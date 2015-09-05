package client

import (
	"net/http"
	"net/http/cookiejar"
)

// Client is NicoClient
type Client struct {
	cl *http.Client
}

// NewClient creates new Client
func NewClient() *Client {
	cl := &Client{}
	cl.cl = &http.Client{}
	cl.cl.Jar, _ = cookiejar.New(nil)
	return cl
}
