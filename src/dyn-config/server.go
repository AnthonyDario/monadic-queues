package main

// A server that hosts a config file.  The file can be read and updated
// remotely

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
    "log"
    "fmt"
)

const CONFIG_FILE = "config.json"

type handler func(http.ResponseWriter, *http.Request)

type readError struct {
	Field string
	Msg   string
}

func (e *readError) Error() string {
    return fmt.Sprintf("error reading %s: %s", e.Field, e.Msg)
}

type readRequest struct {
	Field string
}

type readResponse struct {
	Field string
	Value string
}

// Util
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func makeReadHandler(config *os.File) handler {
	readConfig := func(key string) (string, error) {
        log.Print(fmt.Sprintf("reading key %s", key))

		bs, err := io.ReadAll(config)
		failOnError(err, "Unable to read config file:")

		var js interface{}
		json.Unmarshal(bs, &js)

		m := js.(map[string]interface{})
		v := m[key]

		switch vv := v.(type) {
		case string:
			if vv == "" {
				return vv, &readError{key, "field not populated"}
			}
			return vv, nil
		default:
			return "", &readError{key, "field not a string"}
		}
	}

	// TODO: Turn these failOnError calls into actual failure responses
	return func(w http.ResponseWriter, r *http.Request) {
		bs, err := io.ReadAll(r.Body)
		failOnError(err, "Could not read request body")

        log.Print(fmt.Sprintf("request body: %s", string(bs)))

	    var req readRequest
		err = json.Unmarshal(bs, &req)
		failOnError(err, "Could not unmarshall read request json")

        log.Print(fmt.Sprintf("readRequest.field = %s", req.Field))

		val, err := readConfig(req.Field)
		failOnError(err, "Could not read field")

		// return it in a json?
		ret := readResponse{req.Field, val}
		w.WriteHeader(http.StatusOK)

        retbs, err := json.Marshal(ret)
        failOnError(err, "Could not marshal the return json")

		fmt.Fprintf(w, string(retbs))
		log.Print("Responded OK")
	}
}

func main() {
	config, err := os.OpenFile(CONFIG_FILE, os.O_RDWR|os.O_CREATE, 0644)
    failOnError(err, "Could not open config file")
	defer config.Close()

	// http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/read", makeReadHandler(config))
	log.Print("Starting Log Server")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
