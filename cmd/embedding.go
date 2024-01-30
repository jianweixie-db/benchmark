// Cmd to try
// alias k='kubectl -n <namespace>'
// export N=the POD name
//
// k cp data/shakespeare.txt $N:/tmp/shakespeare.txt
// k cp cmd/embedding.go $N:/tmp/embedding.go
// k exec -ti $N -- bash
// cd /tmp
// go run embedding.go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	kDebugMode       = false
	kBatchSize       = 130
	kRoughTokenCount = 512

	kStreamTotalCount       = 4
	kLoopCountForEachStream = 10

	// Deduced.
	kRoughCharCount = kRoughTokenCount * 6
)

const (
	tplBegin = `{ "input": [`
	tplEnd   = `]}`
)

func main() {
	inputs := readData()
	requests := prepareRowsForInputs(inputs, kStreamTotalCount)

	var wg sync.WaitGroup
	startTime := time.Now()

	for groupId := 0; groupId < kStreamTotalCount; groupId++ {
		wg.Add(1)
		request := requests[groupId]

		go func(req []byte) {
			defer wg.Done()
			for loopId := 0; loopId < kLoopCountForEachStream; loopId++ {
				sendPost(req)
			}
		}(request)
	}

	wg.Wait()
	endTime := time.Now()

	fmt.Printf("Batch Size       : %v\n", kBatchSize)
	fmt.Printf("Parallel Streams : %v\n", kStreamTotalCount)
	fmt.Printf("Chars Per Input  : %v\n", kRoughCharCount)

	fmt.Printf(
		"finished %v streams each having %v loop with %v to run\n",
		kStreamTotalCount,
		kLoopCountForEachStream,
		endTime.Sub(startTime))

	fmt.Printf(
		"this means %v inputs per second\n",
		kLoopCountForEachStream*kBatchSize*kStreamTotalCount/endTime.Sub(startTime).Seconds(),
	)

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

func prepareRowsForInputs(inputs string, totalStreamCount int) [][]byte {
	requests := make([][]byte, 0)

	index := 0

	for streamId := 0; streamId < totalStreamCount; streamId++ {
		rows := make([]string, 0)
		for batchId := 0; batchId < kBatchSize; batchId++ {
			if index+kRoughCharCount > len(inputs) {
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
	}

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
