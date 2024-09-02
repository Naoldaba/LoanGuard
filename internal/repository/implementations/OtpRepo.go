package implementations

import (
	"context"

	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/repository/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoOtpRepository struct {
	collection *mongo.Collection
}

func NewMongoOtpRepository(db *mongo.Database) repository_interface.IOtpRepository {
	return &MongoOtpRepository{
		collection: db.Collection("otp"),
	}
}

func (r *MongoOtpRepository) SaveOtp(ctx context.Context, otp models.OtpEntry) error {
	_, err := r.collection.InsertOne(ctx, otp)
	return err
}

func (r *MongoOtpRepository) FindByOtp(ctx context.Context, otp string) (*models.OtpEntry, error) {
	var otpEntry models.OtpEntry
	err := r.collection.FindOne(ctx, bson.M{"otp": otp}).Decode(&otpEntry)
	if err != nil {
		return nil, err
	}
	return &otpEntry, nil
}