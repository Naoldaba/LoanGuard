package usecases

import (
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/repository/interfaces"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ILoanUsecase interface {
	RequestLoan(userID string, loan *models.Loan) (*models.Loan, error)
	ViewLoanStatus(loanID string) (string, error)
}

type LoanUsecase struct {
	loanRepo repository_interface.ILoanRepository
	logRepo  repository_interface.ILogRepository
}

func NewLoanUsecase(loanRepo repository_interface.ILoanRepository, logRepo repository_interface.ILogRepository) ILoanUsecase {
	return &LoanUsecase{
		loanRepo: loanRepo,
		logRepo: logRepo,
	}
}

func (lu *LoanUsecase) RequestLoan(userID string, loan *models.Loan) (*models.Loan, error) {
	loan.Status = "pending"
	loan.CreatedAt = time.Now()
	id, _ := primitive.ObjectIDFromHex(userID)
	loan.Interest = 0.05 
	
	loan.Total = float32(loan.Amount) + (float32(loan.Amount) * loan.Interest)
	loan.UserId = id

	result, err := lu.loanRepo.RequestLoan(loan)
	if err != nil {
		return nil, err
	}
	Newlog := models.SystemLog{
		Action:  "Loan Application Submission",
		UserID:  loan.UserId,
		LoanID: "Loan ID: " + result.ID.Hex(),
	}
	lu.logRepo.CreateLog(&Newlog)

	return result, nil
}

func (lu *LoanUsecase) ViewLoanStatus(loanID string) (string, error) {
	loanStatus, err := lu.loanRepo.ViewLoanStatus(loanID)
	if err != nil {
		return "", err
	}
	return loanStatus, nil
}