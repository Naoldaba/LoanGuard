package routers

import (
	"LoanGuard/internal/delivery/controllers"
	"LoanGuard/internal/infrastructures/middlewares"

	"github.com/gin-gonic/gin"
)

func CreateAdminRouter(router *gin.Engine, adminController controllers.IAdminController, authMiddleware middlewares.IAuthMiddleware) {
	router.GET("/admin/users", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"),adminController.GetUsers)
	router.DELETE("/admin/users/:id", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"), adminController.DeleteUser)
	router.GET("/admin/loans", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"), adminController.GetLoans)
	router.PATCH("/admin/:id/status", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"), adminController.AcceptOrRejectLoan)
	router.DELETE("/admin/loans/:id", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"), adminController.DeleteLoan)
	router.GET("/admin/logs", authMiddleware.Authentication(), authMiddleware.RoleAuth("ADMIN"), adminController.GetSystemLogs)
}
