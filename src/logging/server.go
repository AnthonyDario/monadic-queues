package main

import (
    "log"
    "net/http"
    "fmt"
	"io"
    "os"
)

const LOG_FILE = "log.log"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// TODO: have the failure return a failure response
func logHandler (w http.ResponseWriter, r *http.Request) {
    // Write the request body to a log file
    // Lets just do the single log file for now...

    if r.Method != "POST" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Incorrect method, this endpoint needs a body")
        return
    }

    reqBody, err := io.ReadAll(r.Body)
    failOnError(err, "Failed to publish a message")

    log.Println(string(reqBody))
}

// Just returning the entire log file now for now. DOES
// NOT SCALE
func getHandler (w http.ResponseWriter, r *http.Request) {
    logFile, err := os.Open(LOG_FILE)
    failOnError(err, "Could not open log file for reading")
    defer logFile.Close()

    bs, err := io.ReadAll(logFile)
    failOnError(err, "Could not read the log file as bytes")

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, string(bs))
}

func main() {
    logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
    failOnError(err, "Could not open log file")
    defer logFile.Close()
    log.SetOutput(logFile)

	http.HandleFunc("/log", logHandler)
    http.HandleFunc("/get", getHandler)
    fmt.Print("Starting Log Server");
	log.Print("Starting Log Server")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
