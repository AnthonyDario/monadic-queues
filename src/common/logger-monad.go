// A writer monad implemented in go.  This keeps track of log messages during
// computation.  A call to "commit" will record the log messages into some log
// sink, currently just the logging server

package common

import (
    "fmt"
    "time"
    "io"
    "log"
	"bytes"
    "net/http"
)

const DOMAIN = "localhost"
const PORT   = "8000"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type LoggerMonad [T any] struct {
    Value T
    Log string
}

// Monadic Functions
// --------------------

// Build a writer monad
func LoggerUnit [T any] (a T) LoggerMonad[T] {
    log := "initialize empty logger monad"
    sendLog(log, DOMAIN, PORT)
    return logTime(a, "initialized logger monad")
}

func LoggerBuild [T any] (a T, msg string) LoggerMonad[T] {
    sendLog(msg, DOMAIN, PORT)
    return logTime(a, msg)
}

// Compose computations using the writer monad
func LoggerBind [T any, U any] (w LoggerMonad[T], f func(T) LoggerMonad [U]) LoggerMonad[U] {
    var w2 = f(w.Value)
    sendLog(w2.Log, DOMAIN, PORT)
    return LoggerMonad[U] {w2.Value, w.Log + "\n" + w2.Log}
}


// Helpers 
// --------------

// prefix our log with the current timestamp
func logTime[T any] (v T, msg string) LoggerMonad[T] {
    t := time.Now()
    return LoggerMonad[T]{v, t.Format(time.RFC3339) + " " + msg}
}

func sendLog (msg string, domain string, port string) {
    body := []byte(msg)
    _, err := http.Post("http://" + domain + ":" + port + "/log",
                          "text/plain",
                          bytes.NewReader(body))
    if err != nil {
        log.Fatalf("Could not commit the logs: %s", err)
	}
}

// We want to be able to commit the value of the writer to the log server
func commit [T any] (w LoggerMonad[T], domain string, port string) {
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

func testLogger () {
    // Build our writer with unit
    var w = LoggerUnit(1) 

    // Our function from int -> LoggerMonad[bool]
    var f = func (i int) LoggerMonad[bool] {
        var isEven = i % 2 == 0
        var log string
        if isEven {
            log = fmt.Sprintf("%d is even", i)
        } else {
            log = fmt.Sprintf("%d is odd", i)
        }
        
        return logTime(i % 2 == 0, log) 
    }

    var g = func (i int) LoggerMonad[int] {
        return logTime(i + 1, "incremented i")
    }
    
    var w2 = LoggerBind(w, g)
    var w3 = LoggerBind(w2, g)
    var w4 = LoggerBind(w3, g)
    var w5 = LoggerBind(w4, f)

    //log.Print(w5.Log)
    commit(w5, "localhost", "8000")
    retrieve()
}
