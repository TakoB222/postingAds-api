package service

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type (
	SignInInput struct {
		Email    string
		Password string
	}

	UserSignUpInput struct {
		FirsName string
		LastName string
		Email    string
		Password string
	}

	RefreshInput struct {
		RefreshToken string `json:"refreshToken"`
	}

	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	Ads struct {
		UserId      string   `json:"user_id"`
		Title       string   `json:"title"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		Price       int      `json:"price"`
		Contacts    Contacts `json:"contacts"`
		Published   bool     `json:"published"`
		ImagesURL   []string `json:"images_url"`
	}

	Contacts struct {
		Name         string `json:"name"`
		Phone_number string `json:"phone_number"`
		Email        string `json:"email"`
		Location     string `json:"location"`
	}
	FtsResponse struct {
		Id string `db:"id"`
		Title string `db:"title"`
	}
)

type Authorization interface {
	SignUp(input UserSignUpInput) (int, error)
	SignIn(input SignInInput) (Tokens, error)
	RefreshSession(input RefreshInput) (Tokens, error)
}

type Admin interface {
	AdminSignIn(input SignInInput) (Tokens, error)
	AdminRefreshSession(input RefreshInput) (Tokens, error)
	AdminGetAllAdsByAdmin() ([]domain.Ad, error)
	AdminGetAd(adId string) (domain.Ad, error)
	AdminDeleteUserAdById(adId string) error
	AdminUpdateAd(adId string, ad Ads) (domain.Ad, error)
}

type Ad interface {
	GetAllAds(userId string) ([]domain.Ad, error)
	CreateAd(userId string, adInput Ads) (int, error)
	GetAdById(userId string, adId string) (domain.Ad, error)
	UpdateAd(userId, adId string, ad Ads) (domain.Ad, error)
	DeleteAd(userId string, adId string) error
	Fts(search_request string)([]repository.FtsResponse, error)
}

type Service struct {
	Authorization
	Admin
	Ad
}

type Dependencies struct {
	Repository   *repository.Repository
	TokenManager *auth.Manager
	Hasher       *hash.SHA1Hasher

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(dep Dependencies) *Service {
	return &Service{
		Authorization: NewAuthService(dep.Repository, dep.TokenManager, dep.Hasher, dep.AccessTokenTTL, dep.RefreshTokenTTL),
		Ad:            NewAdService(dep.Repository),
		Admin:         NewAdminService(dep.Repository, dep.TokenManager, dep.Hasher, dep.AccessTokenTTL, dep.RefreshTokenTTL),
	}
}
