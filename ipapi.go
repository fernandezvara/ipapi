package ipapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiFree = "ip-api.com"
	apiPro  = "pro.ip-api.com"

	defaultFields  = "status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query"
	defaultTimeout = 5 * time.Second
)

type Client struct {
	apiKey     string
	fields     string
	timeout    time.Duration
	httpClient *http.Client
}

type Response struct {
	Query         string  `json:"query,omitempty"`   // IP used for the query
	Status        string  `json:"status,omitempty"`  // success / fail
	Message       string  `json:"message,omitempty"` // Message only appear is status is fail, returning the failure
	Continent     string  `json:"continent,omitempty"`
	ContinentCode string  `json:"continentCode,omitempty"`
	Country       string  `json:"country,omitempty"`
	CountryCode   string  `json:"countryCode,omitempty"`
	Region        string  `json:"region,omitempty"`
	RegionName    string  `json:"regionName,omitempty"`
	City          string  `json:"city,omitempty"`
	District      string  `json:"district,omitempty"`
	Zip           string  `json:"zip,omitempty"`
	Lat           float64 `json:"lat,omitempty"`
	Lon           float64 `json:"lon,omitempty"`
	Timezone      string  `json:"timezone,omitempty"`
	Offset        int     `json:"offset,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	Isp           string  `json:"isp,omitempty"`
	Org           string  `json:"org,omitempty"`
	As            string  `json:"as,omitempty"`
	Asname        string  `json:"asname,omitempty"`
	Reverse       string  `json:"reverse,omitempty"`
	Mobile        bool    `json:"mobile,omitempty"`
	Proxy         bool    `json:"proxy,omitempty"`
	Hosting       bool    `json:"hosting,omitempty"`
}

func New(apiKey string) *Client {

	var (
		client Client = Client{}
	)

	client.apiKey = apiKey
	client.fields = defaultFields
	client.httpClient = http.DefaultClient
	client.timeout = defaultTimeout
	return &client

}

func (c *Client) buildURL(ip string) string {

	var (
		uri    url.URL
		values url.Values = uri.Query()
	)

	values.Add("fields", c.fields)
	uri.Path = fmt.Sprintf("/json/%s", ip)

	if c.apiKey != "" {
		uri.Scheme, uri.Host = "https", apiPro
		values.Add("apiKey", c.apiKey)
	} else {
		uri.Scheme, uri.Host = "http", apiFree
	}

	uri.RawQuery = values.Encode()
	return uri.String()

	// if c.pro {
	// 	if ip != "" {
	// 		return fmt.Sprintf("%s/%s?fields=%s?key=%s", apiPro, ip, c.fields, c.apiKey)
	// 	}
	// 	return fmt.Sprintf("%s?fields=%s?key=%s", apiPro, c.fields, c.apiKey)
	// }
	// return fmt.Sprintf("%s/%s?fields=%s", apiFree, ip, c.fields)
}

func (c *Client) fieldAllowed(fieldName string) bool {

	var allowed = []string{"status", "message", "continent", "continentCode", "country", "countryCode", "region", "regionName", "city", "district", "zip", "lat", "lon", "timezone", "offset", "currency", "isp", "org", "as", "asname", "reverse", "mobile", "proxy", "hosting", "query"}

	for _, value := range allowed {
		if fieldName == value {
			return true
		}
	}
	return false
}

// SetFields allows to change what fields to query to the API,
//   if a field is not allowed by the API won't be added
func (c *Client) SetFields(fields []string) {

	var finalFields []string

	for _, field := range fields {
		if c.fieldAllowed(field) {
			finalFields = append(finalFields, field)
		}
	}
	c.fields = strings.Join(finalFields, ",")

}

// SetTimeout sets a new timeout for API requests
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Client) Query(ctx context.Context, ip string) (response Response, err error) {

	var (
		uri          string = c.buildURL(ip)
		cancel       context.CancelFunc
		httpResponse *http.Response
		httpRequest  *http.Request
	)

	ctx, cancel = context.WithTimeout(ctx, c.timeout)
	defer cancel()

	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return
	}

	httpResponse, err = c.httpClient.Do(httpRequest)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	return

}
