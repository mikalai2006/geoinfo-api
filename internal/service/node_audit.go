package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
)

type NodeAuditService struct {
	repo repository.NodeAudit
}

func NewNodeAuditService(repo repository.NodeAudit) *NodeAuditService {
	return &NodeAuditService{repo: repo}
}

func (s *NodeAuditService) FindNodeAudit(params domain.RequestParams) (domain.Response[model.NodeAudit], error) {
	return s.repo.FindNodeAudit(params)
}

func (s *NodeAuditService) CreateNodeAudit(userID string, nodeAudit *model.NodeAuditInput) (*model.NodeAudit, error) {
	return s.repo.CreateNodeAudit(userID, nodeAudit)
}

func (s *NodeAuditService) UpdateNodeAudit(id string, userID string, data *model.NodeAuditInput) (*model.NodeAudit, error) {
	return s.repo.UpdateNodeAudit(id, userID, data)
}

func (s *NodeAuditService) DeleteNodeAudit(id string) (model.NodeAudit, error) {
	return s.repo.DeleteNodeAudit(id)
}
