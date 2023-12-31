package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LibraryID   primitive.ObjectID     `json:"libraryId" bson:"libraryId" primitive:"true"`
// Library    []Library             `json:"library" bson:"library"`
// SchemaData []ComponentSchemaData `json:"schema_data" bson:"schema_data"`

type ComponentSchema struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty" primitive:"true"`
	ComponentID primitive.ObjectID     `json:"componentId" bson:"componentId" primitive:"true"`
	Data        map[string]interface{} `json:"data" bson:"data"`

	CreatedAt time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type ComponentData struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" primitive:"true"`
	Parent string             `json:"parent" bson:"parent"`
	UID    string             `json:"_uid" bson:"_uid"`
	PageID primitive.ObjectID `json:"pageId" bson:"page_id" primitive:"true"`
	// LayoutID  primitive.ObjectID `json:"layoutId" bson:"layout_id" primitive:"true"`
	Component string `json:"component" bson:"component"`

	Publish   bool                   `json:"publish" bson:"publish" form:"publish"`
	Data      map[string]interface{} `json:"data" bson:"data" form:"data"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type ComponentSchemaData struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" primitive:"true"`
	SchemaID    primitive.ObjectID `json:"schemaId" bson:"schemaId" primitive:"true"`
	ComponentID primitive.ObjectID `json:"componentId" bson:"componentId" primitive:"true"`
	UserID      primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	Name        string             `json:"name" bson:"name"`
	HistoryID   int                `json:"historyId" bson:"historyId"`

	Data      map[string]interface{} `json:"data" bson:"data"`
	Publish   *bool                  `json:"publish" bson:"publish" form:"publish"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type Component struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty" primitive:"true"`
	UserID  primitive.ObjectID `json:"userId" bson:"user_id"  primitive:"true"`
	SpaceID primitive.ObjectID `json:"spaceId" bson:"space_id" primitive:"true"`
	Name    string             `json:"name" bson:"name" `
	Type    string             `json:"type" bson:"type" `
	Title   string             `json:"title" bson:"title" `
	Publish bool               `json:"publish" bson:"publish" `
	Status  string             `json:"status" bson:"status" `
	// IsPage    bool                   `json:"is_page" bson:"is_page" form:"is_page"`
	// IsGlobal  bool                   `json:"is_global" bson:"is_global" form:"is_global"`
	// IsLayout  bool                   `json:"is_layout" bson:"is_layout" form:"is_layout"`
	Tpl       string             `json:"tpl" bson:"tpl" form:"tpl"`
	SortOrder int                `json:"sortOrder" bson:"sort_order" `
	Group     primitive.ObjectID `json:"group" bson:"group" primitive:"true"`
	// Groups    []ComponentGroup                  `json:"groups" bson:"groups"`
	Schema  map[string]map[string]interface{} `json:"schema" bson:"schema" `
	Setting map[string]interface{}            `json:"setting" bson:"setting" `

	Presets   []ComponentPreset `json:"presets" bson:"presets"`
	CreatedAt time.Time         `json:"createdAt" bson:"created_at" `
	UpdatedAt time.Time         `json:"updatedAt" bson:"updated_at"`
}

type ComponentInput struct {
	SpaceID   primitive.ObjectID                `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	UserID    primitive.ObjectID                `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Name      string                            `json:"name" bson:"name" form:"name"`
	Type      string                            `json:"type" bson:"type" form:"type"`
	Title     string                            `json:"title" bson:"title" form:"title"`
	Publish   bool                              `json:"publish" bson:"publish" form:"publish"`
	Status    string                            `json:"status" bson:"status" form:"status"`
	Tpl       string                            `json:"tpl" bson:"tpl" form:"tpl"`
	SortOrder int                               `json:"sortOrder" bson:"sort_order" form:"sortOrder"`
	Schema    map[string]map[string]interface{} `json:"schema" bson:"schema" form:"schema"`
	Group     string                            `json:"group" bson:"group" form:"group"`
	// Groups    []string                          `json:"groups" bson:"groups" form:"groups" primitive:"true"`
	Setting   map[string]interface{} `json:"setting" bson:"setting" form:"setting"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at" form:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updated_at" form:"updatedAt"`
}

type ComponentFind struct {
	SpaceID   string                            `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	UserID    primitive.ObjectID                `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	ID        string                            `json:"id" bson:"_id" form:"id" primitive:"true"`
	Name      string                            `json:"name" bson:"name" form:"name"`
	Type      string                            `json:"type" bson:"type" form:"type"`
	Title     string                            `json:"title" bson:"title" form:"title"`
	Publish   bool                              `json:"publish" bson:"publish" form:"publish"`
	Status    string                            `json:"status" bson:"status" form:"status"`
	Tpl       string                            `json:"tpl" bson:"tpl" form:"tpl"`
	SortOrder int                               `json:"sortOrder" bson:"sort_order" form:"sortOrder"`
	Schema    map[string]map[string]interface{} `json:"schema" bson:"schema" form:"schema"`
	Group     string                            `json:"group" bson:"group" form:"group"`
	// Groups    []string                          `json:"groups" bson:"groups" form:"groups" primitive:"true"`
	Setting   map[string]interface{} `json:"setting" bson:"setting" form:"setting"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at" form:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updated_at" form:"updatedAt"`
}

type ComponentInputByGroup struct {
	Group []string `json:"group[]" bson:"group" form:"group[]"`
}

// func (component *Component) BodyToData() (interface{}, error) {
// 	result := bson.M{}
// 	var tagValue string
// 	elementsFilter := reflect.ValueOf(component)

// 	for i := 0; i < elementsFilter.NumField(); i++ {
// 		typeField := elementsFilter.Type().Field(i)
// 		tag := typeField.Tag

// 		tagValue = tag.Get("bson")

// 		if tagValue == "-" {
// 			continue
// 		}

// 		if elementsFilter.Field(i).Interface() == "" {
// 			continue
// 		}

// 		switch elementsFilter.Field(i).Kind() {
// 		case reflect.String:
// 			value := elementsFilter.Field(i).String()
// 			result[tagValue] = value

// 		case reflect.Bool:
// 			value := elementsFilter.Field(i).Bool()
// 			result[tagValue] = value

// 		case reflect.Int:
// 			value := elementsFilter.Field(i).Int()
// 			result[tagValue] = value
// 		}
// 	}
// 	if len(component.Group) == 0 || component.Group == nil {
// 		result["group"] = []primitive.ObjectID{}
// 		// for i := 0; i < len(component.Group) - 1; i += 1 {
// 		// 	 id, err := primitive.ObjectIDFromHex(component.Group[i])
// 		// 	 if err != nil {
// 		// 		return result, nil
// 		// 	 }
// 		// 	 result["group"][i] = id
// 		// }
// 	} else {
// 		result["group"] = component.Group
// 	}
// 	// if component.IsPage != nil {
// 	// 	result["is_page"] = *component.IsPage
// 	// }
// 	// if component.IsGlobal != nil {
// 	// 	result["is_global"] = *component.IsGlobal
// 	// }
// 	// if component.IsLayout != nil {
// 	// 	result["is_layout"] = *component.IsLayout
// 	// }
// 	// if component.Status != nil {
// 	// 	result["status"] = *component.Status
// 	// }
// 	// if component.Publish != nil {
// 	// 	result["publish"] = *component.Publish
// 	// }

// 	result["updated_at"] = time.Now()
// 	// fmt.Println("user: new data =", result)
// 	return result, nil
// }
