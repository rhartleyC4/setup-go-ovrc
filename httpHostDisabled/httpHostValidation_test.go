package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	invalidHost      = "#Invalid@Host#"
	expectedResponse = "ok"
)

func TestHttpHostValidationDisabled(t *testing.T) {
	// arrange
	server := httptest.NewServer(http.HandlerFunc(testHandler))
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error creating request: %s", err))
		t.FailNow()
	}
	req.Header.Set("Host", invalidHost)
	req.Host = invalidHost
	client := http.Client{}

	// act
	response, err := client.Do(req)

	// assert
	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			fmt.Println(fmt.Sprintf("Error closing body: %v", closeErr))
		}
	}(response.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Request failed: %v", err))
		t.FailNow()
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Request failed: %s", response.Status))
		t.FailNow()
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to read response body, %v", err))
		t.FailNow()
	}
	if string(body) != expectedResponse {
		fmt.Println(fmt.Sprintf("Invalid response body: %s\nExpected: %s", string(body), expectedResponse))
		t.FailNow()
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.Host == invalidHost {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
