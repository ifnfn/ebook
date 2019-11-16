package book

import (
	"ebook/resource/common"

	"github.com/rs/rest-layer/schema"
)

var (
	reviews = schema.Schema{
		Description: "图书书品表",
		Fields: schema.Fields{
			"id":      schema.IDField, // 数据点 ID
			"bookID":  common.StringField,
			"Content": common.StringField,
		},
	}

	accounts = schema.Schema{
		Description: "帐号表",
		Fields: schema.Fields{
			"id":        schema.IDField, // 数据点 ID
			"wechatID":  common.StringField,
			"email":     common.StringField,
			"Favorites": common.StringField,
			"History":   common.StringField,
		},
	}
)
