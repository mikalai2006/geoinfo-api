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
	Nodedata := router.Group("/nodedata")
	Nodedata.GET("/", h.FindNodedata)
	Nodedata.POST("/", middleware.SetUserIdentity, h.CreateNodedata)
	Nodedata.POST("/list/", middleware.SetUserIdentity, h.CreateListNodedata)
	Nodedata.DELETE("/:id", middleware.SetUserIdentity, h.DeleteNodedata)
}

func (h *HandlerV1) CreateNodedata(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

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

	tagIDPrimitive, err := primitive.ObjectIDFromHex(input.TagID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	nodeIDPrimitive, err := primitive.ObjectIDFromHex(input.NodeID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"value", input.Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}}, // {"tag_id", input.TagID},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	fmt.Println(bson.D{{"value", input.Data.Value}, {"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}})

	if len(existNodedata.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, model.ErrNodedataExistValue, nil)
		return
	}

	Nodedata, err := h.services.Nodedata.CreateNodedata(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodedata)
}

func (h *HandlerV1) CreateListNodedata(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

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
		// existNodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		// 	Options: domain.Options{Limit: 1},
		// 	Filter:  bson.D{{"node_id", nodeIDPrimitive}, {"tag_id", tagIDPrimitive}},
		// })
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }

		// if len(existNodedata.Data) == 0 {
		Nodedata, err := h.services.Nodedata.CreateNodedata(userID, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Nodedata)
		// }
		// else {
		// 	fmt.Println("Exist data for ", existNodedata.Data[0])
		// }

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

	user, err := h.services.Nodedata.DeleteNodedata(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
