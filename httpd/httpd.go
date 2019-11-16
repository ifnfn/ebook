package httpd

import (
	"ebook/common"
	"ebook/engines/wechat"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ifnfn/util/config"
)

// Httpd 服务
type Httpd struct {
	echo   *echo.Echo
	kindle *wechat.KindlePush
}

// Run ...
func (h *Httpd) Run(conf string) {
	DefaultLoggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  "${time_rfc3339} ${method}, ${remote_ip}[${city}:${isp}], ${uri}, status=${status}\n",
		// Output: os.Stdout,
	}

	common.Init(conf)

	h.kindle = wechat.NewKindlePush()

	h.echo = echo.New()
	h.echo.Use(AddressInfo()) // 将地址转为地址
	h.echo.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	h.echo.Use(middleware.Recover())
	h.echo.Use(middleware.Gzip())

	// //允许跨站请求
	h.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccessControlRequestMethod},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	h.routeInit()

	go qiniuDocker()

	h.echo.Logger.Fatal(
		h.echo.Start(config.Server.HTTPBindAddress()),
	)
}
