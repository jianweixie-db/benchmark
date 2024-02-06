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
var byteCount = 2_000_000 // 2907082
func main() {
        // Single thread http server
        var mu sync.Mutex
        bytes := make([]byte, byteCount, byteCount)
        for i :=0; i< byteCount; i++ {
                bytes[i] = 'a'
        }
        http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
                mu.Lock()
                defer mu.Unlock()
                //time.Sleep(sleepDuration)
                n, err := w.Write(bytes)
                fmt.Printf("n %v err %v\n", n, err)
                //fmt.Fprintf(w, "%v", bytes)
        })
        log.Fatal(http.ListenAndServe(":80", nil))
}


// package main
// import (
//         "fmt"
//         "log"
//         "net/http"
//         "sync"
// )
// var byteCount = 1_000_000 // 2907082
// func main() {
//         // Single thread http server
//         var mu sync.Mutex
//         bytes := make([]byte, byteCount, byteCount)
//         for i :=0; i< byteCount; i++ {
//                 bytes[i] = 'a'
//         }
//         http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
//                 w.Header().Set("Transfer-Encoding", "identity")
//                 w.Header().Set("Content-Length", fmt.Sprintf("%d", byteCount))
//                 mu.Lock()
//                 defer mu.Unlock()
//                 inSize := r.ContentLength
//                 outSize := 0
//                 for {
//                         n, err := w.Write(bytes)
//                         if err != nil {
//                                 panic(err)
//                         }
//                         outSize += n
//                         if outSize >= byteCount {
//                                 break
//                         }
//                 }
//                 fmt.Printf("in %v out %v\n", inSize, outSize)
//         })
//         log.Fatal(http.ListenAndServe(":80", nil))
// }
