package main

import (
    "log"
    "net/http"
    "fmt"
	"io/ioutil"
    "os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func logHandler (w http.ResponseWriter, r *http.Request) {
    // Write the request body to a log file
    // Lets just do the single log file for now...

    if r.Method != "POST" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Incorrect method, this endpoint needs a body")
    }

    reqBody, err := ioutil.ReadAll(r.Body)

    log.Println(string(reqBody))
    failOnError(err, "Failed to publish a message")
}

func main() {
    LOG_FILE := "log.log"

    logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    failOnError(err, "Could not open log file")
    defer logFile.Close()
    log.SetOutput(logFile)

	http.HandleFunc("/log", logHandler)
    fmt.Print("Starting Log Server");
	log.Print("Starting Log Server")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
