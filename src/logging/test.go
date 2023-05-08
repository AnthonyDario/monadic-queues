package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main() {
	// Send a post body
	json := "This is a new log line"
	body := []byte(json)
	res, err := http.Post("http://localhost:8000/log", "text/plain", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("error with get", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	log.Printf("res body: %s", string(resBody))
}
