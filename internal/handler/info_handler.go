// pendaftaran-backend/internal/handler/info_handler.go

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ulbithebest/BE-pendaftaran/internal/config"
	"github.com/ulbithebest/BE-pendaftaran/internal/model"
	"github.com/ulbithebest/BE-pendaftaran/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var infoCollection = "informations"

// CreateInfoHandler (Admin only)
func CreateInfoHandler(w http.ResponseWriter, r *http.Request) {
	var info model.Information
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	info.ID = primitive.NewObjectID()
	info.CreatedAt = now
	info.UpdatedAt = now

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection(infoCollection)
	_, err := collection.InsertOne(context.TODO(), info)
	if err != nil {
		http.Error(w, `{"error": "Failed to create information"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
}

// GetAllInfoHandler (Untuk semua user yang login)
func GetAllInfoHandler(w http.ResponseWriter, r *http.Request) {
	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection(infoCollection)
	
	// Urutkan berdasarkan yang terbaru
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch information"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var results []model.Information
	if err = cursor.All(context.TODO(), &results); err != nil {
		http.Error(w, `{"error": "Failed to decode information"}`, http.StatusInternalServerError)
		return
	}
	if results == nil {
		results = []model.Information{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// UpdateInfoHandler (Admin only)
func UpdateInfoHandler(w http.ResponseWriter, r *http.Request) {
	infoID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid information ID"}`, http.StatusBadRequest)
		return
	}

	var payload struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":      payload.Title,
			"content":    payload.Content,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection(infoCollection)
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": infoID}, update)
	if err != nil {
		http.Error(w, `{"error": "Failed to update information"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Information updated successfully"})
}

// DeleteInfoHandler (Admin only)
func DeleteInfoHandler(w http.ResponseWriter, r *http.Request) {
	infoID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid information ID"}`, http.StatusBadRequest)
		return
	}

	collection := repository.MongoClient.Database(config.GetConfig().DatabaseName).Collection(infoCollection)
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": infoID})
	if err != nil {
		http.Error(w, `{"error": "Failed to delete information"}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, `{"error": "Information not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Information deleted successfully"})
}