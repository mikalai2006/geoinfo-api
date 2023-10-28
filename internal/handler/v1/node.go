package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) registerNode(router *gin.RouterGroup) {
	node := router.Group("/node")
	node.GET("/", h.FindNode)
	node.POST("/", middleware.SetUserIdentity, h.CreateNode)
	node.POST("/list/", middleware.SetUserIdentity, h.CreateListNode)
	node.PATCH("/:id", middleware.SetUserIdentity, h.UpdateNode)
	node.DELETE("/:id", middleware.SetUserIdentity, h.DeleteNode)
}

func (h *HandlerV1) CreateNode(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *model.Node
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	node, err := h.services.Node.CreateNode(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateListNode(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.Node
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.Node
	for i := range input {
		existOsmID, err := h.services.Node.FindNode(domain.RequestParams{
			Options: domain.Options{Limit: 1},
			Filter:  bson.D{{"osm_id", input[i].OsmID}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		existLatLon := false
		if len(existOsmID.Data) > 0 {
			existLatLon = input[i].Lat == existOsmID.Data[0].Lat && input[i].Lon == existOsmID.Data[0].Lon
			progress := 0
			if existLatLon {
				progress = 100
			}

			_, err := h.services.Ticket.CreateTicket(userID, &model.Ticket{
				Title:       "Double osm object",
				Description: fmt.Sprintf("[osmId]%s[/osmId]: [coords]%v,%v[/coords], [existCoords]%v,%v[/existCoords]", input[i].OsmID, input[i].Lat, input[i].Lon, existOsmID.Data[0].Lat, existOsmID.Data[0].Lon),
				Status:      !existLatLon,
				Progress:    progress,
			})
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			// fmt.Println("Double node:::", input[i].OsmID, input[i].Lat, input[i].Lon)
		}
		if !existLatLon {
			Node, err := h.services.Node.CreateNode(userID, input[i])
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
			result = append(result, Node)
		}
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Node Get all Nodes
// @Security ApiKeyAuth
// @Tags Node
// @Description get all Nodes
// @ModuleID Node
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Node
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Node [get].
func (h *HandlerV1) GetAllNode(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Node{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	nodes, err := h.services.Node.GetAllNode(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// @Summary Find Nodes by params
// @Security ApiKeyAuth
// @Tags Node
// @Description Input params for search Nodes
// @ModuleID Node
// @Accept  json
// @Produce  json
// @Param input query NodeInput true "params for search Node"
// @Success 200 {object} []domain.Node
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Node [get].
func (h *HandlerV1) FindNode(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.NodeInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Nodes, err := h.services.Node.FindNode(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Nodes)
}

func (h *HandlerV1) GetNodeByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateNode(c *gin.Context) {
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
	var input *model.Node
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Node.UpdateNode(id, userID, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteNode(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Node.DeleteNode(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
