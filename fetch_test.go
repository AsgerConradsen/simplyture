package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Unit test of getPage checking if it handles an empty list as response properly.
func TestEmptyResponse(t *testing.T) {
	// Creates a test server that returns and empty json list of invoices
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var invoices Invoices
		jsonBody, _ := json.Marshal(invoices)
		w.Write(jsonBody)
	}))
	defer s.Close()

	invoices, err := getPage(0, s.URL)
	if err != nil {
		t.Error("should return nil err")
	}
	if len(invoices) != 0 {
		t.Error("should return empty slice of invoices")
	}
}
