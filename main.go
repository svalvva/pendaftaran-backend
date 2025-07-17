package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/svalvva/pendaftaran-backend/database" // <-- GANTI INI
)

func main() {
	// Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	// Hubungkan ke database
	database.ConnectDB()

	// Inisialisasi router Gin, jika belum ada
	// Jalankan: go get github.com/gin-gonic/gin
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server backend pendaftaran HIMATIF berjalan!",
		})
	})

	// Menjalankan server di port 8080
	router.Run(":8080")
}