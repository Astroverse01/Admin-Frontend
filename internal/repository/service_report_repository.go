package repository

import (
	"admin-be/internal/database"
	"admin-be/internal/models"
	"context"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type serviceReportRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewServiceReportRepository(db *database.MongoDB) ServiceReportRepository {
	collection := db.GetCollection("serviceReports")
	log.Printf("[Repository] Using collection: %s", collection.Name())
	return &serviceReportRepository{
		db:         db,
		collection: collection,
	}
}

func (r *serviceReportRepository) FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.ServiceReport, error) {
	// Projection: Only fetch required fields for efficiency
	projection := bson.M{
		"serviceType": 1,
		"orderId":     1,
		"astroId":     1,
		"userId":      1,
		"status":      1,
		"createdOn":   1,
		"comment":    1,
		"reportId":    1,
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "createdOn", Value: -1}}).
		SetProjection(projection)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []models.ServiceReport
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *serviceReportRepository) FindByReportID(ctx context.Context, reportID string) (*models.ServiceReport, error) {
	var report models.ServiceReport
	err := r.collection.FindOne(ctx, bson.M{"reportId": reportID}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *serviceReportRepository) FindByOrderIDAndServiceType(ctx context.Context, orderId, serviceType string) (*models.ServiceReport, error) {
	var report models.ServiceReport
	// Use case-insensitive regex match for serviceType to handle case variations
	// Escape special regex characters to prevent issues
	escapedServiceType := regexp.QuoteMeta(serviceType)
	filter := bson.M{
		"orderId": orderId,
		"serviceType": bson.M{
			"$regex":   "^" + escapedServiceType + "$",
			"$options": "i", // case-insensitive
		},
	}
	
	// Debug: Log the filter being used
	log.Printf("[Repository] Searching with filter: orderId=%s, serviceType regex=^%s$ (case-insensitive)", orderId, escapedServiceType)
	
	err := r.collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		log.Printf("[Repository] FindOne error: %v", err)
		return nil, err
	}
	log.Printf("[Repository] Found report: reportID=%s, orderID=%s, serviceType=%s", report.ReportID, report.OrderID, report.ServiceType)
	return &report, nil
}

func (r *serviceReportRepository) FindByOrderID(ctx context.Context, orderId string) (*models.ServiceReport, error) {
	var report models.ServiceReport
	err := r.collection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&report)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *serviceReportRepository) UpdateByReportID(ctx context.Context, reportID string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"reportId": reportID}, bson.M{"$set": update})
	return err
}

func (r *serviceReportRepository) UpdateByOrderID(ctx context.Context, orderId string, update map[string]interface{}) error {
	update["updatedAt"] = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"orderId": orderId}, bson.M{"$set": update})
	return err
}

func (r *serviceReportRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

