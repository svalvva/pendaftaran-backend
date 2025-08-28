package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User sesuai dengan koleksi 'users'
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	NIM        string             `bson:"nim" json:"nim"`
	BirthPlace string             `bson:"birth_place" json:"birth_place"` // <-- TAMBAHKAN INI
	BirthDate  string             `bson:"birth_date" json:"birth_date"`   // <-- TAMBAHKAN INI
	Email      string             `bson:"email" json:"email"`
	PhoneNumber string            `bson:"phone_number" json:"phone_number"`
	Password   string             `bson:"password" json:"password"`
	Role       string             `bson:"role" json:"role"`
}

// Registration sesuai dengan koleksi 'registrations'
type Registration struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Division1        string             `bson:"division1" json:"division1"`           // <-- TAMBAHKAN INI
	Division2        string             `bson:"division2" json:"division2"`  
	Motivation    string             `bson:"motivation" json:"motivation"`
	VisionMission string             `bson:"vision_mission" json:"vision_mission"`
	InterviewSchedule string         `bson:"interview_schedule,omitempty" json:"interview_schedule,omitempty"`
	InterviewLocation string         `bson:"interview_location,omitempty" json:"interview_location,omitempty"`
	CvUrl            string          `bson:"cv_url" json:"cv_url"`                                       // <-- UBAH INI
	CertificateUrl   string          `bson:"certificate_url,omitempty" json:"certificate_url,omitempty"`
	Status        string             `bson:"status" json:"status"`
	Note          string             `bson:"note" json:"note"`
	UpdatedAt     primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

// Struct untuk menggabungkan data Registrasi dan User
type RegistrationDetail struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name          string             `bson:"name" json:"name"` 
	NIM           string             `bson:"nim" json:"nim"`
	Division1        string             `bson:"division1" json:"division1"`           // <-- TAMBAHKAN INI
	Division2        string             `bson:"division2" json:"division2"`  
	Motivation    string             `bson:"motivation" json:"motivation"`
	VisionMission string             `bson:"vision_mission" json:"vision_mission"`
	CvUrl            string          `bson:"cv_url" json:"cv_url"`                                       // <-- UBAH INI
	CertificateUrl   string          `bson:"certificate_url,omitempty" json:"certificate_url,omitempty"`
	Status        string             `bson:"status" json:"status"`
	Note          string             `bson:"note,omitempty" json:"note,omitempty"`
	UpdatedAt     primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

// Information sesuai dengan koleksi 'informations'
type Information struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}