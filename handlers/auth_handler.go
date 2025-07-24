package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	// Ganti svalvva/pendaftaran-backend dengan path modul Anda sendiri jika berbeda
	"github.com/svalvva/pendaftaran-backend/database"
	"github.com/svalvva/pendaftaran-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Register adalah handler untuk mendaftarkan user baru
func Register(c *gin.Context) {
	var user models.User
	// 1. Bind JSON yang masuk ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
		return
	}

	// 2. Cek apakah NIM atau Email sudah ada di database
	collection := database.DB.Collection("users")
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"nim": user.NIM},
			{"email": user.Email},
		},
	}).Decode(&existingUser)

	// Jika tidak ada error, berarti data sudah ada (konflik)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "NIM atau Email sudah terdaftar"})
		return
	}
	// Pastikan error yang terjadi adalah karena dokumen tidak ditemukan, bukan error lain
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa database"})
		return
	}

	// 3. Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user" // Atur role default
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 4. Simpan user ke database
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	// 5. Kirim respons sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

// Login adalah handler untuk proses login user
func Login(c *gin.Context) {
	var payload struct {
		NIM      string `json:"nim"`
		Password string `json:"password"`
	}
	var user models.User

	// 1. Bind data login
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// 2. Cari user berdasarkan NIM
	collection := database.DB.Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"nim": payload.NIM}).Decode(&user)
	if err != nil {
		// Jika user tidak ditemukan atau ada error lain
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIM atau password salah"})
		return
	}

	// 3. Bandingkan password yang di-hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		// Jika password tidak cocok
		c.JSON(http.StatusUnauthorized, gin.H{"error": "NIM atau password salah"})
		return
	}

	// 4. Buat token JWT
	claims := jwt.MapClaims{
		"nim":  user.NIM,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "RAHASIA_JANGAN_DITARUH_DISINI_ASLINYA" // Fallback
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 5. Kirim token sebagai respons
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
	})
}