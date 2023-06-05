package main

import (
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    // return the home page
    if r.Method != "GET" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Incorrect http method, GET required")
        return
    }

    w.write
}

func main() {
    http.HandleFunc("/", homeHandler)
}
