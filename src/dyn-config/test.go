package main

import (
    "fmt"
	"bytes"
	"io"
	"log"
	"net/http"
    "encoding/json"
)

type readRequest struct {
	Field string
}

type Message struct {
    Name string
    Body string
}

func testJson() {
    js := []byte(`{ "Field": "key" }`)

    var tt readRequest
    err := json.Unmarshal(js, &tt)
    if err != nil {
		log.Fatalf("error with posting the read request: %s", err)
    }

    log.Print(fmt.Sprintf("tt.field = %s", tt.Field))
}

func main() {
    testJson()

	json := `{ "field": "key" }`
	body := []byte(json)
    res, err := http.Post("http://localhost:9090/read",
                          "application/json", 
                          bytes.NewReader(body))
    if err != nil {
		log.Fatalf("error with posting the read request: %s", err)
    }

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("impossible to read all body of response: %s", err)
	}
	log.Printf("res body: %s", string(resBody))
}
