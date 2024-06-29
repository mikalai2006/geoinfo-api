package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerNodedata(router *gin.RouterGroup) {
	nodedata := router.Group("/nodedata")
	nodedata.GET("", h.FindNodedata)
	nodedata.POST("", h.CreateNodedata)
	nodedata.POST("/audit", h.AddAudit)
	nodedata.POST("/list", h.CreateListNodedata)
	nodedata.DELETE("/:id", h.DeleteNodedata)

	nodedata.GET("/check_nodes", h.CheckNodes)
}

func (h *HandlerV1) CreateNodedata(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	// var input *model.Nodedata
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.NodedataInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Nodedata, err := h.CreateOrExistNodedata(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodedata)
}

func (h *HandlerV1) CreateListNodedata(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input []*model.NodedataInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	result := []*model.Nodedata{}
	for i := range input {

		nodedata, err := h.CreateOrExistNodedata(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, nodedata)
		// nodeIDPrimitive, err := primitive.ObjectIDFromHex(input[i].NodeID)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// tagIDPrimitive, err := primitive.ObjectIDFromHex(input[i].TagID)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// // existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		// // 	Options: domain.Options{Limit: 1},
		// // 	Filter:  bson.D{{"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}},
		// // })
		// // if err != nil {
		// // 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// // 	return
		// // }

		// existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		// 	Options: domain.Options{Limit: 1},
		// 	Filter:  bson.D{{"data.value", input[i].Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}}, // {"tag_id", input.TagID},
		// })
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// if len(existNodedata.Data) == 0 {
		// 	Nodedata, err := h.services.Nodedata.CreateNodedata(userID, input[i])
		// 	if err != nil {
		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
		// 		return
		// 	}
		// 	result = append(result, Nodedata)
		// } else {
		// 	result = append(result, &existNodedata.Data[0])
		// }
		// // else {
		// // 	fmt.Println("Exist data for ", existNodedata.Data[0])
		// // }

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Nodedata Get all Nodedatas
// @Security ApiKeyAuth
// @Nodedatas Nodedata
// @Description get all Nodedatas
// @ModuleID Nodedata
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Nodedata
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Nodedata [get].
func (h *HandlerV1) GetAllNodedata(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Nodedata{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Nodedatas, err := h.services.Nodedata.GetAllNodedata(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodedatas)
}

// @Summary Find Nodedatas by params
// @Security ApiKeyAuth
// @Nodedatas Nodedata
// @Description Input params for search Nodedatas
// @ModuleID Nodedata
// @Accept  json
// @Produce  json
// @Param input query NodedataInput true "params for search Nodedata"
// @Success 200 {object} []model.Nodedata
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Nodedata [get].
func (h *HandlerV1) FindNodedata(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Nodedata{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Nodedatas, err := h.services.Nodedata.FindNodedata(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodedatas)
}

func (h *HandlerV1) GetNodedataByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateNodedata(c *gin.Context) {

}

func (h *HandlerV1) DeleteNodedata(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	nodedata, err := h.services.Nodedata.DeleteNodedata(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// remove nodedata votes.
	nodedataVotes, err := h.services.NodedataVote.FindNodedataVote(domain.RequestParams{
		Filter: bson.D{{"nodedata_id", nodedata.ID}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(nodedataVotes.Data) > 0 {
		for i, _ := range nodedataVotes.Data {
			_, err := h.services.NodedataVote.DeleteNodedataVote(nodedataVotes.Data[i].ID.Hex())
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
		}
	}

	c.JSON(http.StatusOK, nodedata)
}

func (h *HandlerV1) AddAudit(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var a map[string]json.RawMessage
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.NodedataAuditInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// nodedataIDPrimitive, err := primitive.ObjectIDFromHex(input.NodedataID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"data.value", input.Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}}, // {"tag_id", input.TagID},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// fmt.Println(bson.D{{"data.value", input.Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}})
	// if len(existNodedata.Data) > 0 {
	// 	appG.ResponseError(http.StatusBadRequest, model.ErrNodedataExistValue, nil)
	// 	return
	// }

	Nodedata, err := h.services.Nodedata.AddAudit(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodedata)
}

func (h *HandlerV1) CreateOrExistNodedata(c *gin.Context, input *model.NodedataInput) (*model.Nodedata, error) {
	appG := app.Gin{C: c}

	var result *model.Nodedata

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return result, err
	}

	tagIDPrimitive, err := primitive.ObjectIDFromHex(input.TagID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	nodeIDPrimitive, err := primitive.ObjectIDFromHex(input.NodeID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// check exist node
	existNodes, err := h.services.Node.FindNode(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"_id", nodeIDPrimitive}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	if len(existNodes.Data) == 0 {
		// appG.ResponseError(http.StatusBadRequest, model.ErrNodeNotFound, nil)
		return result, nil
	}

	// check exist nodedata
	existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}, {"data.value", input.Data.Value}}, //  {"tag_id", input.TagID},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	// fmt.Println(existNodedata)

	if len(existNodedata.Data) > 0 {
		// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataExistValue, nil)
		return &existNodedata.Data[0], nil
	}

	result, err = h.services.Nodedata.CreateNodedata(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}

	return result, nil
}

func (h *HandlerV1) CheckNodes(c *gin.Context) {
	appG := app.Gin{C: c}
	var result []string

	// implementation roles for user.
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if !utils.Contains(roles, "admin") {
		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
		return
	}

	allNodes, err := h.services.Node.FindNode(domain.RequestParams{Filter: bson.D{}, Options: domain.Options{Limit: 1000000}})
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	fmt.Println("allNodes", len(allNodes.Data))
	for i, _ := range allNodes.Data {
		res, err := h.services.Node.FindNode(domain.RequestParams{Filter: bson.D{
			{"_id", allNodes.Data[i].ID},
		}, Options: domain.Options{Limit: 1}})
		if err != nil {
			appG.ResponseError(http.StatusUnauthorized, err, nil)
			return
		}
		if len(res.Data) == 0 {
			result = append(result, res.Data[0].ID.Hex())
		}
	}

	c.JSON(http.StatusOK, result)
}
