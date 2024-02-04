package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	sleepDuration = 1000 * time.Millisecond
)

func main() {

	// Single thread http server
	var mu sync.Mutex

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(sleepDuration)
		fmt.Fprintf(w, "Hello, %q", (r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
