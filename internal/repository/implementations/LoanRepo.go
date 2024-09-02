package implementations

import (
	"context"
	"fmt"

	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoLoanRepository struct {
    collection *mongo.Collection
}

func NewMongoLoanRepository(db *mongo.Database) repository_interface.ILoanRepository {
    return &mongoLoanRepository{
        collection: db.Collection("loans"),
    }
}

func (r *mongoLoanRepository) GetAllLoans(status string, order string) ([]models.Loan, error) {
	var loans []models.Loan
	filter := bson.M{}

	if status != "" && status != "all" {
		filter["status"] = status
	}

	findOptions := options.Find()
	if order == "asc" {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}})
	} else {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &loans); err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *mongoLoanRepository) UpdateLoanStatus(loanID string, status string) error {
	Id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return err
	}
    filter := bson.M{"_id": Id}
    update := bson.M{"$set": bson.M{"status": status}}
    _, err = r.collection.UpdateOne(context.Background(), filter, update)
    return err
}

func (r *mongoLoanRepository) DeleteLoan(loanID string) error {
	Id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": Id})
    return err
}

func (r *mongoLoanRepository) RequestLoan(loan *models.Loan) (*models.Loan, error) {
	result, err := r.collection.InsertOne(context.Background(), loan)
	if err != nil {
		return nil, err	
	}

	var created_loan models.Loan
	err = r.collection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&created_loan)
	if err != nil {
		return nil, err
	}
	return &created_loan, nil
}

func (r *mongoLoanRepository) ViewLoanStatus(loanID string) (string, error) {
	Id, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": Id}

	var loan models.Loan
	err = r.collection.FindOne(context.Background(), filter).Decode(&loan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		return "", err
	}
	fmt.Println(loan)
	return loan.Status, nil
}