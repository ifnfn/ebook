package kindlepush

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/ifnfn/util/stores"
	"github.com/ifnfn/util/system"
)

type parserDownload struct {
	qiniu     stores.Store
	qiniuList []stores.Stat
	cookie    string
}

// NewParserDownload ...
func NewParserDownload() crawler.SiteParser {
	q := stores.NewQiniuStore("books", "kindle.ifnfn.com")
	return &parserDownload{
		qiniu:     q,
		qiniuList: q.List(),
	}
}

// Command ...
func (p parserDownload) Command(data interface{}) crawler.Command {
	return crawler.Command{
		Parser: fmt.Sprint(reflect.TypeOf(p)), // 解析器名称
		Data:   data,                          // 数据
	}
}

// Parser ...
func (p *parserDownload) Parser(cmd *crawler.Command) bool {
	bok := cmd.Data.(book.Book)

	if bok.FileName == "" {
		return false
	}

	return true
}

func (p *parserDownload) login(username, password string) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("rememberMe", "false")

	if resp, err := http.PostForm("http://www.kindlepush.com/user/login", data); err == nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "JSESSIONID" {
				p.cookie = cookie.Raw
			}
		}
	}
}

func (p *parserDownload) logout() {
	headers := make(map[string]string)
	headers["Cookie"] = p.cookie
	system.Fetch("http://www.kindlepush.com/user/logout", "GET", headers, nil)
}

func (p *parserDownload) getURL(resourceID string) string {
	headers := map[string]string{
		"Cookie":       p.cookie,
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "Mozilla/5.0 (iPad; CPU OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	}

	body := []byte("resourceId=" + resourceID)

	if data, err := system.Fetch("http://www.kindlepush.com/upload/download/", "POST", headers, body); err == nil {
		switch string(data) {
		case "error_too_many": // 今日下载量已满
			fmt.Println("这本书今天已经下载过了，不能重复下载哦")
		case "error_repeat": // 这本书今天已经下载过了，不能重复下载哦~
			fmt.Println("这本书今天已经下载过了，不能重复下载哦~")
		case "error_param": // 资源不存在
			fmt.Println("资源不存在")
		case "no_login":
			fmt.Println("未登录")
		}

		return string(data)
	}

	return ""
}
