package service

import (
	"github.com/869413421/pg-service/user/pkg/model"
	"github.com/869413421/pg-service/user/pkg/repo"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	key = []byte("pgServiceUserTokenKeySecret")
)

type CustomClaims struct {
	User *model.User
	jwt.StandardClaims
}

type Authble interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *model.User) (string, error)
}

type TokenService struct {
	Repo repo.UserRepositoryInterface
}

func NewTokenService(repo repo.UserRepositoryInterface) *TokenService {
	return &TokenService{Repo: repo}
}

// Decode a token string into a token object
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode a claim into a JWT
func (srv *TokenService) Encode(user *model.User) (string, error) {

	expireToken := time.Now().Add(time.Hour * 72).Unix()

	// Create the Claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "pg.service.user",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(key)
}
