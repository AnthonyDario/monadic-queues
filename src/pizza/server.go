package main

import (
    "net/http"
    "log"
)

func main() {
    fs := http.FileServer(http.Dir("./web"))
    http.Handle("/", fs)
    log.Print("Starting Pizza Server")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
