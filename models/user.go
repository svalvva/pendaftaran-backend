package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    NIM         string             `json:"nim" bson:"nim"`
    Email       string             `json:"email" bson:"email"`
    Password    string             `json:"password" bson:"password"`
    Role        string             `json:"role" bson:"role"`
    // TAMBAHKAN FIELD BARU DI SINI
    BirthPlace  string             `json:"birth_place,omitempty" bson:"birth_place,omitempty"`
    BirthDate   string             `json:"birth_date,omitempty" bson:"birth_date,omitempty"`
}