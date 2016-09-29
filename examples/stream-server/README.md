# Purpose
This directory contains a project that shows how to generate a server with `go-swagger` that can return a stream of newline-delimited JSON bodies.

# How To Use This Project
## Build and Run
(All following instructuctions are to be run in the directory paralle to this file.)

1. Generate the code: `$ swagger generate server -f swagger.yml`
2. Install the project: `$ go install ./...`
3. Run the server: `$ $GOPATH/bin/countdown-server --port=8000`

## See the streaming output
In another terminal window, request some streaming output:
```
$ curl -v http://127.0.0.1:8000/elapse/5
* About to connect() to 127.0.0.1 port 8000 (#0)
*   Trying 127.0.0.1...
* Adding handle: conn: 0x7fdd8400a600
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x7fdd8400a600) send_pipe: 1, recv_pipe: 0
* Connected to 127.0.0.1 (127.0.0.1) port 8000 (#0)
> GET /elapse/5 HTTP/1.1
> User-Agent: curl/7.30.0
> Host: 127.0.0.1:8000
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Sun, 11 Sep 2016 00:54:34 GMT
< Transfer-Encoding: chunked
<
{"remains":5}
{"remains":4}
{"remains":3}
{"remains":2}
{"remains":1}
{"remains":0}
* Connection #0 to host 127.0.0.1 left intact
$
```
## See an error condition
Also in another terminal window, see an error message (not streaming):
```
$ curl -v http://127.0.0.1:8000/elapse/11
* About to connect() to 127.0.0.1 port 8000 (#0)
*   Trying 127.0.0.1...
* Adding handle: conn: 0x7f8582004000
* Adding handle: send: 0
* Adding handle: recv: 0
* Curl_addHandleToPipeline: length: 1
* - Conn 0 (0x7f8582004000) send_pipe: 1, recv_pipe: 0
* Connected to 127.0.0.1 (127.0.0.1) port 8000 (#0)
> GET /elapse/11 HTTP/1.1
> User-Agent: curl/7.30.0
> Host: 127.0.0.1:8000
> Accept: */*
>
< HTTP/1.1 403 Forbidden
< Content-Type: application/json
< Date: Sun, 11 Sep 2016 00:54:48 GMT
< Content-Length: 0
<
* Connection #0 to host 127.0.0.1 left intact
$
```
