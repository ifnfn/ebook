package bookask

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
		Cache: true,
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

	if bok.BookaskID == "" {
		return false
	}

	bookaskURL := fmt.Sprintf("https://www.bookask.com/book/%s.html", bok.BookaskID)
	doc, err := crawler.GetDocument(bookaskURL, nil, p.Cache, p.Proxy)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// 图片
	doc.Find("img.ba_page_prvimg").Each(func(i int, s *goquery.Selection) {
		if href, found := s.Attr("src"); found {
			bok.Image = strings.TrimSpace(href)
			found = true
		}
	})

	// 内容简介
	doc.Find("div.book-text-box#descript p").Each(func(i int, s *goquery.Selection) {
		bok.Content = strings.TrimSpace(s.Text())
	})

	// 目录
	doc.Find("div.book-text-catalog#catalog").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		bok.Topics = append(bok.Topics, html)
	})

	// 其他
	doc.Find("div.book-text-info div.fl div").Each(func(i int, s *goquery.Selection) {
		// println(s.Text())
		text := s.Text()

		if bok.Author == "" {
			if match := regexp.MustCompile(`作者：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Author = match[1]
			}
		}

		if bok.Press == "" {
			if match := regexp.MustCompile(`出版：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Press = match[1]
			}
		}

		if bok.Isbn == "" {
			if match := regexp.MustCompile(`ISBN：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				if bok.Isbn == "" {
					bok.Isbn = match[1]
				}
			}
		}

		if bok.PublicationDate == "" {
			if match := regexp.MustCompile(`出版日期：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.PublicationDate = match[1]
			}
		}
	})

	cmd.Data = bok

	return true
}
