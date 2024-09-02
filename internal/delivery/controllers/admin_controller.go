package controllers

import (
	"LoanGuard/internal/usecases"
	"github.com/gin-gonic/gin"
)


type IAdminController interface{
	GetUsers(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetLoans(ctx *gin.Context)
	AcceptOrRejectLoan(ctx *gin.Context)
	DeleteLoan(ctx *gin.Context)
	GetSystemLogs(ctx *gin.Context)
}

type AdminController struct {
	user_usecase usecases.IUserUsecase
	admin_usecase usecases.IAdminUsecase
}

func NewAdminController(user_usecase usecases.IUserUsecase, admin_usecase usecases.IAdminUsecase) IAdminController{
	return &AdminController{
		user_usecase: user_usecase,
		admin_usecase: admin_usecase,
	}
}


func (uc *AdminController) GetUsers(ctx *gin.Context){
	
	users, err := uc.user_usecase.GetUsers()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"users": users})
}

func (uc *AdminController) DeleteUser(ctx *gin.Context){
	userID := ctx.Param("id")
	err := uc.user_usecase.DeleteUser(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User successfully deleted"})
}

func (uc *AdminController) GetLoans(ctx *gin.Context){
	status := ctx.DefaultQuery("status", "all")
	order := ctx.DefaultQuery("order", "asc")
	
	loans, err := uc.admin_usecase.GetLoans(status, order)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"loans": loans})
}

func (uc *AdminController) AcceptOrRejectLoan(ctx *gin.Context){	
	loanID := ctx.Param("id")
	status := ctx.Query("status")
	err := uc.admin_usecase.AcceptOrRejectLoan(loanID, status)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Loan status updated"})
}

func (uc *AdminController) DeleteLoan(ctx *gin.Context){
	loanID := ctx.Param("id")
	err := uc.admin_usecase.DeleteLoan(loanID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Loan successfully deleted"})
}

func (uc *AdminController) GetSystemLogs(ctx *gin.Context){
	logs, err := uc.admin_usecase.GetSystemLogs()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"logs": logs})
}