package proxy

import (
	"container/list"
	"encoding/json"
	"sync"

	"github.com/ifnfn/util/system"
)

// HTTPProxy ...
type HTTPProxy struct {
	mutex    sync.Mutex
	urlQuery *list.List
}

// NewHTTPProxy 新建 Proxy
func NewHTTPProxy() *HTTPProxy {
	proxy := &HTTPProxy{}

	proxy.urlQuery = list.New()

	return proxy
}

func (p *HTTPProxy) getProxyCache() {
	type data struct {
		Count     int      `json:"count"`
		ProxyList []string `json:"proxy_list"`
	}

	type kuaidaili struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
		Data data   `json:"data"`
	}

	var kuai kuaidaili
	url := "http://ent.kuaidaili.com/api/getproxy/?orderid=936588863967175&num=100&b_pcchrome=1&b_pcie=1&b_pcff=1&carrier=2&protocol=2&method=1&an_an=1&an_ha=1&sp2=1&quality=2&sort=1&format=json&sep=1"

	if body, err := system.HTTPGet(url, nil); err == nil {
		json.Unmarshal(body, &kuai)
		for _, u := range kuai.Data.ProxyList {
			p.urlQuery.PushBack(u)
		}
	}
}

// URL 申请一个 http proxy 地址
func (p *HTTPProxy) URL() string {
	return "http://122.228.25.97:8101"
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.urlQuery.Len() == 0 {
		p.getProxyCache()
	}

	if p.urlQuery.Len() > 0 {
		v := p.urlQuery.Front()
		p.urlQuery.Remove(v)

		return "http://" + v.Value.(string)
	}

	return ""
}

// func GetByProxy(url_addr, proxy_addr string) (*http.Response, error) {
// 	request, _ := http.NewRequest("GET", url_addr, nil)
// 	proxy, err := url.Parse(proxy_addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			Proxy: http.ProxyURL(proxy),
// 		},
// 	}
// 	return client.Do(request)
// }

// // Fetch Httpclient
// func Fetch(urls, method string, headers map[string]string, data []byte) ([]byte, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, urls, bytes.NewReader(data))

// 	if proxyURL, exists := headers["Proxy"]; exists {
// 		delete(headers, "Proxy")
// 		if proxy, err := url.Parse(proxyURL); err == nil {
// 			client.Transport = &http.Transport{
// 				Proxy: http.ProxyURL(proxy),
// 			}
// 		}
// 	}

// 	for k, v := range headers {
// 		req.Header.Add(k, v)
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified || resp.StatusCode == http.StatusNotFound {
// 		return ioutil.ReadAll(resp.Body)
// 	}

// 	if v, e := ioutil.ReadAll(resp.Body); e == nil {
// 		println(string(v))
// 	}
// 	return nil, fmt.Errorf("http error, %d: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode), url)
// }

// // MD5url 将 URL 转为 MD5
// func MD5url(url string) string {
// 	return strings.ToUpper(system.GetMD5([]byte(url)))
// }

// func getCache(url string) ([]byte, error) {
// 	fileName := "./cache/" + MD5url(url)

// 	return ioutil.ReadFile(fileName)
// }

// func saveCache(url string, data []byte) {
// 	fileName := "./cache/" + MD5url(url)

// 	ioutil.WriteFile(fileName, data, 0644)
// }

// // CacheFetch ...
// func CacheFetch(url, method string, headers map[string]string, dody []byte, cache bool) ([]byte, error) {
// 	md5 := MD5url(url)
// 	if cache {
// 		if data, err := getCache(url); err == nil {
// 			println("cache->", md5, url)
// 			return data, err
// 		}
// 	}

// 	data, err := Fetch(url, method, headers, dody)

// 	if err == nil {
// 		saveCache(url, data)
// 	}

// 	println("GET: ", md5, url)

// 	return data, err
// }
