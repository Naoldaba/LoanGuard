package usecases

import (
	"LoanGuard/internal/domain/dtos"
	"LoanGuard/internal/domain/models"
	"LoanGuard/internal/infrastructures/services"
	"LoanGuard/internal/infrastructures/services/email_service"
	"LoanGuard/internal/repository/interfaces"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type IUserUsecase interface {
	Register(user *models.User) (*models.User, error)
	Login(user *dtos.LoginDTO) (string, string, error)
	Logout(token string) error
	RefreshToken(refreshToken string) (string, error)
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	DeleteUser(userID string) error
	UpdateUser(userID string, user *models.User) error
	UpdateProfile(userID string, user *dtos.UpdateProfileDTO, image multipart.File) (*dtos.UpdateProfileDTO, error)
	GetMyProfile(userID string) (*dtos.ProfileDTO, error)
	PromoteUser(userID string) error
	DemoteUser(userID string) error
	VerifyEmailToken(token string) (string, string, error)
}


type UserUsecase struct {
	userRepo          repository_interface.IUserRepository
	passwordService   services.IHashService
	validationService services.IValidationService
	emailService      email_service.IEmailService
	jwtSevices        services.IJWTService
	cloudSvc		  services.ICloudinaryService
	baseUri 	      string
}


func NewUserUsecase(userRepo repository_interface.IUserRepository, passwordService services.IHashService, validationService services.IValidationService, emailService email_service.IEmailService, jwtService services.IJWTService, cloudSvc services.ICloudinaryService, baseUri string) IUserUsecase {
	return &UserUsecase{
		userRepo:          userRepo,
		passwordService:   passwordService,
		validationService: validationService,
		emailService:      emailService,
		jwtSevices:        jwtService,
		cloudSvc:          cloudSvc,
		baseUri: 	       baseUri,			
	}
}

func (u *UserUsecase) Register(user *models.User) (*models.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		user.Role = "ADMIN"
	} else {
		existingUser, _ := u.userRepo.GetUserByEmail(user.Email)
		if existingUser != nil {
			return nil, err
		}
		user.Role = "USER"
	}

	if _, err := u.validationService.ValidatePassword(user.Password); err != nil {
		return nil, err
	}
	if _, err := u.validationService.ValidateEmail(user.Email); err != nil {
		return nil, err
	}

	encryptedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = encryptedPassword
	user.IsVerified = false

	regUser, err := u.userRepo.Register(user)
	if err != nil {
		return nil, err
	}

	verificationToken, err := u.jwtSevices.GenerateVerificationToken(regUser.ID.Hex())
	if err != nil {
		return nil, err
	}
	verificationLink := fmt.Sprintf("%s/users/verify-email?token=%s",u.baseUri, verificationToken)
	user.VerificationToken = verificationToken
	u.userRepo.UpdateUser(regUser.ID.Hex(), user)
	err = u.emailService.SendVerificationEmail(user.Email, verificationLink)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) Login(user *dtos.LoginDTO) (string, string, error) {
	if _, err := u.validationService.ValidateEmail(user.Email); err != nil {
		return "", "", errors.New(err.Error())
	}
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}
	if !u.passwordService.CompareHash(existingUser.Password, user.Password) {
		return "", "", errors.New("invalid password")
	}

	accessToken, _ := u.jwtSevices.GenerateAccessToken(existingUser.ID.Hex(), existingUser.Role)
	refershToken, _ := u.jwtSevices.GenerateRefreshToken(existingUser.ID.Hex(), existingUser.Role)

	existingUser.RefToken = refershToken
	err = u.userRepo.UpdateUser(existingUser.ID.Hex(), existingUser)
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	return accessToken, refershToken, nil
}

func (u *UserUsecase) Logout(token string) error {
	parsedToken, err := u.jwtSevices.ValidateAccessToken(token)
	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	expiration, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration time in token claims")
	}

	remTime := time.Until(time.Unix(int64(expiration), 0))
	if remTime <= 0 {
		return errors.New("token has already expired")
	}

	err = u.userRepo.BlacklistToken(token, remTime)
	if err != nil {
		return err
	}

	existingUser, err := u.userRepo.GetUserByID(claims["user_id"].(string))
	if err != nil {
		return err
	}

	existingUser.RefToken = ""
	err = u.userRepo.UpdateUser(existingUser.ID.Hex(), existingUser)
	if err != nil {
		return errors.New("failed to update user profile")
	}
	return nil
}

