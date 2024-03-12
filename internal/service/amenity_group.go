package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type AmenityGroupService struct {
	repo repository.AmenityGroup
	i18n config.I18nConfig
}

func NewAmenityGroupService(repo repository.AmenityGroup, i18n config.I18nConfig) *AmenityGroupService {
	return &AmenityGroupService{repo: repo, i18n: i18n}
}

func (s *AmenityGroupService) FindAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error) {
	return s.repo.FindAmenityGroup(params)
}

func (s *AmenityGroupService) GetAllAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error) {
	return s.repo.GetAllAmenityGroup(params)
}

func (s *AmenityGroupService) CreateAmenityGroup(userID string, AmenityGroup *model.AmenityGroup) (*model.AmenityGroup, error) {
	return s.repo.CreateAmenityGroup(userID, AmenityGroup)
}

func (s *AmenityGroupService) UpdateAmenityGroup(id string, userID string, AmenityGroup *model.AmenityGroup) (*model.AmenityGroup, error) {
	return s.repo.UpdateAmenityGroup(id, userID, AmenityGroup)
}

func (s *AmenityGroupService) DeleteAmenityGroup(id string) (model.AmenityGroup, error) {
	return s.repo.DeleteAmenityGroup(id)
}
