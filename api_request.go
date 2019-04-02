package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func httpGet(version string) {
	http.Get("http://localhost:8888/" + version + "/pi/10000")
}

func TestGet(version string) {
	var wg sync.WaitGroup

	wg.Add(1000)

	t1 := time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			httpGet(version)
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println(version, " use ", elapsed)
}

func main() {
	//TestGet("v3")
	//TestGet("vgrpc")
	TestGet("v4")
}
