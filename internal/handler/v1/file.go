package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) registerFile(router *gin.RouterGroup) {
	file := router.Group("/file")
	file.GET("/:id", h.GetFile)
	file.GET("/create", h.CreateFile)
}

// @Summary File Get
// @Security ApiKeyAuth
// @Tags File
// @Description get all Files
// @ModuleID File
// @Accept  json
// @Produce  json
// @Success 200 {object} interface{}
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/file [get].
func (h *HandlerV1) GetFile(c *gin.Context) {
	appG := app.Gin{C: c}

	code := c.Param("id")

	nameFile := code
	pathFile := "public/files"

	f, err := os.ReadFile(fmt.Sprintf("%s/%s.json", pathFile, nameFile))
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.Data(http.StatusOK, "application/json", f)
}

func (h *HandlerV1) CreateFile(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)

	// // implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if utils.Contains(roles, "admin") {
	// 	appG.ResponseError(http.StatusUnauthorized, errors.New("admin zone"), nil)
	// 	return
	// }

	counties, err := h.services.Country.FindCountry(domain.RequestParams{Filter: bson.D{}})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	pathFile := "public/files"

	stats := []domain.StatInfo{}
	for i := range counties.Data {
		result := domain.Response[domain.NodeFileItem]{}
		nameFile := counties.Data[i].Code // c.Query("namefile") //
		// f, err := os.Open(fmt.Sprintf("%s/%s.json", pathFile, nameFile))
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }
		// defer f.Close()

		w, err := os.Create(fmt.Sprintf("%s/%s.json", pathFile, nameFile))
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		defer w.Close()
		result, err = h.services.Node.CreateFile(domain.RequestParams{Filter: bson.D{
			{"ccode", bson.D{
				{"$in", bson.A{counties.Data[i].Code}},
			}},
		}})
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		// resp := []map[string]interface{}
		// for j := range result.Data {
		// 	startData := result.Data[j](map[string]interface{})
		// 	delete(startData, "data");
		// 	for j := range result.Data {
		// 		// resp["value"] = result.Data[j].Data
		// 	}
		// }

		b, err := json.Marshal(result.Data)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		w.Write(b)

		stat, err := w.Stat()
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		statData := domain.StatInfo{
			Name:          counties.Data[i].Code,
			CCode:         counties.Data[i].Code,
			Path:          fmt.Sprintf("%s/%s.json", pathFile, nameFile),
			Size:          float64(stat.Size()),
			Count:         float64(result.Total),
			LastUpdatedAt: time.Now(),
		}

		var stateInput domain.CountryInput
		dd := map[string]interface{}{
			"stat": statData,
		}
		stateInputData, err := utils.BindAndValidFromMarshal[domain.CountryInput](dd, stateInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// fmt.Println("stateInputData: ", stateInputData)
		document, err := h.services.Country.UpdateCountry(counties.Data[i].ID.Hex(), stateInputData)
		if err != nil {
			appG.ResponseError(http.StatusInternalServerError, err, nil)
			return
		}
		stats = append(stats, document.Stat)
	}

	c.JSON(http.StatusOK, stats)
}
