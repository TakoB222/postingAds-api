package service

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"time"
)

type AuthService struct {
	repo         repository.User
	tokenManager auth.TokenManager
	hasher       *hash.SHA1Hasher

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAuthService(repo repository.User, tokenManager *auth.Manager, hasher *hash.SHA1Hasher, AccesTokenTTL, RefreshTokenTTL time.Duration) *AuthService {
	return &AuthService{repo: repo, tokenManager: tokenManager, hasher: hasher, AccessTokenTTL: AccesTokenTTL, RefreshTokenTTL: RefreshTokenTTL}
}

func (s *AuthService) SignUp(input UserSignUpInput) (int, error) {
	user := domain.User{
		Email:         input.Email,
		Password_hash: s.hasher.Hash(input.Password),
		First_name:    input.FirsName,
		Last_name:     input.LastName,
		Registered_at: time.Now(),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) SignIn(input SignInInput) (Tokens, error) {
	user, err := s.repo.GetUser(input.Email, s.hasher.Hash(input.Password))
	if err != nil {
		//TODO: create custom repository errors and handle them here
		return Tokens{}, err
	}

	return s.createSession(user.Id)
}

func (s *AuthService) createSession(userId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, s.AccessTokenTTL)

	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		UserId:       userId,
		RefreshToken: res.RefreshToken,
		CreatedAt:    time.Now(),
		ExpiresIn:    time.Now().Add(s.RefreshTokenTTL),
	}

	err = s.repo.SetSession(session)
	if err != nil {
		return Tokens{}, err
	}

	return res, nil
}

func (s *AuthService) RefreshSession(input RefreshInput) (Tokens, error) {
	session, err := s.repo.GetSessionByRefreshToken(input.RefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	err = s.repo.DeleteSessionByUserId(session.UserId)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(session.UserId)
}
