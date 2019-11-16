package sokindle

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/PuerkitoBio/goquery"
	"github.com/ifnfn/util/system"
)

type parserList struct {
	Cache bool // 是否从CACHE 下载
	Proxy bool // 不使用代理
}

// NewParserList ...
func NewParserList() crawler.SiteParser {
	return &parserList{
		Cache: true,  // 是否从CACHE 下载
		Proxy: false, //
	}
}

// Command ...
func (p parserList) Command(data interface{}) crawler.Command {
	cmd := crawler.Command{
		Parser: fmt.Sprint(reflect.TypeOf(p)), // 解析器名称
		Data:   data,                          // 数据
	}
	if data != nil {
		cmd.Data = data
	} else {
		cmd.Data = "https://sokindle.cc/page/1"
	}

	return cmd
}

// Parser ...
func (p *parserList) Parser(cmd *crawler.Command) bool {
	listURL := cmd.Data.(string)
	header := map[string]string{
		"Host":            "book.douban.com",
		"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6",
	}
	doc, err := crawler.GetDocument(listURL, header, p.Cache, p.Proxy)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	var books []book.Book

	doc.Find("div.cardlist > div.span_1_of_4 div.shop-item").Each(func(i int, s *goquery.Selection) {
		bok := book.Book{}
		s.Find("div.thumb-img > a").Each(func(i int, s *goquery.Selection) {
			if href, f := s.Attr("href"); f {
				if match := regexp.MustCompile(`/(\d*).html`).FindStringSubmatch(href); len(match) > 1 {
					bok.SokindleID = match[1]
				}
			}
			if title, f := s.Attr("title"); f {
				bok.Name = strings.TrimSpace(title)
			}
		})
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if match := regexp.MustCompile(`作者：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Author = match[1]
			}

			if match := regexp.MustCompile(`格式：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Formats = match[1]
			}
		})

		bok.ID = system.GetMD5([]byte(bok.Name + bok.Author))[:20] // xid.New().String()
		// println(bok.ID, bok.Formats, bok.Name, bok.Author)

		books = append(books, bok)
	})

	data := make(map[string]interface{})
	data["next"] = ""
	data["books"] = books

	doc.Find("div.pagination li.next-page a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			data["next"] = href
		}
	})
	cmd.Data = data

	return true
}
