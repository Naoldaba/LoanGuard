package controllers

import (
	"LoanGuard/internal/domain/dtos"
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/usecases"
	"strconv"

	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)


type IUserController interface{
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	Logout(c *gin.Context)
	GetUser(ctx *gin.Context)
}

type UserController struct {
	user_usecase usecases.IUserUsecase
}

func NewUserController(user_usecase usecases.IUserUsecase) IUserController{
	return &UserController{
		user_usecase: user_usecase,
	}
}

func (uc *UserController) Register(ctx *gin.Context){
	var user *models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid json format"})
		return
	}
	_, err = uc.user_usecase.Register(user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "verification email sent to your email"})
}

func (uc *UserController) Login(ctx *gin.Context){
	var user *dtos.LoginDTO
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return 
	}
	accTkn, refTkn, err := uc.user_usecase.Login(user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return 
	}

	ctx.JSON(200, gin.H{"accTkn": accTkn, "refTkn": refTkn})
}


func (uc *UserController) RefreshToken(ctx *gin.Context) {
	var refRequest dtos.RefreshTknRequest
	if err := ctx.ShouldBindJSON(&refRequest); err != nil {
		ctx.JSON(400, gin.H{"error": "Refresh token is required"})
		return
	}

	newAccessToken, err := uc.user_usecase.RefreshToken(refRequest.RefreshToken)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"access_token": newAccessToken})
}


func (uc *UserController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(400, gin.H{"error": "Missing token"})
		return
	}

	accTkn, refTkn, err := uc.user_usecase.VerifyEmailToken(token)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Email successfully verified", "access_token": accTkn, "refresh_token": refTkn})
}

func (uc *UserController) Logout(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		c.JSON(400, gin.H{"error": "Missing token"})
		return
	}
	tokenStr, ok := token.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid token format"})
		return
	}
	err := uc.user_usecase.Logout(tokenStr)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully logged out"})
}

func (uc *UserController) UpdateProfile(ctx *gin.Context){
	claims, _ := ctx.Get("claims")
    jwtClaims, ok := claims.(jwt.MapClaims)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
        return
    }
    userID, ok := jwtClaims["user_id"].(string)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
        return
    }

	var updateData dtos.UpdateProfileDTO
    if username := ctx.PostForm("name"); username != "" {
        updateData.Name = username
    }
	if phoneNum := ctx.PostForm("phone_num"); phoneNum != "" {
		updateData.PhoneNum = phoneNum
	}
	if bio := ctx.PostForm("bio"); bio != "" {
		updateData.Bio = bio
	}
	if ageStr := ctx.PostForm("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age format"})
			return
		}
		updateData.Age = age
	}

	image, _, err := ctx.Request.FormFile("profile_picture")
	if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process profile picture"})
		return
	}
	var updatedUser *dtos.UpdateProfileDTO
    if updatedUser, err = uc.user_usecase.UpdateProfile(userID, &updateData, image); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"updaterUser": updatedUser})
}

func (uc *UserController) GetUser(ctx *gin.Context){
	claims, _ := ctx.Get("claims")
	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}
	userID, ok := jwtClaims["user_id"].(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	user, err := uc.user_usecase.GetMyProfile(userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"user": user})
}