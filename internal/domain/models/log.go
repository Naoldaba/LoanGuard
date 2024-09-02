package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemLog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Action    string             `json:"action" bson:"action"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	LoanID    string             `json:"loan_id" bson:"loan_id,omitempty"`
}