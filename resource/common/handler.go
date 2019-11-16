package common

import (
	"ebook/common"

	"github.com/ifnfn/util/config"
	"github.com/rs/rest-layer-mongo"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/resource/testing/mem"

	mgo "gopkg.in/mgo.v2"
)

// Handler 存储器接口
type Handler interface {
	New(collection string) resource.Storer
}

// MemHandler ...
type MemHandler struct {
}

// MemHandlerInit ...
func MemHandlerInit() Handler {
	return &MemHandler{}
}

// New ...
func (h MemHandler) New(collection string) resource.Storer {
	return mem.NewHandler()
}

// MongoHandler ...
type MongoHandler struct {
	session *mgo.Session
	db      string
}

// MongoHandlerInit ...
func MongoHandlerInit() Handler {
	return &MongoHandler{
		db:      config.MongoDB.Database,
		session: common.MgoSession,
	}
}

// New ...
func (h MongoHandler) New(collection string) resource.Storer {
	return mongo.NewHandler(h.session, h.db, collection)
}
