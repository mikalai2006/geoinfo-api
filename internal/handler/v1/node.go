package v1

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	node.GET("", h.FindNode)
	node.POST("", h.CreateNode)
	node.POST("/list/", h.CreateListNode)
	node.PATCH("/:id", h.UpdateNode)
	node.DELETE("/:id", h.DeleteNode)

	node.GET("/zip", h.CreateZip)
	node.GET("/parsekml", h.ParseKML)
	// node.GET("/duplicate", h.FindRepeat)
}

func (h *HandlerV1) CreateNode(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *model.Node
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	node, err := h.CreateOrExistNode(c, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateOrExistNode(c *gin.Context, input *model.Node) (*model.Node, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}

	// // check exist node
	// existNode, err := h.services.Node.FindNode(domain.RequestParams{
	// 	Filter: bson.D{
	// 		{"lat", bson.M{"$lt": node.Lat + 0.001, "$gt": node.Lat - 0.001}},
	// 		{"lon", bson.M{"$lt": node.Lon + 0.001, "$gt": node.Lon - 0.001}},
	// 		{"type", node.Type},
	// 	},
	// 	Options: domain.Options{
	// 		Limit: 1,
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// if len(existNode.Data) > 0 {
	// 	return &existNode.Data[0], nil
	// }
	// check exist node
	existNode, err := h.services.Node.FindNode(domain.RequestParams{
		Filter: bson.D{
			{"lat", bson.M{"$lt": input.Lat + 0.00015, "$gt": input.Lat - 0.00015}},
			{"lon", bson.M{"$lt": input.Lon + 0.00015, "$gt": input.Lon - 0.00015}},
			{"type", input.Type},
		},
		Options: domain.Options{
			Limit: 1,
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}

	// if exist node
	if len(existNode.Data) > 0 {
		//appG.ResponseError(http.StatusBadRequest, errors.New("existSameNode"), nil)
		return &existNode.Data[0], nil
	}

	// Get address.
	pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input.Lat, input.Lon))
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
	r.Header.Add("User-Agent", "a127.0.0.1")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return nil, err
	}
	var bodyResponse domain.ResponseNominatim
	if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
		appG.ResponseError(http.StatusBadRequest, e, nil)
		return nil, err
	}

	var result *model.Node

	if bodyResponse.OsmID != 0 {
		// Check address in to db
		osmID := fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID)
		adrDB, err := h.services.Address.FindAddress(domain.RequestParams{Options: domain.Options{Limit: 1},
			Filter: bson.D{{"osm_id", osmID}}})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return nil, err
		}

		address := &domain.Address{}
		if len(adrDB.Data) > 0 {
			address = &adrDB.Data[0]
		} else {
			address, err = h.services.Address.CreateAddress(userID, &domain.AddressInput{
				OsmID:    osmID,
				Address:  bodyResponse.Address,
				DAddress: bodyResponse.DisplayName,
			})
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return nil, err
			}
		}

		input.OsmID = address.OsmID

		if bodyResponse.Name == "" {
			arrStr := strings.Split(address.DAddress, ",")
			nameNode := ""
			if len(arrStr) >= 2 {
				nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
			} else {
				nameNode = arrStr[0]
			}
			input.Name = strings.TrimSpace(nameNode)
		} else {
			input.Name = bodyResponse.Name
		}

		if ccode, ok := bodyResponse.Address["country_code"]; ok {
			input.CCode = ccode.(string)
		}

		node, err := h.services.Node.CreateNode(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return nil, err
		}

		if len(input.Data) > 0 {
			for i := range input.Data {
				inputNodedata := &model.NodedataInput{
					NodeID:   node.ID.Hex(),
					Data:     input.Data[i].Data,
					TagID:    input.Data[i].TagID.Hex(),
					TagoptID: input.Data[i].TagoptID.Hex(),
				}

				Nodedata, err := h.CreateOrExistNodedata(c, inputNodedata)
				// .services.Nodedata.CreateNodedata(userID, inputNodedata)
				if err != nil {
					appG.ResponseError(http.StatusBadRequest, err, nil)
					return nil, err
				}

				node.Data = append(node.Data, *Nodedata)
			}
		}
		result = node
	}
	//  else {
	// 	fmt.Println("not found osm", bodyResponse.OsmID)
	// }

	return result, nil
}

