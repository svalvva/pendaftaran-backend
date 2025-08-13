package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User sesuai dengan koleksi 'users'
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	NIM      string             `bson:"nim" json:"nim"` // Nomor Induk Mahasiswa
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Role     string             `bson:"role" json:"role"` // "admin" atau "user"
}

// Registration sesuai dengan koleksi 'registrations'
type Registration struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Division  string             `bson:"division" json:"division"`
	Motivation string            `bson:"motivation" json:"motivation"`
	CVPath    string             `bson:"cv_path" json:"cv_path"`
	Status    string             `bson:"status" json:"status"`       // e.g., "pending", "accepted", "rejected"
	Note      string             `bson:"note" json:"note"`           // Catatan dari admin
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}