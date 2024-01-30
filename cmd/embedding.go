//  Cmd to try
//  k cp data/shakespeare.txt $N:/tmp/shakespeare.txt
//  k cp cmd/embedding.go $N:/tmp/embedding.go
package main

import (
        "strings"
        "io/ioutil"
        "bytes"
        "fmt"
        "net/http"
        "time"
)

const (
        kDebugMode = false
        kBatchSize = 150
        kRoughTokenCount = 512
        kLoopCountForEachStream = 4

        // Deduced.
        kRoughCharCount = kRoughTokenCount * 6
)

const (
        tplBegin = `{ "input": [`
        tplEnd = `]}`
)

func main() {
        inputs := readData()
        requests := prepareRowsForInputs(inputs)

        request := requests[0]

        startTime := time.Now()

        for loopId := 0; loopId < kLoopCountForEachStream; loopId++ {
                sendPost(request)
        }

        endTime := time.Now()
        fmt.Printf(
                "finished %v loop with %v to run\n",
                kLoopCountForEachStream,
                endTime.Sub(startTime))

}

func readData() string {
        bs, err := ioutil.ReadFile("data/shakespeare.txt")
        if err == nil {
                return string(bs)
        }

        // To make test easier, try both path.
        bs, err = ioutil.ReadFile("shakespeare.txt")
        if err != nil {
                panic(err)
	}

        return string(bs)
}

func prepareRowsForInputs(inputs string) [][]byte {
        requests:= make([][]byte, 0)

        rows := make([]string, 0)
        index := 0
        for batchId := 0; batchId < kBatchSize; batchId++ {
                if index + kRoughCharCount > len(inputs) {
                        panic("inputs are not enough")
                }

                row := "\"" + string(inputs[index:index+kRoughCharCount]) + "\""
                row = strings.ReplaceAll(row, "\n", "")
                index += kRoughCharCount
                rows = append(rows, row)
        }

        out := strings.Join(rows, ", ")
        // fmt.Printf("%v", out)
        request := []byte(tplBegin + out + tplEnd)

        requests = append(requests, request)
        return requests
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

        if res.StatusCode != 200 {
                panic("status code is not 200")
        }

        if kDebugMode {
                fmt.Printf("res: %v\n", res)

                body, err := ioutil.ReadAll(res.Body)
                if err != nil {
                        panic(err)
                }
                fmt.Printf("body: %v\n", string(body))
        }
}
