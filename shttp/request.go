package shttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	GET     = Method("GET")
	POST    = Method("POST")
	PUT     = Method("PUT")
	DELETE  = Method("DELETE")
	HEAD    = Method("HEAD")
	PATCH   = Method("PATCH")
	OPTIONS = Method("OPTIONS")
	TRACE   = Method("TRACE")

	protocolVersion = "HTTP/1.1"

	httpHeaderUserAgent       = `User-Agent`
	httpHeaderContentType     = `Content-Type`
	httpHeaderContentEncoding = `Content-Encoding`
	httpHeaderContentTypeJson = `application/json`
	httpHeaderContentTypeXml  = `application/xml`
	httpHeaderContentTypeForm = `application/x-www-form-urlencoded`
)

type (
	Method   string
	Header   map[string]string
	Query    map[string]string
	PostForm map[string]string
	Json     map[string]interface{}
)

func (m Method) String() string {
	return string(m)
}

type Request struct {
	req http.Request

	queries  url.Values
	postForm url.Values
	headers  url.Values

	body []byte
}

func NewRequest(req *http.Request) *Request {
	if req == nil {
		return nil
	}
	if len(req.UserAgent()) == 0 {
		req.Header.Set(httpHeaderUserAgent, defaultClientAgent)
	}
	return &Request{
		req:      *req,
		queries:  url.Values{},
		postForm: url.Values{},
		headers:  url.Values{},
		body:     nil,
	}
}

func (r *Request) ContentType(contentType string) {
	r.Header(httpHeaderContentType, contentType)
}

func (r *Request) ContentTypePostForm() {
	r.Header(httpHeaderContentType, httpHeaderContentTypeForm)
}

func (r *Request) ContentTypeJson() {
	r.Header(httpHeaderContentType, httpHeaderContentTypeJson)
}

func (r *Request) ContentTypeXml() {
	r.Header(httpHeaderContentType, httpHeaderContentTypeXml)
}

func (r *Request) Host(h string) {
	r.req.Host = h
	r.req.URL.Host = h
}

func (r *Request) URL() *url.URL {
	return r.req.URL
}

func (r *Request) Url(rawUrl string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	r.req.Host = u.Host
	r.req.URL = u
	return nil
}

func (r *Request) Method(method Method) {
	r.req.Method = method.String()
}

func (r *Request) Get(rawUrl string) error {
	r.Method(GET)
	return r.Url(rawUrl)
}

func (r *Request) Post(rawUrl string) error {
	r.Method(POST)
	return r.Url(rawUrl)
}

func (r *Request) Put(rawUrl string) error {
	r.Method(PUT)
	return r.Url(rawUrl)
}

func (r *Request) Delete(rawUrl string) error {
	r.Method(DELETE)
	return r.Url(rawUrl)
}

func (r *Request) Head(rawUrl string) error {
	r.Method(HEAD)
	return r.Url(rawUrl)
}

func (r *Request) Patch(rawUrl string) error {
	r.Method(PATCH)
	return r.Url(rawUrl)
}

func (r *Request) Trace(rawUrl string) error {
	r.Method(TRACE)
	return r.Url(rawUrl)
}

func (r *Request) Options(rawUrl string) error {
	r.Method(OPTIONS)
	return r.Url(rawUrl)
}

func (r *Request) UserAgent(ua string) {
	r.headers.Set(httpHeaderUserAgent, ua)
}

func (r *Request) Header(key, value string) {
	r.headers.Set(key, value)
}

func (r *Request) HeaderMap(h Header) {
	for k, v := range h {
		r.Header(k, v)
	}
}

func (r *Request) AddHeader(key, value string) {
	r.headers.Add(key, value)
}

func (r *Request) AddCookie(c *http.Cookie) {
	r.req.AddCookie(c)
}

func (r *Request) AddCookies(c []*http.Cookie) {
	for _, cookie := range c {
		r.AddCookie(cookie)
	}
}

func (r *Request) ProtocolVersion(proto string, protoMajor, protoMinor int) {
	if proto == "" {
		r.req.Proto = protocolVersion
		r.req.ProtoMajor = 1
		r.req.ProtoMinor = 0
	} else {
		r.req.Proto = proto
		r.req.ProtoMajor = protoMajor
		r.req.ProtoMinor = protoMinor
	}
}

