package local

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	searchLocalKeywordURL = "https://dapi.kakao.com/v2/local/search/keyword.json"
	appkey                = "61891b84ee488dfa867a03aa50399af7"
)

type Client struct {
	URL    *url.URL
	client *http.Client
}

func NewClient() *Client {
	u, err := url.Parse(searchLocalKeywordURL)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}
	return &Client{
		URL:    u,
		client: client,
	}
}

type Request struct {
	Query  string  `query:"query"`
	X      float64 `query:"x"`
	Y      float64 `query:"y"`
	Radius int     `query:"radius"`
}

func (c *Client) SearchKeyword(r *Request) (*http.Response, error) {
	v := url.Values{}
	v.Set("query", r.Query)
	v.Set("x", fmt.Sprint(r.X))
	v.Set("y", fmt.Sprint(r.Y))
	v.Set("radius", fmt.Sprint(r.Radius))
	body := strings.NewReader(v.Encode())

	httpReq, err := http.NewRequest("POST", c.URL.String(), body)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Authorization", "KakaoAK "+appkey)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != http.StatusOK {
		err := fmt.Errorf("[%d] %s ", httpResp.StatusCode, httpResp.Status)
		return nil, err
	}

	return httpResp, nil
}
