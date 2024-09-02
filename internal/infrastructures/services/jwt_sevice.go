package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type IJWTService interface {
	GenerateAccessToken(userID, role string) (string, error)
	GenerateRefreshToken(userID, role string) (string, error)
	GenerateVerificationToken(userID string) (string, error)
	ValidateAccessToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateVerificationToken(token string) (string, error)
	GetClaimsFromToken(tokenString string) (jwt.MapClaims, bool)
}

type JWTService struct {
	accessSK string
	refreshSK string
	verificationSK string
}

func NewJWTService(accessSK, refreshSK, verificationSK string) IJWTService{
	return &JWTService{
		accessSK: accessSK,
		refreshSK: refreshSK,
		verificationSK: verificationSK,
	}
}

func (jwtservice *JWTService) GenerateAccessToken(userId, role string) (string, error){
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":role,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtservice.accessSK))
}	

func (jwtservice *JWTService) GenerateRefreshToken(userId, role string) (string, error){
	claims := jwt.MapClaims{
		"user_id":userId,
		"role":role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtservice.refreshSK))
}

func (jwtservice *JWTService) GenerateVerificationToken(userId string) (string, error){
	claims := jwt.MapClaims{
		"user_id":userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtservice.verificationSK))
}


func (jwtservice *JWTService) validator(tokenString, secretKey string,) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return token, nil
	}
	return nil, errors.New("invalid Token")
}

func (jwtservice *JWTService) ValidateAccessToken(token string) (*jwt.Token, error) {
	return jwtservice.validator(token, jwtservice.accessSK)
}

func (jwtservice *JWTService) ValidateRefreshToken(token string) (string, error) {
	Token, err := jwtservice.validator(token, jwtservice.refreshSK)
	if err != nil {
		return "", err
	}
	claims, ok := Token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", err
	}

	return userId, nil
}

func (jwtservice *JWTService) ValidateVerificationToken(token string) (string, error) {
	Token, err := jwtservice.validator(token, jwtservice.verificationSK)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", jwt.ErrTokenExpired
		}
		return "", err
	}

	claims, ok := Token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", err
	}
	return userId, nil
}

func (jwtservice *JWTService) GetClaimsFromToken(tokenString string) (jwt.MapClaims, bool) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}