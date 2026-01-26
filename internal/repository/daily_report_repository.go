package repository

import (
	"admin-be/internal/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DailyReportRepository interface {
	GetTodayRecords(ctx context.Context, collectionName string, startTime, endTime time.Time) ([]bson.M, error)
}

type dailyReportRepository struct {
	db *database.MongoDB
}

func NewDailyReportRepository(db *database.MongoDB) DailyReportRepository {
	return &dailyReportRepository{
		db: db,
	}
}

// GetTodayRecords fetches all records from a collection where createdOn is between startTime and endTime
func (r *dailyReportRepository) GetTodayRecords(ctx context.Context, collectionName string, startTime, endTime time.Time) ([]bson.M, error) {
	collection := r.db.GetCollection(collectionName)

	// Filter records where createdOn is between startTime and endTime
	filter := bson.M{
		"createdOn": bson.M{
			"$gte": startTime,
			"$lt":  endTime,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []bson.M{}, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
