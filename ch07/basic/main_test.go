package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	mockUserResp := `{"message":"hello world"}`

	ts := httptest.NewServer(SetupServer())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	responseData, _ := io.ReadAll(resp.Body)
	if string(responseData) != mockUserResp {
		t.Fatalf("Expected %s, got %s", mockUserResp, string(responseData))
	}
}
