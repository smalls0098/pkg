package shttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type (
	Handler func(c *Client, req *Request)
	G       map[string]interface{}
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

	handlers []Handler
}

// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{
	c: http.DefaultClient,
}

func New(handlers ...Handler) *Client {
	return &Client{
		c: &http.Client{
			Transport:     defaultTransport(),
			CheckRedirect: defaultCheckRedirect(),
			Jar:           nil,
			Timeout:       connectTimeout,
		},
		handlers: handlers,
	}
}

func Get(url string, handlers ...Handler) (*Response, error) {
	return DefaultClient.Get(url, handlers...)
}

func Post(url string, handlers ...Handler) (*Response, error) {
	return DefaultClient.Post(url, handlers...)
}

func (c *Client) Get(url string, handlers ...Handler) (*Response, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(c.handlerRequest(httpReq, handlers...))
}

func (c *Client) Get4Bytes(url string, handlers ...Handler) ([]byte, error) {
	resp, err := c.Get(url, handlers...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes()
}

func (c *Client) Get4String(url string, handlers ...Handler) (string, error) {
	resp, err := c.Get(url, handlers...)
	if err != nil {
		return "", err
	}
	return resp.String()
}

func (c *Client) Post(url string, handlers ...Handler) (*Response, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}
	httpReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(c.handlerRequest(httpReq, handlers...))
}

func (c *Client) Post4Bytes(url string, handlers ...Handler) ([]byte, error) {
	resp, err := c.Post(url, handlers...)
	if err != nil {
		return nil, err
	}
	return resp.Bytes()
}

func (c *Client) Post4String(url string, handlers ...Handler) (string, error) {
	resp, err := c.Post(url, handlers...)
	if err != nil {
		return "", err
	}
	return resp.String()
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

func (c *Client) Handler(handler Handler) {
	if c.handlers == nil {
		c.handlers = make([]Handler, 0)
	}
	c.handlers = append(c.handlers, handler)
}

func (c *Client) handlerRequest(httpReq *http.Request, handlers ...Handler) *Request {
	req := NewRequest(httpReq)
	if c.handlers != nil && len(c.handlers) > 0 {
		for _, handler := range c.handlers {
			handler(c, req)
			if c == nil || req == nil {
				panic("client or request is nil")
			}
		}
	}
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
