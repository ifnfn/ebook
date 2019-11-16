package book

import (
	"ebook/resource/common"

	"github.com/rs/rest-layer/resource"
)

// ResourceBind 资源初始化
func ResourceBind(index resource.Index, handler common.Handler) {
	// 图书
	index.Bind("books", books, handler.New("books"), resource.Conf{
		AllowedModes: resource.ReadWrite,
	}).Use(booksHook{})

	// 书评
	index.Bind("reviews", books, handler.New("reviews"), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	// 帐号
	index.Bind("accounts", accounts, handler.New("accounts"), resource.Conf{
		AllowedModes: resource.ReadWrite,
	})
}
