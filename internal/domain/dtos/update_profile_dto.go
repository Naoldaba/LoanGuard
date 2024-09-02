package dtos


type UpdateProfileDTO struct {
	Name   			string `json:"name" bson:"name"`
	ProfilePicture	string `json:"profile_picture" bson:"profile_picture"`
	Age				int    `json:"age" bson:"age"`
	PhoneNum		string `json:"phone_num" bson:"phone_num"`
	Bio				string `json:"bio" bson:"bio"`
}