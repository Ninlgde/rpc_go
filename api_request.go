package main

import (
	"net/http"
	"sync"
)

func httpGet() {
	resp, err := http.Get("http://localhost:8888/pi/10000")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	handle error
	//}

	//fmt.Println(string(body))
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			httpGet()
			wg.Done()
		}()
	}

	wg.Wait()
}
