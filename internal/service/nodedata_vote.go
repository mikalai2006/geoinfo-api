package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type NodedataVoteService struct {
	repo repository.NodedataVote
}

func NewNodedataVoteService(repo repository.NodedataVote) *NodedataVoteService {
	return &NodedataVoteService{repo: repo}
}

func (s *NodedataVoteService) FindNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	return s.repo.FindNodedataVote(params)
}

func (s *NodedataVoteService) GetAllNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	return s.repo.GetAllNodedataVote(params)
}

func (s *NodedataVoteService) CreateNodedataVote(userID string, data *model.NodedataVoteInput) (*model.NodedataVote, error) {
	return s.repo.CreateNodedataVote(userID, data)
}

func (s *NodedataVoteService) UpdateNodedataVote(id string, userID string, data *model.NodedataVote) (*model.NodedataVote, error) {
	return s.repo.UpdateNodedataVote(id, userID, data)
}

func (s *NodedataVoteService) DeleteNodedataVote(id string) (model.NodedataVote, error) {
	return s.repo.DeleteNodedataVote(id)
}
