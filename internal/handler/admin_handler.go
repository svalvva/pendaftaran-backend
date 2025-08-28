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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllRegistrationsDetailHandler(w http.ResponseWriter, r *http.Request) {
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"}, {Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"}, {Key: "as", Value: "userDetails"},
		}}},
		bson.D{{Key: "$unwind", Value: "$userDetails"}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1}, {Key: "user_id", Value: 1}, {Key: "division1", Value: 1},
			{Key: "division2", Value: 1}, {Key: "motivation", Value: 1}, {Key: "vision_mission", Value: 1},
			{Key: "cv_url", Value: 1}, {Key: "certificate_url", Value: 1}, {Key: "status", Value: 1},
			{Key: "note", Value: 1}, {Key: "updated_at", Value: 1}, {Key: "name", Value: "$userDetails.name"},
			{Key: "nim", Value: "$userDetails.nim"},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch registrations"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var results []model.RegistrationDetail
	if err = cursor.All(context.TODO(), &results); err != nil {
		http.Error(w, `{"error": "Failed to decode registrations"}`, http.StatusInternalServerError)
		return
	}
	if results == nil {
		results = []model.RegistrationDetail{}
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

	if payload.Status != "" {
		updateFields["status"] = payload.Status
	}
	if payload.InterviewSchedule != "" {
		updateFields["interview_schedule"] = payload.InterviewSchedule
	}
	if payload.InterviewLocation != "" {
		updateFields["interview_location"] = payload.InterviewLocation
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

// FITUR BARU: Handler untuk update status beberapa pendaftar sekaligus
func BulkUpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct {
        IDs    []string `json:"ids"`
        Status string   `json:"status"`
    }

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    if len(payload.IDs) == 0 {
        http.Error(w, `{"error": "No registration IDs provided"}`, http.StatusBadRequest)
        return
    }

    // Konversi string IDs menjadi BSON ObjectIDs
    var objectIDs []primitive.ObjectID
    for _, idStr := range payload.IDs {
        id, err := primitive.ObjectIDFromHex(idStr)
        if err != nil {
            http.Error(w, `{"error": "Invalid ID format in list"}`, http.StatusBadRequest)
            return
        }
        objectIDs = append(objectIDs, id)
    }

    collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")
    
    // Filter untuk mencari semua dokumen dengan ID yang ada di dalam array
    filter := bson.M{"_id": bson.M{"$in": objectIDs}}
    
    // Data yang akan di-update
    update := bson.M{
        "$set": bson.M{
            "status":     payload.Status,
            "updated_at": primitive.NewDateTimeFromTime(time.Now()),
        },
    }

    // Lakukan operasi UpdateMany
    result, err := collection.UpdateMany(context.TODO(), filter, update)
    if err != nil {
        http.Error(w, `{"error": "Failed to bulk update registrations"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":      "Bulk update successful",
        "updatedCount": result.ModifiedCount,
    })
}


func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("users")

	// Opsi untuk tidak menyertakan field password demi keamanan
	opts := options.Find().SetProjection(bson.M{"password": 0})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch users"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []model.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, `{"error": "Failed to decode users"}`, http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []model.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DeleteRegistrationHandler menghapus data pendaftaran berdasarkan ID
func DeleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Mengambil ID dari parameter URL
	regID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid registration ID"}`, http.StatusBadRequest)
		return
	}

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection("registrations")

	// Menghapus satu dokumen yang cocok dengan ID
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": regID})
	if err != nil {
		http.Error(w, `{"error": "Failed to delete registration"}`, http.StatusInternalServerError)
		return
	}

	// Jika tidak ada dokumen yang terhapus (mungkin ID tidak ditemukan)
	if result.DeletedCount == 0 {
		http.Error(w, `{"error": "Registration not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration deleted successfully"})
}