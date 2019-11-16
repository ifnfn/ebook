package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"ebook/engines/book"
	"ebook/engines/book/robot"
	res "ebook/resource"
	"ebook/util"

	"github.com/ifnfn/util/config"
	"github.com/ifnfn/util/system"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"
)

// ohI6oxC9Dr3IIm8oekNQulg9d49U 三体
// ohI6oxLjpQ9fEn-QG0lkzciUPtNs 朱
// ohI6oxPocYyh1PzSvJB-_WtkannM 姚

// KindlePush 公共号服务
type KindlePush struct {
	msgHandler core.Handler
	msgServer  *core.Server
	msgClient  *core.Client
	BooksRobot robot.Robot
	BookRest   *res.Resource
	Accounts   map[string]*Account
}

// NewKindlePush 新建微信公共号
func NewKindlePush() *KindlePush {
	w := &KindlePush{}

	appID := "wx777625f4ef2809f5"
	appSecret := "5959d3df3ed5c3ebd8ab2b6d76bd6f17"
	token := "daiR8choh7ezahp5ahzumiph2Iphahng"
	encodedAESKey := "r4GjvsIrCywQtI3v1rRqDyxQ5iZwebT7A4HUIjEFxn1"
	w.Init(appID, appSecret, token, encodedAESKey)

	w.BookRest = res.NewResource()
	w.BooksRobot.Init(w.BookRest.Index)
	w.Accounts = make(map[string]*Account)

	if jsonBytes, err := ioutil.ReadFile("accounts.json"); err == nil {
		if err := json.Unmarshal(jsonBytes, &w.Accounts); err != nil {
			fmt.Printf("Could not parse %v", err)
		}
	}

	return w
}

// Save 保存信息
func (wechat *KindlePush) Save() {
	if jsonBytes, err := json.MarshalIndent(wechat.Accounts, "", "    "); err == nil {
		ioutil.WriteFile("accounts.json", jsonBytes, 0644)
	}
}

// Init 微信公共号初始化
func (wechat *KindlePush) Init(appID, appSecret, token, encodedKey string) {
	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(wechat.defaultMsgHandler)
	mux.DefaultEventHandleFunc(wechat.defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, wechat.textMsgHandler)
	mux.MsgHandleFunc(request.MsgTypeImage, wechat.imagetMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, wechat.menuClickEventHandler)

	wechat.msgHandler = mux
	wechat.msgServer = core.NewServer("", appID, token, encodedKey, wechat.msgHandler, nil)

	httpClient := &http.Client{}
	tokenServer := core.NewDefaultAccessTokenServer(appID, appSecret, httpClient)
	if token, err := tokenServer.Token(); err == nil {
		println(token)
	}

	wechat.msgClient = core.NewClient(tokenServer, httpClient)
}

func (wechat *KindlePush) imagetMsgHandler(ctx *core.Context) {
	// image := request.GetImage(ctx.MixedMsg)

	// ret, err := system.HTTPGet(image.PicURL, nil)
	// if err == nil {
	// 	util.UploadImage(ret, stores.NewQiniuStore())
	// }

	ctx.NoneResponse()
}

// GetAcccount 根据ID 得到用户帐户信息
func (wechat *KindlePush) GetAcccount(userID string) *Account {
	if account, found := wechat.Accounts[userID]; found {
		return account
	}

	return nil
}

func (wechat *KindlePush) textMsgHandler(ctx *core.Context) {
	var err error
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)

	message := request.GetText(ctx.MixedMsg)

	// 检查用户帐户
	userID := system.GetMD5([]byte(ctx.MixedMsg.FromUserName))
	account := wechat.GetAcccount(userID)
	if account == nil {
		account = NewAccount(userID)
		wechat.Accounts[userID] = account
	}

	content := message.Content

	// if strings.Contains(content, "http://www.bookask.com") {
	// 	ctx.NoneResponse()
	// 	return
	// }

	// 帮助
	if content == "帮助" {
		text := `从次无限制 Kindle 电子书免费推送。

第一步：关注 『Kindle分享』公众号
第二步：将 kindle@ifnfn.com 添加为你的 Kindle 推送信任邮件
第三步：将您的亚马逊邮箱  “xxxxxx@kindle.cn” 发送到公众号
第四步：通过公众号聊天方式，输入您想到的找的电子书名，公众号回复到的电子书信息，进入详细页面，点南『推送至 Kindle』`

		resp := response.NewText(message.FromUserName, message.ToUserName, message.CreateTime, text)
		ctx.AESResponse(resp, 0, "", nil) // aes密文回复

		return
	}

	// 如果转入的是邮件地址，设置为 kindle 接受地址
	if email := util.GetKindleEmail(content); email != "" {
		account.Email = email
		wechat.Save()

		resp := response.NewText(message.FromUserName, message.ToUserName, message.CreateTime, "您的 Kindle 推送邮件是："+email)
		ctx.AESResponse(resp, 0, "", nil) // aes密文回复

		return
	}

	if err = wechat.lsBooks(ctx, content); err != nil {
		if err = wechat.lsBooks(ctx, ".*"); err != nil {
			resp := response.NewText(message.FromUserName, message.ToUserName, message.CreateTime, "未查找电子书："+content)
			ctx.AESResponse(resp, 0, "", nil) // aes密文回复
		}
	}
}

// BookToArticles ...
func (wechat *KindlePush) BookToArticles(msg *core.MixedMsg, books []book.Book) *response.News {
	count := len(books)
	if count > 6 {
		count = 6
	}

	articles := make([]response.Article, count)

	for i := 0; i < count; i++ {
		desc := []rune(books[i].Content)

		var content string

		if len(desc) > 100 {
			desc = desc[:100]
			content = string(desc[:120]) + " ..."
		} else {
			content = books[i].Content
		}

		userID := system.GetMD5([]byte(msg.FromUserName))
		articles[i] = response.Article{
			Title:       books[i].Name,
			Description: content,
			PicURL:      books[i].Image,
			URL:         config.Server.HTTP(fmt.Sprintf("book/%s?id=%s", books[i].ID, userID)),
		}
	}

	return response.NewNews(
		msg.FromUserName,
		msg.ToUserName,
		msg.CreateTime,
		articles,
	)
}

func (wechat *KindlePush) lsBooks(ctx *core.Context, text string) error {
	if books, found := wechat.BooksRobot.Match(text, 0, 10); found {
		news := wechat.BookToArticles(ctx.MixedMsg, books)

		return ctx.AESResponse(news, time.Now().Unix(), "", nil) // aes 密文回复
	}

	return errors.New("No found")
}

func (wechat *KindlePush) defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func (wechat *KindlePush) menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")

	//ctx.RawResponse(resp) // 明文回复
	ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func (wechat *KindlePush) defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

// Handler Httpd 请求入口
func (wechat *KindlePush) Handler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wechat.msgServer.ServeHTTP(w, r, nil)
		},
	)
}
