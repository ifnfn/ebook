package douban

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/ifnfn/util/system"
	"github.com/PuerkitoBio/goquery"
	"github.com/smashedtoatoms/gofuzz"
)

type parserSearch struct {
	Cache bool // 是否从CACHE 下载
	Proxy bool // 不使用代理
}

// NewParserSearch ...
func NewParserSearch() crawler.SiteParser {
	return &parserSearch{
		Cache: true,  // 是否从CACHE 下载
		Proxy: false, // 不使用代理
	}
}

// Command ...
func (p parserSearch) Command(data interface{}) crawler.Command {
	return crawler.Command{
		Parser: fmt.Sprint(reflect.TypeOf(p)), // 解析器名称
		Data:   data,                          // 数据
	}
}

// Parser ...
func (p *parserSearch) Parser(cmd *crawler.Command) bool {
	bok := cmd.Data.(book.Book)
	foundBook := false
	search := ""
	if bok.Isbn != "" {
		search = bok.Isbn
	} else {
		search = bok.Name
	}

	if search == "" {
		return false
	}

	searchURL := fmt.Sprintf("https://book.douban.com/subject_search?search_text=%s&cat=1001", url.QueryEscape(search))

	doc, err := crawler.GetDocument(searchURL, nil, p.Cache, p.Proxy)
	if err != nil {
		log.Fatal(err)
	}

	var similarity float32
	var title string
	doc.Find("div.article ul.subject-list li.subject-item div.info h2 a").Each(func(i int, s *goquery.Selection) {
		if t, found := s.Attr("title"); found {
			s1 := strings.TrimSpace(t)
			s2 := strings.TrimSpace(s.Find("span").Text())

			if result := Similar(s1, s2, bok.Name); result > similarity {
				if href, found := s.Attr("href"); found {
					if match := regexp.MustCompile(`subject/(\d*)`).FindStringSubmatch(href); len(match) > 1 {
						bok.DoubanID = match[1]
						similarity = result
						title = t
					}
				}
			}
		}
	})

	if similarity > 0.79 {
		cmd.Data = bok
		foundBook = true
		if similarity < 0.80 {
			fmt.Printf("%1.4f, %s|%s\n", similarity, title, bok.Name)
		}
	}
	// if !foundBook && similarity > 0.0 {
	// 	fmt.Printf("%1.4f, %s|%s %v\n", similarity, title, bok.Name, foundBook)
	// }

	if bok.Isbn == "xxx" {
		system.PrintInterface(bok)
		panic("")
	}

	return foundBook
}

func fixChar(s string) string {
	delChar := []string{
		`(：|\(|\)|（|）|　|_|“|”|\"|・|·|《|》|\-|、)`,
		"(套装)", `([全|共].*[卷|册])`, "(套装|珍藏版|白金|纪念版|精装|插图版|插图本|哈佛中国史)",
	}

	for _, c := range delChar {
		reg := regexp.MustCompile(c)
		s = reg.ReplaceAllString(s, "")
	}

	return s
}

// Similar 字符串相似度比较
func Similar(s1, s2, s3 string) float32 {
	if len(s1+s2) == 0 || len(s3) == 0 {
		return 0.0
	}

	s1 = fixChar(s1)
	s2 = fixChar(s2)
	s3 = fixChar(s3)

	var ret float32
	if r, e := gofuzz.Jaro(s1, s3); e == nil && r > ret {
		ret = r
	}
	if r, e := gofuzz.Jaro(s2, s3); e == nil && r > ret {
		ret = r
	}
	if r, e := gofuzz.Jaro(s1+s2, s3); e == nil && r > ret {
		ret = r
	}

	return ret
}
