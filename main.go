package main

import (
	"runtime"

	"ebook/httpd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http := httpd.Httpd{}
	http.Run("/etc/phoenix/config.json")
}
