package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type astroRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewAstroRepository(db *database.MongoDB) AstroRepository {
	return &astroRepository{
		db:         db,
		collection: db.GetCollection("astrologer"),
	}
}

func (r *astroRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.Astrologer, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdOn", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var astros []models.Astrologer
	if err = cursor.All(ctx, &astros); err != nil {
		return nil, err
	}

	return astros, nil
}

func (r *astroRepository) FindByAstroID(ctx context.Context, astroID string) (*models.Astrologer, error) {
	var astro models.Astrologer
	err := r.collection.FindOne(ctx, bson.M{"astroId": astroID}).Decode(&astro)
	if err != nil {
		return nil, err
	}
	return &astro, nil
}

func (r *astroRepository) UpdateByAstroID(ctx context.Context, astroID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"astroId": astroID}, bson.M{"$set": update})
	return err
}

func (r *astroRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}
