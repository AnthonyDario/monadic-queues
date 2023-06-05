// A service wrapping a queue.  It manages the logging that comes in with
// the messages. Committing logs to the log server.
package main

import (
    "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

    c "common"
)

type handler func(http.ResponseWriter, *http.Request)

// Util
// ---------------------
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Interface types
// --------------------
type SendRequest struct {
    Msg string
    Log string
    Dest string
}

// Handlers
// ---------------------
func makeSendHandler(qs map[string]*c.Queue) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Recieved Send Request")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Incorrect method, this endpoint needs a body")
		}

        // Read the message
		bs, err := ioutil.ReadAll(r.Body)
		failOnError(err, "Failed to publish a message")

        // Extract the log line
        var req SendRequest
        err = json.Unmarshal(bs, &req)
        failOnError(err, "Could not unmarshal send Request json")

        // Send the log message to the log queue
        err = qs["log"].Send([]byte(req.Log))
        failOnError(err, "Failed to publish log lines")

        // And the desired message to the message queue
        err = qs[req.Dest].Send([]byte(req.Msg))
        failOnError(err, "Failed to publish queue message")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message Sent")
		log.Print("Finished Sending Message")
	}
}

func main() {
    // Define our queues
	pq := c.Connect("pizza")
    lq := c.Connect("log")
    m := make(map[string]*c.Queue)
    m["pizza"] = &pq
    m["log"] = &lq
	defer pq.Close()
	defer lq.Close()

    // Build our handlers
	http.HandleFunc("/send", makeSendHandler(m))
	log.Print("Starting Server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
