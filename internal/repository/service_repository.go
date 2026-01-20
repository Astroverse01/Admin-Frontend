package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type serviceRepository struct {
	db *database.MongoDB
}

func NewServiceRepository(db *database.MongoDB) ServiceRepository {
	return &serviceRepository{
		db: db,
	}
}

func (r *serviceRepository) FindChatByChatID(ctx context.Context, chatID string) (*models.Chat, error) {
	var chat models.Chat
	collection := r.db.GetCollection("chat")

	// Debug: Check collection name and count documents
	count, _ := collection.CountDocuments(ctx, bson.M{})
	log.Printf("Repository: Querying 'chat' collection (total docs: %d) with chatId: %s", count, chatID)

	// Debug: Try to find any document with this chatId to verify the query
	filter := bson.M{"chatId": chatID}
	log.Printf("Repository: Using filter: %+v", filter)

	// Try to find document first without decoding to see if it exists
	var rawDoc bson.M
	findErr := collection.FindOne(ctx, filter).Decode(&rawDoc)
	if findErr == nil {
		log.Printf("Repository: Document found! Keys in document: %v", getKeys(rawDoc))
		// Now try to decode into Chat struct
		err := collection.FindOne(ctx, filter).Decode(&chat)
		if err != nil {
			log.Printf("Repository: Document exists but decode failed: %v", err)
			return nil, err
		}
		log.Printf("Repository: Successfully decoded chat with chatId: %s", chat.ChatId)
		return &chat, nil
	}

	log.Printf("Repository: Error finding chat with chatId %s: %v (error type: %T)", chatID, findErr, findErr)
	return nil, findErr
}

// Helper function to get keys from a map
func getKeys(m bson.M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (r *serviceRepository) FindIvrByIvrID(ctx context.Context, ivrID string) (*models.IvrCall, error) {
	var ivr models.IvrCall
	err := r.db.GetCollection("ivrCall").FindOne(ctx, bson.M{"ivrId": ivrID}).Decode(&ivr)
	if err != nil {
		return nil, err
	}
	return &ivr, nil
}

func (r *serviceRepository) FindVideoByVideoID(ctx context.Context, videoID string) (*models.VideoCall, error) {
	var video models.VideoCall
	err := r.db.GetCollection("videoCall").FindOne(ctx, bson.M{"videoId": videoID}).Decode(&video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *serviceRepository) UpdateChatByChatID(ctx context.Context, chatID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.db.GetCollection("chat").UpdateOne(ctx, bson.M{"chatId": chatID}, bson.M{"$set": update})
	return err
}

func (r *serviceRepository) UpdateIvrByIvrID(ctx context.Context, ivrID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.db.GetCollection("ivrCall").UpdateOne(ctx, bson.M{"ivrId": ivrID}, bson.M{"$set": update})
	return err
}

func (r *serviceRepository) UpdateVideoByVideoID(ctx context.Context, videoID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.db.GetCollection("videoCall").UpdateOne(ctx, bson.M{"videoId": videoID}, bson.M{"$set": update})
	return err
}
