package httpx

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
)

type client struct {
	client        *http.Client
	defalutHeader map[string]string
}

type Client interface {
	clientPrivate()
	Do(string, ...optsFunc) ([]byte, error)
	AddHeader(map[string]string)
	DoTest(api string, optsF ...optsFunc) ([]byte, error)
}

func New(optsF ...optsFunc) Client {
	o := &cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(o)
	opt := &opts{
		deadline:    7,
		httpVersion: "1",
	}
	for _, optF := range optsF {
		optF(opt)
	}
	var clientTran = &http.Transport{}
	if opt.proxy != "" {
		clientTran.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(opt.proxy)
		}

	}
	if opt.httpVersion == "2" {
		http2.ConfigureTransport(clientTran)
	}

	// netAddr := &net.TCPAddr{Port: 29321}
	// dialer := &net.Dialer{LocalAddr: netAddr}
	// clientTran.DialContext = dialer.DialContext
	return &client{
		client: &http.Client{
			Jar:       jar,
			Transport: clientTran,
			Timeout:   opt.deadline * time.Second,
		},
		defalutHeader: opt.header,
	}

}

func (c *client) clientPrivate() {}

func (c *client) AddHeader(header map[string]string) {
	for key, v := range header {
		c.defalutHeader[key] = v
	}
}

func (c *client) Do(api string, optsF ...optsFunc) ([]byte, error) {
	opt := &opts{
		method: "GET",
	}
	for _, optF := range optsF {
		optF(opt)
	}
	// 新的实例就设置代理
	// if c.client.Transport == nil {
	if opt.proxy != "" {
		var clientTran = &http.Transport{}
		clientTran.Proxy = func(_ *http.Request) (*url.URL, error) {

			return url.Parse(opt.proxy)
		}
		c.client.Transport = clientTran
		defer func() { c.client.Transport = nil }()
	}

	// }
	req, err := http.NewRequest(opt.method, api, opt.body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.defalutHeader {
		req.Header.Set(k, v)
	}
	for k, v := range opt.header {
		req.Header.Set(k, v)
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	switch rsp.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, _ = gzip.NewReader(rsp.Body)
	case "deflate":
		bodyReader = flate.NewReader(rsp.Body)
	default:
		bodyReader = rsp.Body
	}
	return io.ReadAll(bodyReader)

}

type opts struct {
	httpVersion string
	method      string
	proxy       string
	header      map[string]string
	body        io.ReadCloser
	deadline    time.Duration
}

type optsFunc func(*opts)

func SetMethod(method string) optsFunc {
	return func(o *opts) {
		o.method = method
	}
}
func SetProxy(proxy string) optsFunc {
	return func(o *opts) {
		o.proxy = proxy
	}
}

// 默认 1, 如果启用 http2 ,请输入 2
func SetHttp2(version string) optsFunc {
	return func(o *opts) {
		o.httpVersion = version
	}
}

func SetDeadLine(deadline time.Duration) optsFunc {
	return func(o *opts) {
		o.deadline = deadline
	}
}

func SetHeader(header map[string]string) optsFunc {
	return func(o *opts) {
		o.header = header
	}
}

func SetBody(body io.ReadCloser) optsFunc {
	return func(o *opts) {
		o.body = body
	}
}

func SetParams(params map[string]string) optsFunc {
	return func(o *opts) {
		uv := url.Values{}
		var body io.ReadCloser
		for k, v := range params {
			uv.Set(k, v)
		}
		body = ioutil.NopCloser(strings.NewReader(uv.Encode()))
		o.body = body
	}
}

func (c *client) DoTest(api string, optsF ...optsFunc) ([]byte, error) {
	opt := &opts{
		method: "GET",
	}
	for _, optF := range optsF {
		optF(opt)
	}
	// 新的实例就设置代理
	// if c.client.Transport == nil {
	if opt.proxy != "" {
		var clientTran = &http.Transport{}
		clientTran.Proxy = func(_ *http.Request) (*url.URL, error) {

			return url.Parse(opt.proxy)
		}
		c.client.Transport = clientTran
		defer func() { c.client.Transport = nil }()
	}
	req, err := http.NewRequest(opt.method, api, opt.body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.defalutHeader {
		req.Header.Set(k, v)
	}

	for k, v := range opt.header {
		req.Header.Set(k, v)
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	switch rsp.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, _ = gzip.NewReader(rsp.Body)
	case "deflate":
		bodyReader = flate.NewReader(rsp.Body)
	default:
		bodyReader = rsp.Body
	}
	return io.ReadAll(bodyReader)

}
func SetJson(params map[string]interface{}) optsFunc {
	return func(o *opts) {
		pStr, _ := json.Marshal(params)
		o.body = ioutil.NopCloser(strings.NewReader(string(pStr)))
	}
}

func GenUrlParamsStr(params map[string]string) string {
	var uv = url.Values{}
	if len(params) > 0 {
		for k, v := range params {
			uv.Set(k, v)
		}
	}
	return uv.Encode()
}
