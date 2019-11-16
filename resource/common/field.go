package common

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ifnfn/util/system"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/rest-layer/schema/query"
)

var (
	// Now is a field hook handler that returns the current time, to be used in
	// schema with OnInit and OnUpdate.
	Now = func(ctx context.Context, value interface{}) interface{} {
		return time.Now()
		// return time.Now().Unix()
		// return fmt.Sprintf("%d", ntime)
	}

	// CreatedField 创建时候字段
	CreatedField = schema.Field{
		Description: "The time at which the item has been inserted",
		Required:    true,
		ReadOnly:    true,
		OnInit:      Now,
		Sortable:    true,
		Validator:   &schema.Time{},
	}

	// UpdatedField is a common schema field configuration for "updated" fields. It stores
	// the current date each time the item is modified.
	UpdatedField = schema.Field{
		Description: "The time at which the item has been last updated",
		Required:    true,
		ReadOnly:    true,
		OnInit:      Now,
		OnUpdate:    Now,
		Sortable:    true,
		Validator:   &schema.Time{},
	}

	// IDField ID
	IDField = schema.Field{
		Required:   true,
		Filterable: true,
		Sortable:   true,
		Validator: &schema.String{
			// This regexp matches a base32 id
			Regexp: "^[0-9a-v]{20}$",
		},
	}

	// StringField Field
	StringField = schema.Field{
		Required:   false,
		Filterable: true,
		Validator:  &schema.String{},
	}

	// FloatField Field
	FloatField = schema.Field{
		Required:   false,
		Filterable: true,
		Validator:  &schema.Float{},
	}

	// IntegerField Field
	IntegerField = schema.Field{
		Required:   false,
		Filterable: true,
		Validator:  &schema.Integer{},
	}

	// BoolField Field
	BoolField = schema.Field{
		Required:   true,
		Filterable: true,
		Default:    false,
		Validator:  &schema.Bool{},
	}
)

// InsertDB 插入数据
func InsertDB(ctx context.Context, res *resource.Resource, payload map[string]interface{}) (*resource.Item, error) {
	changes, base := res.Validator().Prepare(ctx, payload, nil, false)
	doc, errs := res.Validator().Validate(changes, base)
	if len(errs) > 0 {
		system.PrintInterface(errs)
		system.PrintInterface(doc)

		return nil, errors.New(system.StructToString(errs))
	}

	item, err := resource.NewItem(doc)
	if err != nil {
		return nil, err
	}

	if err := res.Insert(ctx, []*resource.Item{item}); err != nil {
		return nil, err
	}

	return item, nil
}

// PatchDB 更新数据
func PatchDB(ctx context.Context, res *resource.Resource, original *resource.Item, payload map[string]interface{}) (*resource.Item, error) {
	changes, base := res.Validator().Prepare(ctx, payload, &original.Payload, false)
	doc, errs := res.Validator().Validate(changes, base)
	if len(errs) > 0 {
		system.PrintInterface(errs)
		system.PrintInterface(doc)

		return nil, errors.New(system.StructToString(errs))
	}

	item, err := resource.NewItem(doc)
	if err != nil {
		return nil, err
	}

	if err := res.Update(ctx, item, original); err != nil {
		return nil, err
	}

	return item, nil
}

// InsertOrPatchDB 不存在则插入，存在则补充
func InsertOrPatchDB(ctx context.Context, qx *query.Query, res *resource.Resource, payload map[string]interface{}) (*resource.Item, error) {
	if item, err := res.Find(ctx, qx); err == nil {
		if item.Total == 0 {
			return InsertDB(ctx, res, payload)
		} else if item.Total == 1 {
			return PatchDB(ctx, res, item.Items[0], payload)
		}
	}

	return InsertDB(ctx, res, payload)
}

// Struct2Map 将类转成MAP
func Struct2Map(obj interface{}) map[string]interface{} {
	var out map[string]interface{}
	if v, err := json.Marshal(obj); err == nil {
		err := json.Unmarshal(v, &out)
		if err != nil {
			println(err.Error())
		}
	}

	return out
}

// Map2Struct 将 map 转为类
func Map2Struct(data map[string]interface{}, out interface{}) {
	if v, err := json.Marshal(data); err == nil {
		json.Unmarshal(v, out)
	}
}
