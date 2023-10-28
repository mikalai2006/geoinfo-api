package service

import (
	"time"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"github.com/mikalai2006/geoinfo-api/pkg/auths"
	"github.com/mikalai2006/geoinfo-api/pkg/hasher"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, data *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
}
type Address interface {
	FindAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	GetAllAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	CreateAddress(userID string, address *domain.AddressInput) (*domain.Address, error)
}

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	SignIn(input *domain.SignInInput) (domain.ResponseTokens, error)
	ExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	CreateSession(auth *domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userID string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
	RemoveRefreshTokens(refreshToken string) (string, error)
}

type Track interface {
	FindTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	GetAllTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	CreateTrack(userID string, track *domain.Track) (*domain.Track, error)
}

type Node interface {
	FindNode(params domain.RequestParams) (domain.Response[model.Node], error)
	GetAllNode(params domain.RequestParams) (domain.Response[model.Node], error)
	CreateNode(userID string, node *model.Node) (*model.Node, error)
	UpdateNode(id string, userID string, data *model.Node) (*model.Node, error)
	DeleteNode(id string) (model.Node, error)
}
type Nodedata interface {
	FindNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error)
	GetAllNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error)
	CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error)
	UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error)
	DeleteNodedata(id string) (model.Nodedata, error)
}

type Review interface {
	FindReview(params domain.RequestParams) (domain.Response[domain.Review], error)

	GetAllReview(params domain.RequestParams) (domain.Response[domain.Review], error)
	CreateReview(userID string, review *domain.Review) (*domain.Review, error)
}
type User interface {
	GetUser(id string) (model.User, error)
	FindUser(params domain.RequestParams) (domain.Response[model.User], error)
	CreateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(id string) (model.User, error)
	UpdateUser(id string, user *model.User) (model.User, error)
	Iam(userID string) (model.User, error)
}

type Image interface {
	CreateImage(userID string, data *model.ImageInput) (model.Image, error)
	GetImage(id string) (model.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[model.Image], error)
	DeleteImage(id string) (model.Image, error)
}

type Apps interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Tag interface {
	FindTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	CreateTag(userID string, tag *model.Tag) (*model.Tag, error)
	UpdateTag(id string, userID string, data *model.Tag) (*model.Tag, error)
	DeleteTag(id string) (model.Tag, error)
}
type Tagopt interface {
	FindTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error)
	GetAllTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error)
	CreateTagopt(userID string, tagopt *model.TagoptInput) (*model.Tagopt, error)
	UpdateTagopt(id string, userID string, data *model.TagoptInput) (*model.Tagopt, error)
	DeleteTagopt(id string) (model.Tagopt, error)
}
type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	DeleteTicket(id string) (model.Ticket, error)
}

type Like interface {
	FindLike(params domain.RequestParams) (domain.Response[model.Like], error)
	CreateLike(userID string, like *model.LikeInput) (*model.Like, error)
	UpdateLike(id string, userID string, data *model.Like) (*model.Like, error)
	DeleteLike(id string) (model.Like, error)
}
type Amenity interface {
	FindAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error)
	GetAllAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error)
	CreateAmenity(userID string, Amenity *model.Amenity) (*model.Amenity, error)
	UpdateAmenity(id string, userID string, data *model.Amenity) (*model.Amenity, error)
	DeleteAmenity(id string) (model.Amenity, error)
}

type Services struct {
	Action
	Address
	Amenity
	Authorization
	Apps
	Image
	Review
	User
	Track
	Node
	Nodedata
	Tag
	Tagopt
	Ticket

	Like
}

type ConfigServices struct {
	Repositories           *repository.Repositories
	Hasher                 hasher.PasswordHasher
	TokenManager           auths.TokenManager
	OtpGenerator           utils.Generator
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeLength int
	I18n                   config.I18nConfig
}

func NewServices(cfgService *ConfigServices) *Services {
	return &Services{
		Authorization: NewAuthService(
			cfgService.Repositories.Authorization,
			cfgService.Hasher,
			cfgService.TokenManager,
			cfgService.RefreshTokenTTL,
			cfgService.AccessTokenTTL,
			cfgService.OtpGenerator,
			cfgService.VerificationCodeLength,
		),
		Action:   NewActionService(cfgService.Repositories.Action, cfgService.I18n),
		Address:  NewAddressService(cfgService.Repositories.Address, cfgService.I18n),
		Amenity:  NewAmenityService(cfgService.Repositories.Amenity, cfgService.I18n),
		Review:   NewReviewService(cfgService.Repositories.Review),
		Apps:     NewAppsService(cfgService.Repositories, cfgService.I18n),
		Image:    NewImageService(cfgService.Repositories.Image),
		User:     NewUserService(cfgService.Repositories.User),
		Track:    NewTrackService(cfgService.Repositories.Track),
		Node:     NewNodeService(cfgService.Repositories.Node),
		Nodedata: NewNodedataService(cfgService.Repositories.Nodedata),
		Tag:      NewTagService(cfgService.Repositories.Tag),
		Tagopt:   NewTagoptService(cfgService.Repositories.Tagopt),
		Ticket:   NewTicketService(cfgService.Repositories.Ticket),

		Like: NewLikeService(cfgService.Repositories.Like),
	}
}
