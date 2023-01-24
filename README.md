# h2cli
Simple http 2.0 (h2) client in go


## Installation
```shell
go install github.com/fortio/h2cli@latest
```

## Running

Options/flags:
```shell
$ h2cli -h
Usage of h2cli:
  -cacert Path
    	Path to a custom CA certificate file instead standard internet/system CAs
  -h2
    	use HTTP2 (default true)
  -url string
    	URL to fetch (default "https://localhost:8080/debug")
```

### Internet example

Standard
```shell
$ h2cli -url https://httpbin.org/get
12:57:47 Response code 200, proto HTTP/2.0
```
```json
{
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/2.0", 
    "X-Amzn-Trace-Id": "Root=1-63d0464b-583c8c4a5b5d0a7f2f2e585a"
  }, 
  "origin": "xxx.xxx.xxx.xxx", 
  "url": "https://httpbin.org/get"
}
```

### Fortio example

With
```
fortio server -cert cert-tmp/server.crt -key cert-tmp/server.key &
h2cli -cacert cert-tmp/ca.crt
```
gets (in stderr)
```
13:02:09 Debug: GET /debug HTTP/2.0 [::1]:54370 ()  "Go-http-client/2.0"
13:02:09 Response code 200, proto HTTP/2.0
```
stdout
```
Φορτίο version dev  go1.19.5 arm64 darwin echo debug server up for 11.8s on K6922C3FGJ - request from [::1]:54370

GET /debug HTTP/2.0

headers:

Host: localhost:8080
Accept-Encoding: gzip
User-Agent: Go-http-client/2.0

body:

```
