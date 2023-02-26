# h2cli
Simple http 2.0 (h2) client in go including h2c (http2 over clear text, not requiring TLS)


## Installation
```shell
go install github.com/fortio/h2cli@latest
```

## Running

Options/flags:
```
h2cli dev usage:
	h2cli [flags]
flags:
  -cacert Path
    	Path to a custom CA certificate file instead standard internet/system CAs
  -h2
    	use HTTP2 (default true)
  -loglevel level
    	log level, one of [Debug Verbose Info Warning Error Critical Fatal] (default Info)
  -method string
    	HTTP method to use (default "GET")
  -quiet
    	Quiet mode, sets log level to warning
  -url string
    	URL to fetch (default "https://debug.fortio.org")
```

### Internet example

Standard

```
$ go run . -url https://debug.fortio.org
12:50:58 I GET on https://debug.fortio.org
12:50:58 I Response code 200, proto HTTP/2.0, size 366
Φορτίο version 1.50.1 h1:5FSttAHQsyAsi3dzxDmSByfzDYByrWY/yw53bqOg+Kc= go1.19.6 amd64 linux (in fortio.org/proxy 1.10.0)
Debug server on ol1 up for 15h38m51.4s
Request from [2600:1700:xxx]:56772 https TLS_AES_128_GCM_SHA256

GET / HTTP/2.0

headers:

Host: debug.fortio.org
Accept-Encoding: gzip
User-Agent: Go-http-client/2.0

body:

```
or
```shell
$ h2cli -url https://httpbin.org/get
12:52:03 I GET on https://httpbin.org/get
12:52:04 I Response code 200, proto HTTP/2.0, size 270
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

### H2C examples

h2c against debug.fortio.org is the new default for `go run .`
```
% go run .
15:28:59 I h2c GET on http://debug.fortio.org
15:29:00 I Response code 200, proto HTTP/2.0, size 334
Φορτίο version 1.52.0 h1:xHVOXkR3k5V5DvVM7/byfIHff3ia613qunnm+7O0EuQ= go1.19.6 arm64 linux (in fortio.org/proxy 1.11.1)
Debug server on a1 up for 19h4m25s
Request from [2600:1700:xxx]:52833

GET / HTTP/2.0

headers:

Host: debug.fortio.org
Accept-Encoding: gzip
User-Agent: Go-http-client/2.0

body:
```

```
% go run . -url localhost:8080/debug
15:31:16 I h2c GET on http://localhost:8080/debug
15:31:16 I Response code 200, proto HTTP/2.0, size 245
Φορτίο version dev  go1.19.6 arm64 darwin echo debug server up for 26m39.3s on MacBook-Air.local - request from [::1]:53010

GET /debug HTTP/2.0

headers:

Host: localhost:8080
Accept-Encoding: gzip
User-Agent: Go-http-client/2.0

body:
```
