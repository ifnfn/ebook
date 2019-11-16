package httpd

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Httpd) routeInit() {
	h.echo.GET("/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	h.echo.GET("/book/:bookID", h.getBook)
	h.echo.GET("/book/push_kindle", h.pushKindle)
	h.echo.GET("/book/favorite", h.favorite)

	h.echo.Static("/bootstrap", "bootstrap")
	h.echo.Static("/js", "js")

	h.echo.Any("/wechat/kd_callback", echo.WrapHandler(h.kindle.Handler()))

	h.echo.Any("/474888328294966zhuZHG/*", echo.WrapHandler(
		http.StripPrefix("/api", h.kindle.BookRest.Handler()),
	))

	h.echo.Any("/api/*", echo.WrapHandler(
		http.StripPrefix("/api", h.kindle.BookRest.Handler()),
	))

	// h.echo.Any("/handler", echo.WrapHandler(
	// 	http.HandlerFunc(
	// 		func(w http.ResponseWriter, r *http.Request) {
	// 			qiniuHealth(w, r)
	// 		},
	// 	),
	// ))
}
