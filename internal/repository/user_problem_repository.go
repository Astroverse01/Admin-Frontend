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

type userProblemRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewUserProblemRepository(db *database.MongoDB) UserProblemRepository {
	return &userProblemRepository{
		db:         db,
		collection: db.GetCollection("userProblem"),
	}
}

func (r *userProblemRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.UserProblem, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var problems []models.UserProblem
	if err = cursor.All(ctx, &problems); err != nil {
		return nil, err
	}

	return problems, nil
}

func (r *userProblemRepository) FindByProblemID(ctx context.Context, problemID string) (*models.UserProblem, error) {
	var problem models.UserProblem
	err := r.collection.FindOne(ctx, bson.M{"problemId": problemID}).Decode(&problem)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *userProblemRepository) UpdateByProblemID(ctx context.Context, problemID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"problemId": problemID}, bson.M{"$set": update})
	return err
}

func (r *userProblemRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

