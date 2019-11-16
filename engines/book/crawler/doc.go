package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"ebook/proxy"

	"github.com/PuerkitoBio/goquery"
	"github.com/ifnfn/util/system"
)

var defaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
}

var httpProxy = proxy.NewHTTPProxy()

// GetDocument ...
func GetDocument(url string, headers map[string]string, cache bool, proxy bool) (*goquery.Document, error) {
	var err error
	var data []byte
	retry := 1
	if proxy {
		retry = 10
	}

	if url == "" {
		return nil, errors.New("URL is blank")
	}

	tmpHeaders := make(map[string]string)
	for k, v := range defaultHeaders {
		tmpHeaders[k] = v
	}
	for k, v := range headers {
		tmpHeaders[k] = v
	}

	for i := 0; i < retry; i++ {
		if proxy {
			tmpHeaders["Proxy"] = httpProxy.URL()
		}
		data, err = system.CacheFetch(url, "GET", tmpHeaders, nil, cache)
		if err == nil {
			return toDocument(url, data)
		}
		println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", err.Error())
	}

	return nil, err
}

// PostDocument ...
func PostDocument(url string, headers map[string]string, dody []byte, cache bool) (*goquery.Document, error) {
	var err error
	var data []byte

	if url == "" {
		return nil, errors.New("URL is blank")
	}

	tmpHeaders := make(map[string]string)
	for k, v := range defaultHeaders {
		tmpHeaders[k] = v
	}

	for k, v := range headers {
		tmpHeaders[k] = v
	}
	// tmpHeaders["Proxy"] = httpProxy.URL()

	data, err = system.CacheFetch(url, "POST", tmpHeaders, dody, cache)
	if err == nil {
		return toDocument(url, data)
	}

	return nil, err
}

func toDocument(url string, data []byte) (*goquery.Document, error) {
	if match := regexp.MustCompile(`<html(.*)`).FindStringSubmatch(string(data)); len(match) <= 1 {
		return nil, fmt.Errorf("cache/%s %s 不是 HTML 文件", system.MD5url(url), url)
	}

	io := bytes.NewBuffer(data)
	return goquery.NewDocumentFromReader(io)
}
