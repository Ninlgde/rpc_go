# Ninlgde Asyn rpc server
非常简单的分布式异步rpc server实现
依赖etcd做服务发现

* 添加了grpc的实现
* 添加了v4客户端，类似stream？
* 添加了v5客户端，链接池

~~造轮子使我快乐~~

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
    go run server_vgrpc_main.go -port=18080
    go run server_vgrpc_main.go -port=18081
    go run server_vgrpc_main.go -port=18082
    ```
    
3. start http aip server
    ```sh
    go run http_api_server.go
    ```
    
4. go and test
    
    - http://localhost:8888/v5/ping/helloworld
    - http://localhost:8888/vgrpc/ping/helloworld
    - http://localhost:8888/v5/pi/1000
    - http://localhost:8888/vgrpc/pi/1000
    
5. benchmark

    1. wrk install
    
    ```git clone https://github.com/wg/wrk```
    
    ```make & make install```
    
    2. go and test
    
    ```
    wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/v4/pi/10000
    wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/v5/pi/10000
    wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/vgrpc/pi/10000
    ```
    
## Run on k8s(minikube)
run rpc server on minikube
fork from: https://github.com/tinrab/kubernetes-go-grpc-tutorial.git

1. build pb
   ```protoc -I . --go_out=plugins=grpc:. *.proto```
   
2. build docker file
   ```
   docker build -t local/service -f Dockerfile.service .
   docker build -t local/apis -f Dockerfile.apis .
   ```
   
3. apply yaml
   ```
   kubectl apply -f calc.yaml
   kubectl apply -f api.yaml
   ```

4. test
   ```
   curl $(minikube service api-service --url)/ping/hello
   curl $(minikube service api-service --url)/pi/20000
   curl $(minikube service api-service --url)/gcd/333/653
   ```
    
## TODO

- api server的错误码

- 灰度发布

- 心跳
