package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"github.com/mikalai2006/geoinfo-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerV1 struct {
	db           *mongo.Database
	repositories *repository.Repositories
	services     *service.Services
	oauth        config.OauthConfig
	i18n         config.I18nConfig
	imageConfig  config.IImageConfig
}

func NewHandler(services *service.Services, repositories *repository.Repositories, db *mongo.Database, oauth *config.OauthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig) *HandlerV1 {
	return &HandlerV1{
		repositories: repositories,
		db:           db,
		services:     services,
		oauth:        *oauth,
		i18n:         *i18n,
		imageConfig:  *imageConfig,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{

		h.registerAuth(v1)
		oauth := v1.Group("/oauth")
		h.registerVkOAuth(oauth)
		h.registerGoogleOAuth(oauth)

		h.registerReview(v1)
		h.RegisterUser(v1)
		h.RegisterApp(v1)
		h.RegisterImage(v1)
		h.registerAddress(v1)
		h.registerTrack(v1)
		h.registerNode(v1)
		h.registerGql(v1)
		h.registerTag(v1)
		h.registerTagopt(v1)
		h.registerTicket(v1)
		h.registerLike(v1)
		h.registerAmenity(v1)
		h.registerNodedata(v1)
		h.registerAction(v1)

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "v1",
			})
		})
	}
}
