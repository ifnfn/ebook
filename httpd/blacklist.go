package httpd

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

var blackListRegex []*regexp.Regexp

// Blacklist 黑名单中间件
func Blacklist() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			remoteAddr := c.RealIP()

			addr := GetIPAddress(remoteAddr)

			for _, b := range blackListRegex {
				if b.MatchString(remoteAddr) {
					println("Backlist ADDR:", remoteAddr, b.String())
					return nil
				}
				if b.MatchString(addr.City) {
					println("Backlist City:", addr.City, b.String())
					return nil
				}
			}

			return next(c)
		}
	}
}

func init() {
	if data, err := ioutil.ReadFile("blacklist.txt"); err == nil {
		for _, s := range strings.Split(string(data), "\n") {
			if len(s) > 0 {
				println("add blacklist:", s)
				blackListRegex = append(blackListRegex, regexp.MustCompile(s))
			}
		}
	}
}
