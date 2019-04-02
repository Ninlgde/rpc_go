package v3_0

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestClient_Rpc(t *testing.T) {
	client := NewClient()

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				client.Rpc("ping", "ireader "+strconv.Itoa(i))
				time.Sleep(200 * time.Millisecond)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestClient_RpcStream(t *testing.T) {
	client := NewClient()

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			conn := client.discover.FindAvailableConn()
			defer conn.Close()
			for j := 0; j < 10; j++ {
				client.RpcStream(conn, "ping", "ireader "+strconv.Itoa(i))
				time.Sleep(200 * time.Millisecond)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
