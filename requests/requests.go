package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HttpRequest struct {
	Client      *http.Client
	Endpoint    string
	Auth        string
	JSON        []byte
	VerbHTTP    string
	Form        url.Values
	ContentType string
}

type HttpError struct {
	StatusCode int
	Message    string
	ReqBody    string
}

var (
	sharedClient *http.Client
	once         sync.Once
)

func getSharedClient() *http.Client {
	once.Do(func() {
		sharedClient = &http.Client{}
	})
	return sharedClient
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v - %v", e.Message, e.ReqBody)
}

func (args HttpRequest) Do() (*http.Response, error) {
	request, err := args.createRequest()
	if err != nil {
		return nil, err
	}

	response, err := args.executeRequest(request)
	if err != nil {
		return nil, err
	}
	return checkStatusCode(response)
}

func parseJson(resp *http.Response) (map[string]any, error) {
	defer resp.Body.Close()
	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response JSON: %w", err)
	}
	return result, nil
}

func ParseJson(resp *http.Response) (map[string]any, error) {
	return parseJson(resp)
}

func ParseJsonFatal(resp *http.Response) map[string]any {
	result, err := parseJson(resp)
	if err != nil {
		log.Fatalf("err decoding response JSON : '%v'", err)
	}
	return result
}

func (args HttpRequest) setHeaders(request *http.Request) {
	if strings.Contains(args.Auth, "Bearer") {
		request.Header.Set("Authorization", args.Auth)
	} else {
		request.SetBasicAuth(args.Auth, "")
	}
	request.Header.Set("Content-Type", args.ContentType)
}

func (args *HttpRequest) createRequest() (*http.Request, error) {
	args.Client = getSharedClient()
	var buf io.Reader
	if args.ContentType == "application/x-www-form-urlencoded" {
		buf = strings.NewReader(args.Form.Encode())
	} else {
		buf = bytes.NewBuffer(args.JSON)
	}
	request, err := http.NewRequest(args.VerbHTTP, args.Endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("err creating http request : %w", err)
	}
	args.setHeaders(request)
	return request, nil
}

func (args *HttpRequest) executeRequest(request *http.Request) (*http.Response, error) {
	response, err := args.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("err executing http request : %w", err)
	}
	return response, nil
}

func checkStatusCode(response *http.Response) (*http.Response, error) {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if err := handleNon2XXStatusCode(response); err != nil {
			return nil, err
		}
		return response, nil
	}
	return response, nil
}

func handleNon2XXStatusCode(response *http.Response) error {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("err reading http response body : %w", err)
	}
	bodyStr := string(bodyBytes)

	fmt.Printf("response : %+v\n", response)
	fmt.Printf("response.StatusCode : %+v\n", response.StatusCode)
	fmt.Println("bodyStr [handleNon2xxStatusCode()] : ", bodyStr)

	return &HttpError{
		StatusCode: response.StatusCode,
		Message:    fmt.Sprintf("[REQUEST.FAIL]-'%s'-'%s'", response.Status, bodyStr),
		ReqBody:    bodyStr,
	}
}
