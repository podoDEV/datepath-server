package direction

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	directionURL = "https://maps.googleapis.com/maps/api/directions/json"
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
	u, err := url.Parse(directionURL)
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
	Origin        string `query:"origin"`
	Destination   string `query:"destination"`
	DepartureTime string `query:"departure_time"`
}

type Response struct {
	GeocodedWaypoints []struct {
		GeocoderStatus string   `json:"geocoder_status"`
		PlaceID        string   `json:"place_id"`
		Types          []string `json:"types"`
	} `json:"geocoded_waypoints"`
	Routes []struct {
		Bounds struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"bounds"`
		Copyrights string `json:"copyrights"`
		Legs       []struct {
			ArrivalTime struct {
				Text     string `json:"text"`
				TimeZone string `json:"time_zone"`
				Value    int    `json:"value"`
			} `json:"arrival_time"`
			DepartureTime struct {
				Text     string `json:"text"`
				TimeZone string `json:"time_zone"`
				Value    int    `json:"value"`
			} `json:"departure_time"`
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			EndAddress  string `json:"end_address"`
			EndLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"end_location"`
			StartAddress  string `json:"start_address"`
			StartLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"start_location"`
			Steps []struct {
				Distance struct {
					Text  string `json:"text"`
					Value int    `json:"value"`
				} `json:"distance"`
				Duration struct {
					Text  string `json:"text"`
					Value int    `json:"value"`
				} `json:"duration"`
				EndLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"end_location"`
				HTMLInstructions string `json:"html_instructions"`
				Polyline         struct {
					Points string `json:"points"`
				} `json:"polyline"`
				StartLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"start_location"`
				Steps []struct {
					Distance struct {
						Text  string `json:"text"`
						Value int    `json:"value"`
					} `json:"distance"`
					Duration struct {
						Text  string `json:"text"`
						Value int    `json:"value"`
					} `json:"duration"`
					EndLocation struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"end_location"`
					Polyline struct {
						Points string `json:"points"`
					} `json:"polyline"`
					StartLocation struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"start_location"`
					TravelMode string `json:"travel_mode"`
				} `json:"steps,omitempty"`
				TravelMode     string `json:"travel_mode"`
				TransitDetails struct {
					ArrivalStop struct {
						Location struct {
							Lat float64 `json:"lat"`
							Lng float64 `json:"lng"`
						} `json:"location"`
						Name string `json:"name"`
					} `json:"arrival_stop"`
					ArrivalTime struct {
						Text     string `json:"text"`
						TimeZone string `json:"time_zone"`
						Value    int    `json:"value"`
					} `json:"arrival_time"`
					DepartureStop struct {
						Location struct {
							Lat float64 `json:"lat"`
							Lng float64 `json:"lng"`
						} `json:"location"`
						Name string `json:"name"`
					} `json:"departure_stop"`
					DepartureTime struct {
						Text     string `json:"text"`
						TimeZone string `json:"time_zone"`
						Value    int    `json:"value"`
					} `json:"departure_time"`
					Headsign string `json:"headsign"`
					Headway  int    `json:"headway"`
					Line     struct {
						Agencies []struct {
							Name string `json:"name"`
							URL  string `json:"url"`
						} `json:"agencies"`
						Color     string `json:"color"`
						Name      string `json:"name"`
						ShortName string `json:"short_name"`
						TextColor string `json:"text_color"`
						Vehicle   struct {
							Icon string `json:"icon"`
							Name string `json:"name"`
							Type string `json:"type"`
						} `json:"vehicle"`
					} `json:"line"`
					NumStops int `json:"num_stops"`
				} `json:"transit_details,omitempty"`
			} `json:"steps"`
			TrafficSpeedEntry []interface{} `json:"traffic_speed_entry"`
			ViaWaypoint       []interface{} `json:"via_waypoint"`
		} `json:"legs"`
		OverviewPolyline struct {
			Points string `json:"points"`
		} `json:"overview_polyline"`
		Summary       string        `json:"summary"`
		Warnings      []string      `json:"warnings"`
		WaypointOrder []interface{} `json:"waypoint_order"`
	} `json:"routes"`
	Status string `json:"status"`
}

func (c *Client) Direction(r *Request) (*http.Response, error) {
	q := url.Values{}
	q.Set("origin", r.Origin)
	q.Set("destination", r.Destination)
	q.Set("departure_time", r.DepartureTime)
	q.Set("mode", "transit")
	q.Set("language", "ko")
	q.Set("key", c.appKey)

	c.URL.RawQuery, _ = url.QueryUnescape(q.Encode())
	directionURL := c.URL.String()

	httpReq, err := http.NewRequest("GET", directionURL, nil)
	if err != nil {
		return nil, err
	}

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
