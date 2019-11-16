package douban

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

	urls := fmt.Sprintf("https://book.douban.com/subject_search?search_text=%s&cat=1001", bok.Isbn)

	doc, err := crawler.GetDocument(urls, nil, p.Cache, p.Proxy)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.article ul.subject-list li.subject-item div.info h2 a").Each(func(i int, s *goquery.Selection) {
		if href, found := s.Attr("href"); found {
			if match := regexp.MustCompile(`subject/(\d*)`).FindStringSubmatch(href); len(match) > 1 {
				bok.DoubanID = match[1]
				cmd.Data = bok
				foundBook = true
			}
		}
	})

	return foundBook
}
