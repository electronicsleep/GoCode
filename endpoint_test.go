package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootEndpointsCheckHandler(t *testing.T) {
	endpoint_array := [...]string{"/", "/about", "/history"}
	check_array := [...]string{"GoCode", "About", "History"}
	count := 0
	expected := ""
	for _, endpoint := range endpoint_array {
		fmt.Println("endpoint: ", endpoint)
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(templateHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: (%v : %v)", status, http.StatusOK)
		} else {
			fmt.Printf("handler returned correct status code: (%v : %v)\n", status, http.StatusOK)
		}

		fmt.Println("count: ", count)
		expected = check_array[count]
		fmt.Println("expected: ", expected)
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("handler did not find expected string body: (expected: %v, endpoint: %v)", expected, endpoint)
		} else {
			fmt.Printf("handler contained expected string: (expected: %v, endpoint: %v)\n", expected, endpoint)
		}
		count += 1
	}
}
