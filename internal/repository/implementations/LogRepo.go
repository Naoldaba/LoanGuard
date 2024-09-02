package implementations

import (
    "context"

    "LoanGuard/internal/domain/models"
	"LoanGuard/internal/repository/interfaces"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type mongoLogRepository struct {
    collection *mongo.Collection
}

func NewMongoLogRepository(db *mongo.Database) repository_interface.ILogRepository {
    return &mongoLogRepository{
        collection: db.Collection("logs"),
    }
}

func (r *mongoLogRepository) CreateLog(log *models.SystemLog) error {
    _, err := r.collection.InsertOne(context.Background(), log)
    return err
}

func (r *mongoLogRepository) GetAllLogs() ([]models.SystemLog, error) {
    var logs []models.SystemLog
    cursor, err := r.collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    if err := cursor.All(context.Background(), &logs); err != nil {
        return nil, err
    }
    return logs, nil
}