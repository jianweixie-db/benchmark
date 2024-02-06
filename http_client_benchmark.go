package main

import (
        "io"
        "log"
        "io/ioutil"
        "fmt"
        "net/http"
	"time"
        "bytes"
)

var (
        URL = "http://ec2-54-237-91-215.compute-1.amazonaws.com/test"
        INPUT_SIZE = 400_000
)

func main() {

        var buf bytes.Buffer

        for i := 0; i < INPUT_SIZE; i++ {
                buf.WriteByte('a')
        }

        req, err := http.NewRequest("POST", URL, &buf)
        if err != nil {
                panic(err)
        }

        // explicity setting the content-length will disable the "transfer-encoding: chunked" header
        req.ContentLength = int64(INPUT_SIZE)
        req.TransferEncoding = []string{"identity"}

        client := &http.Client{}

	startTime := time.Now()
        res, err := client.Do(req)
        if err != nil {
                panic(err)
        }
        n, err2:= io.Copy(ioutil.Discard, res.Body)

        //bytes, err := ioutil.ReadAll(res.Body)
        defer res.Body.Close()
	endTime := time.Now()

        fmt.Printf("took %v\n", endTime.Sub(startTime))
        fmt.Printf("res %v err %v\n", res, err)

        log.Printf("n %v err2 %v", n, err2)
}
