package robot

import (
	"ebook/engines/book"
	"ebook/engines/book/crawler"
	"ebook/engines/book/parser/bookask"
	"ebook/engines/book/parser/douban"
	"ebook/engines/book/parser/kindlepush"
	"ebook/engines/book/parser/sokindle"

	"github.com/rs/rest-layer/resource"
)

// Robot ...
type Robot struct {
	Library
	task               *crawler.TaskGroup
	SokindleDetailed   crawler.SiteParser
	SokindleList       crawler.SiteParser
	SokindleDownload   crawler.SiteParser
	BookaskIsbn        crawler.SiteParser
	BookaskDetailed    crawler.SiteParser
	DoubanSearch       crawler.SiteParser
	DoubanIsbn         crawler.SiteParser
	DoubanDetailed     crawler.SiteParser
	KindlepushList     crawler.SiteParser
	KindlepushDetailed crawler.SiteParser
}

// Init ...
//
//
//                 SokindleDownload
//                       ^
//                       |
//                     Save                                      Save                                 Save
//                       ^               (isbn != "")             ^                                    ^
//                       |           +-> doubanIsbn  ---+         |                                    |
//                       |           |                  |         |                                    |
// sokindleList -> sokindleDetailed--+                  +-> DoubanDetailed ---> BookaskIsbn -> BookaskDetailed
//                                   |                  |         ^
//                                   +-> oubanSearch ---+         |
//                                                                |
// KindlepushList -> KindlepushDetailed --------------------------+
//
func (r *Robot) Init(index resource.Index) {
	r.Library = NewLibraryDB(index)
	r.Library.Load()
	r.task = crawler.NewTaskGroup(1)

	// 增加 sokindle 网站解析器
	r.SokindleList = r.task.SetParser(sokindle.NewParserList(), r.sokindleListCmd)
	r.SokindleDetailed = r.task.SetParser(sokindle.NewParserDetailed(), func(cmd crawler.Command) {
		bok := cmd.Data.(book.Book)
		bok.Finish = true
		r.Set(bok)

		// 解析豆瓣
		if bok.Isbn != "" {
			r.Command(r.DoubanIsbn, bok)
		} else {
			r.Command(r.DoubanSearch, bok)
		}
	})
	r.SokindleDownload = r.task.SetParser(sokindle.NewParserDownload(), func(cmd crawler.Command) {
		// bok := cmd.Data.(book.Book)
		// println("download:", bok.Name)
	})

	// 增加 douban 网站解析器
	r.DoubanIsbn = r.task.SetParser(douban.NewParserISBN(), r.doubanSearchCmd)
	r.DoubanSearch = r.task.SetParser(douban.NewParserSearch(), r.doubanSearchCmd)
	r.DoubanDetailed = r.task.SetParser(douban.NewParserDetailed(), func(cmd crawler.Command) {
		bok := cmd.Data.(book.Book)

		if bok.Content == "" && bok.Isbn != "" { // 如何豆瓣没有找到
			// 解析书问
			r.Command(r.BookaskIsbn, bok)
		} else {
			r.Set(bok) // 写入数据库
		}
	})

	// 增加 kindlepush 网站解析器
	r.KindlepushList = r.task.SetParser(kindlepush.NewParserList(), func(cmd crawler.Command) {
		data := cmd.Data.(map[string]interface{})
		books := data["books"].([]book.Book)
		count := 0
		for _, bok := range books {
			if _, found := r.Find("", bok.ID); !found {
				// r.Command(r.DoubanSearch, bok)//  启动 Douban 搜索
				r.Command(r.KindlepushDetailed, bok)
				count++
			}
		}
		// r.Command(r.SokindleList, data["next"])
	})
	r.KindlepushDetailed = r.task.SetParser(kindlepush.NewParserDetailed(), func(cmd crawler.Command) {
		bok := cmd.Data.(book.Book)
		r.Command(r.DoubanSearch, bok) //  启动 Douban 搜索
		// r.Set(bok)
	})

	// 增加 bookask 网站解析器
	r.BookaskIsbn = r.task.SetParser(bookask.NewParserISBN(), func(cmd crawler.Command) {
		bok := cmd.Data.(book.Book)
		r.Command(r.BookaskDetailed, bok)
	})
	r.BookaskDetailed = r.task.SetParser(bookask.NewParserDetailed(), func(cmd crawler.Command) {
		r.Set(cmd.Data.(book.Book))
	})
}

// Run ...
func (r *Robot) Run() {
	r.Command(r.SokindleList, nil)
	// r.Command(r.KindlepushList, nil)
}

// Wait 等待任务结束
func (r *Robot) Wait() {
	r.task.Wait()
	r.Save()
}

func (r *Robot) sokindleListCmd(cmd crawler.Command) {
	data := cmd.Data.(map[string]interface{})
	books := data["books"].([]book.Book)
	count := 0
	for _, bok := range books {
		println(bok.SokindleID, bok.Name)
		if _, found := r.Find("SokindleID", bok.SokindleID); !found {
			r.Command(r.SokindleDetailed, bok)
			count++
		}
	}

	if count > 0 {
		r.Command(r.SokindleList, data["next"])
	}
}

func (r *Robot) doubanSearchCmd(cmd crawler.Command) {
	bok := cmd.Data.(book.Book)
	// 解析豆瓣详细
	r.Command(r.DoubanDetailed, bok)
}

// Command 增加命令
func (r *Robot) Command(paser crawler.SiteParser, data interface{}) {
	r.task.AddCommand(paser.Command(data))
}
