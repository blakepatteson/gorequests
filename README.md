# Go `requests` Package

A simple and intuitive HTTP client for Go, inspired by Python's popular [`requests`](https://pypi.org/project/requests/) module.

This package aims to provide an easy-to-use interface for making HTTP requests in Go, abstracting away much of the boilerplate of the standard library.

## Installation

To install the `requests` package, use the standard `go get`:

```bash
go get github.com/blakepatteson/gorequests
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
		log.Fatalf("err with 'GET' request : %v", err)
	}

	body, err := requests.ParseJson(resp)
	if err != nil {
		fmt.Printf("err parsing JSON : %v", err)
		// do something with the error
    }
	fmt.Printf("Response Body : %+v\n", body)
}
```

## Making a POST request with JSON

```go
func main() {
	resp, err := requests.HttpRequest{
		VerbHTTP: "POST",
		Endpoint: "https://api.example.com/data",
		JSON:     []byte(`{"key": "value","anotherKey":"anotherValue"}`),
	}.Do()
	if err != nil {
		fmt.Printf("err with post request : '%v'", err)
	}
	fmt.Printf("resp : '%v'\n", resp)
	// if you just want the json, call the 'Fatal' variant
	// will crash if the parseJson err != nil (for less boilerplate)
	fmt.Printf("Response Body : '%+v\n", requests.ParseJsonFatal(resp))
}

```

## Setting Auth

For bearer token :

```
req.Auth = "Bearer YOUR_TOKEN_HERE"
```

For basic auth :

```
req.Auth = "username:password"
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
