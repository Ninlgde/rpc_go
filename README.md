# Ninlgde Asyn rpc server
非常简单的分布式异步rpc server实现
依赖etcd做服务发现

## 安装
1. Download and install it:

```sh
$ go get -u github.com/Ninlgde/rpc_go
```

2. Import it in your code:

```go
import "github.com/Ninlgde/rpc_go/v3.0"
```

## Quick start
1. start etcd
    * docker 
    ```sh
    sudo docker run -p 2379:2379 -v /etc/ssl/certs/:/etc/ssl/certs/ elcolio/etcd
    ```
    
    * macos
    ```sh
    brew install etcd
    etcd
    ```
    
2. start server
    ```sh
    go run server_v3_main.go -port=8080
    go run server_v3_main.go -port=8081
    go run server_v3_main.go -port=8082
    ```
    
3. start http aip server
    ```sh
    go run http_api_server.go
    ```
    
4. go and test
    
    http://localhost:8888/ping/helloworld
    
    http://localhost:8888/pi/1000