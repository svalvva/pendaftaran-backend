package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/svalvva/pendaftaran-backend/database" 
	"github.com/svalvva/pendaftaran-backend/models"   
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