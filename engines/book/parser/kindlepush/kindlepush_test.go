package kindlepush

import (
	"fmt"
	"testing"
)

func TestGetGoodsIDFromURL(t *testing.T) {
	p := parserDownload{}
	p.login("zzgmtv@163.com", "780227CNSCZD")

	u := p.getURL("6251")
	fmt.Println(u)
	p.logout()
}
