package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) RegisterApp(router *gin.RouterGroup) {
	lang := router.Group("/lang")
	lang.POST("", h.SetUserFromRequest, h.createLanguage)
	lang.GET("", h.findLanguage)
	lang.GET("/:id", h.getLanguage)
	lang.PATCH("/:id", h.SetUserFromRequest, h.updateLanguage)
	lang.DELETE("/:id", h.SetUserFromRequest, h.deleteLanguage)
}

func (h *HandlerV1) createLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.LanguageInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Apps.CreateLanguage(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) getLanguage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Apps.GetLanguage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) findLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.LanguageInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.services.Apps.FindLanguage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *HandlerV1) updateLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.LanguageInput
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.services.Apps.UpdateLanguage(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) deleteLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}
	// var input domain.Page
	// if err := c.BindJSON(&input); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)

	// 	return
	// }

	document, err := h.services.Apps.DeleteLanguage(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
