package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/svalvva/pendaftaran-backend/database"
	"github.com/svalvva/pendaftaran-backend/handlers" // Sekarang akan digunakan
)

func main() { // Inisialisasi server Gin dan koneksi database
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	database.ConnectDB() // Inisialisasi koneksi database
	router := gin.Default() // Inisialisasi router Gin

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server backend pendaftaran HIMATIF berjalan!",
		})
	})

    // --- TAMBAHKAN BLOK INI ---
	// Membuat grup untuk API
	api := router.Group("/api")
	{
		// Menambahkan route untuk registrasi
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login) // Proses login user
	}
    // -------------------------

	// Menjalankan server di port 8080
	router.Run(":8080")
}