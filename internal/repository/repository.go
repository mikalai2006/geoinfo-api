package repository

import (
	"reflect"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, tag *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
	GqlGetActions(params domain.RequestParams) ([]*model.Action, error)
}

type Address interface {
	FindAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	GetAllAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	CreateAddress(userID string, address *domain.AddressInput) (*domain.Address, error)
	DeleteAddress(id string) (model.Address, error)
	GqlGetAdresses(params domain.RequestParams) ([]*model.Address, error)
}

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	GetAuth(id string) (domain.Auth, error)
	CheckExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	GetByCredentials(auth *domain.SignInInput) (domain.Auth, error)
	SetSession(authID primitive.ObjectID, session domain.Session) error
	VerificationCode(userID string, code string) error
	RefreshToken(refreshToken string) (domain.Auth, error)
	RemoveRefreshToken(refreshToken string) (string, error)
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
	GetNodedata(id string) (*model.Nodedata, error)
	CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error)
	UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error)
	DeleteNodedata(id string) (model.Nodedata, error)
	GqlGetNodedatas(params domain.RequestParams) ([]*model.Nodedata, error)

	AddAudit(userID string, data *model.NodedataAuditInput) (*model.Nodedata, error)
	FindAudits(params domain.RequestParams) (domain.Response[model.NodedataAudit], error)
}

type NodedataVote interface {
	FindNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error)
	GetAllNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error)
	GetNodedataVote(id string) (*model.NodedataVote, error)
	CreateNodedataVote(userID string, data *model.NodedataVoteInput) (*model.NodedataVote, error)
	UpdateNodedataVote(id string, userID string, data *model.NodedataVote) (*model.NodedataVote, error)
	DeleteNodedataVote(id string) (model.NodedataVote, error)
	GqlGetNodedataVote(params domain.RequestParams) ([]*model.NodedataVote, error)
}

type Review interface {
	FindReview(params domain.RequestParams) (domain.Response[model.Review], error)
	GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error)
	CreateReview(userID string, review *model.Review) (*model.Review, error)
	UpdateReview(id string, userID string, data *model.ReviewInput) (*model.Review, error)
	DeleteReview(id string) (*model.Review, error)

	GqlGetReviews(params domain.RequestParams) ([]*model.Review, error)
	GqlGetCountReviews(params domain.RequestParams) (*model.ReviewInfo, error)
}

type User interface {
	GetUser(id string) (model.User, error)
	FindUser(params domain.RequestParams) (domain.Response[model.User], error)
	CreateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(id string) (model.User, error)
	UpdateUser(id string, user *model.User) (model.User, error)
	Iam(userID string) (model.User, error)

	SetStat(userID string, statData model.UserStat) (model.User, error)

	GqlGetUsers(params domain.RequestParams) ([]*model.User, error)
}

type Image interface {
	CreateImage(userID string, data *model.ImageInput) (model.Image, error)
	GetImage(id string) (model.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[model.Image], error)
	DeleteImage(id string) (model.Image, error)

	GqlGetImages(params domain.RequestParams) ([]*model.Image, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Currency interface {
	CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error)
	GetCurrency(id string) (domain.Currency, error)
	FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error)
	UpdateCurrency(id string, data interface{}) (domain.Currency, error)
	DeleteCurrency(id string) (domain.Currency, error)
}

type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Tag interface {
	FindTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	CreateTag(userID string, tag *model.Tag) (*model.Tag, error)
	UpdateTag(id string, userID string, data *model.Tag) (*model.Tag, error)
	DeleteTag(id string) (model.Tag, error)
	GqlGetTags(params domain.RequestParams) ([]*model.Tag, error)
}

type Tagopt interface {
	FindTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error)
	GetAllTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error)
	CreateTagopt(userID string, tag *model.TagoptInput) (*model.Tagopt, error)
	UpdateTagopt(id string, userID string, data *model.TagoptInput) (*model.Tagopt, error)
	DeleteTagopt(id string) (model.Tagopt, error)
	GqlGetTagopts(params domain.RequestParams) ([]*model.Tagopt, error)
}

type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error)
	DeleteTicket(id string) (model.Ticket, error)
	GqlGetTickets(params domain.RequestParams) ([]*model.Ticket, error)
}
type Like interface {
	FindLike(params domain.RequestParams) (domain.Response[model.Like], error)
	CreateLike(userID string, like *model.LikeInput) (*model.Like, error)
	UpdateLike(id string, userID string, data *model.Like) (*model.Like, error)
	DeleteLike(id string) (model.Like, error)
	GqlGetIamLike(userID string, nodeID string) (*model.Like, error)
	GqlGetLikes(params domain.RequestParams) ([]*model.Like, error)
}

