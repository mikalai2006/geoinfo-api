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

func (h *HandlerV1) registerNodeAudit(router *gin.RouterGroup) {
	nodeAudit := router.Group("/node_audit")
	nodeAudit.POST("", h.CreateNodeAudit)
	nodeAudit.POST("/list", h.CreateListNodeAudit)
	nodeAudit.GET("", h.FindNodeAudit)
	nodeAudit.PATCH("/:id", h.UpdateNodeAudit)
	nodeAudit.DELETE("/:id", h.DeleteNodeAudit)
}

func (h *HandlerV1) CreateNodeAudit(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *model.NodeAuditInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// node, err := h.services.NodeAudit.CreateNodeAudit(userID, input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	node, err := h.CreateOrExistNodeAudit(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateListNodeAudit(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input []*model.NodeAuditInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.NodeAudit
	for i := range input {
		nodeAudit, err := h.CreateOrExistNodeAudit(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, nodeAudit)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Find NodeAudits by params
// @Security ApiKeyAuth
// @Tags NodeAudit
// @Description Input params for search NodeAudits
// @ModuleID NodeAudit
// @Accept  json
// @Produce  json
// @Param input query NodeAudit true "params for search NodeAudit"
// @Success 200 {object} []domain.NodeAudit
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/node_audit [get].
func (h *HandlerV1) FindNodeAudit(c *gin.Context) {
	appG := app.Gin{C: c}

	// authData, err := middleware.GetAuthFromCtx(c)
	// fmt.Println("auth ", authData.Roles)

	params, err := utils.GetParamsFromRequest(c, model.NodeAudit{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println(params)
	Nodes, err := h.services.NodeAudit.FindNodeAudit(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodes)
}

func (h *HandlerV1) UpdateNodeAudit(c *gin.Context) {
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
	// var a map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON[model.Node](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// fmt.Println(data)
	var input *model.NodeAuditInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.NodeAudit.UpdateNodeAudit(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteNodeAudit(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	// implementation roles for user.
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if !utils.Contains(roles, "admin") {
		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
		return
	}

	node, err := h.services.NodeAudit.DeleteNodeAudit(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateOrExistNodeAudit(c *gin.Context, input *model.NodeAuditInput) (*model.NodeAudit, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	// nodeIDPrimitive, err := primitive.ObjectIDFromHex(string(input.NodeID))
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return nil, err
	// }
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}

	var result *model.NodeAudit

	// check exist node
	existNode, err := h.services.Node.FindNode(domain.RequestParams{
		Filter: bson.D{
			{"_id", input.NodeID},
		},
		Options: domain.Options{
			Limit: 1,
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	if len(existNode.Data) == 0 {
		//appG.ResponseError(http.StatusBadRequest, errors.New("not found node"), nil)
		return result, nil
	}

	// check exist node audit
	existNodeAudit, err := h.services.NodeAudit.FindNodeAudit(domain.RequestParams{
		Filter: bson.D{
			{"node_id", input.NodeID},
			{"user_id", userIDPrimitive},
		},
		Options: domain.Options{
			Limit: 1,
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	if len(existNodeAudit.Data) > 0 {
		// //appG.ResponseError(http.StatusBadRequest, errors.New("existSameNode"), nil)
		// update node audit.
		id := &existNodeAudit.Data[0].ID
		result, err = h.services.NodeAudit.UpdateNodeAudit(id.Hex(), userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}

		return result, nil
	} else {
		// create node audit.
		result, err = h.services.NodeAudit.CreateNodeAudit(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return result, err
		}
	}

	return result, nil
}
