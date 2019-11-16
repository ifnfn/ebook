package res

import (
	"context"
	"log"
	"net/http"

	"ebook/resource/book"
	"ebook/resource/common"

	"github.com/ifnfn/util/config"
	"github.com/justinas/alice"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	"github.com/rs/xaccess"
	"github.com/rs/xlog"
)

// Resource 资源总表
type Resource struct {
	Index  resource.Index
	Stores common.Handler
}

// NewResource 创建资源
func NewResource() *Resource {
	res := &Resource{}
	res.Init()

	return res
}

// Init 资源初始化
func (r *Resource) Init() {
	// Create a REST API resource index
	r.Index = resource.NewIndex()

	db := config.Server.Database

	db = "mem"
	// db = "mongo"

	if db == "mongodb" {
		r.Stores = common.MongoHandlerInit()
	} else {
		r.Stores = common.MemHandlerInit()
	}

	book.ResourceBind(r.Index, r.Stores)
}

// GetResource 从上下文中获取资源对对象
func GetResource(ctx context.Context) *Resource {
	if db, ok := ctx.Value("resource").(*Resource); ok {
		return db
	}

	return nil
}

func ctxHandler(db *Resource) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract context from request
			ctx := r.Context()

			// // Store it into the request's context
			ctx = context.WithValue(ctx, interface{}("resource"), db)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// Handler 资源服务开始
func (r *Resource) Handler() http.Handler {
	// Create API HTTP handler for the resource graph
	api, err := rest.NewHandler(r.Index)
	if err != nil {
		log.Fatalf("Invalid API configuration: %s", err)
		return nil
	}

	// Setup logger
	c := alice.New()
	c = c.Append(xlog.NewHandler(xlog.Config{}))
	c = c.Append(xaccess.NewHandler())

	// resource.LoggerLevel = resource.LogLevelDebug
	// resource.Logger = func(ctx context.Context, level resource.LogLevel, msg string, fields map[string]interface{}) {
	// 	xlog.FromContext(ctx).OutputF(xlog.Level(level), 2, msg, fields)
	// }

	c = c.Append(ctxHandler(r))
	return c.Then(api)
}
