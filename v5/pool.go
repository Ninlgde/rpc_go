package v5_0

import (
	"fmt"
	"net"
	"sync"
)

type Pool struct {
	min     int
	max     int
	cap     int
	conns   []net.Conn
	busy    map[net.Conn]struct{}
	address string
	sync.Mutex
}

type maxError struct {
	error
}

func NewPool(address string, min, max int) (pool *Pool) {
	conns := make([]net.Conn, min)
	cap := 0
	for cap < min {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Fail to create conn ", err)
			continue
		}
		conns[cap] = conn
		cap++
	}
	pool = &Pool{
		min:     min,
		max:     max,
		cap:     min,
		conns:   conns,
		busy:    make(map[net.Conn]struct{}),
		address: address,
	}
	return
}

func DisposePool(pool *Pool) {
	pool.Lock()
	defer pool.Unlock()
	// 将池里的链接全部销毁，外面的
	for i := range pool.conns {
		pool.conns[i].Close()
	}
	pool.conns = nil

	for k := range pool.busy {
		k.Close()
	}
	pool.busy = nil
}

func (pool *Pool) Get() (conn net.Conn, err error) {
	pool.Lock()
	defer pool.Unlock()
	if len(pool.conns) > 0 {
		conn = pool.conns[0]
		pool.conns = pool.conns[1:]
		pool.busy[conn] = struct{}{}
		return
	}
	if pool.cap < pool.max {
		conn, err = pool.newConn()
		if err == nil {
			pool.cap++
		}
		return
	}
	return nil, maxError{}
}

func (pool *Pool) Put(conn net.Conn) {
	pool.Lock()
	defer pool.Unlock()
	delete(pool.busy, conn)
	pool.conns = append(pool.conns, conn)
}

func (pool *Pool) newConn() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", pool.address)
	if err != nil {
		fmt.Println("Fail to create conn ", err)
		return
	}
	return
}
