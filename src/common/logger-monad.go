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

type LoggerMonad [T any] struct {
    Value T
    Log string
    Sink string
}

// Monadic Functions
// --------------------

// Build a writer monad
func LoggerUnit [T any] (a T, msg string, sink string) LoggerMonad[T] {
    sendLog(msg, sink)
    return logTime(a, msg, sink)
}

// Compose computations using the writer monad
func LoggerBind [T any, U any] (w LoggerMonad[T], f func(T, string) LoggerMonad [U]) LoggerMonad[U] {
    var w2 = f(w.Value, w.Sink)
    sendLog(w2.Log, w2.Sink)
    return LoggerMonad[U] {w2.Value, w.Log + "\n" + w2.Log, w2.Sink}
}

// Helpers 
// --------------

// prefix our log with the current timestamp
func logTime[T any] (v T, msg string, sink string) LoggerMonad[T] {
    t := time.Now()
    return LoggerMonad[T]{v, t.Format(time.RFC3339) + " " + msg, sink}
}

func sendLog (msg string, sink string) {
    body := []byte(msg)
    _, err := http.Post(sink + "/log", "text/plain", bytes.NewReader(body))
    if err != nil {
        log.Fatalf("Could not commit the logs: %s", err)
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
    // Our function from int -> LoggerMonad[bool]
    var f = func (i int, sink string) LoggerMonad[bool] {
        var isEven = i % 2 == 0
        var log string
        if isEven {
            log = fmt.Sprintf("%d is even", i)
        } else {
            log = fmt.Sprintf("%d is odd", i)
        }
        
        return logTime(i % 2 == 0, log, sink)
    }

    var g = func (i int, sink string) LoggerMonad[int] {
        return logTime(i + 1, "incremented i", sink)
    }
    
    var w = LoggerUnit(1, "initialized LoggerMonad", "http://localhost:8000")
    var w2 = LoggerBind(w, g)
    var w3 = LoggerBind(w2, g)
    var w4 = LoggerBind(w3, g)
    var w5 = LoggerBind(w4, f)

    log.Print(w5.Log)
    retrieve()
}
