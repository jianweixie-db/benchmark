package main

import (
        "bytes"
        "fmt"
        "net/http"
)

type Post struct {
	Input  string `json:"input"`
}

func main() {
        // HTTP endpoint
        posturl := "http://localhost:8080/v2/embeddings"

        // JSON body
        body := []byte(`{
                "input": "hello world"
        }`)

        // Create a HTTP post request
        r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
        if err != nil {
                panic(err)
	}
        fmt.Printf("res: %v\n", r)
        res, err := http.DefaultClient.Do(r)
        if err != nil {
                panic(err)
	}
        fmt.Printf("res: %v", res.Body)
}
