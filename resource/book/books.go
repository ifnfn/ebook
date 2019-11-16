package book

import (
	"context"
	"ebook/resource/common"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema"
	"github.com/rs/rest-layer/schema/query"
)

var (
	// Define a datapoint resource schema
	books = schema.Schema{
		Description: "图书基本信息表",
		Fields: schema.Fields{
			"id":          common.IDField,     // 数据点 ID
			"Author":      common.StringField, // 作者
			"Name":        common.StringField, // 书名
			"Isbn":        common.StringField, // 号号
			"Content":     common.StringField, // 内容简介
			"AuthorIntro": common.StringField, // 作者简介
			"Topics": { //  目录表
				Required:   false,
				Filterable: false,
				Validator: &schema.Array{
					ValuesValidator: &schema.String{},
				},
			},
			"Reviews": schema.Field{
				Required:   false,
				Filterable: false,
				Validator: &schema.Array{
					ValuesValidator: &schema.String{},
				},
			},
			"Category":        common.StringField,  // 分类
			"Tags":            common.StringField,  // 标签
			"Press":           common.StringField,  // 出版社
			"PublicationDate": common.StringField,  // 出版日期
			"Image":           common.StringField,  // 图片
			"ReviewCount":     common.IntegerField, // 评论次数
			"Rating":          common.FloatField,   // 豆瓣评分
			"Meta":            common.IntegerField, // 评价人数
			"Count":           common.IntegerField, // 下载次数
			"Time":            common.StringField,  // 上架时间
			"Formats":         common.StringField,  // 文件格式
			"BookaskID":       common.StringField,  // 书问页
			"DoubanID":        common.StringField,  // 豆瓣页
			"SokindleID":      common.StringField,  // sokindle ID
			"KindlePushID":    common.StringField,  // kindlepush 网页
			"Baidupwd":        common.StringField,  // 百度下载密码
			"BaiduURL":        common.StringField,  // 百度下载
			"FileName":        common.StringField,  // 文件名
			"Size":            common.StringField,  // 文件大小
			"Finish":          common.BoolField,    // 资源已准备好
			"Created":         common.CreatedField,
			"Update":          common.UpdatedField,
		},
	}
)

type booksHook struct {
}

func (a booksHook) onFind(ctx context.Context, q *query.Query, list **resource.ItemList, err *error) {
}
