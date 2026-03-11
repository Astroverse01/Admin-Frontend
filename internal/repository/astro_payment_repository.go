package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type astroPaymentRepository struct {
	collection *mongo.Collection
}

func NewAstroPaymentRepository(db *database.MongoDB) AstroPaymentRepository {
	return &astroPaymentRepository{
		collection: db.GetCollection("astroPayment"),
	}
}

func (r *astroPaymentRepository) Create(ctx context.Context, payment *models.AstroPayment) error {
	_, err := r.collection.InsertOne(ctx, payment)
	return err
}

