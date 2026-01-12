package repository

import (
	"astro-admin/internal/database"
	"astro-admin/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type astroProblemRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewAstroProblemRepository(db *database.MongoDB) AstroProblemRepository {
	return &astroProblemRepository{
		db:         db,
		collection: db.GetCollection("astroProblems"),
	}
}

func (r *astroProblemRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.AstroProblem, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var problems []models.AstroProblem
	if err = cursor.All(ctx, &problems); err != nil {
		return nil, err
	}

	return problems, nil
}

func (r *astroProblemRepository) FindByProblemID(ctx context.Context, problemID string) (*models.AstroProblem, error) {
	var problem models.AstroProblem
	err := r.collection.FindOne(ctx, bson.M{"problemId": problemID}).Decode(&problem)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *astroProblemRepository) UpdateByProblemID(ctx context.Context, problemID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"problemId": problemID}, bson.M{"$set": update})
	return err
}

func (r *astroProblemRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}
