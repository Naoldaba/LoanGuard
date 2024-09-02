package usecases

import (
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/repository/interfaces"
	"errors"
)
	
type IAdminUsecase interface {
    GetLoans(status, order string) ([]models.Loan, error)
    AcceptOrRejectLoan(loanID string, status string) error
    DeleteLoan(loanID string) error
    GetSystemLogs() ([]models.SystemLog, error)
}

type adminUseCase struct {
    loanRepo repository_interface.ILoanRepository
    logRepo  repository_interface.ILogRepository
}

func NewAdminUsecase(loanRepo repository_interface.ILoanRepository, logRepo repository_interface.ILogRepository) IAdminUsecase {
    return &adminUseCase{loanRepo: loanRepo, logRepo: logRepo}
}

func (uc *adminUseCase) GetLoans(status string, order string) ([]models.Loan, error) {
    loans, err := uc.loanRepo.GetAllLoans(status, order)
    if err != nil {
        return nil, err
    }
    return loans, nil
}

func (uc *adminUseCase) AcceptOrRejectLoan(loanID string, status string) error {
    if status != "approved" && status != "rejected" {
        return errors.New("invalid status")
    }

    err := uc.loanRepo.UpdateLoanStatus(loanID, status)
    if err != nil {
        return err
    }

    log := models.SystemLog{
        Action: "Loan " + status + ": " + loanID,
    }
    uc.logRepo.CreateLog(&log)

    return nil
}

func (uc *adminUseCase) DeleteLoan(loanID string) error {
    err := uc.loanRepo.DeleteLoan(loanID)
    if err != nil {
        return err
    }

    log := models.SystemLog{
        Action: "Loan deleted: " + loanID,
    }
    uc.logRepo.CreateLog(&log)

    return nil
}

func (uc *adminUseCase) GetSystemLogs() ([]models.SystemLog, error) {
    logs, err := uc.logRepo.GetAllLogs()
    if err != nil {
        return nil, err
    }
    return logs, nil
}
