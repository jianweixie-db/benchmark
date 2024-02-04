package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	kStreamTotalCount       = 4
	kLoopCountForEachStream = 10
	kDebugMode              = true
)

func main() {
	var wg sync.WaitGroup
	startTime := time.Now()

	for groupId := 0; groupId < kStreamTotalCount; groupId++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for loopId := 0; loopId < kLoopCountForEachStream; loopId++ {
				sendRequst()
			}
		}()
	}

	wg.Wait()
	endTime := time.Now()

	fmt.Printf("Parallel Streams : %v\n", kStreamTotalCount)

	fmt.Printf(
		"finished %v streams each having %v loop with %v to run\n",
		kStreamTotalCount,
		kLoopCountForEachStream,
		endTime.Sub(startTime))

	fmt.Printf(
		"this means %v inputs per second\n",
		kLoopCountForEachStream*kStreamTotalCount/endTime.Sub(startTime).Seconds(),
	)

}

var fixedBody = []byte("{\"inputs\": []} ")

func sendRequst() {
	// HTTP endpoint
	posturl := "http://localhost:8080/test"

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(fixedBody))
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
