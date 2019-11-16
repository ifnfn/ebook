package kindlepush

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/PuerkitoBio/goquery"
)

type parserDetailed struct {
	Cache bool // 是否从CACHE 下载
	Proxy bool // 不使用代理
}

// NewParserDetailed ...
func NewParserDetailed() crawler.SiteParser {
	return &parserDetailed{
		Cache: true, // 是否从CACHE 下载
		Proxy: false,
	}
}

// Command ...
func (p parserDetailed) Command(data interface{}) crawler.Command {
	return crawler.Command{
		Parser: fmt.Sprint(reflect.TypeOf(p)), // 解析器名称
		Data:   data,                          // 数据
	}
}

// Parser ...
func (p *parserDetailed) Parser(cmd *crawler.Command) bool {
	bok := cmd.Data.(book.Book)

	if bok.KindlePushID == "" {
		return false
	}
	kindlePushURL := "http://www.kindlepush.com/book/" + bok.KindlePushID
	doc, err := crawler.GetDocument(kindlePushURL, nil, p.Cache, p.Proxy)

	if err != nil {
		log.Fatal(err)
		return false
	}

	doc.Find("div#cover-img img").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("src"); exists {
			bok.Image = href
		}
	})

	doc.Find("article#book-content").Each(func(i int, s *goquery.Selection) {
		bok.Content = strings.TrimSpace(s.Text())
	})

	// 类型：『人文社科』
	doc.Find("div.cnt div.data h5").Each(func(i int, s *goquery.Selection) {
		if match := regexp.MustCompile(`类型：『(.*)』`).FindStringSubmatch(s.Text()); len(match) > 1 {
			bok.Category = strings.TrimSpace(match[1])
		}
	})

	cmd.Data = bok

	return true
}
