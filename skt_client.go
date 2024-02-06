package main

import (
        "net"
	"time"
        "fmt"
        "io"
        "log"
)

var (
        URL = "ec2-54-84-67-103.compute-1.amazonaws.com:80"
        BytesToWrite = 100
        LoopsToCall = 10
)

func main() {

	startTime := time.Now()
        for i:=0; i < LoopsToCall; i++ {
                onecall()
        }
	endTime := time.Now()
        dur := endTime.Sub(startTime)
        fmt.Printf("took %v ave %v \n", dur, dur.Seconds()/float64(LoopsToCall))
}

func onecall() {
        conn, err := net.Dial("tcp", URL)
        if err != nil {
                panic(err)
        }

        bytes := make([]byte, BytesToWrite, BytesToWrite)

        conn.Write(bytes)
        n, err:=io.ReadFull(conn, bytes)
        defer conn.Close()
        if err != nil {
                panic(err)
        }

        log.Printf("n %v", n)
}
