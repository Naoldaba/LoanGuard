package usecases

import (
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/infrastructures/services"
	"LoanGuard/internal/infrastructures/services/email_service"
	"LoanGuard/internal/repository/interfaces"
	"context"
	"errors"
	"fmt"
	"time"
)

type IOtpUsecase interface {
	GenerateAndSendOtp(ctx context.Context, email string) error
	ValidateOtp(ctx context.Context, otp string) (*models.OtpEntry, error)
	ResetPassword(ctx context.Context, userID string, newPassword string) error
}

type OtpUsecase struct {
	otpRepo       repository_interface.IOtpRepository
	userRepo      repository_interface.IUserRepository
	emailSvc      email_service.IEmailService
	passSvc       services.IHashService
	baseUri       string
	validationSvc services.IValidationService
}

func NewOtpUseCase(otpRepo repository_interface.IOtpRepository, userRepo repository_interface.IUserRepository, emailSvc email_service.IEmailService, passSvc services.IHashService, baseUri string, validationSvc services.IValidationService) IOtpUsecase {
	return &OtpUsecase{
		otpRepo:       otpRepo,
		userRepo:      userRepo,
		emailSvc:      emailSvc,
		baseUri:   baseUri,
		passSvc:       passSvc,
		validationSvc: validationSvc,
	}
}

func (u *OtpUsecase) GenerateAndSendOtp(ctx context.Context, email string) error {
	var otpEntry models.OtpEntry
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	otp := services.GenerateOTP()
	expiresAt := time.Now().Add(10 * time.Minute)

	otpEntry.ExpiresAt = expiresAt
	otpEntry.UserID = user.ID.Hex()
	otpEntry.OTP = otp

	err = u.otpRepo.SaveOtp(ctx, otpEntry)
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/users/password-update?otp=%s", u.baseUri, otp)

	err = u.emailSvc.SendResetEmail(user.Email, resetLink)
	if err != nil {
		return err
	}

	return nil
}

func (u *OtpUsecase) ValidateOtp(ctx context.Context, otp string) (*models.OtpEntry, error) {
	otpEntry, err := u.otpRepo.FindByOtp(ctx, otp)
	if err != nil {
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(otpEntry.ExpiresAt) {
		return nil, errors.New("OTP has expired")
	}

	return otpEntry, nil
}

func (u *OtpUsecase) ResetPassword(ctx context.Context, userID string, newPassword string) error {
	var newUser *models.User
	newUser, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if _, err := u.validationSvc.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := u.passSvc.HashPassword(newPassword)
	if err != nil {
		return err
	}

	newUser.Password = hashedPassword
	err = u.userRepo.UpdateUser(userID, newUser)
	if err != nil {
		return err
	}

	return nil
}
