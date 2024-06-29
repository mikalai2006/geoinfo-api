package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) registerTicket(router *gin.RouterGroup) {
	Ticket := router.Group("/ticket")
	Ticket.GET("/", h.FindTicket)
	Ticket.POST("/", h.CreateTicket)
	Ticket.POST("/list/", h.CreateListTicket)
	Ticket.DELETE("/:id", h.DeleteTicket)
}

func (h *HandlerV1) CreateTicket(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *model.Ticket
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Ticket, err := h.services.Ticket.CreateTicket(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Ticket)
}

func (h *HandlerV1) CreateListTicket(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.Ticket
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.Ticket
	for i := range input {
		Ticket, err := h.services.Ticket.CreateTicket(userID, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, Ticket)

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Ticket Get all Tickets
// @Security ApiKeyAuth
// @Tickets Ticket
// @Description get all Tickets
// @ModuleID Ticket
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Ticket
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Ticket [get].
func (h *HandlerV1) GetAllTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Ticket{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tickets, err := h.services.Ticket.GetAllTicket(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tickets)
}

// @Summary Find Tickets by params
// @Security ApiKeyAuth
// @Tickets Ticket
// @Description Input params for search Tickets
// @ModuleID Ticket
// @Accept  json
// @Produce  json
// @Param input query TicketInput true "params for search Ticket"
// @Success 200 {object} []model.Ticket
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Ticket [get].
func (h *HandlerV1) FindTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.TicketInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tickets, err := h.services.Ticket.FindTicket(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tickets)
}

func (h *HandlerV1) GetTicketByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateTicket(c *gin.Context) {

}

func (h *HandlerV1) DeleteTicket(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Ticket.DeleteTicket(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
