package service

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"time"
)

type AdminService struct {
	repo         repository.Admin
	tokenManager auth.TokenManager
	hasher       *hash.SHA1Hasher

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAdminService(repo repository.Admin, tokenManager *auth.Manager, hasher *hash.SHA1Hasher, AccesTokenTTL, RefreshTokenTTL time.Duration) *AdminService {
	return &AdminService{repo: repo, tokenManager: tokenManager, hasher: hasher, AccessTokenTTL: AccesTokenTTL, RefreshTokenTTL: RefreshTokenTTL}
}

func (s *AdminService) AdminSignIn(input SignInInput) (Tokens, error) {
	id, err := s.repo.GetAdminId(input.Email, s.hasher.Hash(input.Password))
	if err != nil {
		//TODO: create custom repository errors and handle them here
		return Tokens{}, err
	}

	return s.createSession(id)
}

func (s *AdminService) createSession(adminId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(adminId, s.AccessTokenTTL)

	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.AdminSession{
		AdminId:      adminId,
		RefreshToken: res.RefreshToken,
		CreatedAt:    time.Now(),
		ExpiresIn:    time.Now().Add(s.RefreshTokenTTL),
	}

	err = s.repo.SetAdminSession(session)
	if err != nil {
		return Tokens{}, err
	}

	return res, nil
}

func (s *AdminService) AdminRefreshSession(input RefreshInput) (Tokens, error) {
	session, err := s.repo.GetAdminSessionByRefreshToken(input.RefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if err = s.repo.DeleteAdminSessionByAdminId(session.AdminId); err != nil {
		return Tokens{}, err
	}

	return s.createSession(session.AdminId)
}

func (s *AdminService) AdminGetAllAdsByAdmin() ([]domain.Ad, error) {
	ads, err := s.repo.GetAllAdsByAdmin()
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (s *AdminService) AdminGetAd(adId string) (domain.Ad, error) {
	ad, err := s.repo.GetAd(adId)
	if err != nil {
		return domain.Ad{}, err
	}

	return ad, nil
}

func (s *AdminService) AdminDeleteUserAdById(adId string) error {
	if err := s.repo.AdminDeleteAd(adId); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AdminUpdateAd(adId string, ad Ads) error {
	if err := s.repo.AdminUpdateAd(adId, repository.Ads{
		Title:       ad.Title,
		Category:    ad.Category,
		Description: ad.Description,
		Price:       ad.Price,
		Contacts: repository.Contacts{
			Name:         ad.Contacts.Name,
			Phone_number: ad.Contacts.Phone_number,
			Email:        ad.Contacts.Email,
			Location:     ad.Contacts.Location,
		},
		Published: ad.Published,
		ImagesURL: ad.ImagesURL,
	}); err != nil {
		return err
	}

	return nil
}
