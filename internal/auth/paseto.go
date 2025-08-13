// internal/auth/paseto.go
package auth

import (
	"errors"
	"time"

	"github.com/o1egl/paseto/v2"
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Kita tetap gunakan struct ini untuk kemudahan di bagian lain aplikasi,
// meskipun kita tidak akan menyimpannya langsung ke dalam token.
type PasetoPayload struct {
	UserID primitive.ObjectID `json:"user_id"`
	NIM    string             `json:"nim"`
	Role   string             `json:"role"`
}

// GenerateToken membuat token Paseto baru untuk user (VERSI BARU)
func GenerateToken(userID primitive.ObjectID, nim, role string) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	cfg := config.GetConfig()

	jsonToken := paseto.JSONToken{
		Audience:   "himatif-app",
		Issuer:     "himatif-api",
		Jti:        nim,
		Subject:    userID.Hex(),
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  now,
	}

	// PERBAIKAN DI SINI: Panggil .Set() secara langsung
	jsonToken.Set("role", role)

	key := []byte(cfg.PasetoSecretKey)
	return paseto.NewV2().Encrypt(key, jsonToken, nil)
}

// VerifyToken memverifikasi token Paseto (VERSI BARU)
func VerifyToken(tokenString string) (*PasetoPayload, error) {
	cfg := config.GetConfig()
	key := []byte(cfg.PasetoSecretKey)

	var jsonToken paseto.JSONToken
	err := paseto.NewV2().Decrypt(tokenString, key, &jsonToken, nil)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Validasi klaim waktu
	if err := jsonToken.Validate(); err != nil {
		return nil, err
	}

	// Ambil 'role' dari klaim custom
	var role string
	if err := jsonToken.Get("role", &role); err != nil {
		return nil, errors.New("token missing role claim")
	}

	// Ambil UserID dari 'Subject' dan konversi kembali ke ObjectID
	userID, err := primitive.ObjectIDFromHex(jsonToken.Subject)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	// Buat ulang struct PasetoPayload untuk dikembalikan agar
	// bagian lain dari aplikasi tidak perlu diubah.
	payload := &PasetoPayload{
		UserID: userID,
		NIM:    jsonToken.Jti,
		Role:   role,
	}

	return payload, nil
}
