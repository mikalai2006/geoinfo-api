package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id string) (model.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) FindUser(params domain.RequestParams) (domain.Response[model.User], error) {
	return s.repo.FindUser(params)
}

func (s *UserService) CreateUser(userID string, user *model.User) (*model.User, error) {
	return s.repo.CreateUser(userID, user)
}

func (s *UserService) DeleteUser(id string) (model.User, error) {
	return s.repo.DeleteUser(id)
}

func (s *UserService) UpdateUser(id string, user *model.User) (model.User, error) {
	return s.repo.UpdateUser(id, user)
}

func (s *UserService) Iam(userID string) (model.User, error) {
	return s.repo.Iam(userID)
}
