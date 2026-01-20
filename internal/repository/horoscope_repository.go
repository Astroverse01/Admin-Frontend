package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type horoscopeRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewHoroscopeRepository(db *database.MongoDB) HoroscopeRepository {
	return &horoscopeRepository{
		db:         db,
		collection: db.GetCollection("horoscope"),
	}
}

func (r *horoscopeRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.Horoscope, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "CreatedOn", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var horoscopes []models.Horoscope
	if err = cursor.All(ctx, &horoscopes); err != nil {
		return nil, err
	}

	return horoscopes, nil
}

func (r *horoscopeRepository) FindByHoroscopeID(ctx context.Context, horoscopeID string) (*models.Horoscope, error) {
	var horoscope models.Horoscope
	err := r.collection.FindOne(ctx, bson.M{"horoscopeId": horoscopeID}).Decode(&horoscope)
	if err != nil {
		return nil, err
	}
	return &horoscope, nil
}

func (r *horoscopeRepository) FindBySignNameAndDate(ctx context.Context, signName, date string) (*models.Horoscope, error) {
	var horoscope models.Horoscope
	err := r.collection.FindOne(ctx, bson.M{"signName": signName, "date": date}).Decode(&horoscope)
	if err != nil {
		return nil, err
	}
	return &horoscope, nil
}

func (r *horoscopeRepository) Create(ctx context.Context, horoscope *models.Horoscope) error {
	_, err := r.collection.InsertOne(ctx, horoscope)
	return err
}

func (r *horoscopeRepository) UpdateByHoroscopeID(ctx context.Context, horoscopeID string, update map[string]interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"horoscopeId": horoscopeID}, bson.M{"$set": update})
	return err
}

func (r *horoscopeRepository) DeleteByHoroscopeID(ctx context.Context, horoscopeID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"horoscopeId": horoscopeID})
	return err
}

func (r *horoscopeRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}