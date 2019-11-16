package book

import (
	"ebook/resource/common"
	"path"

	"github.com/ifnfn/util/config"
	"github.com/ifnfn/util/stores"
)

// Book 电子图书
type Book struct {
	ID              string   `json:"id"`
	Author          string   // 作者
	Name            string   // 书名
	Isbn            string   // 号号
	Content         string   // 内容简介
	AuthorIntro     string   // 作者简介
	Topics          []string // 目录表
	Category        string   // 分类
	Press           string   // 出版社
	PublicationDate string   // 出版日期
	Image           string   // 图片
	ReviewCount     int      // 评论次数
	Rating          float64  // 豆瓣评分
	Meta            int      // 评价人数
	Tags            string   // 标签
	Reviews         []string // 书评
	Count           int      // 下载次数
	Time            string   // 上架时间
	Formats         string   // 文件格式
	DoubanID        string   // 豆瓣页
	SokindleID      string   // Sokindle ID
	KindlePushID    string   //
	BookaskID       string   // 书问
	Baidupwd        string   // 百度下载密码
	BaiduURL        string   // 百度下载
	FileName        string   // 文件名
	Size            string   // 文件大小
	Finish          bool     // 完成
}

// Fill 从 MAP转为类
func (bok *Book) Fill(m map[string]interface{}) error {
	common.Map2Struct(m, bok)

	return nil
}

// Map 转为 MAP
func (bok Book) Map() map[string]interface{} {
	payload := common.Struct2Map(bok)
	if payload["Topics"] == nil {
		delete(payload, "Topics")
	}

	if payload["Reviews"] == nil {
		delete(payload, "Reviews")
	}

	return payload
}

// Merge 将参数指定的 Book中未空的字段，覆盖原字段数据
func (bok *Book) Merge(b Book) {
	if b.Author != "" {
		bok.Author = b.Author
	}

	if b.Name != "" {
		bok.Name = b.Name
	}

	if b.Isbn != "" {
		bok.Isbn = b.Isbn
	}

	if b.Content != "" {
		bok.Content = b.Content
	}

	if b.AuthorIntro != "" {
		bok.AuthorIntro = b.AuthorIntro
	}

	if len(b.Topics) >= 0 {
		bok.Topics = b.Topics
	}

	if b.Category != "" {
		bok.Category = b.Category
	}

	if b.Press != "" {
		bok.Press = b.Press
	}

	if b.PublicationDate != "" {
		bok.PublicationDate = b.PublicationDate
	}

	if b.Image != "" {
		bok.Image = b.Image
	}

	if b.ReviewCount > 0 {
		bok.ReviewCount = b.ReviewCount
	}

	if b.Rating > 0.0 {
		bok.Rating = b.Rating
	}

	if b.Meta > 0 {
		bok.Meta = b.Meta
	}

	if b.Tags != "" {
		bok.Tags = b.Tags
	}

	if len(b.Reviews) > 0 {
		bok.Reviews = b.Reviews
	}

	if b.Count > 0 {
		bok.Count = b.Count
	}

	if b.Time != "" {
		bok.Time = b.Time
	}

	if b.Formats != "" {
		bok.Formats = b.Formats
	}

	if b.DoubanID != "" {
		bok.DoubanID = b.DoubanID
	}

	if b.BookaskID != "" {
		bok.BookaskID = b.BookaskID
	}

	if b.SokindleID != "" {
		bok.SokindleID = b.SokindleID
	}

	if b.KindlePushID != "" {
		bok.KindlePushID = b.KindlePushID
	}

	if b.Baidupwd != "" {
		bok.Baidupwd = b.Baidupwd
	}

	if b.BaiduURL != "" {
		bok.BaiduURL = b.BaiduURL
	}

	if b.FileName != "" {
		bok.FileName = b.FileName
	}

	if b.Size != "" {
		bok.Size = b.Size
	}
	if b.Finish {
		bok.Finish = b.Finish
	}
}

// DownloadFile ...
func (bok *Book) DownloadFile() string {
	fileName := bok.FileName
	fstore := stores.NewLocalFileStore(config.Server.UploadPath)
	if _, err := fstore.Stat(fileName); err != nil { // 如果本地缓冲不存在
		qstore := stores.NewQiniuStore("books", "kindle.ifnfn.com")

		err := stores.Copy(fileName, qstore, fstore)
		if err != nil {
			println(err.Error())

			return ""
		}
	}

	if _, err := fstore.Stat(fileName); err == nil { // 如果本地缓冲不存在
		return path.Join(config.Server.UploadPath, fileName)
	}

	return ""
}
