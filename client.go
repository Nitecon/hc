package hc

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	isReadBody bool
	httpClient *http.Client
	resp       *http.Response
	body       []byte
}

func New() *Client {
	client := &http.Client{}
	return &Client{
		httpClient: client,
	}
}

func (c *Client) Get(url string) error {
	c.isReadBody = false
	c.body = nil
	resp, err := c.httpClient.Get(url)
	c.resp = resp
	return err
}

func (c *Client) Post(url string, data []byte) error {
	c.isReadBody = false
	c.body = nil
	resp, err := c.doRequest("POST", url, data)
	c.resp = resp
	return err
}

func (c *Client) Put(url string, data []byte) error {
	c.isReadBody = false
	c.body = nil
	resp, err := c.doRequest("PUT", url, data)
	c.resp = resp
	return err
}

func (c *Client) Delete(url string) error {
	c.isReadBody = false
	c.body = nil
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

func (c *Client) Header() http.Header {
	if c.resp != nil {
		return c.resp.Header
	}
	return nil
}

func (c *Client) ReadBody() []byte {
	if c.body == nil && !c.isReadBody {
		body, err := ReadResponseBody(c.resp)
		c.isReadBody = true
		if err != nil {
			return nil
		}
		c.body = body
	}
	return c.body
}

func (c *Client) ReadString() string {
	if c.body == nil && !c.isReadBody {
		body, err := ReadResponseBody(c.resp)
		c.isReadBody = true
		if err != nil {
			return ""
		}
		return string(body)
	}
	return string(c.body)
}

func (c *Client) ReadJson(v interface{}) error {
	if c.body == nil && !c.isReadBody {
		body, err := ReadResponseBody(c.resp)
		c.isReadBody = true
		if err != nil {
			return err
		}
		c.body = body
	}
	if c.body == nil {
		return fmt.Errorf("body is nil")
	}
	return json.Unmarshal(c.body, v)
}

func (c *Client) doRequest(method, url string, data []byte) (*http.Response, error) {
	// If the data is not empty, compress it.
	if len(data) > 0 {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(data); err != nil {
			return nil, err
		}
		if err := gz.Close(); err != nil {
			return nil, err
		}

		data = buf.Bytes()
	}

	// Create a bytes.Reader from the compressed data.
	reader := bytes.NewReader(data)

	// Set the Content-Encoding header to gzip.
	req, err := http.NewRequest(method, url, reader)
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
