package dtos

type ResetPassword struct {
	NewPassword string `bson:"new_password" json:"new_password"`
}