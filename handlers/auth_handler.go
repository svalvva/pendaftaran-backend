package handlers

import (
	"context"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/svalvva/pendaftaran-backend/database"
	"github.com/svalvva/pendaftaran-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register adalah handler untuk mendaftarkan user baru
func Register(c *gin.Context) {
	var user models.User
	// 1. Bind JSON yang masuk ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	// 2. Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user" // Atur role default

	// 3. Simpan user ke database
	collection := database.DB.Collection("users")
	// Cek apakah NIM atau Email sudah ada
	existingUser := collection.FindOne(context.TODO(), bson.M{
		"$or": []bson.M{
			{"nim": user.NIM},
			{"email": user.Email},
		},
	})
	if existingUser.Err() == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "NIM atau Email sudah terdaftar"})
		return
	}

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan user"})
		return
	}

	// 4. Kirim respons sukses
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
        // MATA-MATA 1: Cari tahu kenapa gagal mencari user
        fmt.Printf("Gagal mencari user dengan NIM %s. Error: %v\n", payload.NIM, err) 
        
        c.JSON(http.StatusUnauthorized, gin.H{"error": "NIM atau password salah"})
        return
    }

    // MATA-MATA 2: Cek apakah user ditemukan dan apa hash passwordnya
    fmt.Printf("User ditemukan! Hash dari DB: %s\n", user.Password)


    // 3. Bandingkan password yang di-hash
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
    if err != nil {
        // MATA-MATA 3: Cari tahu kenapa password gagal dicocokkan
        fmt.Printf("Gagal mencocokkan password. Error: %v\n", err)
        
        c.JSON(http.StatusUnauthorized, gin.H{"error": "NIM atau password salah"})
        return
    }


	// ---- BAGIAN BARU: MEMBUAT TOKEN JWT ----
	// 4. Buat token JWT
	// Semua paket yang tadinya "not used" sekarang kita pakai di sini
	claims := jwt.MapClaims{
		"nim":  user.NIM,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Dapatkan secret key dari environment variable (butuh paket "os")
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "RAHASIA_JANGAN_DITARUH_DISINI_ASLINYA" // Fallback jika tidak ada di env
	}
	
	// Tanda tangani token dengan secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}
	// ---------------------------------------------

	// 5. Kirim token sebagai respons
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
	})
}