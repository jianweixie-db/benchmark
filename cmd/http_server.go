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

        bytes = make([]byte, 2907082, 2907082)
        for i :=0; i< 2907082; i++ {
                bytes[i] = 'a'
        }

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		//time.Sleep(sleepDuration)
		fmt.Fprintf(w, "%v", bytes)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
