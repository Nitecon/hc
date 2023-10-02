package hc

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"golang.org/x/net/http2"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	resp       *http.Response
	body       []byte
}

func New() *Client {
	transport := &http2.Transport{
		AllowHTTP: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	return &Client{
		httpClient: client,
	}
}

func (c *Client) Get(url string) error {
	resp, err := c.httpClient.Get(url)
	c.resp = resp
	return err
}

func (c *Client) Post(url string, data []byte) error {
	resp, err := c.doRequest("POST", url, data)
	c.resp = resp
	return err
}

func (c *Client) Put(url string, data []byte) error {
	resp, err := c.doRequest("PUT", url, data)
	c.resp = resp
	return err
}

func (c *Client) Delete(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	c.resp = resp
	return err
}

func (c *Client) PostJson(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Post(url, jsonData)
}

func (c *Client) PutJson(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Put(url, jsonData)
}

func (c *Client) Status() int {
	if c.resp != nil {
		return c.resp.StatusCode
	}
	return 0
}

func (c *Client) StatusText() string {
	if c.resp != nil {
		return c.resp.Status
	}
	return ""
}

func (c *Client) ReadJson(v interface{}) error {
	if c.body == nil {
		body, err := ReadResponseBody(c.resp)
		if err != nil {
			return err
		}
		c.body = body
	}
	return json.Unmarshal(c.body, v)
}

func (c *Client) doRequest(method, url string, data []byte) (*http.Response, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	return c.httpClient.Do(req)
}

func ReadResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		return DecompressGzip(resp)
	}
	return io.ReadAll(resp.Body)
}

func DecompressGzip(resp *http.Response) ([]byte, error) {
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
