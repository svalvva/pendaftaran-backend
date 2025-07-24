package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User sekarang memiliki NIM dan Role
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FullName  string             `bson:"full_name" json:"full_name"` // Hapus binding required agar tidak error saat login
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	NIM       string             `bson:"nim" json:"nim" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required,min=6"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}