package shttp_test

import (
	"net/url"
	"testing"

	"github.com/smalls0098/pkg/shttp"
)

func Test_Client_Get(t *testing.T) {
	response, err := shttp.Get("http://baidu.com", func(c *shttp.Client, req *shttp.Request) {
		req.Header("test", "1")
		req.Query("test", "1")

		p, _ := url.Parse("http://127.0.0.1:9900")
		c.ProxyUrl(p)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

func Test_Client_Post(t *testing.T) {
	response, err := shttp.Post("http://baidu.com", func(c *shttp.Client, req *shttp.Request) {
		req.Header("test", "1")
		req.Query("test", "1")
		req.PostForm("test", "2")
		p, _ := url.Parse("http://127.0.0.1:9900")
		c.ProxyUrl(p)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.String())
}