func (u *UserUsecase) RefreshToken(refreshTok string) (string, error) {
	userId, err := u.jwtSevices.ValidateRefreshToken(refreshTok)
	if err != nil {
		return "", errors.New(err.Error())
	}

	existingUser, err := u.userRepo.GetUserByID(userId)
	if err != nil {
		return "", errors.New("user not found")
	}

	if existingUser.RefToken != refreshTok {
		return "", errors.New("invalid token")
	}
	accessToken, _ := u.jwtSevices.GenerateAccessToken(existingUser.ID.Hex(), existingUser.Role)
	return accessToken, nil
}

func (u *UserUsecase) VerifyEmailToken(token string) (string, string, error) {
	userId, err := u.jwtSevices.ValidateVerificationToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, ok := u.jwtSevices.GetClaimsFromToken(token)
			if ok {
				userId, _ := claims["user_id"].(string)
				user, err := u.userRepo.GetUserByID(userId)
				if err != nil {
					return "", "", errors.New("invalid token")
				}
				if user.IsVerified {
					return "", "", errors.New("user already verified")
				}

				verificationToken, _ := u.jwtSevices.GenerateVerificationToken(user.ID.Hex())
				verificationLink := fmt.Sprintf("%s/users/verify-email?token=%s", u.baseUri, verificationToken)

				user.VerificationToken = verificationToken
				err = u.userRepo.UpdateUser(user.ID.Hex(), user)
				if err != nil {
					return "", "", err
				}

				err = u.emailService.SendVerificationEmail(user.Email, verificationLink)
				if err != nil {
					return "", "", err
				}

				return "", "", errors.New("verification token expired. A new verification email has been sent")
			}
		}
		return "", "", errors.New(err.Error())
	}

	user, err := u.userRepo.GetUserByID(userId)
	if err != nil {
		return "", "", errors.New("invalid token")
	} else if user.IsVerified {
		return "", "", errors.New("user already verified")
	} else if token != user.VerificationToken {
		return "", "", errors.New("invalid token")
	}
	user.IsVerified = true
	user.VerificationToken = ""

	accessToken, _ := u.jwtSevices.GenerateAccessToken(user.ID.Hex(), user.Role)
	refershToken, _ := u.jwtSevices.GenerateRefreshToken(user.ID.Hex(), user.Role)

	user.RefToken = refershToken
	err = u.userRepo.UpdateUser(user.ID.Hex(), user)
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	return accessToken, refershToken, nil
}

func (u *UserUsecase) GetUserByID(userID string) (*models.User, error) {
	return u.userRepo.GetUserByID(userID)
}

func(u *UserUsecase) GetMyProfile(userID string) (*dtos.ProfileDTO, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &dtos.ProfileDTO{
		ID: user.ID.Hex(),
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		PhoneNum: user.PhoneNum,
		Bio: user.Bio,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (u *UserUsecase) GetUserByEmail(email string) (*models.User, error) {
	if _, err := u.validationService.ValidateEmail(email); err != nil {
		return nil, err
	}
	return u.userRepo.GetUserByEmail(email)
}

func (u *UserUsecase) DeleteUser(userID string) error {
	return u.userRepo.DeleteUser(userID)
}

func (u *UserUsecase) GetUsers() ([]*models.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *UserUsecase) UpdateUser(userID string, user *models.User) error {
	if _, err := u.validationService.ValidateEmail(user.Email); err != nil {
		return err
	}
	return u.userRepo.UpdateUser(userID, user)
}

func (u *UserUsecase) UpdateProfile(userID string, user *dtos.UpdateProfileDTO, image multipart.File) (*dtos.UpdateProfileDTO, error){
	existingUser, err := u.userRepo.GetUserByID(userID)
    if err != nil {
        return nil,err
    }

    if image != nil {
        profilePictureURL, err := u.cloudSvc.UploadProfilePicture(image)
        if err != nil {
            return nil,err
        }
        user.ProfilePicture = profilePictureURL
    }

    if user.Name != "" {
        existingUser.Name = user.Name
    }
    if user.PhoneNum != "" {
        existingUser.PhoneNum = user.PhoneNum
    }
    if user.Bio != "" {
        existingUser.Bio = user.Bio
    }
    if user.ProfilePicture != "" {
        existingUser.ProfilePicture = user.ProfilePicture
    }

    return u.userRepo.UpdateUserProfile(userID, existingUser)
}

func (u *UserUsecase) PromoteUser(userID string) error {
	return u.userRepo.PromoteUser(userID)
}

func (u *UserUsecase) DemoteUser(userID string) error {
	return u.userRepo.DemoteUser(userID)
}
