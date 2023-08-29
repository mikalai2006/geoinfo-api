package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) registerTrack(router *gin.RouterGroup) {
	track := router.Group("/track")
	track.GET("/", h.FindTrack)
	track.POST("/", middleware.SetUserIdentity, h.CreateTrack)
}

func (h *HandlerV1) CreateTrack(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.Track
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Track, err := h.services.Track.CreateTrack(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Track)
}

// @Summary Track Get all Tracks
// @Security ApiKeyAuth
// @Tags Track
// @Description get all Tracks
// @ModuleID Track
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Track
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Track [get].
func (h *HandlerV1) GetAllTrack(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Track{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tracks, err := h.services.Track.GetAllTrack(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tracks)
}

// @Summary Find Tracks by params
// @Security ApiKeyAuth
// @Tags Track
// @Description Input params for search Tracks
// @ModuleID Track
// @Accept  json
// @Produce  json
// @Param input query TrackInput true "params for search Track"
// @Success 200 {object} []domain.Track
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Track [get].
func (h *HandlerV1) FindTrack(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.TrackInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tracks, err := h.services.Track.FindTrack(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tracks)
}

func (h *HandlerV1) GetTrackByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateTrack(c *gin.Context) {

}

func (h *HandlerV1) DeleteTrack(c *gin.Context) {

}
