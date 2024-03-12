package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) registerReview(router *gin.RouterGroup) {
	review := router.Group("/review")
	review.GET("/", h.FindReview)
	review.POST("", middleware.SetUserIdentity, h.CreateReview)
}

func (h *HandlerV1) CreateReview(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *model.Review
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	review, err := h.services.Review.CreateReview(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary Review Get all reviews
// @Security ApiKeyAuth
// @Tags review
// @Description get all reviews
// @ModuleID review
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Review
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/review [get].
func (h *HandlerV1) GetAllReview(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Review{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	reviews, err := h.services.Review.GetAllReview(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// @Summary Find reviews by params
// @Security ApiKeyAuth
// @Tags review
// @Description Input params for search reviews
// @ModuleID review
// @Accept  json
// @Produce  json
// @Param input query ReviewInput true "params for search review"
// @Success 200 {object} []domain.Review
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/review [get].
func (h *HandlerV1) FindReview(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.ReviewInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	reviews, err := h.services.Review.FindReview(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *HandlerV1) GetReviewByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateReview(c *gin.Context) {

}

func (h *HandlerV1) DeleteReview(c *gin.Context) {

}
