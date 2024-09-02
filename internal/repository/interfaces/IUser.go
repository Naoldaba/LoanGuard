package repository_interface

import (
	"LoanGuard/internal/domain/dtos"
	"LoanGuard/internal/domain/models"
	"time"
)

type IUserRepository interface {
	Register(user *models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user *models.User) error
	UpdateUserProfile(userID string, updateData *models.User) (*dtos.UpdateProfileDTO, error)
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	UpdatePassword(userID string, hashedPassword string) error
	BlacklistToken(token string, remainingTime time.Duration) error
}
