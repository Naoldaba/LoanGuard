package routers

import (
	"LoanGuard/internal/delivery/controllers"
	"LoanGuard/internal/infrastructures/middlewares"

	"github.com/gin-gonic/gin"
)

func CreateLoanRouter(router *gin.Engine, loanController controllers.ILoanController, authMiddleware middlewares.IAuthMiddleware) {
	router.POST("/loan", authMiddleware.Authentication(),loanController.RequestLoan)
	router.GET("/loan/:id", authMiddleware.Authentication(), loanController.ViewLoanStatus)
}
