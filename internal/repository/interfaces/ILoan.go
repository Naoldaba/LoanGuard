package repository_interface

import (
	"LoanGuard/internal/domain/models"
)

type ILoanRepository interface {
	GetAllLoans(status string, order string) ([]models.Loan, error)
	UpdateLoanStatus(loanID string, status string) error
	DeleteLoan(loanID string) error
	RequestLoan(loan *models.Loan) (*models.Loan, error)
	ViewLoanStatus(loanID string) (string, error)
}
