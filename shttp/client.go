package shttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type (
	Middleware     func(c *Client, req *Request, resp *Response) error
	RequestHandler func(c *Client, req *Request)
	G              map[string]interface{}
	Option         func(*Client)
)

// VERSION sHttp version
const VERSION = "0.1.0"

var (
	defaultClientAgent = fmt.Sprintf(`sHttp %s`, VERSION)

	connectTimeout   = 5 * time.Second
	readWriteTimeout = 5 * time.Second
)

type Client struct {
	c *http.Client

	middlewares []Middleware
}

var DefaultClient = &Client{
	c:           http.DefaultClient,
	middlewares: make([]Middleware, 0),
}

func New(opts ...Option) *Client {
	options := &Client{
		c: &http.Client{
			Transport:     defaultTransport(),
			CheckRedirect: defaultCheckRedirect(),
			Jar:           nil,
			Timeout:       connectTimeout,
		},
		middlewares: make([]Middleware, 0),
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithMiddleware(middleware Middleware) Option {
	return func(opts *Client) {
		opts.Use(middleware)
	}
}

func (c *Client) Use(middleware Middleware) *Client {
	if c.middlewares == nil {
		c.middlewares = make([]Middleware, 0)
	}
	c.middlewares = append(c.middlewares, middleware)
	return c
}

func (c *Client) ProxyUrl(proxy *url.URL) {
	proxyFunc := func(r *http.Request) (*url.URL, error) {
		return proxy, nil
	}
	if t, ok := c.c.Transport.(*http.Transport); ok {
		t.Proxy = proxyFunc
	} else {
		c.c.Transport = &http.Transport{
			Proxy: proxyFunc,
		}
	}
}

func (c *Client) Transport(t *http.Transport) {
	c.c.Transport = t
}

func (c *Client) TLSClientConfig(tls *tls.Config) {
	if t, ok := c.c.Transport.(*http.Transport); ok {
		t.TLSClientConfig = tls
	} else {
		c.c.Transport = &http.Transport{
			TLSClientConfig: tls,
		}
	}
}

func (c *Client) Timeout(connectTimeout time.Duration, readWriteTimeout time.Duration) {
	dialCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, err
	}
	if t, ok := c.c.Transport.(*http.Transport); ok {
		t.DialContext = dialCtx
	} else {
		c.c.Transport = &http.Transport{
			DialContext: dialCtx,
		}
	}
}

func (c *Client) Request(url string, method Method, body io.Reader, handlers ...RequestHandler) (*Response, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	httpReq, err := http.NewRequest(method.String(), url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(c.handlerRequest(httpReq, handlers...))
}

func (c *Client) Get(url string, handlers ...RequestHandler) (*Response, error) {
	return c.Request(url, GET, nil, handlers...)
}

func (c *Client) Get4Bytes(url string, handlers ...RequestHandler) ([]byte, error) {
	resp, err := c.Get(url, handlers...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes()
}

func (c *Client) Get4String(url string, handlers ...RequestHandler) (string, error) {
	resp, err := c.Get(url, handlers...)
	if err != nil {
		return "", err
	}
	return resp.String()
}

func (c *Client) Post(url string, handlers ...RequestHandler) (*Response, error) {
	return c.Request(url, POST, nil, handlers...)
}

func (c *Client) Post4Bytes(url string, handlers ...RequestHandler) ([]byte, error) {
	resp, err := c.Post(url, handlers...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes()
}

func (c *Client) Post4String(url string, handlers ...RequestHandler) (string, error) {
	resp, err := c.Post(url, handlers...)
	if err != nil {
		return "", err
	}
	return resp.String()
}

func (c *Client) handlerRequest(httpReq *http.Request, handlers ...RequestHandler) *Request {
	req := NewRequest(httpReq)
	if handlers != nil && len(handlers) > 0 {
		for _, handler := range handlers {
			handler(c, req)
			if c == nil || req == nil {
				panic("client or request is nil")
			}
		}
	}
	return req
}

func (c *Client) Do(req *Request) (*Response, error) {
	var err error
	if len(c.middlewares) > 0 {
		for _, m := range c.middlewares {
			err = m(c, req, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	httpReq, err := req.buildRequest()
	if err != nil {
		return nil, err
	}
	httpResp, err := c.c.Do(httpReq)
	if err != nil {
		return nil, err
	}
	resp := NewResponse(httpResp)
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	if len(c.middlewares) > 0 {
		for _, m := range c.middlewares {
			err = m(c, req, resp)
			if err != nil {
				return nil, err
			}
		}
	}
	return resp, nil
}

func defaultTransport() *http.Transport {
	return &http.Transport{
		// No validation for https certification of the server in default.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, connectTimeout)
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(readWriteTimeout))
			return conn, err
		},
		MaxIdleConnsPerHost: 100,
		DisableKeepAlives:   true,
	}
}

func defaultCheckRedirect() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("stopped after 5 redirects")
		}
		return nil
	}
}

func Get(url string, handlers ...RequestHandler) (*Response, error) {
	return DefaultClient.Get(url, handlers...)
}

func Post(url string, handlers ...RequestHandler) (*Response, error) {
	return DefaultClient.Post(url, handlers...)
}
