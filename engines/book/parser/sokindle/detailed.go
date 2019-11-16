package sokindle

import (
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/PuerkitoBio/goquery"
	"github.com/ifnfn/util/system"
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
	if bok.SokindleID == "" {
		return false
	}
	sokindleURL := fmt.Sprintf("https://sokindle.com/books/%s.html", bok.SokindleID)

	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Mobile Safari/537.36"
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Host"] = "sokindle.com"
	headers["Origin"] = "https://sokindle.com"
	headers["Referer"] = sokindleURL

	body := []byte("e_secret_key=523523")

	doc, err := crawler.PostDocument(sokindleURL, headers, body, p.Cache)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// 基本信息
	doc.Find("div.book-info div.book-left").Each(func(i int, s *goquery.Selection) {
		s.Find("div.bookpic img").Each(func(i int, s *goquery.Selection) {
			if href, f := s.Attr("src"); f {
				bok.Image = href
			}
		})

		s.Find("div.bookinfo li").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if match := regexp.MustCompile(`书名：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Name = match[1]
			}

			if match := regexp.MustCompile(`作者：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Author = match[1]
			}

			if match := regexp.MustCompile(`格式：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Formats = match[1]
			}

			if match := regexp.MustCompile(`浏览：(\d*)`).FindStringSubmatch(text); len(match) > 1 {
				if v, err := strconv.Atoi(match[1]); err == nil {
					bok.Count = int(v)
				}
			}

			if match := regexp.MustCompile(`标签：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Category = match[1]
			}

			if match := regexp.MustCompile(`时间：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Time = match[1]
			}

			if match := regexp.MustCompile(`ISBN：(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Isbn = match[1]
				if strings.Contains(bok.Isbn, "ISBN:") {
					if match := regexp.MustCompile(`ISBN: (.*)`).FindStringSubmatch(text); len(match) > 1 {
						bok.Isbn = match[1]
					}
				}
			}
		})
	})

	// 下载
	doc.Find("table.dltable a").Each(func(i int, s *goquery.Selection) {
		if href, f := s.Attr("href"); f {
			if match := regexp.MustCompile(`url=(http://pan.baidu.com/.*)`).FindStringSubmatch(href); len(match) > 1 {
				bok.BaiduURL = match[1]
			}
		}
	})

	// 下载密码
	doc.Find("div.e-secret strong").Each(func(i int, s *goquery.Selection) {
		if match := regexp.MustCompile(`提取密码：(.*)`).FindStringSubmatch(s.Text()); len(match) > 1 {
			bok.Baidupwd = match[1]
		}
	})

	// 文件大小
	doc.Find("i.fa-th-large").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Parent().Text())
		if match := regexp.MustCompile(`文件大小：(.*)`).FindStringSubmatch(text); len(match) > 1 {
			bok.Size = match[1]
		}
	})

	// qiniu 文件名
	doc.Find("i.fa-list-alt").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Parent().Text())
		if match := regexp.MustCompile(`文件名称：(.*)`).FindStringSubmatch(text); len(match) > 1 {
			bok.FileName = fmt.Sprintf("%s: %s", bok.SokindleID, match[1])
		}
		ext := path.Ext(bok.FileName)
		if ext == "" {
			ext = "." + strings.ToLower(bok.Formats)
		}

		bok.FileName = bok.SokindleID + "_" + strings.ToUpper(system.GetMD5([]byte(bok.Name))) + ext
	})

	level := 0
	doc.Find("article.article-content > *").Each(func(i int, s *goquery.Selection) {
		name := goquery.NodeName(s)
		text := s.Text()

		if name == "h2" {
			if text == "内容简介" {
				level = 1
			} else if text == "作者简介" {
				level = 2
			} else {
				level = 3
			}
		} else if name == "p" {
			if len(s.Find("table").Nodes) == 0 {
				if level == 1 {
					bok.Content += text + "\n"
				} else if level == 2 {
					bok.AuthorIntro += text + "\n"
				}
			}
		}
	})

	cmd.Data = bok

	return true
}
