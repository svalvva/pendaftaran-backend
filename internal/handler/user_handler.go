package handler

import (
	// Pustaka Standar Go
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// Pustaka Pihak Ketiga
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	// Paket Internal Proyek Anda (sesuaikan path jika perlu)
	"github.com/ulbithebest/BE-pendaftaran/internal/auth"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/model"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"
)

// RegisterHandler menangani pendaftaran user baru
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user" // Role default

	collection := repository.MongoClient.Database("himatif_db").Collection("users")

	// Cek apakah NIM sudah ada
	count, _ := collection.CountDocuments(context.TODO(), bson.M{"nim": user.NIM})
	if count > 0 {
		http.Error(w, `{"error": "NIM already registered"}`, http.StatusConflict)
		return
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, `{"error": "Failed to register user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}

// LoginHandler menangani login user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		NPM      string `json:"NPM"` // Sesuaikan dengan field name di login.html
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	var user model.User
	collection := repository.MongoClient.Database("himatif_db").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"nim": creds.NPM}).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid NIM or password"}`, http.StatusUnauthorized)
		return
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, `{"error": "Invalid NIM or password"}`, http.StatusUnauthorized)
		return
	}

	// Generate Paseto token
	token, err := auth.GenerateToken(user.ID, user.NIM, user.Role)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"role":  user.Role,
	})
}

// SubmitRegistrationHandler menangani pengiriman formulir dari form.html
func SubmitRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil payload dari token untuk mengetahui siapa yang mengirim
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "Failed to get user data from token"}`, http.StatusInternalServerError)
		return
	}

	// 2. Batasi ukuran request body (misal 10MB) untuk keamanan
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, `{"error": "File size exceeds limit"}`, http.StatusBadRequest)
		return
	}

	// 3. Ambil file CV dari form
	file, handler, err := r.FormFile("cv")
	if err != nil {
		log.Printf("Error retrieving CV file: %v", err)
		http.Error(w, `{"error": "CV file is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. Buat path untuk menyimpan file di folder ./uploads
	// Nama file: <user_id>_<nama_file_asli>.pdf
	ext := filepath.Ext(handler.Filename)
	if ext != ".pdf" {
		http.Error(w, `{"error": "CV must be a PDF file"}`, http.StatusBadRequest)
		return
	}
	fileName := fmt.Sprintf("%s_%s", payload.UserID.Hex(), handler.Filename)
	filePath := filepath.Join("uploads", fileName)

	// Buat file tujuan di server
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"error": "Failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Salin konten file yang di-upload ke file tujuan
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, `{"error": "Failed to copy file content"}`, http.StatusInternalServerError)
		return
	}

	// 5. Buat entri pendaftaran di database
	registration := model.Registration{
		UserID:     payload.UserID,
		Division:   r.FormValue("division"),
		Motivation: r.FormValue("motivation"),
		CVPath:     filePath, // Simpan path-nya, bukan file-nya
		Status:     "pending", // Status awal
		Note:       "",
		UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}

	collection := repository.MongoClient.Database("himatif_db").Collection("registrations")
	_, err = collection.InsertOne(context.TODO(), registration)
	if err != nil {
		log.Printf("Failed to insert registration: %v", err)
		http.Error(w, `{"error": "Failed to submit registration"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration submitted successfully"})
}

// GetUserProfileHandler untuk mengambil data user yang sedang login
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "Failed to get user data from token"}`, http.StatusInternalServerError)
		return
	}

	var user model.User
	collection := repository.MongoClient.Database("himatif_db").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"_id": payload.UserID}).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	// Jangan kirim password ke frontend
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetAllRegistrationsHandler untuk admin melihat semua data pendaftaran
func GetAllRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	var registrations []model.Registration

	collection := repository.MongoClient.Database("himatif_db").Collection("registrations")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch registrations"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &registrations); err != nil {
		http.Error(w, `{"error": "Failed to decode registrations"}`, http.StatusInternalServerError)
		return
	}

	if registrations == nil {
		registrations = []model.Registration{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}