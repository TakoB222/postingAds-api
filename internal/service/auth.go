package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"time"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	accessTokenTTL   = 30 * time.Minute
	refreshTokenTTL =  30 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
	tokenManager auth.TokenManager
	hasher *hash.SHA1Hasher
}

func NewAuthService(repo repository.Authorization, tokenManager *auth.Manager, hasher *hash.SHA1Hasher) *AuthService {
	return &AuthService{repo: repo, tokenManager: tokenManager, hasher: hasher}
}

func (s *AuthService) SignUp(input UserSignUpInput) (int, error) {
	user := domain.User{
		Email: input.Email,
		Password_hash: s.hasher.Hash(input.Password),
		First_name: input.FirsName,
		Last_name: input.LastName,
		Registered_at: time.Now(),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) SignIn(input UserSignInInput)(Tokens, error){
	user, err := s.repo.GetUser(input.Email, s.hasher.Hash(input.Password))
	fmt.Println("userid - ", user.Id)
	if err != nil {
		//TODO: create custom repository errors and handle them here
		return Tokens{}, err
	}

	return s.createSession(user.Id, input.Ua, input.Ip)
}

func (s *AuthService) createSession(userId string, ua string, ip string)(Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		UserId: userId,
		RefreshToken: res.RefreshToken,
		CreatedAt: time.Now(),
		ExpiresIn: time.Now().Add(refreshTokenTTL),
		UA: ua,
		Ip: ip,
	}

	err = s.repo.SetSession(session)
	if err != nil {
		return Tokens{}, err
	}

	return res, nil
}

func (s *AuthService) RefreshSession(refreshToken string)(Tokens, error) {
	session, err := s.repo.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	err = s.repo.DeleteSessionByUserId(session.UserId)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(session.UserId, session.UA,session.Ip)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

