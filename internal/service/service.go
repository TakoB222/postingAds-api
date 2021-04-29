package service

import (
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"time"
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

	RefreshInput struct {
		RefreshToken string `json:"refreshToken"`
		Ua string `json:"ua"`
		Ip string `json:"ip"`
	}

	Tokens struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

type Authorization interface {
	SignUp(input UserSignUpInput)(int,error)
	SignIn(input UserSignInInput)(Tokens,error)
	RefreshSession(input RefreshInput)(Tokens, error)
}

type Admin interface {
}

type Ad interface {

}

type Service struct {
	Authorization
	Admin
	Ad
}

type Dependencies struct {
	Repository *repository.Repository
	TokenManager *auth.Manager
	Hasher *hash.SHA1Hasher

	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(dep Dependencies) *Service {
	return &Service{
		Authorization:NewAuthService(dep.Repository, dep.TokenManager, dep.Hasher),
	}
}
