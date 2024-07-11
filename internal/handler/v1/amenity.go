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

func (h *HandlerV1) registerAmenity(router *gin.RouterGroup) {
	Amenity := router.Group("/amenity")
	Amenity.GET("", h.FindAmenity)
	Amenity.POST("", h.SetUserFromRequest, h.CreateAmenity)
	Amenity.POST("/list/", h.SetUserFromRequest, h.CreateListAmenity)
	Amenity.PATCH("/:id", h.SetUserFromRequest, h.UpdateAmenity)
	Amenity.DELETE("/:id", h.SetUserFromRequest, h.DeleteAmenity)
}

func (h *HandlerV1) CreateAmenity(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// // var input *model.Amenity
	// // if er := c.BindJSON(&input); er != nil {
	// // 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// // 	return
	// // }
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.Amenity](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage //  map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.Amenity](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// find exist amenity.
	existAmenity, err := h.services.Amenity.FindAmenity(domain.RequestParams{Filter: bson.D{
		{"type", data.Type},
	}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(existAmenity.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("exist amenity"), nil)
		return
	}

	Amenity, err := h.services.Amenity.CreateAmenity(userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Amenity)
}

func (h *HandlerV1) CreateListAmenity(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.Amenity
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.Amenity
	for i := range input {
		existOsmID, err := h.services.Amenity.FindAmenity(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"key", input[i].Key}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		if len(existOsmID.Data) == 0 {
			Amenity, err := h.services.Amenity.CreateAmenity(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, Amenity)
		}

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Amenity Get all Amenitys
// @Security ApiKeyAuth
// @Amenitys Amenity
// @Description get all Amenitys
// @ModuleID Amenity
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Amenity
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Amenity [get].
func (h *HandlerV1) GetAllAmenity(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Amenity{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Amenitys, err := h.services.Amenity.GetAllAmenity(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Amenitys)
}

// @Summary Find Amenitys by params
// @Security ApiKeyAuth
// @Amenitys Amenity
// @Description Input params for search Amenitys
// @ModuleID Amenity
// @Accept  json
// @Produce  json
// @Param input query AmenityInput true "params for search Amenity"
// @Success 200 {object} []model.Amenity
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Amenity [get].
func (h *HandlerV1) FindAmenity(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.AmenityInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Amenitys, err := h.services.Amenity.FindAmenity(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Amenitys)
}

func (h *HandlerV1) GetAmenityByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAmenity(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// // var input model.AmenityInput
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
	// data, er := utils.BindJSON[model.Amenity](a)
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
	data, er := utils.BindJSON2[model.Amenity](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Amenity.UpdateAmenity(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteAmenity(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Amenity.DeleteAmenity(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
