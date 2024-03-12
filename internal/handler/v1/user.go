package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
)

func (h *HandlerV1) RegisterUser(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.POST("", middleware.SetUserIdentity, h.CreateUser)
	user.GET("", h.FindUser)
	user.GET("/:id", h.GetUser)
	user.DELETE("/:id", middleware.SetUserIdentity, h.DeleteUser)
	user.PATCH("/:id", middleware.SetUserIdentity, h.UpdateUser)
}

// @Summary Get user by Id
// @Tags user
// @Description get user info
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [get].
func (h *HandlerV1) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	user, err := h.services.User.GetUser(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// type InputUser struct {
// 	domain.RequestParams
// 	model.User
// }

// @Summary Find few users
// @Security ApiKeyAuth
// @Tags user
// @Description Input params for search users
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param input query model.UserInput true "params for search users"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user [get].
func (h *HandlerV1) FindUser(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.UserInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.services.User.FindUser(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *model.User
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, er := h.services.User.CreateUser(userID, input)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [delete].
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	// var input model.User
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	// 	return
	// }

	user, err := h.services.User.DeleteUser(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Update user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body model.User true "body for update user"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [put].
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	fmt.Println("UpdateUser: ", id)

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	var input model.User
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println("UpdateUser input: ", input)

	user, err := h.services.User.UpdateUser(id, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var imageInput = &model.ImageInput{}
	imageInput.Service = "user"
	imageInput.ServiceID = user.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "user"

	// fmt.Println("imageInput: ", imageInput)

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	var result []model.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}

	c.JSON(http.StatusOK, user)
}
