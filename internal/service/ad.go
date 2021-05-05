package service

import (
	"github.com/TakoB222/postingAds-api/internal/domain"
	"github.com/TakoB222/postingAds-api/internal/repository"
)

type AdService struct {
	repo repository.Ad
}

func NewAdService(repo repository.Ad) *AdService {
	return &AdService{repo: repo}
}

func (s *AdService) GetAllAds(userId string) ([]domain.Ad, error) {
	ads, err := s.repo.GetAllAdsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (s *AdService) CreateAd(userId string, adInput Ads) (int, error) {
	adId, err := s.repo.CreateAd(userId, repository.Ads{
		Title:       adInput.Title,
		Category:    adInput.Category,
		Description: adInput.Description,
		Price:       adInput.Price,
		Contacts:    repository.Contacts(adInput.Contacts),
		Published:   adInput.Published,
		ImagesURL:   adInput.ImagesURL,
	})
	if err != nil {
		return 0, err
	}
	return adId, nil
}

func (s *AdService) GetAdById(userId string, adId string)([]domain.Ad, error) {
	ad, err := s.repo.GetAdById(userId, adId)
	if err != nil {
		return []domain.Ad{}, err
	}

	return ad, nil
}

func (s *AdService) UpdateAd(userId string, adId string, ad Ads) error{
	err := s.repo.UpdateAd(userId, adId, repository.Ads{
		Title: ad.Title,
		Category: ad.Category,
		Description: ad.Description,
		Price: ad.Price,
		Contacts: repository.Contacts(ad.Contacts),
		Published: ad.Published,
		ImagesURL: ad.ImagesURL,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AdService) DeleteAd (userId string, adId string) error {
	if err := s.repo.DeleteAd(userId, adId); err != nil {
		return err
	}

	return nil
}
