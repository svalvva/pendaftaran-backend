package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua variabel konfigurasi aplikasi
type Config struct {
	MongoURI      string
	DatabaseName  string
	PasetoSecretKey string
	ServerPort    string
	CloudinaryCloudName string
    CloudinaryApiKey    string
    CloudinaryApiSecret string

}

var appConfig *Config

// LoadConfig memuat konfigurasi dari file .env
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	appConfig = &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		DatabaseName:  os.Getenv("MONGO_DATABASE"),
		PasetoSecretKey: os.Getenv("PASETO_SECRET_KEY"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
    	CloudinaryApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
    	CloudinaryApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
    log.Println("Configuration loaded successfully")
}

// GetConfig mengembalikan instance konfigurasi yang sudah dimuat
func GetConfig() *Config {
	if appConfig == nil {
		LoadConfig()
	}
	return appConfig
}