package util

import (
	"fmt"
	"testing"

	"ebook/common"
	"ebook/engines/book"
)

func TestA(x *testing.T) {
	common.Init("/Users/zhuzhg/works/golang/src/ebook/config.json")
	bok := book.Book{
		ID:       "100",
		Name:     "布达佩斯往事",
		Content:  `布达佩斯往事`,
		FileName: "100: 把时间当作朋友.mobi",
	}

	Sendkindle("zzgmtv@163.com", bok)
}

func TestSendxxx(x *testing.T) {
	v := GetKindleEmail("我的邮件是  aaaaaaaaaa@kindle.cn")
	fmt.Println(v)

	v = GetKindleEmail("我的邮件是  中@163.com")
	println(v)
}
