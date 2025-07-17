package main

import (
	"log"

	"github.com/gin-contrib/cors" // Pastikan import ini ada
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/svalvva/pendaftaran-backend/database"
	"github.com/svalvva/pendaftaran-backend/handlers"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	// Hubungkan ke database
	database.ConnectDB()

	// Inisialisasi router Gin
	router := gin.Default()

	// ----------------------------------------------------
	// PENTING: Konfigurasi CORS diletakkan di sini,
	// TEPAT SETELAH router dibuat dan SEBELUM route didefinisikan.
	// ----------------------------------------------------
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5500", "http://127.0.0.1:5500"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
	// ----------------------------------------------------


	// Definisikan route Anda seperti biasa
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server backend pendaftaran HIMATIF berjalan!",
		})
	})

	// Membuat grup untuk API
	api := router.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}

	// Menjalankan server di port 8080
	log.Println("Listening and serving HTTP on :8080")
	router.Run(":8080")
}