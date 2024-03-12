package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type NodedataService struct {
	repo repository.Nodedata
}

func NewNodedataService(repo repository.Nodedata) *NodedataService {
	return &NodedataService{repo: repo}
}

func (s *NodedataService) FindNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	return s.repo.FindNodedata(params)
}

func (s *NodedataService) GetAllNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	return s.repo.GetAllNodedata(params)
}

func (s *NodedataService) CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error) {
	return s.repo.CreateNodedata(userID, data)
}

func (s *NodedataService) UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error) {
	return s.repo.UpdateNodedata(id, userID, data)
}

func (s *NodedataService) DeleteNodedata(id string) (model.Nodedata, error) {
	return s.repo.DeleteNodedata(id)
}

func (s *NodedataService) AddAudit(userID string, data *model.NodedataAuditInput) (*model.Nodedata, error) {
	return s.repo.AddAudit(userID, data)
}
