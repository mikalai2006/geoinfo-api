package service

import (
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type ReviewService struct {
	repo repository.Review
}

func NewReviewService(repo repository.Review) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) FindReview(params domain.RequestParams) (domain.Response[domain.Review], error) {
	return s.repo.FindReview(params)
}

func (s *ReviewService) GetAllReview(params domain.RequestParams) (domain.Response[domain.Review], error) {
	return s.repo.GetAllReview(params)
}

func (s *ReviewService) CreateReview(userID string, review *domain.Review) (*domain.Review, error) {
	return s.repo.CreateReview(userID, review)
}
