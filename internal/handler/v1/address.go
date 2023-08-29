package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) registerAddress(router *gin.RouterGroup) {
	address := router.Group("/address")
	address.GET("/", h.FindAddress)
	address.POST("/", middleware.SetUserIdentity, h.CreateAddress)
}

func (h *HandlerV1) CreateAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("language")
	if lang == "" {
		lang = h.i18n.Default
	}
	fmt.Println("lang", lang)

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.AddressInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	if input.Lang == "" {
		input.Lang = lang
	}

	address, err := h.services.Address.CreateAddress(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, address)
}

// @Summary Address Get all Address
// @Security ApiKeyAuth
// @Tags address
// @Description get all Address
// @ModuleID address
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Address
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/address [get].
func (h *HandlerV1) GetAllAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Address{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	addresses, err := h.services.Address.GetAllAddress(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// @Summary Address by params
// @Security ApiKeyAuth
// @Tags address
// @Description Input params for search Addresses
// @ModuleID address
// @Accept  json
// @Produce  json
// @Param input query AddressInput true "params for search Address"
// @Success 200 {object} []domain.Address
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/address [get].
func (h *HandlerV1) FindAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.AddressInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	addresses, err := h.services.Address.FindAddress(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *HandlerV1) GetAddressByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAddress(c *gin.Context) {

}

func (h *HandlerV1) DeleteAddress(c *gin.Context) {

}
