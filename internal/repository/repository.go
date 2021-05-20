package repository

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	//TODO: create custom errors
	ErrorEmptyResult = "empty result"
)

type (
	Ads struct {
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

type User interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email, password_hash string) (domain.User, error)
	GetSessionByRefreshToken(refreshToken string) (domain.Session, error)
	DeleteSessionByUserId(userId string) error
	SetSession(session domain.Session) error
}

type Admin interface {
	GetAdminId(email, password_hash string) (string, error)
	GetAdminSessionByRefreshToken(refrehsToken string) (domain.AdminSession, error)
	DeleteAdminSessionByAdminId(adminId string) error
	SetAdminSession(session domain.AdminSession) error
	GetAllAdsByAdmin() ([]domain.Ad, error)
	GetAd(adId string) (domain.Ad, error)
	AdminDeleteAd(adId string) error
	AdminUpdateAd(adId string, ad Ads) error
}

type Ad interface {
	GetAllAdsByUserId(userId string) ([]domain.Ad, error)
	CreateAd(userId string, input Ads) (int, error)
	GetAdById(userId string, adId string) (domain.Ad, error)
	UpdateAd(userId string, adId string, ad Ads) error
	DeleteAd(userId string, adId string) error
	SearchAdByRequest(search_request string)([]FtsResponse, error)
}

type Repository struct {
	User
	Admin
	Ad
}

func NewRepositories(db *sqlx.DB) *Repository {
	return &Repository{
		User:  NewAuthRepository(db),
		Ad:    NewAdRepository(db),
		Admin: NewAdminRepository(db),
	}
}