type Amenity interface {
	FindAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error)
	GetAllAmenity(params domain.RequestParams) (domain.Response[model.Amenity], error)
	CreateAmenity(userID string, Amenity *model.Amenity) (*model.Amenity, error)
	UpdateAmenity(id string, userID string, data *model.Amenity) (*model.Amenity, error)
	DeleteAmenity(id string) (model.Amenity, error)
	GqlGetAmenitys(params domain.RequestParams) ([]*model.Amenity, error)
}

type AmenityGroup interface {
	FindAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error)
	GetAllAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error)
	CreateAmenityGroup(userID string, AmenityGroup *model.AmenityGroup) (*model.AmenityGroup, error)
	UpdateAmenityGroup(id string, userID string, data *model.AmenityGroup) (*model.AmenityGroup, error)
	DeleteAmenityGroup(id string) (model.AmenityGroup, error)
	GqlGetAmenityGroups(params domain.RequestParams) ([]*model.AmenityGroup, error)
}

type Repositories struct {
	Action
	Address
	Amenity
	AmenityGroup
	Authorization
	Lang
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

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		Action:        NewActionMongo(mongodb, i18n),
		Address:       NewAddressMongo(mongodb, i18n),
		Amenity:       NewAmenityMongo(mongodb, i18n),
		AmenityGroup:  NewAmenityGroupMongo(mongodb, i18n),
		Authorization: NewAuthMongo(mongodb),
		Lang:          NewLangMongo(mongodb, i18n),
		Country:       NewCountryMongo(mongodb, i18n),
		Currency:      NewCurrencyMongo(mongodb, i18n),
		Image:         NewImageMongo(mongodb, i18n),
		Review:        NewReviewMongo(mongodb, i18n),
		User:          NewUserMongo(mongodb, i18n),
		Track:         NewTrackMongo(mongodb, i18n),
		Node:          NewNodeMongo(mongodb, i18n),
		NodeVote:      NewNodeVoteMongo(mongodb, i18n),
		NodeAudit:     NewNodeAuditMongo(mongodb, i18n),
		Nodedata:      NewNodedataMongo(mongodb, i18n),
		NodedataVote:  NewNodedataVoteMongo(mongodb, i18n),
		Tag:           NewTagMongo(mongodb, i18n),
		Tagopt:        NewTagoptMongo(mongodb, i18n),
		Ticket:        NewTicketMongo(mongodb, i18n),

		Like: NewLikeMongo(mongodb, i18n),
	}
}

// func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
// 	var opts *options.FindOptions
// 	if pagination != nil {
// 		opts = &options.FindOptions{
// 			Skip:  pagination.GetSkip(),
// 			Limit: pagination.GetLimit(),
// 		}
// 	}

// 	return opts
// }

func createFilter[V any](filterData *V) any {
	var filter V

	filterReflect := reflect.ValueOf(filterData)
	// fmt.Println("========== filterReflect ===========")
	// fmt.Println("struct > ", filterReflect)
	// fmt.Println("struct type > ", filterReflect.Type())
	filterIndirectData := reflect.Indirect(filterReflect)
	// fmt.Println("filter data > ", filterIndirectData)
	// fmt.Println("filter numField > ", filterIndirectData.NumField())
	dataFilter := bson.M{}

	var tagJSON, tagPrimitive string
	for i := 0; i < filterIndirectData.NumField(); i++ {
		field := filterIndirectData.Field(i)
		if field.Kind() == reflect.Ptr {
			field = reflect.Indirect(field)
		}
		typeField := filterIndirectData.Type().Field(i)
		tag := typeField.Tag
		// tagBson = tag.Get("bson")
		tagJSON = tag.Get("json")
		tagPrimitive = tag.Get("primitive")
		switch field.Kind() {
		case reflect.String:
			value := field.String()
			if tagPrimitive == "true" {
				id, _ := primitive.ObjectIDFromHex(value)
				// fmt.Println("===== string add ", tag, value)
				dataFilter[tagJSON] = id
			} else {
				dataFilter[tagJSON] = value
			}

		case reflect.Bool:
			value := field.Bool()
			dataFilter[tagJSON] = value

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := field.Int()
			dataFilter[tagJSON] = value

		default:

		}

		// fmt.Println(tagBson, tagJSON, tagPrimitive, fmt.Sprintf("[%s]", field), field.Kind(), field)
	}

	// structure := reflect.ValueOf(&filter)
	// fmt.Println("========== filter ===========")
	// fmt.Println("struct > ", structure)
	// fmt.Println("struct type > ", structure.Type())
	// fmt.Println("filter data > ", reflect.Indirect(structure))
	// fmt.Println("filter numField > ", reflect.Indirect(structure).NumField())

	// fmt.Println("========== result ===========")
	// fmt.Println("dataFilter > ", dataFilter)
	return filter
}
