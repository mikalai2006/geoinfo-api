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
)

func (h *HandlerV1) registerAmenityGroup(router *gin.RouterGroup) {
	AmenityGroup := router.Group("/amenitygroup")
	AmenityGroup.GET("", h.FindAmenityGroup)
	AmenityGroup.POST("", h.CreateAmenityGroup)
	AmenityGroup.POST("/list/", h.CreateListAmenityGroup)
	AmenityGroup.PATCH("/:id", h.UpdateAmenityGroup)
	AmenityGroup.DELETE("/:id", h.DeleteAmenityGroup)
}

func (h *HandlerV1) CreateAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// // var input *model.AmenityGroup
	// // if er := c.BindJSON(&input); er != nil {
	// // 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// // 	return
	// // }
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.AmenityGroup](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.AmenityGroup](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	AmenityGroup, err := h.services.AmenityGroup.CreateAmenityGroup(userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, AmenityGroup)
}

func (h *HandlerV1) CreateListAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.AmenityGroup
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.AmenityGroup
	for i := range input {
		existOsmID, err := h.services.AmenityGroup.FindAmenityGroup(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"title", input[i].Title}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		if len(existOsmID.Data) == 0 {
			AmenityGroup, err := h.services.AmenityGroup.CreateAmenityGroup(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, AmenityGroup)
		}

	}

	c.JSON(http.StatusOK, result)
}

// @Summary AmenityGroup Get all AmenityGroups
// @Security ApiKeyAuth
// @AmenityGroups AmenityGroup
// @Description get all AmenityGroups
// @ModuleID AmenityGroup
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.AmenityGroup
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/AmenityGroup [get].
func (h *HandlerV1) GetAllAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.AmenityGroup{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	AmenityGroups, err := h.services.AmenityGroup.GetAllAmenityGroup(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, AmenityGroups)
}

// @Summary Find AmenityGroups by params
// @Security ApiKeyAuth
// @AmenityGroups AmenityGroup
// @Description Input params for search AmenityGroups
// @ModuleID AmenityGroup
// @Accept  json
// @Produce  json
// @Param input query AmenityGroupInput true "params for search AmenityGroup"
// @Success 200 {object} []model.AmenityGroup
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/AmenityGroup [get].
func (h *HandlerV1) FindAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.AmenityGroupInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	AmenityGroups, err := h.services.AmenityGroup.FindAmenityGroup(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, AmenityGroups)
}

func (h *HandlerV1) GetAmenityGroupByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// // var input model.AmenityGroupInput
	// // data, err := utils.BindAndValid(c, &input)
	// // if err != nil {
	// // 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// // 	return
	// // }
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.AmenityGroup](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// // fmt.Println(data)

	fmt.Println("hello1")
	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.AmenityGroup](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.AmenityGroup.UpdateAmenityGroup(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteAmenityGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.AmenityGroup.DeleteAmenityGroup(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
