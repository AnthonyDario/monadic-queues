package main

import (
    "net/http"
    "log"
    "fmt"

    c "common"
)

func getKey(key string, q map[string][]string) c.LoggerMonad[string] {
    var log string
    val, err := q[key]
    if !err {
        log = key + " not selected"
    } else {
        log = val[0] + " selected for " + key
    }
    return c.LoggerBuild(val[0], log)
}

func writeBadRequest(w http.ResponseWriter, msg string) func(error) c.LoggerMonad[error] {
    return func(e error) c.LoggerMonad[error] {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, msg)
        return c.LoggerBuild(e, msg + ": " + e.Error())
    }
}

func order(w http.ResponseWriter, r *http.Request) {

    // initialize our monad
    query := r.URL.Query()

    lm := getKey("toppings", query)
    if lm.Value == "" {
        writeBadRequest(w, lm.Log)
        return
    }
    toppings := lm.Value

    lm = getKey("size", query)
    if lm.Value == "" {
        writeBadRequest(w, lm.Log)
        return
    }
    size := lm.Value

    lm = getKey("username", query)
    if lm.Value == "" {
        writeBadRequest(w, lm.Log)
        return
    }
    name := lm.Value

    // add the order to the database
    // return a success page

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("order successful with toppings " + toppings + "size: " + size + "name: " + name))
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/order", order)
    log.Print("Starting Pizza Server")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
