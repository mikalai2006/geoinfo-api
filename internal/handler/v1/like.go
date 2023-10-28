package v1

import (
	"errors"
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

func (h *HandlerV1) registerLike(router *gin.RouterGroup) {
	Like := router.Group("/like")
	Like.GET("", h.FindLike)
	Like.POST("", middleware.SetUserIdentity, h.CreateLike)
	Like.PATCH("/:id", middleware.SetUserIdentity, h.UpdateLike)
	Like.DELETE("/:id", middleware.SetUserIdentity, h.DeleteLike)
}

func (h *HandlerV1) CreateLike(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var data *model.LikeInput
	if er := c.BindJSON(&data); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// Check exists likefor node and user
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
	}
	nodeIDPrimitive, err := primitive.ObjectIDFromHex(data.NodeID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
	}
	LikeExist, err := h.services.Like.FindLike(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"node_id", nodeIDPrimitive}, {"user_id", userIDPrimitive}},
	},
	)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(LikeExist.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, model.ErrLikeExist, nil)
		return
	}

	Like, err := h.services.Like.CreateLike(userID, data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Like)
}

// @Summary Find Likes by params
// @Security ApiKeyAuth
// @Likes Like
// @Description Input params for search Likes
// @ModuleID Like
// @Accept  json
// @Produce  json
// @Param input query LikeInput true "params for search Like"
// @Success 200 {object} []model.Like
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Like [get].
func (h *HandlerV1) FindLike(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.LikeInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Likes, err := h.services.Like.FindLike(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Likes)
}

func (h *HandlerV1) GetLikeByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateLike(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	var data *model.Like
	if er := c.BindJSON(&data); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.services.Like.UpdateLike(id, userID, data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteLike(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Like.DeleteLike(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
