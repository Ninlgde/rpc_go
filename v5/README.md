## connection pool client

连接池客户端

```sh
 π ~/Desktop/source/thirdOpenSource ❯ wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/v5/pi/1000
Running 30s test @ http://127.0.0.1:8888/v5/pi/1000
  144 threads and 3000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.82ms    2.08ms 132.39ms   89.10%
    Req/Sec    10.73k     5.14k   28.68k    75.70%
  Latency Distribution
     50%    1.17ms
     75%    2.41ms
     90%    4.10ms
     99%    9.45ms
  759507 requests in 48.06s, 118.06MB read
Requests/sec:  15804.51
Transfer/sec:      2.46MB
```
比grpc---
```
π wrk master ❯ wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/vgrpc/pi/1000
Running 30s test @ http://127.0.0.1:8888/vgrpc/pi/1000
  144 threads and 3000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.03ms   10.73ms 549.00ms   92.34%
    Req/Sec     3.38k     3.93k   25.02k    91.42%
  Latency Distribution
     50%    2.65ms
     75%    8.38ms
     90%   14.90ms
     99%   35.70ms
  718402 requests in 47.12s, 111.67MB read
Requests/sec:  15245.03
Transfer/sec:      2.37MB
```