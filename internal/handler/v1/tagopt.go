package v1

import (
	"encoding/json"
	"errors"
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

func (h *HandlerV1) registerTagopt(router *gin.RouterGroup) {
	Tagopt := router.Group("/tagopt")
	Tagopt.GET("/", h.FindTagopt)
	Tagopt.POST("/", middleware.SetUserIdentity, h.CreateTagopt)
	Tagopt.POST("/list/", middleware.SetUserIdentity, h.CreateListTagopt)
	Tagopt.PATCH("/:id", middleware.SetUserIdentity, h.UpdateTagopt)
	Tagopt.DELETE("/:id", middleware.SetUserIdentity, h.DeleteTagopt)
}

func (h *HandlerV1) CreateTagopt(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// var input *model.Tagopt
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	input, er := utils.BindJSON2[model.TagoptInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	existTagopt, err := h.services.Tagopt.FindTagopt(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"value", input.Value}}, // {"tag_id", input.TagID},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if len(existTagopt.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, model.ErrTagOptExistValue, nil)
		return
	}

	Tagopt, err := h.services.Tagopt.CreateTagopt(userID, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tagopt)
}

func (h *HandlerV1) CreateListTagopt(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.TagoptInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	result := []*model.Tagopt{}
	for i := range input {
		existTagopt, err := h.services.Tagopt.FindTagopt(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"value", input[i].Value}, {"tag_id", input[i].TagID}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		if len(existTagopt.Data) == 0 {
			Tagopt, err := h.services.Tagopt.CreateTagopt(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, Tagopt)
		}

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Tagopt Get all Tagopts
// @Security ApiKeyAuth
// @Tagopts Tagopt
// @Description get all Tagopts
// @ModuleID Tagopt
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Tagopt
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Tagopt [get].
func (h *HandlerV1) GetAllTagopt(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Tagopt{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tagopts, err := h.services.Tagopt.GetAllTagopt(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tagopts)
}

// @Summary Find Tagopts by params
// @Security ApiKeyAuth
// @Tagopts Tagopt
// @Description Input params for search Tagopts
// @ModuleID Tagopt
// @Accept  json
// @Produce  json
// @Param input query TagoptInput true "params for search Tagopt"
// @Success 200 {object} []model.Tagopt
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Tagopt [get].
func (h *HandlerV1) FindTagopt(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.TagoptInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tagopts, err := h.services.Tagopt.FindTagopt(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tagopts)
}

func (h *HandlerV1) GetTagoptByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateTagopt(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// var input model.TagInput
	// data, err := utils.BindAndValid(c, &input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	var a map[string]json.RawMessage // map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	data, er := utils.BindJSON2[model.TagoptInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.services.Tagopt.UpdateTagopt(id, userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteTagopt(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Tagopt.DeleteTagopt(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
