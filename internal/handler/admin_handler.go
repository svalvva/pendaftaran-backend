// internal/handler/admin_handler.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/model" // <-- PERBAIKAN 1: Tambahkan import model
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllRegistrationsDetailHandler(w http.ResponseWriter, r *http.Request) {
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "userDetails"},
		}}},
		bson.D{{Key: "$unwind", Value: "$userDetails"}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1}, {Key: "user_id", Value: 1}, {Key: "division", Value: 1},
			{Key: "motivation", Value: 1}, {Key: "vision_mission", Value: 1}, {Key: "cv_path", Value: 1},
			{Key: "status", Value: 1}, {Key: "note", Value: 1}, {Key: "updated_at", Value: 1},
			{Key: "name", Value: "$userDetails.name"}, {Key: "nim", Value: "$userDetails.nim"},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch registrations"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var results []model.RegistrationDetail // <-- PERBAIKAN 2: Gunakan model.RegistrationDetail
	if err = cursor.All(context.TODO(), &results); err != nil {
		http.Error(w, `{"error": "Failed to decode registrations"}`, http.StatusInternalServerError)
		return
	}

	if results == nil {
		results = []model.RegistrationDetail{} // <-- PERBAIKAN 3: Gunakan model.RegistrationDetail
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func UpdateRegistrationDetailsHandler(w http.ResponseWriter, r *http.Request) {
	regID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid registration ID"}`, http.StatusBadRequest)
		return
	}

	var payload model.Registration 
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")
	updateFields := bson.M{}

	// Hanya update field yang dikirim oleh admin
	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}
	if payload.InterviewSchedule != "" {
		updateFields["interview_schedule"] = payload.InterviewSchedule
	}
    
    // --- PERBAIKAN DI SINI: Tambahkan logika untuk lokasi wawancara ---
	if payload.InterviewLocation != "" {
		updateFields["interview_location"] = payload.InterviewLocation
	}
    
	if payload.Division != "" {
		updateFields["division"] = payload.Division
	}
	updateFields["updated_at"] = primitive.NewDateTimeFromTime(time.Now())
	
	update := bson.M{"$set": updateFields}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": regID}, update)
	if err != nil {
		http.Error(w, `{"error": "Failed to update registration"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration updated successfully"})
}