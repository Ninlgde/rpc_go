## connection pool client

连接池客户端

```sh
π wrk master ❯ wrk -t144 -c3000 -d30s -T30s --latency http://127.0.0.1:8888/v5/pi/1000
Running 30s test @ http://127.0.0.1:8888/v5/pi/1000
  144 threads and 3000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.86ms    1.61ms  65.59ms   87.90%
    Req/Sec    10.54k     2.70k   28.20k    85.22%
  Latency Distribution
     50%    1.56ms
     75%    2.48ms
     90%    3.74ms
     99%    7.21ms
  575365 requests in 41.05s, 89.44MB read
Requests/sec:  14015.37
Transfer/sec:      2.18MB
```
比grpc略慢。。回头调优
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