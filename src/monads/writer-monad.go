// A writer monad implemented in go.  This keeps track of log messages during
// computation.  A call to "commit" will record the log messages into some log
// sink, currently just the logging server

package main

import (
    "fmt"
    "time"
    "io"
    "log"
	"bytes"
    "net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Writer [T any] struct {
    Value T
    Log string
}

// Build a writer monad
func unit [T any] (a T) Writer[T] {
    return Writer[T] {a, ""}
}

// Compose computations using the writer monad
func bind [T any, U any] (w Writer[T], f func(T) Writer [U]) Writer[U] {
    var w2 = f(w.Value)
    return Writer[U] {w2.Value, w.Log + "\n" + w2.Log}
}

// prefix our log with the current timestamp
func logTime[T any] (v T, msg string) Writer[T] {
    t := time.Now()
    return Writer[T]{v, t.Format(time.RFC3339) + " " + msg}
}

// We want to be able to commit the value of the writer to the log server
func commit [T any] (w Writer[T], domain string, port string) {
    // Call the log server with our log message

    body := []byte(w.Log)
    res, err := http.Post("http://" + domain + ":" + port + "/log",
                          "text/plain",
                          bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Could not commit the logs: %s", err)
	}

	// Send a post body
	_, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Could not read log-server response: %s ", err)
	}
}

func retrieve () {
    res, err := http.Get("http://localhost:8000/get")
    failOnError(err, "Could not retrieve the log file from the server")

    resBody, err := io.ReadAll(res.Body)
    failOnError(err, "Could not read log file response")
    log.Printf("res body:\n%s", string(resBody))
}

func main () {
    // Build our writer with unit
    var w = unit(1) 

    // Our function from int -> Writer[bool]
    var f = func (i int) Writer[bool] {
        var isEven = i % 2 == 0
        var log string
        if isEven {
            log = fmt.Sprintf("%d is even", i)
        } else {
            log = fmt.Sprintf("%d is odd", i)
        }
        
        return logTime(i % 2 == 0, log) 
    }

    var g = func (i int) Writer[int] {
        return logTime(i + 1, "incremented i")
    }
    
    var w2 = bind(w, g)
    var w3 = bind(w2, g)
    var w4 = bind(w3, g)
    var w5 = bind(w4, f)

    //log.Print(w5.Log)
    commit(w5, "localhost", "8000")
    retrieve()
}
