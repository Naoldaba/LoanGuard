package controllers

import (
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ILoanController interface{
	RequestLoan(ctx *gin.Context)
	ViewLoanStatus(ctx *gin.Context)
}

type LoanController struct {
	loanUsecase usecases.ILoanUsecase
}

func NewLoanController(loanUsecase usecases.ILoanUsecase) ILoanController{
	return &LoanController{
		loanUsecase: loanUsecase,
	}
}

func (lc *LoanController) RequestLoan(ctx *gin.Context){
	var loan *models.Loan
	err := ctx.ShouldBindJSON(&loan)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid json format"})
		return
	}

	claims, _ := ctx.Get("claims")
    jwtClaims, ok := claims.(jwt.MapClaims)
    if !ok {
        ctx.JSON(400, gin.H{"error": "Failed to parse claims"})
        return
    }
    userID, _ := jwtClaims["user_id"].(string)

	result, err := lc.loanUsecase.RequestLoan(userID, loan)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"loan_request": result})
}

func (lc *LoanController) ViewLoanStatus(ctx *gin.Context){
	loanId := ctx.Param("id")
	loanStatus, err := lc.loanUsecase.ViewLoanStatus(loanId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"loanStatus": loanStatus})
}