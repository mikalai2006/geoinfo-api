package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type TagoptService struct {
	repo repository.Tagopt
}

func NewTagoptService(repo repository.Tagopt) *TagoptService {
	return &TagoptService{repo: repo}
}

func (s *TagoptService) FindTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error) {
	return s.repo.FindTagopt(params)
}

func (s *TagoptService) GetAllTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error) {
	return s.repo.GetAllTagopt(params)
}

func (s *TagoptService) CreateTagopt(userID string, tag *model.TagoptInput) (*model.Tagopt, error) {
	return s.repo.CreateTagopt(userID, tag)
}

func (s *TagoptService) UpdateTagopt(id string, userID string, data *model.TagoptInput) (*model.Tagopt, error) {
	return s.repo.UpdateTagopt(id, userID, data)
}

func (s *TagoptService) DeleteTagopt(id string) (model.Tagopt, error) {
	return s.repo.DeleteTagopt(id)
}
