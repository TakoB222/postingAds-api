package service

import (
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
)


type(
	UserSignInInput struct {
		Email    string
		Password string
		Ua string //user-agent
		Ip string
	}

	UserSignUpInput struct {
		FirsName string
		LastName string
		Email    string
		Password string
	}

	Tokens struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

type Users interface {
	SignUp(input UserSignUpInput)(int,error)
	SignIn(input UserSignInInput)(Tokens,error)
	RefreshSession(refreshToken string)(Tokens, error)
}

type Admin interface {
}

type Ad interface {

}

type Service struct {
	Users
	Admin
	Ad
}

type Dependencies struct {
	Repository *repository.Repository
	TokenManager *auth.Manager
	Hasher *hash.SHA1Hasher
}

func NewServices(dep Dependencies) *Service {
	return &Service{
		Users:NewAuthService(dep.Repository, dep.TokenManager, dep.Hasher),
	}
}
