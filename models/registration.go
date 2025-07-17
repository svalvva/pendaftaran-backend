package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Registration merepresentasikan data pendaftaran di collection 'registrations'
type Registration struct {
    ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
    Division   string             `json:"division" bson:"division"`
    Motivation string             `json:"motivation" bson:"motivation"`
    CVPath     string             `json:"cv_path" bson:"cv_path"`
    Status     string             `json:"status" bson:"status"` // "menunggu", "lulus", "tidak_lulus"
    Note       string             `json:"note,omitempty" bson:"note,omitempty"` // Catatan dari admin
    UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}