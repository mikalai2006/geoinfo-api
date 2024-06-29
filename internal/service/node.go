package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type NodeService struct {
	repo        repository.Node
	userService *UserService
}

func NewNodeService(repo repository.Node, userService *UserService) *NodeService {
	return &NodeService{repo: repo, userService: userService}
}

func (s *NodeService) FindNode(params domain.RequestParams) (domain.Response[model.Node], error) {
	return s.repo.FindNode(params)
}

func (s *NodeService) CreateFile(params domain.RequestParams) (domain.Response[domain.NodeFileItem], error) {
	return s.repo.CreateFile(params)
}

func (s *NodeService) FindForKml(params domain.RequestParams) (domain.Response[domain.Kml], error) {
	return s.repo.FindForKml(params)
}

func (s *NodeService) GetAllNode(params domain.RequestParams) (domain.Response[model.Node], error) {
	return s.repo.GetAllNode(params)
}

func (s *NodeService) CreateNode(userID string, node *model.Node) (*model.Node, error) {
	result, err := s.repo.CreateNode(userID, node)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(userID, model.UserStat{Node: 1})
	}

	return result, err
}

func (s *NodeService) UpdateNode(id string, userID string, data *model.Node) (*model.Node, error) {
	return s.repo.UpdateNode(id, userID, data)
}

func (s *NodeService) DeleteNode(id string) (model.Node, error) {
	result, err := s.repo.DeleteNode(id)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(result.UserID.Hex(), model.UserStat{Node: -1})
	}

	return result, err
}
