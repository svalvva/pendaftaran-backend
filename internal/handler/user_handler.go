// internal/handler/user_handler.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ulbithebest/BE-pendaftaran/internal/auth"
	"github.com/ulbithebest/BE-pendaftaran/internal/config" // <-- PASTIKAN CONFIG DI-IMPORT
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/model"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user" 

	// PERBAIKAN: Gunakan nama database dari config, bukan hardcode
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("users")

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		NPM      string `json:"NPM"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	var user model.User
	// PERBAIKAN: Gunakan nama database dari config, bukan hardcode
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"nim": creds.NPM}).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid NIM or password"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, `{"error": "Invalid NIM or password"}`, http.StatusUnauthorized)
		return
	}

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

func SubmitRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "Failed to get user data from token"}`, http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error": "File size exceeds limit"}`, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("cv")
	if err != nil {
		http.Error(w, `{"error": "CV file is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".pdf" {
		http.Error(w, `{"error": "CV must be a PDF file"}`, http.StatusBadRequest)
		return
	}
	fileName := fmt.Sprintf("%s_%s", payload.UserID.Hex(), handler.Filename)
	filePath := filepath.Join("uploads", fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"error": "Failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, `{"error": "Failed to copy file content"}`, http.StatusInternalServerError)
		return
	}

	registration := model.Registration{
		UserID:        payload.UserID,
		Division:      r.FormValue("division"),
		Motivation:    r.FormValue("motivation"),
		VisionMission: r.FormValue("vision_mission"),
		CVPath:        filePath,
		Status:        "pending",
		Note:          "",
		UpdatedAt:     primitive.NewDateTimeFromTime(time.Now()),
	}

	// PERBAIKAN: Gunakan nama database dari config, bukan hardcode
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")
	_, err = collection.InsertOne(context.TODO(), registration)
	if err != nil {
		http.Error(w, `{"error": "Failed to submit registration"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration submitted successfully"})
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "Failed to get user data from token"}`, http.StatusInternalServerError)
		return
	}

	var user model.User
	// PERBAIKAN: Gunakan nama database dari config, bukan hardcode
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"_id": payload.UserID}).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserRegistrationHandler mengambil detail pendaftaran milik user yang sedang login
func GetUserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "User data not found in token"}`, http.StatusInternalServerError)
		return
	}

	var registration model.Registration
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")
	
	// Cari pendaftaran berdasarkan user_id dari token
	err := collection.FindOne(context.TODO(), bson.M{"user_id": payload.UserID}).Decode(&registration)
	if err != nil {
		// Jika tidak ditemukan, itu bukan error. Kirim respons kosong.
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			return
		}
		http.Error(w, `{"error": "Failed to fetch registration data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registration)
}