# Go `requests` Package

A simple and intuitive HTTP client for Go, inspired by Python's popular `requests` module. This package aims to provide an easy-to-use interface for making HTTP requests in Go.

## Installation

To install the `requests` package, use the standard `go get`:

```bash
go get github.com/blakepatteson/requests
```

## Usage

Making a GET request

```go
package main

import (
	"fmt"
	"log"
	"github.com/blakepatteson/requests"
)

func main() {
	resp, err := requests.HttpRequest{
		VerbHTTP:    "GET",
		Endpoint:    "https://api.example.com/data",
	}.Do()
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	body, err := requests.ParseJson(resp)
    if err != nil{
        log.Fatalf("Error : %v", err)
    }
	fmt.Printf("Response : %+v\n", body)
}
```

## Making a POST request with JSON

```go
out, err := requests.HttpRequest{
	VerbHTTP:    "POST",
	Endpoint:    "https://api.example.com/data",
	JSON:        []byte(`{"key": "value"}`),
	ContentType: "application/json",
}.Do()
fmt.Printf("out : '%v'", out)
```

## Setting Auth

For bearer token :

```
args.Auth = "Bearer YOUR_TOKEN_HERE"
```

For basic auth :

```
args.Auth = "username:password"
```

## Testing

The package comes with a set of basic tests to ensure functionality. To run the tests, navigate to the package directory and execute:

```
go test .
```

## Contributing

Feel free to contribute to this project by opening issues or submitting pull requests.

## License

This project is licensed under the [MIT license](/LICENSE).
