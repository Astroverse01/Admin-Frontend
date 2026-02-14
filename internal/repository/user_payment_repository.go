package repository

import (
	"admin-be/internal/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userPaymentRepository struct {
	collection *mongo.Collection
}

func NewUserPaymentRepository(db *database.MongoDB) UserPaymentRepository {
	return &userPaymentRepository{
		collection: db.GetCollection("userPayment"),
	}
}

// SumAmountsByUserIDs returns total amount per userId for the given userIDs (one aggregation query).
func (r *userPaymentRepository) SumAmountsByUserIDs(ctx context.Context, userIDs []string) (map[string]float64, error) {
	if len(userIDs) == 0 {
		return map[string]float64{}, nil
	}
	pipe := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"userId": bson.M{"$in": userIDs}}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$userId",
			"total": bson.M{"$sum": "$amount"},
		}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]float64)
	for _, id := range userIDs {
		result[id] = 0
	}
	for cursor.Next(ctx) {
		var doc struct {
			ID    string  `bson:"_id"`
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		result[doc.ID] = doc.Total
	}
	return result, cursor.Err()
}
