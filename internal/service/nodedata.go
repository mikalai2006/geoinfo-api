package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type NodedataService struct {
	repo                repository.Nodedata
	userService         *UserService
	nodedataVoteService *NodedataVoteService
}

func NewNodedataService(repo repository.Nodedata, userService *UserService, nodedataVoteService *NodedataVoteService) *NodedataService {
	return &NodedataService{repo: repo, userService: userService, nodedataVoteService: nodedataVoteService}
}

func (s *NodedataService) FindNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	return s.repo.FindNodedata(params)
}

func (s *NodedataService) GetAllNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	return s.repo.GetAllNodedata(params)
}

func (s *NodedataService) CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error) {
	result, err := s.repo.CreateNodedata(userID, data)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(userID, model.UserStat{Nodedata: 1})
	}

	return result, err
}

func (s *NodedataService) UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error) {
	return s.repo.UpdateNodedata(id, userID, data)
}

func (s *NodedataService) DeleteNodedata(id string) (model.Nodedata, error) {
	result := model.Nodedata{}

	removedNodedata, err := s.repo.GetNodedata(id)
	if err != nil {
		return result, err
	}

	// find all nodedata_vote for remove.
	nodedataVotes, err := s.nodedataVoteService.FindNodedataVote(domain.RequestParams{
		Filter:  bson.M{"nodedata_id": removedNodedata.ID},
		Options: domain.Options{Limit: 1000},
	})
	if err != nil {
		return result, err
	}

	for i, _ := range nodedataVotes.Data {
		_, err := s.nodedataVoteService.DeleteNodedataVote(nodedataVotes.Data[i].ID.Hex())
		if err != nil {
			return result, err
		}
		// fmt.Println("Remove nodedata: ", nodedata.Data[i].ID.Hex())
	}

	result, err = s.repo.DeleteNodedata(id)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(result.UserID.Hex(), model.UserStat{Nodedata: -1})
	}

	return result, err
}

func (s *NodedataService) AddAudit(userID string, data *model.NodedataAuditInput) (*model.Nodedata, error) {
	return s.repo.AddAudit(userID, data)
}
