package v2_0

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func newClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}
	defer conn.Close()
	client := Client{conn}
	for i := 0; i < 10; i++ {
		client.Rpc("ping", "ireader "+strconv.Itoa(i))
		//fmt.Println(out, result)
		//time.Sleep(time.Second * 1)
	}
}

func BenchmarkClient(b *testing.B) {
	t1 := time.Now()
	var wg sync.WaitGroup

	wg.Add(500)
	for i := 0; i < 500; i++ {
		go func(i int) {
			newClient()
			wg.Done()
			//fmt.Println(i)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(t1)
	fmt.Println("all finished, use ", elapsed)
}
