package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	httpGet = "GET"
)

func TestServe(t *testing.T) {
	// Create a request to be passed to the serve function
	req := httptest.NewRequest(httpGet, "/", nil)
	w := httptest.NewRecorder()

	// Capture the current environment variables
	originalEnv := os.Environ()

	// Restore the environment variables at the end of the test
	defer func() {
		os.Clearenv()
		for _, keyval := range originalEnv {
			keyval := strings.SplitN(keyval, "=", 2)
			os.Setenv(keyval[0], keyval[1])
		}
	}()

	// Set some environment variables for testing
	os.Setenv("TEST_VAR_1", "value1")
	os.Setenv("TEST_VAR_2", "value2")

	// Call the serve function
	serve(w, req)

	// Check the HTTP response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse the response body as JSON
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing JSON response: %v", err)
	}

	// Check the response for specific environment variables
	expectedVars := map[string]string{
		"TEST_VAR_1": "value1",
		"TEST_VAR_2": "value2",
	}

	for key, expectedValue := range expectedVars {
		actualValue, ok := response[key]
		if !ok {
			t.Errorf("Expected environment variable %s not found in response", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected value %s for %s, got %s", expectedValue, key, actualValue)
		}
	}
}
