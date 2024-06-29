package service

import (
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodedataVoteService struct {
	repo         repository.NodedataVote
	repoUser     repository.User
	repoNodedata repository.Nodedata
}

func NewNodedataVoteService(repo repository.NodedataVote, repoUser repository.User, repoNodedata repository.Nodedata) *NodedataVoteService {
	return &NodedataVoteService{repo: repo, repoUser: repoUser, repoNodedata: repoNodedata}
}

func (s *NodedataVoteService) FindNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	return s.repo.FindNodedataVote(params)
}

func (s *NodedataVoteService) GetAllNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	return s.repo.GetAllNodedataVote(params)
}

func (s *NodedataVoteService) CreateNodedataVote(userID string, data *model.NodedataVoteInput) (*model.NodedataVote, error) {
	var result *model.NodedataVote

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	nodedataIDPrimitive, err := primitive.ObjectIDFromHex(data.NodedataID)
	if err != nil {
		return nil, err
	}
	existNodedataVote, err := s.repo.FindNodedataVote(domain.RequestParams{
		Filter:  bson.M{"nodedata_id": nodedataIDPrimitive, "user_id": userIDPrimitive},
		Options: domain.Options{Limit: 1},
	})
	if err != nil {
		return nil, err
	}
	if len(existNodedataVote.Data) > 0 {
		// fmt.Println("Update nodedata Vote")
		updateNodedataVote := &model.NodedataVote{
			// UserID:     userIDPrimitive,
			// NodedataID: nodedataIDPrimitive,
			Value: data.Value,
		}
		return s.UpdateNodedataVote(existNodedataVote.Data[0].ID.Hex(), userID, updateNodedataVote)
		// result, err = s.repo.UpdateNodedataVote(existNodedataVote.Data[0].ID.Hex(), userID, updateNodedataVote)
		// if err != nil {
		// 	return nil, err
		// }
	} else {
		// fmt.Println("Create nodedata Vote")
		result, err = s.repo.CreateNodedataVote(userID, data)
		if err != nil {
			return nil, err
		}

		// set user stat.
		if err == nil {
			// set for user.
			statFragment := model.UserStat{}
			if result.Value > 0 {
				statFragment.NodedataLike = 1
			} else {
				statFragment.NodedataDLike = 1
			}
			_, _ = s.repoUser.SetStat(userID, statFragment)

			// set for author.
			nodedata, err := s.repoNodedata.GetNodedata(result.NodedataID.Hex())
			if err == nil {
				statFragmentAuthor := model.UserStat{}
				if result.Value > 0 {
					statFragmentAuthor.NodedataAuthorLike = 1
				} else {
					statFragmentAuthor.NodedataAuthorDLike = 1
				}
				_, _ = s.repoUser.SetStat(nodedata.UserID.Hex(), statFragmentAuthor)
			}
		}

	}
	return result, err
}

func (s *NodedataVoteService) UpdateNodedataVote(id string, userID string, data *model.NodedataVote) (*model.NodedataVote, error) {
	result, err := s.repo.UpdateNodedataVote(id, userID, data)

	// set user stat.
	if err == nil {
		// set for user.
		statFragment := model.UserStat{}
		if result.Value > 0 {
			statFragment.NodedataLike = 1
			statFragment.NodedataDLike = -1
		} else {
			statFragment.NodedataLike = -1
			statFragment.NodedataDLike = 1
		}
		_, _ = s.repoUser.SetStat(userID, statFragment)

		// set for author.
		nodedata, err := s.repoNodedata.GetNodedata(result.NodedataID.Hex())
		if err == nil {
			statFragmentAuthor := model.UserStat{}
			if result.Value > 0 {
				statFragmentAuthor.NodedataAuthorLike = 1
				statFragmentAuthor.NodedataAuthorDLike = -1
			} else {
				statFragmentAuthor.NodedataAuthorLike = -1
				statFragmentAuthor.NodedataAuthorDLike = 1
			}
			_, _ = s.repoUser.SetStat(nodedata.UserID.Hex(), statFragmentAuthor)
		}
	}

	return result, err
}

func (s *NodedataVoteService) DeleteNodedataVote(id string) (model.NodedataVote, error) {
	result, err := s.repo.DeleteNodedataVote(id)

	// set user stat.
	if err == nil {
		// set for user.
		statFragment := model.UserStat{}
		if result.Value > 0 {
			statFragment.NodedataLike = -1
		} else {
			statFragment.NodedataDLike = -1
		}
		_, _ = s.repoUser.SetStat(result.UserID.Hex(), statFragment)

		// set for author.
		nodedata, err := s.repoNodedata.GetNodedata(result.NodedataID.Hex())
		if err == nil {
			statFragmentAuthor := model.UserStat{}
			if result.Value > 0 {
				statFragmentAuthor.NodedataAuthorLike = -1
			} else {
				statFragmentAuthor.NodedataAuthorDLike = -1
			}
			_, _ = s.repoUser.SetStat(nodedata.UserID.Hex(), statFragmentAuthor)
		}
	}

	return result, err
}
