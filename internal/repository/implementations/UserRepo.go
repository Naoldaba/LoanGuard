package implementations

import (
	"LoanGuard/internal/domain/dtos"
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/infrastructures/services"
	"LoanGuard/internal/repository/interfaces"
	"errors"
	"time"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection  *mongo.Collection
	redisClient services.ICacheService
}

func NewMongoUserRepository(db *mongo.Database, redisClient services.ICacheService) repository_interface.IUserRepository {
	return &MongoUserRepository{
		collection:  db.Collection("users"),
		redisClient: redisClient,
	}
}

func (r *MongoUserRepository) Register(user *models.User) (*models.User, error) {
	if user.ID == primitive.NilObjectID {
		user.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoUserRepository) BlacklistToken(token string, remainingTime time.Duration) error {
	err := r.redisClient.BlacklistTkn(token, remainingTime)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoUserRepository) GetUserByID(id string) (*models.User, error) {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var user models.User
	err = r.collection.FindOne(context.Background(), bson.M{"_id": user_id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *MongoUserRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	cursor, err := repo.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) DeleteUser(id string) error {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": user_id})
	return err
}

func (r *MongoUserRepository) UpdateUser(id string, user *models.User) error {
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": user_id}, bson.M{"$set": user})
	return err
}

func (r *MongoUserRepository) PromoteUser(userID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": user_id}, bson.M{"$set": bson.M{"role": "admin"}})
	return err
}

func (r *MongoUserRepository) DemoteUser(userID string) error {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": user_id}, bson.M{"$set": bson.M{"role": "user"}})
	return err
}

func (r *MongoUserRepository) UpdatePassword(userID string, hashedPassword string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"password": hashedPassword}})
	return err
}


func (r *MongoUserRepository) UpdateUserProfile(userID string, updateData *models.User) (*dtos.UpdateProfileDTO, error) {
	user_id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
    filter := bson.M{"_id": user_id}
    update := bson.M{"$set": updateData,}
    _, err = r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	
	var updatedUser *dtos.UpdateProfileDTO
	err = r.collection.FindOne(context.Background(), filter).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}
    return updatedUser, err
}