package main

import (
	"log"
	"net/http"
	// "strings"

	// Pastikan path import ini sesuai dengan nama modul di go.mod Anda
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/handler"
	"github.com/ulbithebest/BE-pendaftaran/internal/middleware"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// 1. Muat Konfigurasi dari file .env
	cfg := config.GetConfig()

	// 2. Hubungkan ke Database MongoDB
	repository.ConnectDB(cfg)

	// 3. Inisialisasi Router menggunakan Chi
	r := chi.NewRouter()

	// 4. Setup Middleware Global
	r.Use(chiMiddleware.Logger)    // Middleware untuk mencatat (log) setiap request yang masuk
	r.Use(chiMiddleware.Recoverer) // Middleware untuk menangani panic dan menjaga server tetap hidup

	// Setup CORS (Cross-Origin Resource Sharing)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5500", "http://127.0.0.1:5500"}, // Sesuaikan dengan alamat frontend Anda
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}, // TAMBAHKAN PATCH
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 5. Definisikan Routes (Endpoint API)

	// Routes publik yang bisa diakses tanpa login/token
	r.Post("/register", handler.RegisterHandler)
	r.Post("/login", handler.LoginHandler)

	// Group routes yang memerlukan otentikasi (wajib ada token Paseto)
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Semua di dalam grup ini akan dilindungi oleh middleware otentikasi

		// --- Routes untuk user biasa ---
		r.Get("/user/profile", handler.GetUserProfileHandler)
		r.Post("/user/registration", handler.SubmitRegistrationHandler)
		r.Get("/user/my-registration", handler.GetUserRegistrationHandler)

		// --- TAMBAHAN: File Server untuk CV (dilindungi otentikasi) ---
		// Ini akan membuat file di folder /uploads bisa diakses via URL
		fileServer := http.FileServer(http.Dir("./uploads"))
		r.Handle("/uploads/*", http.StripPrefix("/api/uploads/", fileServer))

		// --- Routes khusus admin ---
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminOnlyMiddleware) // Perlindungan tambahan, hanya untuk admin
			
			// --- PERUBAHAN ENDPOINT ADMIN ---
			r.Get("/registrations-with-details", handler.GetAllRegistrationsDetailHandler)
			// r.Patch("/registrations/{id}/status", handler.UpdateRegistrationStatusHandler)
			r.Patch("/registrations/{id}", handler.UpdateRegistrationDetailsHandler)
		})
	})

	// 6. Jalankan Server HTTP
	log.Printf("✅ Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, r); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}