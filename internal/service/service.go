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
	DeleteAddress(id string) (model.Address, error)
}

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	GetAuth(id string) (domain.Auth, error)
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

	FindForKml(params domain.RequestParams) (domain.Response[domain.Kml], error)
	CreateFile(params domain.RequestParams) (domain.Response[domain.NodeFileItem], error)
}

type NodeAudit interface {
	CreateNodeAudit(userID string, nodeAudit *model.NodeAuditInput) (*model.NodeAudit, error)
	FindNodeAudit(params domain.RequestParams) (domain.Response[model.NodeAudit], error)
	UpdateNodeAudit(id string, userID string, data *model.NodeAuditInput) (*model.NodeAudit, error)
	DeleteNodeAudit(id string) (model.NodeAudit, error)
}

type NodeVote interface {
	CreateNodeVote(userID string, data *model.NodeVote) (*model.NodeVote, error)
	FindNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error)
	UpdateNodeVote(id string, userID string, data *model.NodeVoteInput) (*model.NodeVote, error)
	DeleteNodeVote(id string) (model.NodeVote, error)
}

type Nodedata interface {
	FindNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error)
	GetAllNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error)
	CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error)
	UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error)
	DeleteNodedata(id string) (model.Nodedata, error)

	AddAudit(userID string, data *model.NodedataAuditInput) (*model.Nodedata, error)
}

type NodedataVote interface {
	FindNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error)
	GetAllNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error)
	CreateNodedataVote(userID string, data *model.NodedataVoteInput) (*model.NodedataVote, error)
	UpdateNodedataVote(id string, userID string, data *model.NodedataVote) (*model.NodedataVote, error)
	DeleteNodedataVote(id string) (model.NodedataVote, error)
}

type Review interface {
	CreateReview(userID string, review *model.Review) (*model.Review, error)
	FindReview(params domain.RequestParams) (domain.Response[model.Review], error)
	UpdateReview(id string, userID string, data *model.ReviewInput) (*model.Review, error)
	DeleteReview(id string) (*model.Review, error)

	GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error)
}
type User interface {
	GetUser(id string) (model.User, error)
	FindUser(params domain.RequestParams) (domain.Response[model.User], error)
	CreateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(id string) (model.User, error)
	UpdateUser(id string, user *model.User) (model.User, error)
	Iam(userID string) (model.User, error)

	SetStat(id string, statData model.UserStat) (model.User, error)
}

type Image interface {
	CreateImage(userID string, data *model.ImageInput) (model.Image, error)
	GetImage(id string) (model.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[model.Image], error)
	DeleteImage(id string) (model.Image, error)
}
type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Currency interface {
	CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error)
	GetCurrency(id string) (domain.Currency, error)
	FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error)
	UpdateCurrency(id string, data interface{}) (domain.Currency, error)
	DeleteCurrency(id string) (domain.Currency, error)
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
type AmenityGroup interface {
	FindAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error)
	GetAllAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error)
	CreateAmenityGroup(userID string, AmenityGroup *model.AmenityGroup) (*model.AmenityGroup, error)
	UpdateAmenityGroup(id string, userID string, data *model.AmenityGroup) (*model.AmenityGroup, error)
	DeleteAmenityGroup(id string) (model.AmenityGroup, error)
}

type Services struct {
	Action
	Address
	Amenity
	AmenityGroup
	Authorization
	Apps
	Country
	Currency
	Image
	Review
	User
	Track
	Node
	NodeAudit
	NodeVote
	Nodedata
	NodedataVote
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
	ImageService           config.IImageConfig
}

func NewServices(cfgService *ConfigServices) *Services {
	Authorization := NewAuthService(
		cfgService.Repositories.Authorization,
		cfgService.Hasher,
		cfgService.TokenManager,
		cfgService.RefreshTokenTTL,
		cfgService.AccessTokenTTL,
		cfgService.OtpGenerator,
		cfgService.VerificationCodeLength,
	)
	Action := NewActionService(cfgService.Repositories.Action, cfgService.I18n)
	Address := NewAddressService(cfgService.Repositories.Address, cfgService.I18n)
	Amenity := NewAmenityService(cfgService.Repositories.Amenity, cfgService.I18n)
	AmenityGroup := NewAmenityGroupService(cfgService.Repositories.AmenityGroup, cfgService.I18n)
	User := NewUserService(cfgService.Repositories.User)
	// Review := NewReviewService(cfgService.Repositories.Review)
	Apps := NewAppsService(cfgService.Repositories, cfgService.I18n)
	Country := NewCountryService(cfgService.Repositories, cfgService.I18n)
	Currency := NewCurrencyService(cfgService.Repositories, cfgService.I18n)
	Image := NewImageService(cfgService.Repositories.Image, cfgService.ImageService)
	Track := NewTrackService(cfgService.Repositories.Track)
	Node := NewNodeService(cfgService.Repositories.Node, User)
	NodeAudit := NewNodeAuditService(cfgService.Repositories.NodeAudit)
	NodeVote := NewNodeVoteService(cfgService.Repositories.NodeVote)
	// Nodedata:=     NewNodedataService(cfgService.Repositories.Nodedata, cfgService.Repositories.User, cfgService.Repositories.NodedataVote, *Services)
	NodedataVote := NewNodedataVoteService(cfgService.Repositories.NodedataVote, cfgService.Repositories.User, cfgService.Repositories.Nodedata)
	Tag := NewTagService(cfgService.Repositories.Tag)
	Tagopt := NewTagoptService(cfgService.Repositories.Tagopt)
	Ticket := NewTicketService(cfgService.Repositories.Ticket)
	Like := NewLikeService(cfgService.Repositories.Like)

	return &Services{
		Authorization: Authorization,
		Action:        Action,
		Address:       Address,
		Amenity:       Amenity,
		AmenityGroup:  AmenityGroup,
		User:          User,
		Apps:          Apps,
		Like:          Like,
		Country:       Country,
		Image:         Image,
		Track:         Track,
		Node:          Node,
		NodeAudit:     NodeAudit,
		NodeVote:      NodeVote,
		Nodedata:      NewNodedataService(cfgService.Repositories.Nodedata, User, NodedataVote),
		Tag:           Tag,
		Tagopt:        Tagopt,
		Ticket:        Ticket,
		Review:        NewReviewService(cfgService.Repositories.Review, User),
		Currency:      Currency,
		NodedataVote:  NodedataVote,
	}
}
