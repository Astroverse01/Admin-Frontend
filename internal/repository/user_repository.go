package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewUserRepository(db *database.MongoDB) UserRepository {
	return &userRepository{
		db:         db,
		collection: db.GetCollection("user"),
	}
}

func (r *userRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.User, error) {
	log.Printf("[UserRepository.FindAll] Called with filter: %v, skip: %d, limit: %d", filter, skip, limit)
	log.Printf("[UserRepository.FindAll] Collection name: %s", r.collection.Name())

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	log.Println("[UserRepository.FindAll] Executing Find query on MongoDB")
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("[UserRepository.FindAll] Error executing Find query: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	log.Println("[UserRepository.FindAll] Cursor obtained, decoding results")
	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("[UserRepository.FindAll] Error decoding cursor results: %v", err)
		return nil, err
	}

	log.Printf("[UserRepository.FindAll] Successfully decoded %d users", len(users))
	return users, nil
}

func (r *userRepository) FindByUserID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateByUserID(ctx context.Context, userID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"userId": userID}, bson.M{"$set": update})
	return err
}

func (r *userRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	log.Printf("[UserRepository.Count] Called with filter: %v", filter)
	log.Printf("[UserRepository.Count] Collection name: %s", r.collection.Name())

	log.Println("[UserRepository.Count] Executing CountDocuments query on MongoDB")
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("[UserRepository.Count] Error executing CountDocuments query: %v", err)
		return 0, err
	}

	log.Printf("[UserRepository.Count] Successfully counted documents: %d", count)
	return count, err
}
