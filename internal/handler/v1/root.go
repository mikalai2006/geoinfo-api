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
	auth         config.AuthConfig
	i18n         config.I18nConfig
	imageConfig  config.IImageConfig
}

func NewHandler(services *service.Services, repositories *repository.Repositories, db *mongo.Database, oauth *config.OauthConfig, auth *config.AuthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig) *HandlerV1 {
	return &HandlerV1{
		repositories: repositories,
		db:           db,
		services:     services,
		oauth:        *oauth,
		auth:         *auth,
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

		h.registerAmenity(v1)
		h.registerAmenityGroup(v1)
		h.RegisterApp(v1)
		h.RegisterCurrency(v1)
		h.RegisterCountry(v1)
		h.registerTag(v1)
		h.registerTagopt(v1)

		authenticated := v1.Group("", h.SetUserFromRequest)
		{
			h.registerFile(authenticated)
			h.registerAction(authenticated)
			h.registerAddress(authenticated)
			h.RegisterImage(authenticated)
			h.registerGql(authenticated)
			h.registerLike(authenticated)
			h.registerNode(authenticated)
			h.registerNodeVote(authenticated)
			h.registerNodeAudit(authenticated)
			h.registerNodedata(authenticated)
			h.registerNodedataVote(authenticated)
			h.registerReview(authenticated)
			h.registerTicket(authenticated)
			h.registerTrack(authenticated)
			h.RegisterUser(authenticated)
		}

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "v1",
			})
		})
	}
}
