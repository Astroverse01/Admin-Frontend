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

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(db *database.MongoDB) BlogRepository {
	return &blogRepository{
		collection: db.GetCollection("blogs"),
	}
}

func (r *blogRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.Blog, error) {
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdOn", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog
	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}

func (r *blogRepository) FindByBlogID(ctx context.Context, blogID string) (*models.Blog, error) {
	var blog models.Blog
	if err := r.collection.FindOne(ctx, bson.M{"blogId": blogID}).Decode(&blog); err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) FindBySlug(ctx context.Context, slug string) (*models.Blog, error) {
	var blog models.Blog
	if err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&blog); err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *blogRepository) Create(ctx context.Context, blog *models.Blog) error {
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

func (r *blogRepository) UpdateByBlogID(ctx context.Context, blogID string, update map[string]interface{}) error {
	update["updatedOn"] = time.Now().UTC()
	_, err := r.collection.UpdateOne(ctx, bson.M{"blogId": blogID}, bson.M{"$set": update})
	return err
}

func (r *blogRepository) DeleteByBlogID(ctx context.Context, blogID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"blogId": blogID})
	return err
}

func (r *blogRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

