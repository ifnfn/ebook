package common

import (
	mgo "gopkg.in/mgo.v2"

	"github.com/ifnfn/util/config"
	"github.com/ifnfn/util/system"
)

// Config
var (
	MgoSession *mgo.Session
)

// Init 配置初始化
func Init(configFile string) {
	config.NewConfig(configFile)
	MgoSession, _ = system.NewMongoClient()

	config.Qiniu.AccessKey = "xxxxxxxxxxxxxxxxxx"
	config.Qiniu.SecretKey = "xxxxxxxxxxxxxxx"
	config.Qiniu.Domain = "xxxxxxxxxxxxx"
	config.Qiniu.Bucket = "xxxxxxxxxxxxx"

	config.Wechat.AppID = "xxxxxxxxxxxxxxxxxx"
	config.Wechat.AppSecret = "xxxxxxxxxxxxxxxxxx"
	config.Wechat.Token = "xxxxxxxxxxxxxxxxxx"
	config.Wechat.EncodedAESKey = "xxxxxxxxxxxxxxxxxx"
}
