// internal/handler/user_handler.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ulbithebest/BE-pendaftaran/internal/auth"
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
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

	if len(user.PhoneNumber) < 10 || len(user.PhoneNumber) > 13 {
		http.Error(w, `{"error": "Nomor telepon harus antara 10 hingga 13 digit."}`, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user"

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("users")

	// Cek duplikasi NIM
	count, _ := collection.CountDocuments(context.TODO(), bson.M{"nim": user.NIM})
	if count > 0 {
		http.Error(w, `{"error": "NIM already registered"}`, http.StatusConflict)
		return
	}

	// Simpan user ke database
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
	// 1. Dapatkan data user dari token
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "User data not found"}`, http.StatusInternalServerError)
		return
	}

	// 2. Parse form (maks 10MB)
	if err := r.ParseMultipartForm(20 << 20); err != nil { // Batas total 20MB
		http.Error(w, `{"error": "Request size exceeds limit"}`, http.StatusBadRequest)
		return
	}

	division1 := r.FormValue("division1")
	division2 := r.FormValue("division2")
	if division1 == division2 {
		http.Error(w, `{"error": "Pilihan Divisi 1 dan 2 tidak boleh sama."}`, http.StatusBadRequest)
		return
	}

	// 3. Setup koneksi ke Cloudinary
	cfg := config.GetConfig()
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryCloudName, cfg.CloudinaryApiKey, cfg.CloudinaryApiSecret)
	if err != nil {
		http.Error(w, `{"error": "Failed to connect to Cloudinary"}`, http.StatusInternalServerError)
		return
	}
	ctx := context.Background()

	// --- LOGIKA UPLOAD CV (BANYAK FILE) ---
	cvFiles := r.MultipartForm.File["cv"]
	if len(cvFiles) == 0 {
		http.Error(w, `{"error": "CV file is required"}`, http.StatusBadRequest)
		return
	}

	var cvUrls []string
	for i, fileHeader := range cvFiles {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, `{"error": "Failed to open CV file"}`, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		publicID := fmt.Sprintf("himatif-registrations/%s_cv_%d_%d", payload.NIM, time.Now().Unix(), i)
		uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
			PublicID:     publicID,
			ResourceType: "image",
			Overwrite:    api.Bool(true),
		})
		if err != nil {
			log.Printf("Cloudinary CV upload error: %v", err)
			http.Error(w, `{"error": "Failed to upload one of the CVs"}`, http.StatusInternalServerError)
			return
		}
		cvUrls = append(cvUrls, uploadResult.SecureURL)
	}

	// --- LOGIKA UPLOAD SERTIFIKAT (BANYAK FILE, OPSIONAL) ---
	certFiles := r.MultipartForm.File["certificate"]
	var certificateUrls []string
	if len(certFiles) > 0 {
		for i, fileHeader := range certFiles {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, `{"error": "Failed to open certificate file"}`, http.StatusInternalServerError)
				return
			}
			defer file.Close()

			publicID := fmt.Sprintf("himatif-registrations/%s_cert_%d_%d", payload.NIM, time.Now().Unix(), i)
			uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
				PublicID:     publicID,
				ResourceType: "image",
				Overwrite:    api.Bool(true),
			})
			if err != nil {
				log.Printf("Cloudinary certificate upload error: %v", err)
				http.Error(w, `{"error": "Failed to upload one of the certificates"}`, http.StatusInternalServerError)
				return
			}
			certificateUrls = append(certificateUrls, uploadResult.SecureURL)
		}
	}

	// Simpan ke database
	registration := model.Registration{
		UserID:          payload.UserID,
		Division1:       division1,
		Division2:       division2,
		Motivation:      r.FormValue("motivation"),
		VisionMission:   r.FormValue("vision_mission"),
		CvUrls:          cvUrls,
		CertificateUrls: certificateUrls,
		Status:          "pending",
		Note:            "",
		UpdatedAt:       primitive.NewDateTimeFromTime(time.Now()),
	}

	collection := repository.MongoClient.Database(cfg.DatabaseName).Collection("registrations")
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

func GetUserRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetPayloadFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "User data not found in token"}`, http.StatusInternalServerError)
		return
	}

	var registration model.Registration
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")

	err := collection.FindOne(context.TODO(), bson.M{"user_id": payload.UserID}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, `{"error": "Failed to fetch registration data"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registration)
}