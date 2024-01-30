//  k cp data/shakespeare.txt $N:/tmp/shakespeare.txt
//  k cp cmd/embedding.go $N:/tmp/embedding.go
package main

import (
        "strings"
        "io/ioutil"
        "bytes"
        "fmt"
        "net/http"
)

const (
        kBatchSize = 4
        kRoughTokenCount = 512
)

const (
        tplBegin = `{ "input": [`
        tplEnd = `]}`
)

func main() {

        inputs := readData()

        rows := make([]string, 0)
        index := 0
        for batchId := 0; batchId < kBatchSize; batchId++ {
                row := "\"" + string(inputs[index:index+kRoughTokenCount*6]) + "\""
                 row = strings.ReplaceAll(row, "\n", "")
                index += kRoughTokenCount*6
                rows = append(rows, row)
        }

        out := strings.Join(rows, ", ")
        // fmt.Printf("%v", out)

        request := []byte(tplBegin + out + tplEnd)

        sendPost(request)

}

func readData() string {
        bs, err := ioutil.ReadFile("data/shakespeare.txt")
        if err != nil {
                panic(err)
	}

        return string(bs)
}

func sendPost(body []byte) {
        // HTTP endpoint
        posturl := "http://localhost:8080/v2/embeddings"

        // Create a HTTP post request
        r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
        if err != nil {
                panic(err)
	}
        // fmt.Printf("res: %v\n", r)
        res, err := http.DefaultClient.Do(r)
        if err != nil {
                panic(err)
	}
        fmt.Printf("res: %v", res)
        //bs, err := ioutil.ReadAll(res.Body)
        //fmt.Printf("res: %v", string(bs))
}
