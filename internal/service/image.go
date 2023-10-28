package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type ImageService struct {
	repo repository.Image
}

func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) FindImage(params domain.RequestParams) (domain.Response[model.Image], error) {
	return s.repo.FindImage(params)
}

func (s *ImageService) GetImage(id string) (model.Image, error) {
	return s.repo.GetImage(id)
}

func (s *ImageService) GetImageDirs(id string) ([]interface{}, error) {
	return s.repo.GetImageDirs(id)
}
func (s *ImageService) CreateImage(userID string, image *model.ImageInput) (model.Image, error) {
	return s.repo.CreateImage(userID, image)
}

func (s *ImageService) DeleteImage(id string) (model.Image, error) {
	return s.repo.DeleteImage(id)
}