func (r *Request) Query(key, value string) {
	if r.queries == nil {
		r.queries = url.Values{}
	}
	r.queries.Set(key, value)
}

func (r *Request) QueryMap(m Query) {
	for k, v := range m {
		r.Query(k, v)
	}
}

func (r *Request) AddQuery(key, value string) {
	if r.queries == nil {
		r.queries = url.Values{}
	}
	r.queries.Add(key, value)
}

func (r *Request) AddQueryMap(m Query) {
	for k, v := range m {
		r.AddQuery(k, v)
	}
}

func (r *Request) PostForm(key, value string) {
	if r.postForm == nil {
		r.postForm = url.Values{}
	}
	r.postForm.Set(key, value)
}

func (r *Request) PostFormMap(pf PostForm) {
	for k, v := range pf {
		r.PostForm(k, v)
	}
}

func (r *Request) AddPostForm(key, value string) {
	if r.postForm == nil {
		r.postForm = url.Values{}
	}
	r.postForm.Add(key, value)
}

func (r *Request) AddPostFormMap(pf PostForm) {
	for k, v := range pf {
		r.AddPostForm(k, v)
	}
}

func (r *Request) BodyJSON(obj interface{}) error {
	if r.body == nil && obj != nil {
		bs, err := jsonMarshal(obj)
		if err != nil {
			return err
		}
		r.BodyJSON4Bytes(bs)
	}
	return nil
}

func (r *Request) GetBody() []byte {
	return r.body
}

func (r *Request) BodyJSON4Str(json string) {
	r.BodyJSON4Bytes([]byte(json))
}

func (r *Request) BodyJSON4Bytes(json []byte) {
	if r.body == nil && len(json) > 0 {
		r.body = json
		r.Header(httpHeaderContentType, httpHeaderContentTypeJson)
	}
}

func (r *Request) Body(content []byte) {
	if r.body == nil {
		r.body = content
	}
}

func (r *Request) RequestInfo() string {
	cookieStr := strings.Builder{}
	cookieStr.WriteString("CookieInfo:\n")
	for i, cookie := range r.req.Cookies() {
		if i != 0 {
			cookieStr.WriteString("; ")
		}
		cookieStr.WriteString(cookie.String())
	}

	headerStr := strings.Builder{}
	headerStr.WriteString("headerInfo:\n")
	for k, v := range r.headers {
		headerStr.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, " ")))
	}

	bodyStr := strings.Builder{}
	bodyStr.WriteString("BodyInfo:\n")
	if len(r.body) > 0 {
		bodyStr.Write(r.body)
	}
	for k, v := range r.postForm {
		bodyStr.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, " ")))
	}

	queriesStr := strings.Builder{}
	queriesStr.WriteString("queriesInfo:\n")
	for k, v := range r.queries {
		queriesStr.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, " ")))
	}

	return fmt.Sprintf("\nURL:%s\n\n%s\n%s\n%s\n\n%s", r.req.URL.String(), queriesStr.String(), headerStr.String(), cookieStr.String(), bodyStr.String())
}

func jsonMarshal(obj interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	err := jsonEncoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (r *Request) buildFormBody() {
	// build POST/PUT/PATCH/DELETE url and body
	if !(r.req.Method == POST.String() || r.req.Method == PUT.String() ||
		r.req.Method == PATCH.String() || r.req.Method == DELETE.String()) || r.body != nil {
		return
	}
	// with params
	if len(r.postForm) == 0 {
		return
	}
	postForm := r.postForm.Encode()
	if len(postForm) > 0 {
		buf := []byte(postForm)
		r.body = buf
		r.ContentTypePostForm()
	}
}

func (r *Request) buildRequest() (*http.Request, error) {
	r.buildFormBody()
	if r.body != nil && len(r.body) > 0 {
		r.req.ContentLength = int64(len(r.body))
		r.req.Body = io.NopCloser(bytes.NewReader(r.body))
	}

	if r.queries != nil && len(r.queries) > 0 {
		uqs := r.req.URL.Query()
		for k, vs := range r.queries {
			for _, v := range vs {
				uqs.Add(k, v)
			}
		}
		r.req.URL.RawQuery = uqs.Encode()
	}

	if r.headers != nil && len(r.headers) > 0 {
		if r.req.Header == nil {
			r.req.Header = http.Header{}
		}
		for k, v := range r.headers {
			r.req.Header[k] = v
		}
	}

	return &r.req, nil
}
