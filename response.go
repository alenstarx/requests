package requests

import (
	"io/ioutil"
	"net/http"
	"compress/gzip"
	"encoding/json"
	"io"
)

type Response struct {
	*http.Response
	err error
}

func (r *Response) Err() error {
	return r.err
}

func (r *Response) GetCookie() map[string]string {
	cs := r.Cookies()
	if len(cs) == 0 {
		return nil
	}
	cookie := make(map[string]string)
	for _, v := range cs {
		cookie[v.Name] = v.Value
	}
	return cookie
}
func (r *Response) GetHeader() map[string]string {
	var header map[string]string
	for k, values := range r.Header {
		value := ""
		for _, v := range values {
			if value != "" {
				value = value + ";" + v
			} else {
				value = v
			}
		}
		if header == nil {
			header = make(map[string]string)
		}
		header[k] = value
	}
	return header
}

func (r *Response) Bytes() []byte {
	if r.err != nil || r.Body == nil {
		return nil
	}

	var body []byte
	if r.Header.Get("Content-Encoding") == "gzip" {
		var reader io.Reader
		reader, r.err = gzip.NewReader(r.Body)
		if r.err != nil {
			return nil
		}
		body, r.err = ioutil.ReadAll(reader)
		return body
	}

	body, r.err = ioutil.ReadAll(r.Body)
	return body
}

func (r *Response) BindJson(obj interface{}) error {
	body := r.Bytes()
	if r.err != nil {
		return r.err
	}
	return json.Unmarshal(body, obj)
}

func (r *Response) StoreCookie(req *Request) *Response{
	c := r.GetCookie()
	if c != nil {
		req.SetCookie(c)
	}
	return nil
}

func (r *Response) DumpToFile(filename string) error {
	// TODO
	return nil
}

func (r *Response) SaveToFile(filename string) error {
	// TODO
	return nil
}