package dtos


type ProfileDTO struct {
	ID				string `json:"id" bson:"_id"`
	Name   			string `json:"name" bson:"name"`
	Email			string `json:"email" bson:"email"`
	Role			string `json:"role" bson:"role"`
	ProfilePicture	string `json:"profile_picture" bson:"profile_picture"`
	Age				int    `json:"age" bson:"age"`
	PhoneNum		string `json:"phone_num" bson:"phone_num"`
	Bio				string `json:"bio" bson:"bio"`
}