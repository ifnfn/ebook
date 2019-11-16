package bookask

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/PuerkitoBio/goquery"
)

type parserISBN struct {
	Cache bool // 是否从CACHE 下载
	Proxy bool // 不使用代理
}

// NewParserISBN ...
func NewParserISBN() crawler.SiteParser {
	return &parserISBN{
		Cache: true,  // 是否从CACHE 下载
		Proxy: false, // 不使用代理
	}
}

// Command ...
func (p parserISBN) Command(data interface{}) crawler.Command {
	return crawler.Command{
		Parser: fmt.Sprint(reflect.TypeOf(p)), // 解析器名称
		Data:   data,                          // 数据
	}
}

// Parser ...
func (p *parserISBN) Parser(cmd *crawler.Command) bool {
	bok := cmd.Data.(book.Book)
	foundBook := false

	if bok.Isbn == "" {
		return false
	}

	url := fmt.Sprintf("http://www.bookask.com/s/kw_%s.html", bok.Isbn)
	doc, err := crawler.GetDocument(url, nil, p.Cache, p.Proxy)

	if err != nil {
		log.Fatal(err)
		return false
	}

	cmd.Data = nil
	doc.Find("div.s-tittle a.am-text-truncate").Each(func(i int, s *goquery.Selection) {
		if href, found := s.Attr("href"); found {
			if match := regexp.MustCompile(`/(\d*).html`).FindStringSubmatch(href); len(match) > 1 {
				bok.BookaskID = match[1]
				cmd.Data = bok
				foundBook = true
			}
		}
	})

	return foundBook
}
