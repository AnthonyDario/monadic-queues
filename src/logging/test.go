package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Send a post body
	json := "This is a new log line"
	body := []byte(json)
	res, err := http.Post("http://localhost:8000/log", "text/plain", bytes.NewReader(body))
    failOnError(err, "error with posting the log: %s")

	resBody, err := io.ReadAll(res.Body)
    failOnError(err, "impossible to read all body of response")
	log.Printf("res body: %s", string(resBody))

    // Get the log file
    res, err = http.Get("http://localhost:8000/get")
    failOnError(err, "Could not retrieve the log file from the server")

    resBody, err = io.ReadAll(res.Body)
    failOnError(err, "Could not read log file response")
    log.Printf("res body:\n%s", string(resBody))
}
