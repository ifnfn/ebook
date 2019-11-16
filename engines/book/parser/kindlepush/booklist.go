package kindlepush

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/ifnfn/util/system"
	"github.com/PuerkitoBio/goquery"
)

// parserList ...
type parserList struct {
	Cache bool // 是否从CACHE 下载
	Proxy bool // 不使用代理
}

// NewParserList ...
func NewParserList() crawler.SiteParser {
	return &parserList{
		Cache: true,  // 是否从CACHE 下载
		Proxy: false, // 不使用代理
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
		cmd.Data = "http://www.kindlepush.com/category/-1/1/1"
	}

	return cmd
}

// Parser ...
func (p *parserList) Parser(cmd *crawler.Command) bool {
	listURL := cmd.Data.(string)
	doc, err := crawler.GetDocument(listURL, nil, p.Cache, p.Proxy)

	if err != nil {
		log.Println(err.Error())
		return false
	}

	var books []book.Book

	doc.Find("div.info > div.wrap").Each(func(i int, s *goquery.Selection) {
		bok := book.Book{}
		// 详细网址与书名
		s.Find("a.title").Each(func(i int, s *goquery.Selection) {
			bok.Name = strings.TrimSpace(s.Text())
			if href, f := s.Attr("href"); f {
				if match := regexp.MustCompile(`book/(\d*)`).FindStringSubmatch(href); len(match) > 1 {
					bok.KindlePushID = match[1]
				}
			}
		})

		// 豆瓣评分
		s.Find("div.u-stargrade span").Each(func(i int, s *goquery.Selection) {
			// println(s.Text())
			if match := regexp.MustCompile(`\((.*)\)`).FindStringSubmatch(s.Text()); len(match) > 1 {
				if v, err := strconv.ParseFloat(match[1], 64); err == nil {
					bok.Rating = v
				}
			}
		})

		// 作者
		s.Find("div.u-author span").Each(func(i int, s *goquery.Selection) {
			bok.Author = strings.TrimSpace(s.Text())
		})

		bok.ID = system.GetMD5([]byte(bok.Name + bok.Author))[:20] // xid.New().String()
		fmt.Printf("%s %s %s %3.1f\n", bok.ID, bok.Name, bok.Author, bok.Rating)

		books = append(books, bok)
	})

	data := make(map[string]interface{})
	data["next"] = ""
	data["books"] = books

	doc.Find("div.u-page a.next").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			data["next"] = href
		}
	})
	cmd.Data = data

	return true
}
