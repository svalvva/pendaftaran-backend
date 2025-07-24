package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User mendefinisikan struktur untuk data pengguna di dalam database.
// Ini adalah satu-satunya model yang Anda butuhkan untuk pengguna.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName  string             `bson:"full_name" json:"full_name" binding:"required"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Password  string             `bson:"password" json:"password" binding:"required,min=6"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}