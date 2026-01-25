package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type feedbackRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewFeedbackRepository(db *database.MongoDB) FeedbackRepository {
	return &feedbackRepository{
		db:         db,
		collection: db.GetCollection("feedback"),
	}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *models.Feedback) error {
	_, err := r.collection.InsertOne(ctx, feedback)
	return err
}

func (r *feedbackRepository) BulkCreate(ctx context.Context, feedbacks []*models.Feedback) error {
	if len(feedbacks) == 0 {
		return nil
	}

	// Convert []*models.Feedback to []interface{} for InsertMany
	docs := make([]interface{}, len(feedbacks))
	for i, feedback := range feedbacks {
		docs[i] = feedback
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

