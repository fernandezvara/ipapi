[![Build status](https://img.shields.io/github/workflow/status/fernandezvara/ipapi/Test?style=flat-square)](https://github.com/fernandezvara/ipapi/actions?workflow=Test)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square&logo=go)](https://pkg.go.dev/github.com/fernandezvara/ipapi)
[![GoreportCard](https://goreportcard.com/badge/github.com/fernandezvara/ipapi?style=flat-square)](https://goreportcard.com/report/github.com/fernandezvara/ipapi)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

# ip-api.com Go Client

This is a client for **ip-api.com** API service that uses its JSON endpoints. It works for the free and professional API, requiring an API key to identify the customer.

**IMPORTANT NOTE:**

_If you are using a free account there is a request limit per minute, as the time of this writting is 45 req/min. Currently the library do no manage this limits and you can have the service banned for your IP._

## Install package

```bash
    go get github.com/fernandezvara/ipapi
```

## Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fernandezvara/ipapi"
)

func main() {

	// create a client for the API. Adding an API key will use the Pro API
	client := ipapi.New("")

	// queries use contexts to manage timeouts, so just pass your current context if
	// you have one, or set a new one for the query
	ctx := context.Background()

	// response holds the information from the API (for the fields requested)
	response, err := client.Query(ctx, "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Query)   // IP requested
	fmt.Println(response.Status)  // success or fail
	fmt.Println(response.Message) // Message is filled when status of the request is fail,
	                              // in this case it states that IP is from a reserved range

	// require a different set of fields
	client.SetFields([]string{"status", "message", "query", "country", "regionName", "city", "zip", "lat", "lon"}, false)

	// set a different timeout
	client.SetTimeout(10 * time.Second)

	// querying without IP address it will return the IP that originates the request
	response, err = client.Query(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %#v\n", response)

}
```

# License

[The MIT License (MIT)](LICENSE)
