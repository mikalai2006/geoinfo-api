package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type LikeService struct {
	repo repository.Like
}

func NewLikeService(repo repository.Like) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) FindLike(params domain.RequestParams) (domain.Response[model.Like], error) {
	return s.repo.FindLike(params)
}

func (s *LikeService) CreateLike(userID string, like *model.LikeInput) (*model.Like, error) {
	return s.repo.CreateLike(userID, like)
}

func (s *LikeService) UpdateLike(id string, userID string, Like *model.Like) (*model.Like, error) {
	return s.repo.UpdateLike(id, userID, Like)
}

func (s *LikeService) DeleteLike(id string) (model.Like, error) {
	return s.repo.DeleteLike(id)
}
