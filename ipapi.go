package ipapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiFree = "ip-api.com"
	apiPro  = "pro.ip-api.com"

	defaultFields  = "status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query"
	defaultTimeout = 5 * time.Second
)

var allFields = map[string]int{
	"country":       1,
	"countryCode":   2,
	"region":        4,
	"regionName":    8,
	"city":          16,
	"zip":           32,
	"lat":           64,
	"lon":           128,
	"timezone":      256,
	"isp":           512,
	"org":           1024,
	"as":            2048,
	"reverse":       4096,
	"query":         8192,
	"status":        16384,
	"message":       32768,
	"mobile":        65536,
	"proxy":         131072, // 262144 missing
	"district":      524288,
	"continent":     1048576,
	"continentCode": 2097152,
	"asname":        4194304,
	"currency":      8388608,
	"hosting":       16777216,
	"offset":        33554432,
}

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

}

func (c *Client) fieldAllowed(fieldName string) (ok bool) {

	_, ok = allFields[fieldName]
	return

}

// SetFields allows to change what fields to query to the API,
//   if a field is not allowed by the API won't be added
func (c *Client) SetFields(fields []string, numeric bool) {

	var finalFields []string

	for _, field := range fields {
		if c.fieldAllowed(field) && !exists(finalFields, field) {
			finalFields = append(finalFields, field)
		}
	}

	if numeric {
		var total int
		for _, k := range finalFields {
			total = total + allFields[k]
		}
		c.fields = strconv.Itoa(total)
		return
	}

	c.fields = strings.Join(finalFields, ",")

}

func exists(items []string, item string) bool {

	for _, i := range items {
		if item == i {
			return true
		}
	}

	return false

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
