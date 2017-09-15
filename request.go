package requests

import (
	"bytes"
	"crypto/tls"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"log"
	"io"
	"encoding/json"
)

type Request struct {
	*http.Request
	client *http.Client
	err    error
}

func initialization(r *Request) {
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Header.Set("Accept-Lanaguage", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	r.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:55.0) Gecko/20100101 Firefox/55.0")
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("DNT", "1")
}

func reset(r *Request) {
	r.Header.Del("Cookie")
}
func NewRequest() *Request {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Fatalln("http.NewRequest:", err.Error())
	}
	r := &Request{
		req,
		&http.Client{},
		nil,
	}
	return r
}

func (r *Request) Err() error {
	return r.err
}

func (r *Request) DelCookie() *Request {
	r.Header.Del("Cookie")
	return r
}

func (r *Request) SetCookie(cookie map[string]string) *Request {
	for k, v := range cookie {
		r.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return r
}
func (r *Request) SetUserAgent(userAgent string) *Request {
	r.Header.Set("User-Agent", userAgent)
	return r
}

func (r *Request) SetHeader(header map[string]string) *Request {
	for k, v := range header {
		r.Header.Set(k, v)
	}
	return r
}
func (r *Request) EnableSPDY() *Request {
	var tr *http.Transport
	if r.client.Transport == nil {
		tr = &http.Transport{}
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http2.ConfigureTransport(tr)
	r.client.Transport = tr
	return r
}

func (r *Request) SetProxy(proxy string) *Request {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		r.err = err
		return r
	}
	var tr *http.Transport
	if r.client.Transport == nil {
		tr = &http.Transport{}
	}
	tr.Proxy = http.ProxyURL(proxyUrl)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r.client.Transport = tr
	return r
}

func (r *Request) PostJson(obj interface{}) *Response {
	r.Method = "POST"
	r.Header.Set("Content-Type", "application/json")
	return r.do()
}
func (r *Request) PostForm(form map[string]string) *Response {
	r.Method = "POST"
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r.do()
}
func (r *Request) Get() *Response {
	r.Method = "GET"
	return r.do()
}
func (r *Request) Post() *Response {
	r.Method = "POST"
	return r.do()
}
func (r *Request) Put() *Response {
	r.Method = "PUT"
	return r.do()
}
func (r *Request) Patch() *Response {
	r.Method = "PATCH"
	return r.do()
}
func (r *Request) Delete() *Response {
	r.Method = "DELETE"
	return r.do()
}
func (r *Request) Head() *Response {
	r.Method = "HEAD"
	return r.do()
}
func (r *Request) Options() *Response {
	r.Method = "OPTIONS"
	return r.do()
}

func (r *Request) SetUrl(rawurl string) *Request {
	r.URL, r.err = url.Parse(rawurl)
	if r.err != nil {
		return r
	}
	return r
}
func (r *Request) SetParams(params map[string]string) *Request {
	q := r.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()
	r.err = nil
	return r
}

func (r *Request) SetPayload(payload interface{}) *Request {
	switch payload.(type) {
	case string:
		reader := strings.NewReader(payload.(string))
		r.ContentLength = int64(reader.Len())
		snapshot := *reader
		r.GetBody = func() (io.ReadCloser, error) {
			rd := snapshot
			return ioutil.NopCloser(&rd), nil
		}
	case []byte:
		reader := bytes.NewReader(payload.([]byte))
		r.ContentLength = int64(reader.Len())
		snapshot := *reader
		r.GetBody = func() (io.ReadCloser, error) {
			rd := snapshot
			return ioutil.NopCloser(&rd), nil
		}
	case map[string]string:
		m := payload.(map[string]string)
		values := url.Values{}
		for k, v := range m {
			values.Add(k, v)
		}
		form := values.Encode()
		reader := strings.NewReader(form)
		r.ContentLength = int64(reader.Len())
		snapshot := *reader
		r.GetBody = func() (io.ReadCloser, error) {
			rd := snapshot
			return ioutil.NopCloser(&rd), nil
		}
	default:
		body, err := json.Marshal(payload)
		if err != nil {
			panic("not support parameter type")
		}
		reader := bytes.NewReader(body)
		r.ContentLength = int64(reader.Len())
		snapshot := *reader
		r.GetBody = func() (io.ReadCloser, error) {
			rd := snapshot
			return ioutil.NopCloser(&rd), nil
		}
	}
	r.err = nil
	return r
}

func (r *Request) do() *Response {
	var resp *http.Response
	resp, r.err = r.client.Do(r.Request)
	//defer resp.Body.Close()
	return &Response{resp, r.err}
}
