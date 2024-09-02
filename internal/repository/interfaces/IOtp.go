package repository_interface

import (
	"LoanGuard/internal/domain/models"
	"context"
)

type IOtpRepository interface {
	SaveOtp(ctx context.Context, otp models.OtpEntry) error
	FindByOtp(ctx context.Context, otp string) (*models.OtpEntry, error)
}