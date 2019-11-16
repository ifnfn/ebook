package httpd

import (
	"github.com/labstack/echo"
	"github.com/ifnfn/util/system"
)

// Address 位置
type Address struct {
	Region string
	City   string
	ISP    string
}

var ipInfo map[string]Address

// GetIPAddress 获取I地址位置信息
func GetIPAddress(ip string) Address {
	if _, exists := ipInfo[ip]; !exists {
		u := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
		if data, err := system.HTTPGetJSON(u, nil); err == nil {
			a := data.(map[string]interface{})
			if int(a["code"].(float64)) == 0 {
				data := a["data"].(map[string]interface{})
				if data["city"].(string) != "" {
					ipInfo[ip] = Address{
						data["region"].(string),
						data["city"].(string),
						data["isp"].(string),
					}
				}
			}
		}
	}

	return ipInfo[ip]
}

// AddressInfo 转地址
func AddressInfo() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			addr := GetIPAddress(c.RealIP())
			c.Set("Region", addr.Region)
			c.Set("City", addr.City)
			c.Set("ISP", addr.ISP)

			return next(c)
		}
	}
}

func init() {
	ipInfo = make(map[string]Address)
}
