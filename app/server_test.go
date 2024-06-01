package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHome(t *testing.T) {
	go main()

	met := "GET"
	url := "http://localhost:4221/"
	req, _ := http.NewRequest(met, url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Got no response from sending %s request to %s", met, url)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf(fmt.Sprintf("Could not connect to route. Expected 200, got %s", res.Status))
	}
}
