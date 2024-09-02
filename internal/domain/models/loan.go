package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID             primitive.ObjectID 		`json:"id" bson:"_id,omitempty"`
	Amount         int    					`json:"amount" bson:"amount"`
	Interest       float32    			    `json:"interest" bson:"interest"`
	Total          float32    				`json:"total" bson:"total"`
	Status         string 					`json:"status" bson:"status"`
	LoanPurpose    string 					`json:"loan_purpose" bson:"loan_purpose"`
	UserId         primitive.ObjectID 		`json:"userId" bson:"userId"`
	CreatedAt 	   time.Time                `jaon:"created_at" bson:"created_at"`
}