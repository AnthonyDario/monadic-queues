package main

import (
    "strconv"
    "os"
    "time"
    "bytes"
    "errors"
    "net/http"
    "log"
    "fmt"

    c "common"
)

type handler func(http.ResponseWriter, *http.Request)

type Result struct {
    Value string
    Error error
}

// Helpers 
// ---------------

// Grab the value of the key from the url query
func getKey(key string, q map[string][]string, sink string) c.LoggerMonad[Result] {
    var log string
    var err error
    var v string
    val, prs := q[key]
    if !prs {
        log = key + " not selected"
        err = errors.New(key + " not selected")
        v = ""
    } else {
        log = val[0] + " selected for " + key
        err = nil
        v = val[0]
    }

    return c.LoggerUnit(Result {v, err}, log, sink)
}

// Return a failure response
func writeBadRequest(w http.ResponseWriter, msg string, sink string) func(error) c.LoggerMonad[error] {
    return func(e error) c.LoggerMonad[error] {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, msg)
        return c.LoggerUnit(e, msg + ": " + e.Error(), sink)
    }
}

// Build a url from a domain and a port
func buildUrl(domain string, port string) string {
    return "http://" + domain + ":" + port
}

// The pipeline for an order using monads. Very confusing
func orderPipeline(query map[string][]string) c.ConfigMonad[c.LoggerMonad[Result]] {

    buildKeyConfigMonad := func (key string) c.ConfigMonad[c.LoggerMonad[Result]] {
        return c.ConfigMonad[c.LoggerMonad[Result]] {
            func (env map[string]string) c.LoggerMonad[Result] {
                sink := buildUrl(env["LogDomain"], env["LogPort"])
                return getKey(key, query, sink)
            },
        }
    }

    // want a function from a logger monad to another logger monad, getting a key
    buildKeyLoggerMonad := func (key string) func (c.LoggerMonad[Result]) c.LoggerMonad[Result] {
        return func (lm c.LoggerMonad[Result]) c.LoggerMonad[Result] {
            return c.LoggerBind(lm, func(r Result, sink string) c.LoggerMonad[Result] {
                err := r.Error
                if err != nil {
                    return lm
                } else {
                    return getKey(key, query, sink)
                }
            })
        }
    }

    buildWrapGetKey := func (key string) func(c.LoggerMonad[Result]) c.ConfigMonad[c.LoggerMonad[Result]] {
        return func (lm c.LoggerMonad[Result]) c.ConfigMonad[c.LoggerMonad[Result]] {
            return c.ConfigUnit(buildKeyLoggerMonad(key)(lm))
        }
    }

    toppings := buildKeyConfigMonad("toppings")
    size := c.ConfigBind(toppings, buildWrapGetKey("size")) 
    name := c.ConfigBind(size, buildWrapGetKey("username")) 

    return name
}

func order(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    c.RunConfig(orderPipeline(query))
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("hello"))
}

// The non-monadic approach
// -------------------------------------

func sendLog (msg string) {
    body := []byte(msg)
    _, err := http.Post("http://localhost:8000/log", "text/plain", bytes.NewReader(body))
    if err != nil {
        log.Fatalf("Could not commit the logs: %s", err)
	}
}

// Return a failure response
func writeFail(w http.ResponseWriter, msg string) {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, msg)
}

func orderFree(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()

    toppings, prs := q["toppings"]
    if !prs {
        sendLog("toppings not selected")
        writeFail(w, "toppings not selected")
        return
    }

    size, prs := q["size"]
    if !prs {
        sendLog("size not selected")
        writeFail(w, "size not selected")
        return
    }

    name, prs := q["username"]
    if !prs {
        sendLog("name not selected")
        writeFail(w, "name not selected")
        return
    }
    
    msg := "order successful with toppings " + toppings[0] + " size: " + size[0] + " name: " + name[0]
    sendLog(msg)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(msg))
}

func instrument(f func (http.ResponseWriter, *http.Request)) handler {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now() 
        f(w, r)

        end := time.Now()
        elapsed := end.Sub(start)
        data, _ := os.OpenFile("data.data", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        defer data.Close()
        msg := strconv.Itoa(int(elapsed.Milliseconds())) + "\n"
        data.WriteString(msg)
    }
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    //http.HandleFunc("/order", order)
    http.HandleFunc("/order", instrument(orderFree))
    log.Print("Starting Pizza Server")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
