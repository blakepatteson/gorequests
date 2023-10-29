package requests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"message": "success"}`)
		case "/post":
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"key": "value"}` {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, `{"message": "bad request body"}`)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"message": "post success"}`)
		case "/headers":
			auth := r.Header.Get("Authorization")
			if auth != "Bearer token123" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"message": "header success"}`)
		case "/badrequest":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, `{"message": "bad request"}`)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{"message": "not found"}`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	tests := []struct {
		name          string
		endpoint      string
		verbHTTP      string
		headers       string
		json          []byte
		wantErrorCode int
	}{
		{"success", server.URL + "/success", "GET", "", nil, http.StatusOK},
		{"POST request", server.URL + "/post", "POST", "", []byte(`{"key": "value"}`), http.StatusOK},
		{"bad POST request", server.URL + "/post", "POST", "", []byte(`{"wrong": "value"}`), http.StatusBadRequest},
		{"header success", server.URL + "/headers", "GET", "Bearer token123", nil, http.StatusOK},
		{"header fail", server.URL + "/headers", "GET", "Bearer wrongtoken", nil, http.StatusUnauthorized},
		{"bad request", server.URL + "/badrequest", "GET", "", nil, http.StatusBadRequest},
		{"not found", server.URL + "/notfound", "GET", "", nil, http.StatusNotFound},
		{"unknown error", server.URL + "/unknown", "GET", "", nil, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := HttpRequest{
				VerbHTTP:    tt.verbHTTP,
				Endpoint:    tt.endpoint,
				Auth:        tt.headers,
				JSON:        tt.json,
				ContentType: "application/json",
			}
			resp, err := args.Do()

			if err != nil && tt.wantErrorCode != http.StatusOK {
				if !strings.Contains(err.Error(), strconv.Itoa(tt.wantErrorCode)) {
					t.Errorf("Expected error containing status code '%v', got: '%v'", tt.wantErrorCode, err)
				}
			} else if err == nil && resp == nil {
				t.Errorf("Expected a response, got nil")
			} else if err == nil && resp.StatusCode != tt.wantErrorCode {
				t.Errorf("Expected status code '%v', got: '%v'", tt.wantErrorCode, resp.StatusCode)
			}
		})
	}
}
