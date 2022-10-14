package shttp

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	resp http.Response

	body []byte
}

func NewResponse(resp *http.Response) *Response {
	if resp == nil {
		return nil
	}
	return &Response{
		resp: *resp,
	}
}

func (r *Response) Response() *http.Response {
	return &r.resp
}

func (r *Response) Ok() bool {
	return r.resp.StatusCode == http.StatusOK
}

func (r *Response) Cookie() []*http.Cookie {
	return r.resp.Cookies()
}

func (r *Response) Bytes() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	bs, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	defer r.resp.Body.Close()

	// handle gzip
	contentEncode := r.resp.Header.Get(httpHeaderContentEncoding)
	if strings.Contains(contentEncode, "gzip") {
		reader, err := gzip.NewReader(bytes.NewReader(bs))
		if err != nil {
			return nil, err
		}
		bs, err = io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	}

	r.body = bs
	return bs, nil
}

func (r *Response) String() (string, error) {
	if bs, err := r.Bytes(); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}
