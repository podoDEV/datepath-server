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
)

type Config struct {
	AppKey string `toml:"app_key"`
}

type Client struct {
	URL    *url.URL
	client *http.Client
	appKey string
}

func NewClient(conf *Config) *Client {
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
		appKey: conf.AppKey,
	}
}

type Request struct {
	Query  string  `query:"query"`
	X      float64 `query:"x"`
	Y      float64 `query:"y"`
	Radius int     `query:"radius"`
}

type Response struct {
	Documents []Document `json:"documents"`
	Meta      struct {
		IsEnd         bool `json:"is_end"`
		PageableCount int  `json:"pageable_count"`
		SameName      struct {
			Keyword        string   `json:"keyword"`
			Region         []string `json:"region"`
			SelectedRegion string   `json:"selected_region"`
		} `json:"same_name"`
		TotalCount int `json:"total_count"`
	} `json:"meta"`
}

type Document struct {
	AddressName       string `json:"address_name"`
	CategoryGroupCode string `json:"category_group_code"`
	CategoryGroupName string `json:"category_group_name"`
	CategoryName      string `json:"category_name"`
	Distance          string `json:"distance"`
	Id                string `json:"id"`
	Phone             string `json:"phone"`
	PlaceName         string `json:"place_name"`
	PlaceURL          string `json:"place_url"`
	RoadAddressName   string `json:"road_address_name"`
	X                 string `json:"x"`
	Y                 string `json:"y"`
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
	httpReq.Header.Set("Authorization", "KakaoAK "+c.appKey)

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
