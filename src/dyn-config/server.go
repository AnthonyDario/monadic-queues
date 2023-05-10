/* 
 * A server that hosts a config file.  The file can be read and updated
 * remotely
 */

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
    "log"
    "fmt"
    "sync"
)

const CONFIG_FILE = "config.json"

type handler func(http.ResponseWriter, *http.Request)

// Multiple requests can attempt to write to the config concurrently,
// gotta make it thread safe
type SafeConfig struct {
    mu   sync.Mutex
    conf map[string]interface{}
}

func (c *SafeConfig) write(key string, val string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // update the local copy
    c.conf[key] = val

    // write to disk
	config, err := os.OpenFile(CONFIG_FILE, os.O_RDWR|os.O_CREATE, 0644)
    failOnError(err, "Could not open config file")
	defer config.Close()

    bs, err := json.Marshal(c.conf)
    failOnError(err, "Unable to marshal updated config json")
    io.WriteString(config, string(bs))
}

// Don't really need to be thread safe for reading
func (c *SafeConfig) read(key string) (string, error) {
    v := c.conf[key]
    switch vv := v.(type) {
        case string:
            return vv, nil
        default:
            return "", &ReadError{key, "value of wrong type"}
    }
}

// Reading Types
type ReadError struct {
	Field string
	Msg   string
}

func (e *ReadError) Error() string {
    return fmt.Sprintf("error reading %s: %s", e.Field, e.Msg)
}

type ReadRequest struct {
	Field string
}

type ReadResponse struct {
	Field string
	Value string
}

// Writing Types
type WriteError struct {
    Field string
    Value string
    Msg   string
}

func (e *WriteError) Error() string {
    return fmt.Sprintf(`error writing "%s": "%s" - %s`, 
                       e.Field, e.Value, e.Msg)
}

type WriteRequest struct {
    Field string
    Value string
}

type WriteResponse struct {
    success bool
}

// Util
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func loadConfig() SafeConfig {
	config, err := os.OpenFile(CONFIG_FILE, os.O_RDWR|os.O_CREATE, 0644)
    failOnError(err, "Could not open config file")
	defer config.Close()

    bs, err := io.ReadAll(config)
    failOnError(err, "Unable to read config file")

    var js interface{}
    json.Unmarshal(bs, &js)

    conf := js.(map[string]interface{})
    return SafeConfig{conf: conf}
}

// Handlers
func makeWriteHandler(config *SafeConfig) handler {
    // TODO: failOnError -> return a failure response
    return func(w http.ResponseWriter, r *http.Request) {
        bs, err := io.ReadAll(r.Body)
        failOnError(err, "Could not read write request body")

        var req WriteRequest
        err = json.Unmarshal(bs, &req)
        failOnError(err, "Could not unmarshal write request json")

        config.write(req.Field, req.Value)

        ret := WriteResponse{true}
        retbs, err := json.Marshal(ret)
        failOnError(err, "Could not marshal the return json")

        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, string(retbs))
        log.Print("Responded OK to write request")
    }
}

func makeReadHandler(config *SafeConfig) handler {
	// TODO: Turn these failOnError calls into actual failure responses
	return func(w http.ResponseWriter, r *http.Request) {
        
		bs, err := io.ReadAll(r.Body)
		failOnError(err, "Could not read request body")

	    var req ReadRequest
		err = json.Unmarshal(bs, &req)
		failOnError(err, "Could not unmarshal read request json")

        val, err := config.read(req.Field)
        failOnError(err, "Could not read from config")

		ret := ReadResponse{req.Field, val}
        retbs, err := json.Marshal(ret)
        failOnError(err, "Could not marshal the return json")

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(retbs))
		log.Print("Responded OK to read request")
	}
}

func main() {
    config := loadConfig()

	http.HandleFunc("/read", makeReadHandler(&config))
	http.HandleFunc("/write", makeWriteHandler(&config))
	log.Print("Starting Log Server")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
