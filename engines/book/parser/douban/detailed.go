package douban

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
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
		Cache: true,  // 是否从CACHE 下载
		Proxy: false, // 不使用代理
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

	if bok.DoubanID == "" {
		return false
	}

	headers := map[string]string{
		"Host":   "book.douban.com",
		"Accept": "*/*",
	}
	doubanURL := fmt.Sprintf("https://book.douban.com/subject/%s/", bok.DoubanID)
	doc, err := crawler.GetDocument(doubanURL, headers, p.Cache, p.Proxy)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if len(doc.Find("html.ua-webkit").Nodes) > 0 {
		p.Parserx(&bok, doc)
	} else {
		p.ParserM(&bok, doc)
	}

	p.review(&bok)
	cmd.Data = bok

	return true
}

func (p *parserDetailed) Parserx(bok *book.Book, doc *goquery.Document) {
	doc.Find("div.article div.indent").Each(func(i int, s *goquery.Selection) {
		// 图片
		if attr, found := s.Find("div#mainpic a.nbg").Attr("href"); found {
			bok.Image = strings.TrimSpace(attr)
		}
		// 图片
		s.Find("div#info").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			if match := regexp.MustCompile(`出版社: (.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Press = match[1]
			}
			if match := regexp.MustCompile(`ISBN: (.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.Isbn = match[1]
			}
			if match := regexp.MustCompile(`出版年:(.*)`).FindStringSubmatch(text); len(match) > 1 {
				bok.PublicationDate = match[1]
			}
		})
	})

	doc.Find("div#interest_sectl").Each(func(i int, s *goquery.Selection) {
		// 豆瓣评分
		s.Find("strong.rating_num ").Each(func(i int, s *goquery.Selection) {
			if v, err := strconv.ParseFloat(strings.TrimSpace(s.Text()), 64); err == nil {
				bok.Rating = v
			}
		})

		s.Find("a.rating_people span").Each(func(i int, s *goquery.Selection) {
			if v, err := strconv.Atoi(strings.TrimSpace(s.Text())); err == nil {
				bok.ReviewCount = v
			}
		})
	})

	doc.Find(`div.related_info`).Each(func(i int, s *goquery.Selection) {
		// 内容简介
		bok.Content = s.Find(`div.indent > span.hidden p`).Text()

		if bok.Content == "" {
			bok.Content = doc.Find(`div.indent div.intro`).Text()
		}
		// 作者简介
		bok.AuthorIntro = s.Find(`div.indent > div > div.intro`).Text()
	})

	// 标签
	var tags []string
	doc.Find(`div#db-tags-section > div.indent span a`).Each(func(i int, s *goquery.Selection) {
		if text := strings.TrimSpace(s.Text()); text != "" {
			tags = append(tags, text)
		}
	})
	bok.Tags = strings.Join(tags, ",")
}

// Parser ...
func (p *parserDetailed) ParserM(bok *book.Book, doc *goquery.Document) {
	// 图片
	if attr, found := doc.Find("div.article div#mainpic a.nbg").Attr("href"); found {
		bok.Image = strings.TrimSpace(attr)
	}

	// 豆瓣评分
	doc.Find("strong.rating_num ").Each(func(i int, s *goquery.Selection) {
		if v, err := strconv.ParseFloat(strings.TrimSpace(s.Text()), 64); err == nil {
			bok.Rating = v
		}
	})

	doc.Find("a.rating_people span").Each(func(i int, s *goquery.Selection) {
		if v, err := strconv.Atoi(strings.TrimSpace(s.Text())); err == nil {
			bok.ReviewCount = v
		}
	})

	doc.Find(`div.related_info`).Each(func(i int, s *goquery.Selection) {
		// 内容简介
		bok.Content = s.Find(`div.indent > span.hidden p`).Text()

		if bok.Content == "" {
			bok.Content = doc.Find(`div.indent div.intro`).Text()
		}
		// 作者简介
		bok.AuthorIntro = s.Find(`div.indent > div > div.intro`).Text()
	})

	// 标签
	var tags []string
	doc.Find(`div#db-tags-section > div.indent span a`).Each(func(i int, s *goquery.Selection) {
		if text := strings.TrimSpace(s.Text()); text != "" {
			tags = append(tags, text)
		}
	})
	bok.Tags = strings.Join(tags, ",")
}

func (p *parserDetailed) review(bok *book.Book) bool { // 解析图书详细信息
	found := false

	headers := map[string]string{
		"Host":       "book.douban.com",
		"Accept":     "*/*",
		"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	}
	reviewsURL := fmt.Sprintf("https://m.douban.com/book/subject/%s/reviews", bok.DoubanID)
	doc, err := crawler.GetDocument(reviewsURL, headers, true, false)
	if err != nil {
		log.Fatal(err)
		return false
	}

	doc.Find("section.review-list li").Each(func(i int, s *goquery.Selection) {
		s.Find("div.info").Remove()
		if html, err := s.Html(); err == nil {
			bok.Reviews = append(bok.Reviews, html)
			found = true
		}
	})

	return found
}
