package service

import 	"github.com/TakoB222/postingAds-api/internal/repository"

type Services struct {
	User  User
	Admin Admin
	Ads   Ad
}

type User interface {
	SignIn()
	SignUp()
	RefreshSession()
}

type Admin interface {
}

type Ad interface {
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{}
}
