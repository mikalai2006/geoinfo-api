package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type AmenityService struct {
	repo repository.Amenity
	i18n config.I18nConfig
}

func NewAmenityService(repo repository.Amenity, i18n config.I18nConfig) *AmenityService {
	return &AmenityService{repo: repo, i18n: i18n}
}

func (s *AmenityService) FindAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error) {
	return s.repo.FindAmenity(params)
}

func (s *AmenityService) GetAllAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error) {
	return s.repo.GetAllAmenity(params)
}

func (s *AmenityService) CreateAmenity(userID string, Amenity *model.Amenity) (*model.Amenity, error) {
	return s.repo.CreateAmenity(userID, Amenity)
}

func (s *AmenityService) UpdateAmenity(id string, userID string, Amenity *model.Amenity) (*model.Amenity, error) {
	return s.repo.UpdateAmenity(id, userID, Amenity)
}

func (s *AmenityService) DeleteAmenity(id string) (model.Amenity, error) {
	return s.repo.DeleteAmenity(id)
}
