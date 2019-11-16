package sokindle

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"reflect"

	"ebook/engines/book"
	"ebook/engines/book/crawler"

	"github.com/ifnfn/util/stores"
)

type parserDownload struct {
	qiniu     stores.Store
	qiniuList []stores.Stat
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

	exists := false

	for _, f := range p.qiniuList {
		if f.Name == bok.FileName {
			exists = true
			break
		}
	}

	if exists {
		return false
	}

	if _, err := p.qiniu.Stat(bok.FileName); err == nil { // 如果七牛上文件存在
		return false
	}

	destPath := "books/" + bok.FileName
	run := exec.Command("bd", "download", bok.BaiduURL, "--secret="+bok.Baidupwd, "--dir="+destPath)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if run.Run() != nil { // 百度云下载失败
		return false
	}

	fstore := stores.NewLocalFileStore(destPath)
	dirList := fstore.List()

	// 如果找到 『.aria2』的文件，表示還沒有下載完
	downloadFinish := true
	for _, fi := range dirList {
		if path.Ext(fi.Name) == ".aria2" {
			downloadFinish = false
			println("下载未完成：", fi.Name)
			break
		}
	}

	if downloadFinish {
		for _, fi := range dirList {
			if value, err := fstore.Get(fi.Name); err == nil {
				println(fi.Name)
				defer value.Close()
				if p.qiniu.Save(bok.FileName, value) == nil {
					os.RemoveAll(destPath)
					return true
				}
				break
			}
		}
	}

	return false
}