func (h *HandlerV1) CreateListNode(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
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
		// existOsmID, err := h.services.Node.FindNode(domain.RequestParams{
		// 	Options: domain.Options{Limit: 1},
		// 	Filter:  bson.D{{"osm_id", input[i].OsmID}},
		// })
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }

		// existLatLon := false
		// if len(existOsmID.Data) > 0 {
		// 	existLatLon = input[i].Lat == existOsmID.Data[0].Lat && input[i].Lon == existOsmID.Data[0].Lon
		// 	progress := 0
		// 	if existLatLon {
		// 		progress = 100
		// 	}

		// 	_, err := h.services.Ticket.CreateTicket(userID, &model.Ticket{
		// 		Title:       "Double osm object",
		// 		Description: fmt.Sprintf("[osmId]%s[/osmId]: [coords]%v,%v[/coords], [existCoords]%v,%v[/existCoords]", input[i].OsmID, input[i].Lat, input[i].Lon, existOsmID.Data[0].Lat, existOsmID.Data[0].Lon),
		// 		Status:      !existLatLon,
		// 		Progress:    progress,
		// 	})
		// 	if err != nil {
		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
		// 		return
		// 	}
		// 	// fmt.Println("Double node:::", input[i].OsmID, input[i].Lat, input[i].Lon)
		// }
		// if !existLatLon {

		// // Get address.
		// pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input[i].Lat, input[i].Lon))
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
		// r.Header.Add("User-Agent", "a127.0.0.1")

		// resp, err := http.DefaultClient.Do(r)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// defer resp.Body.Close()

		// bytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// var bodyResponse domain.ResponseNominatim
		// if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
		// 	appG.ResponseError(http.StatusBadRequest, e, nil)
		// 	return
		// }

		// if bodyResponse.OsmID != 0 {
		// 	address, err := h.services.Address.CreateAddress(userID, &domain.AddressInput{
		// 		OsmID:    fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID),
		// 		Address:  bodyResponse.Address,
		// 		DAddress: bodyResponse.DisplayName,
		// 	})
		// 	if err != nil {
		// 		appG.ResponseError(http.StatusBadRequest, err, nil)
		// 		return
		// 	}

		// 	input[i].OsmID = address.OsmID

		// 	if bodyResponse.Name == "" {
		// 		arrStr := strings.Split(address.DAddress, ",")
		// 		nameNode := ""
		// 		if len(arrStr) >= 2 {
		// 			nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
		// 		} else {
		// 			nameNode = arrStr[0]
		// 		}
		// 		input[i].Name = strings.TrimSpace(nameNode)
		// 	} else {
		// 		input[i].Name = bodyResponse.Name
		// 	}

		// 	if ccode, ok := bodyResponse.Address["country_code"]; ok {
		// 		input[i].CCode = ccode.(string)
		// 	}
		// }

		// Node, err := h.services.Node.CreateNode(userID, input[i])
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }

		// if len(input[i].Data) > 0 {
		// 	for j := range input[i].Data {
		// 		inputNodedata := &model.NodedataInput{
		// 			NodeID:   Node.ID.Hex(),
		// 			Data:     input[i].Data[j].Data,
		// 			TagID:    input[i].Data[j].TagID.Hex(),
		// 			TagoptID: input[i].Data[j].TagoptID.Hex(),
		// 		}

		// 		Nodedata, err := h.services.Nodedata.CreateNodedata(userID, inputNodedata)
		// 		if err != nil {
		// 			appG.ResponseError(http.StatusBadRequest, err, nil)
		// 			return
		// 		}
		// 		Node.Data = append(Node.Data, *Nodedata)
		// 	}
		// }

		Node, err := h.CreateOrExistNode(c, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		result = append(result, Node)
		// }
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

	authData, err := middleware.GetAuthFromCtx(c)
	fmt.Println("auth ", authData.Roles)

	params, err := utils.GetParamsFromRequest(c, model.NodeInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	fmt.Println(params)
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

	node, err := h.services.Node.DeleteNode(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// find all images for remove.
	images, err := h.services.Image.FindImage(domain.RequestParams{
		Filter: bson.D{
			{"service", "node"},
			{"service_id", node.ID.Hex()},
		},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i, _ := range images.Data {
		_, err := h.services.Image.DeleteImage(images.Data[i].ID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("Remove image", image.ID)
	}

	// find all nodedata for remove.
	nodedata, err := h.services.Nodedata.FindNodedata(domain.RequestParams{
		Filter: bson.D{
			{"node_id", node.ID},
		},
		Options: domain.Options{Limit: 1000},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i, _ := range nodedata.Data {
		_, err := h.services.Nodedata.DeleteNodedata(nodedata.Data[i].ID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("Remove nodedata: ", nodedata.Data[i].ID.Hex())
	}

	// find all reviews for remove.
	reviews, err := h.services.Review.FindReview(domain.RequestParams{
		Filter: bson.D{
			{"node_id", node.ID},
		},
		Options: domain.Options{Limit: 1000},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i, _ := range reviews.Data {
		_, err := h.services.Review.DeleteReview(reviews.Data[i].ID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("Remove review: ", reviews.Data[i].ID.Hex())
	}

	// find all audits for remove.
	nodeaudits, err := h.services.NodeAudit.FindNodeAudit(domain.RequestParams{
		Filter: bson.D{
			{"node_id", node.ID},
		},
		Options: domain.Options{Limit: 1000},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i, _ := range nodeaudits.Data {
		_, err := h.services.NodeAudit.DeleteNodeAudit(nodeaudits.Data[i].ID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("Remove nodeaudits: ", nodeaudits.Data[i].ID.Hex())
	}

	// Remove address.
	nodeAlsoOsmID, err := h.services.Node.FindNode(domain.RequestParams{
		Filter: bson.D{{"osm_id", node.OsmID}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(nodeAlsoOsmID.Data) == 0 {
		addr, err := h.services.Address.FindAddress(domain.RequestParams{
			Filter: bson.D{{"osm_id", node.OsmID}},
		})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		if len(addr.Data) > 0 {
			_, err = h.services.Address.DeleteAddress(addr.Data[0].ID.Hex())
			if err != nil {
				appG.ResponseError(http.StatusBadRequest, err, nil)
				return
			}
		}
	}

	c.JSON(http.StatusOK, node)
}

func (h *HandlerV1) CreateZip(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)

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

	// var input *model.Node
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }

	params, err := utils.GetParamsFromRequest(c, model.NodeInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Nodes, err := h.services.Node.FindForKml(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	country := c.Query("ccode")
	if country == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("not found country"), nil)
		return
	}

	// get current language.
	var currentLocale utils.LocaleItem = utils.Locales[0]
	for _, v := range utils.Locales {
		if v.Code == params.Lang {
			currentLocale = v
		}
	}

	pathDefault := "public\\default\\"
	archive, err := os.Create(fmt.Sprintf("public/kml/%s(%s).zip", country, currentLocale.Code))
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	walker := func(p string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", p)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		relPath := strings.TrimPrefix(p, filepath.Dir(pathDefault))
		f, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(pathDefault, walker)
	if err != nil {
		panic(err)
	}

	groups := map[string][]domain.Kml{}
	for _, v := range Nodes.Data {
		if err, ok := groups[v.Type]; ok {
			if err != nil && len(groups[v.Type]) == 0 {
				groups[v.Type] = []domain.Kml{}
			}
		}
		groups[v.Type] = append(groups[v.Type], v)
	}

	kml := domain.KMLSchema{
		XMLNS: "http://www.opengis.net/kml/2.2",
		GX:    "http://www.google.com/kml/ext/2.2",
		KML:   "http://www.opengis.net/kml/2.2",
		Atom:  "http://www.w3.org/2005/Atom",
	}

	// Create extends schema.
	schemaFields := []domain.KMLFileSchemaSimpleField{}
	schemaFields = append(schemaFields, domain.KMLFileSchemaSimpleField{
		Type:        "string",
		Name:        "coo",
		DisplayName: currentLocale.Coordinates,
	})
	// schemaFields = append(schemaFields, domain.KMLFileSchemaSimpleField{
	// 	Type:        "string",
	// 	Name:        fmt.Sprintf("%s", "description"),
	// 	DisplayName: "Description",
	// })
	schemaFields = append(schemaFields, domain.KMLFileSchemaSimpleField{
		Type:        "string",
		Name:        "cre",
		DisplayName: currentLocale.Created,
	})
	schemaFields = append(schemaFields, domain.KMLFileSchemaSimpleField{
		Type:        "string",
		Name:        "cou",
		DisplayName: currentLocale.Country,
	})
	schemaFields = append(schemaFields, domain.KMLFileSchemaSimpleField{
		Type:        "string",
		Name:        "aut",
		DisplayName: currentLocale.Authors,
	})
	kml.Document.KMLFileSchema = domain.KMLFileSchema{
		Name:        "__managed_schema",
		ID:          "__managed_schema",
		SimpleField: schemaFields,
	}

	// Create placemarkers.
	KMLGX := []domain.KMLGX{}
	stylesMap := []domain.StyleMap{}
	for _, v := range groups {
		// Create style group folder.
		bgColor := "#ffffff"
		if _, ok := v[0].Amenity.Props["bgColor"]; ok {
			bgColor = v[0].Amenity.Props["bgColor"].(string)
		}
		color := strings.Replace(bgColor, "#", "", -1)

		var indexIconKml int64 = 2000
		if _, ok := v[0].Amenity.Props["iconkml"]; ok {
			indexIconKml = int64(v[0].Amenity.Props["iconkml"].(float64))
		}
		href := fmt.Sprintf("https://earth.google.com/earth/rpc/cc/icon?color=%s&id=%d&scale=4", color, indexIconKml)

		highlightStyle := domain.KMLGX{
			ID:       fmt.Sprintf("highlight_%s", v[0].Type),
			StyleURL: "https://earth.google.com/balloon_components/base/1.0.26.0/card_template.kml#main",
			Style: domain.KMLGXStyle{
				IconStyle: domain.KMLGXIconStyle{
					Icon: domain.KMLGXIcon{
						Href: href,
					},
					HotSpot: domain.HotSpot{
						X:      "64",
						Y:      "128",
						XUnits: "pixels",
						YUnits: "insetPixels",
					},
					Scale: "1.5",
				},
				LabelStyle: domain.KMLLabelStyle{
					Scale: "1.3",
				},
				LineStyle: domain.KMLLineStyle{
					Color: color,
					Width: "4",
				},
				PolyStyle: domain.KMLPolyStyle{
					Color: "ffffff",
				},
			},
		}
		normalStyle := domain.KMLGX{
			ID:       fmt.Sprintf("normal_%s", v[0].Type),
			StyleURL: "https://earth.google.com/balloon_components/base/1.0.26.0/card_template.kml#main",
			Style: domain.KMLGXStyle{
				IconStyle: domain.KMLGXIconStyle{
					Icon: domain.KMLGXIcon{
						Href: href,
					},
					HotSpot: domain.HotSpot{
						X:      "64",
						Y:      "128",
						XUnits: "pixels",
						YUnits: "insetPixels",
					},
					Scale: "1",
				},
				LabelStyle: domain.KMLLabelStyle{
					Scale: "0",
				},
				LineStyle: domain.KMLLineStyle{
					Color: color,
					Width: "4",
				},
				PolyStyle: domain.KMLPolyStyle{
					Color: "ffffff",
				},
				// BalloonStyle: domain.KMLBalloonStyle{
				// 	Text: "<![CDATA[<B>Name: $[name]</B><br>$[description]<P>Country: $[country]<br/>Authors $[authors]]]>",
				// },
			},
		}
		KMLGX = append(KMLGX, normalStyle)
		KMLGX = append(KMLGX, highlightStyle)

		styleMap := domain.StyleMap{
			ID: fmt.Sprintf("style_%s", v[0].Type),
			Pair: []domain.KMLPair{
				{
					KEY:      "normal",
					StyleURL: fmt.Sprintf("#normal_%s", v[0].Type),
				},
				{
					KEY:      "highlight",
					StyleURL: fmt.Sprintf("#highlight_%s", v[0].Type),
				},
			},
		}
		stylesMap = append(stylesMap, styleMap)
	}

	fmt.Println(len(KMLGX), len(stylesMap))
	kml.Document.KMLGX = KMLGX
	kml.Document.StyleMap = stylesMap
	kml.Document.ID = fmt.Sprintf("doc_%s_%s", country, currentLocale.Code)

	groupsFolder := []domain.KMLFolder{}
	for _, v := range groups {
		// fmt.Println(len(v))
		group := domain.KMLFolder{
			Metadata: domain.KMLMetadata{
				Igoicon: domain.KMLIgoicon{
					Filename: fmt.Sprintf("%s.bmp", v[0].Type),
				},
			},
			ID:   v[0].Amenity.Type,
			Name: v[0].Amenity.Title,
			// StyleURL: fmt.Sprintf("#style_%s", v[0].Type),
		}
		for _, g := range v {
			descriptionArrayForName := []string{}
			descriptionArray := []string{}
			contributorsArray := map[string]bool{
				fmt.Sprintf("%s(%s)", g.User.Name, g.User.Login): true,
			}
			if g.Data != nil {
				for _, t := range g.Data {
					if t.Data.Value == "yes" {
						descriptionArray = append(descriptionArray, fmt.Sprintf("<div><span><b>%s</b>,</span></div>", t.Tag.Title))
						descriptionArrayForName = append(descriptionArrayForName, fmt.Sprintf("%s", t.Tag.Title))
					} else {
						descriptionArray = append(descriptionArray, fmt.Sprintf("<div><span><b>%s</b>:%s,</span></div>", t.Tag.Title, t.Data.Value))
						descriptionArrayForName = append(descriptionArrayForName, fmt.Sprintf("%s:%s", t.Tag.Title, t.Data.Value))
					}

					contributorsArray[fmt.Sprintf("%s(%s)", t.User.Name, t.User.Login)] = true
				}
			}

			keysContributors := make([]string, len(contributorsArray))
			i := 0
			for k := range contributorsArray {
				keysContributors[i] = k
				i++
			}

			description := fmt.Sprintf(
				"<![CDATA[<div style=\"font-size:14px\"><div>%s</div><small>%s: %s(%s: %s)</small><small>%s: %s</small></div>]]>",
				strings.Join(descriptionArray, " "),
				currentLocale.Created,
				g.CreatedAt.Format("2006-01-02"),
				currentLocale.Updated,
				g.UpdatedAt.Format("2006-01-02"),
				currentLocale.Authors,
				strings.Join(keysContributors, ", "),
			)

			group.Placemark = append(group.Placemark, domain.KMLPlacemark{
				ID:   fmt.Sprintf("%s_%s", g.Amenity.Type, g.ID.Hex()),
				Name: fmt.Sprintf("%s - %s ::: %s", g.Amenity.Title, g.Name, strings.Join(descriptionArrayForName, ", ")),
				Description: domain.KMLPlacemarkDescription{
					Description: description,
				},
				// ExtendedData: domain.ExtendedData{
				// 	Data: []domain.ExtendedDataData{
				// 		{Name: "name",
				// 			Value: fmt.Sprintf("%s - %s", g.Amenity.Title, g.Name),
				// 		},
				// 		{Name: "description",
				// 			Value: description,
				// 		},
				// 	},
				// },
				ExtendedData: domain.ExtendedData{
					SchemaData: domain.SchemaData{
						SchemaUrl: "#__managed_schema",
						SimpleData: []domain.SimpleData{
							{
								Name:  "coo",
								Value: fmt.Sprintf("%f,%f", g.Lat, g.Lon),
							},
							{
								Name:  "aut",
								Value: strings.Join(keysContributors, ", "),
							},
							{
								Name:  "cou",
								Value: country,
							},
							{
								Name:  "cre",
								Value: fmt.Sprintf("%s(%s: %s)", g.CreatedAt.Format("2006-01-02"), currentLocale.Updated, g.UpdatedAt.Format("2006-01-02")),
							},
						},
					},
				},
				Point: domain.KMLPoint{
					Coordinates: fmt.Sprintf("%f,%f,%d", g.Lon, g.Lat, 0),
				},
				StyleURL: fmt.Sprintf("#style_%s", v[0].Type),
			})
		}
		groupsFolder = append(groupsFolder, group)
	}
	kml.Document.Name = "POIs"
	kml.Document.Folder = domain.KMLParentFolder{
		ID:          country,
		Name:        fmt.Sprintf("%s-%s(%s)", "POIs", strings.ToUpper(country), currentLocale.Code),
		Folder:      groupsFolder,
		Description: "Поддержка - Viber, WhatsApp +00000000000",
	}
	// kml.Document.KMLGX = []domain.KMLGX{}
	// kml.Document.StyleMap = []domain.StyleMap{}

	// kml.Document.Name = fmt.Sprintf("%s-%s", "POIS", country)

	fileNameKml := fmt.Sprintf("pois-%s(%s).kml", country, currentLocale.Code)
	w2, err := zipWriter.Create(fileNameKml)
	if err != nil {
		panic(err)
	}

	w3 := &bytes.Buffer{}
	w3.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))

	enc := xml.NewEncoder(w3)
	enc.Indent("", "    ")
	if err := enc.Encode(kml); err != nil {
		panic(err)
	}
	_, err = io.Copy(w2, w3)
	if err != nil {
		panic(err)
	}

	// if err := doc.WriteIndent(w2, "", "  "); err != nil {
	// 	panic(err)
	// }

	// c.JSON(http.StatusOK, Nodes)
	c.JSON(http.StatusOK, len(Nodes.Data))
}

func (h *HandlerV1) ParseKML(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)

	// implementation roles for user.
	roles, err := middleware.GetRoles(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if utils.Contains(roles, "admin") {
		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
		return
	}

	// var input *model.Node
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }

	// params, err := utils.GetParamsFromRequest(c, domain.NodeInputData{}, &h.i18n)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// Nodes, err := h.services.Node.FindForKml(params)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	pathFile := "public/parsekml"
	nameFile := c.Query("namefile") // "prom"
	f, err := os.Open(fmt.Sprintf("%s/%s.kml", pathFile, nameFile))
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer f.Close()

	w, err := os.Create(fmt.Sprintf("%s/%s.json", pathFile, nameFile))
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer w.Close()

	var kml domain.KMLParseSchema
	enc := xml.NewDecoder(f)
	if err := enc.Decode(&kml); err != nil {
		panic(err)
	}

	fields := []domain.KMLParseField{}
	for _, v := range kml.Document.Folder {
		for _, p := range v.Placemark {
			fie := domain.KMLParseField{}
			for _, g := range p.ExtendedData.SchemaData.SimpleData {
				switch g.Name {
				case "str:UlU=":
					fie.Description = g.Value
				case "str:TGF0aXR1ZGU=":
					fie.Lat = g.Value
				case "str:TG9uZ2l0dWRl":
					fie.Lon = g.Value
				case "str:UG9pbnQgYXV0aG9y":
					fie.Author = g.Value
				case "str:Q2F0ZWdvcnk=":
					fie.Name = g.Value
				}
			}
			fields = append(fields, fie)
		}
	}
	b, err := json.Marshal(fields)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	w.Write(b)

	c.JSON(http.StatusOK, fields)
}

// func (h *HandlerV1) FindRepeat(c *gin.Context) {
// 	appG := app.Gin{C: c}
// 	// userID, err := middleware.GetUID(c)

// 	// implementation roles for user.
// 	roles, err := middleware.GetRoles(c)
// 	if err != nil {
// 		appG.ResponseError(http.StatusUnauthorized, err, nil)
// 		return
// 	}
// 	if !utils.Contains(roles, "admin") {
// 		appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
// 		return
// 	}

// 	pathFile := "public"
// 	nameFile := "doublenode" // "prom"
// 	w, err := os.Create(fmt.Sprintf("%s/%s.json", pathFile, nameFile))
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}
// 	defer w.Close()

// 	allNodes, err := h.services.Node.FindNode(domain.RequestParams{
// 		Filter:  bson.D{},
// 		Options: domain.Options{Limit: 10000000},
// 	})

// 	fields := []interface{}{}
// 	for i := 0; i < 1000; i++ {
// 		doublenodes, err := h.services.Node.FindNode(domain.RequestParams{
// 			Filter: bson.D{
// 				{"lat", bson.M{"$lt": allNodes.Data[i].Lat + 0.0001, "$gt": allNodes.Data[i].Lat - 0.0001}},
// 				{"lon", bson.M{"$lt": allNodes.Data[i].Lon + 0.0001, "$gt": allNodes.Data[i].Lon - 0.0001}},
// 				{"type", allNodes.Data[i].Type},
// 			},
// 			Options: domain.Options{Limit: 20},
// 		})
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return
// 		}

// 		if len(doublenodes.Data) > 1 {
// 			fields = append(fields, doublenodes.Data)
// 		}
// 	}

// 	b, err := json.Marshal(fields)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	w.Write(b)

// 	c.JSON(http.StatusOK, fields)
// }
