package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID     				primitive.ObjectID `json:"id" bson:"_id"`
	Name   				string             `json:"name" bson:"name"`
	Email  				string             `json:"email" bson:"email"`
	Password 			string			   `json:"password" bson:"password"`
	Role   				string             `json:"role" bson:"role"`
	Age    				int                `json:"age" bson:"age"`
	PhoneNum			string             `json:"phone_num" bson:"phone_num"`
	Bio   				string             `json:"bio" bson:"bio"`
	ProfilePicture		string             `json:"profile_picture" bson:"profile_picture"`
	IsVerified			bool 			   `json:"is_verified" bson:"is_verified"`
	VerificationToken	string			   `json:"-" bson:"verification_token"`	
	RefToken			string 			   `json:"-" bson:"refresh_token"`	
}
