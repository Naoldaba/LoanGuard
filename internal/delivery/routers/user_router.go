package routers

import (
	"LoanGuard/internal/delivery/controllers"
	"LoanGuard/internal/infrastructures/middlewares"

	"github.com/gin-gonic/gin"
)

func CreateUserRouter(router *gin.Engine, userController controllers.IUserController, otpController controllers.IOTPController, authMiddleware middlewares.IAuthMiddleware) {
	//auth
	router.POST("/users/sign-up", userController.Register)
	router.POST("/users/sign-in", userController.Login)
	router.GET("/users/sign-out", authMiddleware.Authentication(), userController.Logout)
	router.GET("/users/verify-email", userController.VerifyEmail)
	router.POST("/users/token/refresh", userController.RefreshToken)
	router.POST("/users/password-update", otpController.ResetPassword)
	router.POST("/users/password-reset", otpController.ForgotPassword)

	//user
	router.POST("/users/profile-update", authMiddleware.Authentication(), userController.UpdateProfile)
	router.GET("/users/profile", authMiddleware.Authentication(), userController.GetUser)
}
