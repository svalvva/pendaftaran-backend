package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User merepresentasikan data pengguna di collection 'users'
type User struct {
    ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name     string             `json:"name" bson:"name"`
    NIM      string             `json:"nim" bson:"nim"`
    Email    string             `json:"email" bson:"email"`
    Password string 			`json:"password" bson:"password"` // Tanda - agar password tidak dikirim dalam JSON
    Role     string             `json:"role" bson:"role"`
}