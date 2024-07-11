package v1

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerNodedataVote(router *gin.RouterGroup) {
	NodedataVote := router.Group("/nodedatavote")
	NodedataVote.GET("", h.FindNodedataVote)
	NodedataVote.POST("", h.CreateNodedataVote)
	NodedataVote.POST("/list", h.CreateNodedataVoteList)
	NodedataVote.PATCH("/:id", h.UpdateNodedataVote)
	NodedataVote.DELETE("/:id", h.DeleteNodedataVote)
}

func (h *HandlerV1) CreateNodedataVote(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	// var input *model.NodedataVote
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// var a map[string]json.RawMessage //  map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// input, er := utils.BindJSON2[model.NodedataVoteInput](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }

	var input *model.NodedataVoteInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	nodedataVote, err := h.CreateOrExistNodedataVote(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, nodedataVote)
}

func (h *HandlerV1) CreateNodedataVoteList(c *gin.Context) {
	appG := app.Gin{C: c}

	var input []*model.NodedataVoteInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.NodedataVote
	for i := range input {
		nodedataVote, err := h.CreateOrExistNodedataVote(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		result = append(result, nodedataVote)
	}

	c.JSON(http.StatusOK, result)
}

func (h *HandlerV1) CreateOrExistNodedataVote(c *gin.Context, input *model.NodedataVoteInput) (*model.NodedataVote, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *model.NodedataVote

	nodedataIDPrimitive, err := primitive.ObjectIDFromHex(input.NodedataID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// check exist nodedata
	existNodedatas, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"_id", nodedataIDPrimitive}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	if len(existNodedatas.Data) == 0 {
		// appG.ResponseError(http.StatusBadRequest, errors.New("not found nodedata"), nil)
		return result, nil
	}

	input.NodedataUserID = existNodedatas.Data[0].UserID
	input.NodeID = existNodedatas.Data[0].NodeID

	// check exist vote
	existNodedataVote, err := h.services.NodedataVote.FindNodedataVote(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter: bson.D{
			{"value", input.Value},
			{"nodedata_id", nodedataIDPrimitive},
			{"nodedata_user_id", input.NodedataUserID},
			{"node_id", input.NodeID},
			{"user_id", userIDPrimitive},
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	if len(existNodedataVote.Data) > 0 {
		// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
		return &existNodedataVote.Data[0], nil
	}

	result, err = h.services.NodedataVote.CreateNodedataVote(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// update counter votes nodedata.
	votes, err := h.services.NodedataVote.FindNodedataVote(domain.RequestParams{
		Filter: bson.D{{"nodedata_id", result.NodedataID}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	like := 0
	dlike := 0
	for i, _ := range votes.Data {
		if votes.Data[i].Value > 0 {
			like += 1
		} else {
			dlike += 1
		}
	}
	status := 100
	if dlike > 5 {
		status = -1
	}
	_, err = h.services.Nodedata.UpdateNodedata(result.NodedataID.Hex(), userID, &model.Nodedata{
		Like:   int64(like),
		Dlike:  int64(math.Abs(float64(dlike))),
		Status: int64(status),
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}

// @Summary NodedataVote Get all NodedataVotes
// @Security ApiKeyAuth
// @NodedataVotes NodedataVote
// @Description get all NodedataVotes
// @ModuleID NodedataVote
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.NodedataVote
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/NodedataVote [get].
func (h *HandlerV1) GetAllNodedataVote(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.NodedataVote{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	NodedataVotes, err := h.services.NodedataVote.GetAllNodedataVote(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, NodedataVotes)
}

// @Summary Find NodedataVotes by params
// @Security ApiKeyAuth
// @NodedataVotes NodedataVote
// @Description Input params for search NodedataVotes
// @ModuleID NodedataVote
// @Accept  json
// @Produce  json
// @Param input query NodedataVoteInput true "params for search NodedataVote"
// @Success 200 {object} []model.NodedataVote
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/NodedataVote [get].
func (h *HandlerV1) FindNodedataVote(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.NodedataVote{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	NodedataVotes, err := h.services.NodedataVote.FindNodedataVote(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, NodedataVotes)
}

func (h *HandlerV1) UpdateNodedataVote(c *gin.Context) {

	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	var input *model.NodedataVote
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.NodedataVote.UpdateNodedataVote(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteNodedataVote(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.NodedataVote.DeleteNodedataVote(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
