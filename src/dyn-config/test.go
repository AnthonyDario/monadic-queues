/*
 * Tests for the dynamic config server
 */
package main

import (
    "fmt"
	"bytes"
	"io"
	"log"
	"net/http"
)

type readRequest struct {
	Field string
}

type Message struct {
    Name string
    Body string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

/*
func testJson() {
    js := []byte(`{ "Field": "key" }`)

    var tt readRequest
    err := json.Unmarshal(js, &tt)
    if err != nil {
		log.Fatalf("error with posting the read request: %s", err)
    }

    log.Print(fmt.Sprintf("tt.field = %s", tt.Field))
}
*/

func testReading(field string) {
	js := fmt.Sprintf(`{ "field": "%s" }`, field)
	body := []byte(js)
    res, err := http.Post("http://localhost:9090/read",
                          "application/json", 
                          bytes.NewReader(body))
    failOnError(err, "Error with posting the read request")

	resBody, err := io.ReadAll(res.Body)
    failOnError(err, "Impossible to read entire body of response")

	log.Printf("read response body: %s", string(resBody))
}

func testReadAll() {
    res, err := http.Get("http://localhost:9090/readAll")
    failOnError(err, "Error with posting the read request")

	resBody, err := io.ReadAll(res.Body)
    failOnError(err, "Impossible to read entire body of response")

	log.Printf("read response body: %s", string(resBody))
}

func testWriting(field string, value string) {
    js := fmt.Sprintf("{ \"Field\": \"%s\", \"Value\": \"%s\" }", field, value)
    body := []byte(js)
    res, err := http.Post("http://localhost:9090/write",
                          "application/json",
                          bytes.NewReader(body))
    failOnError(err, "Error with posting the write request")

    resBody, err := io.ReadAll(res.Body)
    failOnError(err, "Impossible to read entire body of response")

    log.Printf("Write response body: %s", string(resBody))
}

func main() {
    //testJson()
    testReading("LogDomain")
    testWriting("newfield", "newvalue")
    testWriting("newnew", "oldold")
    testWriting("dockertest", "dockersuccess")
    testReading("dockertest")
    testReadAll()
}
