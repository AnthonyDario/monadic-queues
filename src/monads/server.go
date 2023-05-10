package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

// Util
// ---------------------
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Handlers
// ---------------------
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func makeSendHandler(q Queue) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Recieved Send Request")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Incorrect method, this endpoint needs a body")
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		failOnError(err, "Failed to publish a message")

        err = q.send([]byte(reqBody))
		failOnError(err, "Failed to publish a message")

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message Sent")
		log.Print("Finished Sending Message")
	}
}

func main() {
	q := connect("pizza")
	defer q.Close()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/send", makeSendHandler(q))
	log.Print("Starting Server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
