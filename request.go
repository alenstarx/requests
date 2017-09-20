package requests

import (
	"bytes"
	"crypto/tls"
	"golang.org/x/net/http2"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"io"
	"errors"
)

type Request struct {
	*http.Client
	URL    *url.URL
	Body   io.Reader
	Method string
	Header map[string]string
	err    error
}

func initialization(r *Request) {
	r.Header["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	r.Header["Accept-Lanaguage"] = "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3"
	r.Header["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64; rv:55.0) Gecko/20100101 Firefox/55.0"
	r.Header["Accept-Encoding"] = "gzip"
	r.Header["DNT"] = "1"
}

func reset(r *Request) {
	if _, ok := r.Header["Cookie"]; ok {
		delete(r.Header, "Cookie")
	}
}

func NewRequest() *Request {
	r := &Request{
		&http.Client{},
		nil,
		nil,
		"GET",
		make(map[string]string),
		nil,
	}
	initialization(r)
	return r
}

func (r *Request) Err() error {
	return r.err
}

func (r *Request) DelCookie() *Request {
	if _, ok := r.Header["Cookie"]; ok {
		delete(r.Header, "Cookie")
	}
	return r
}

func (r *Request) SetCookie(cookie map[string]string) *Request {
	for k, v := range cookie {
		if c, ok := r.Header["Cookie"]; ok {
			if len(c) > 0 {
				r.Header["Cookie"] = c + "; " + k + "=" + v
			} else {
				r.Header["Cookie"] = k + "=" + v
			}
		} else {
			r.Header["Cookie"] = k + "=" + v
		}
	}
	return r
}

func (r *Request) AddCookie(name, value string) *Request {
	if c, ok := r.Header["Cookie"]; ok {
		if len(c) > 0 {
			r.Header["Cookie"] = c + "; " + name + "=" + value
		} else {
			r.Header["Cookie"] = name + "=" + value
		}
	} else {
		r.Header["Cookie"] = name + "=" + value
	}

	return r
}

func (r *Request) SetUserAgent(userAgent string) *Request {
	r.Header["User-Agent"] = userAgent
	return r
}

func (r *Request) SetHeader(header map[string]string) *Request {
	for k, v := range header {
		r.Header[k] = v
	}
	return r
}
func (r *Request) EnableSPDY() *Request {
	var tr *http.Transport
	if r.Transport == nil {
		tr = &http.Transport{}
	}
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http2.ConfigureTransport(tr)
	r.Transport = tr
	return r
}

func (r *Request) SetProxy(proxy string) *Request {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		r.err = err
		return r
	}
	var tr *http.Transport
	if r.Transport == nil {
		tr = &http.Transport{}
	}
	tr.Proxy = http.ProxyURL(proxyUrl)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r.Transport = tr
	return r
}

func (r *Request) JSON(obj interface{}) *Response {
	r.Header["Content-Type"] = "application/json"
	r.SetPayload(obj)
	return r.POST()
}
func (r *Request) FORM(form map[string]string) *Response {
	r.Header["Content-Type"] = "application/x-www-form-urlencoded"
	r.SetPayload(form)
	return r.POST()
}
func (r *Request) GET() *Response {
	r.Method = "GET"
	return r.do()
}
func (r *Request) POST() *Response {
	r.Method = "POST"
	r.URL.RawQuery = ""
	return r.do()
}
func (r *Request) PUT() *Response {
	r.Method = "PUT"
	return r.do()
}
func (r *Request) PATCH() *Response {
	r.Method = "PATCH"
	return r.do()
}
func (r *Request) DELETE() *Response {
	r.Method = "DELETE"
	return r.do()
}
func (r *Request) HEAD() *Response {
	r.Method = "HEAD"
	return r.do()
}
func (r *Request) OPTIONS() *Response {
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
func (r *Request) SetParam(param map[string]string) *Request {
	if r.URL == nil {
		r.err = errors.New("invalid URL")
		return r
	}
	q := r.URL.Query()
	for k, v := range param {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()
	return r
}

func (r *Request) SetPayload(payload interface{}) *Request {
	switch payload.(type) {
	case string:
		r.Body = strings.NewReader(payload.(string))
		r.err = nil
	case []byte:
		r.Body = bytes.NewReader(payload.([]byte))
		r.err = nil
	case map[string]string:
		m := payload.(map[string]string)
		values := url.Values{}
		for k, v := range m {
			values.Add(k, v)
		}
		form := values.Encode()
		r.Body = strings.NewReader(form)
		r.err = nil
	default:
		body, err := json.Marshal(payload)
		if err != nil {
			r.err = err
			break
		}
		r.Body = bytes.NewReader(body)
		r.err = nil
	}
	return r
}

func (r *Request) do() *Response {
	if r.err != nil {
		return &Response{nil, r.err}
	}
	var resp *http.Response
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		r.err = err
		return &Response{nil, r.err}
	}
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	resp, r.err = r.Do(req)
	return &Response{resp, r.err}
}
