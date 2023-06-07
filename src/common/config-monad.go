// The ConfigMonad is a slightly modified version of the reader monad

package common

import (
    "io"
    "encoding/json"
    "log"
    "net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// The config monad contains the current value of computation and the
// configuration
type ConfigMonad [T any] struct {
    F func(map[string]string) (T)
}

// The monadic functions
// ---------------------------

// Initialize the config monad with the current config
func ConfigUnit [T any] (v T) ConfigMonad[T] {
    return ConfigMonad[T] {
        func (env map[string]string) T {
            return v
        },
    }
}

// Bind modifies the monad, applying the given function, f, to its result
// bind : m a -> (a -> m b) -> mb
func ConfigBind[T any, U any] (m ConfigMonad[T], f func(T) ConfigMonad[U]) ConfigMonad[U] {
    return ConfigMonad[U] { 
        func (env map[string]string) U {
            val := m.F(env)
            return f(val).F(env)
        },
    }
}

// Helpers
// -----------

// Loads the config 
func initialConfig() map[string]string {
    res, err := http.Get("http://localhost:9090/readAll")
    failOnError(err, "Error with posting the read request")

	resBody, err := io.ReadAll(res.Body)
    failOnError(err, "Impossible to read entire body of response")

    var m map[string]string
    err = json.Unmarshal(resBody, &m)
    failOnError(err, "Failed to unmarshal the json config")

    return m
}

// Runs the config and provides the output value(we know what the environment
// is already so don't need to provide it)
func RunConfig [T any] (m ConfigMonad[T]) T {
    return m.F(initialConfig())
}

func testConfig() {
    // Reads the domain and port from the config and returns the url
    m := ConfigMonad[string] {
        func (m map[string]string) string {
            return "https://" + m["LogDomain"] + ":" + m["LogPort"]
        },
    }

    url := RunConfig(m)
    log.Print(url)

}
