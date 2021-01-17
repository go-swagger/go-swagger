# Streaming with go-swagger

## Purpose

This directory contains a project that shows how to generate with `go-swagger`:
1. a server that can return a stream of newline-delimited JSON bodies.
2. a client that can read this stream.

## Build and run a streaming server
(All following instructuctions are to be run from the directory parallel to this file.)

1. Generate the code: `$ swagger generate server -f swagger.yml`
2. Install the project: `$ go install ./...`
3. Run the server: `$ $GOPATH/bin/countdown-server --port=8000`

### See the streaming output
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
### See an error condition
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

## Build and run a streaming client

The client library in this folder has been generated with `swagger generate client -f swagger.yml --skip-models`
(assuming the models where built in the previous step).

A sample client using this library is provided here: `elapsed_client.go`.

This client reads asynchronously from the stream of json produced by the server above (from `http://localhost:8000`)
unmarshals and print the result. The client maintains the connection to receive chunks, for up to 7 seconds.

This program takes the start of the countdown as a command line argument.

### Try it

```
go run elapsed_client.go 5

2021/01/17 14:58:22 asking server for countdown timings: 5
2021/01/17 14:58:22 received countdown mark - raw: {"remains":5}
2021/01/17 14:58:22 received countdown mark - remaining: 5
2021/01/17 14:58:23 received countdown mark - raw: {"remains":4}
2021/01/17 14:58:23 received countdown mark - remaining: 4
2021/01/17 14:58:24 received countdown mark - raw: {"remains":3}
2021/01/17 14:58:24 received countdown mark - remaining: 3
2021/01/17 14:58:25 received countdown mark - raw: {"remains":2}
2021/01/17 14:58:25 received countdown mark - remaining: 2
2021/01/17 14:58:26 received countdown mark - raw: {"remains":1}
2021/01/17 14:58:26 received countdown mark - remaining: 1
2021/01/17 14:58:27 received countdown mark - raw: {"remains":0}
2021/01/17 14:58:27 received countdown mark - remaining: 0
2021/01/17 14:58:27 response complete
2021/01/17 14:58:27 EOF
```

```
go run elapsed_client.go 8

2021/01/17 14:58:31 asking server for countdown timings: 8
2021/01/17 14:58:31 received countdown mark - raw: {"remains":8}
2021/01/17 14:58:31 received countdown mark - remaining: 8
2021/01/17 14:58:32 received countdown mark - raw: {"remains":7}
2021/01/17 14:58:32 received countdown mark - remaining: 7
2021/01/17 14:58:33 received countdown mark - raw: {"remains":6}
2021/01/17 14:58:33 received countdown mark - remaining: 6
2021/01/17 14:58:34 received countdown mark - raw: {"remains":5}
2021/01/17 14:58:34 received countdown mark - remaining: 5
2021/01/17 14:58:35 received countdown mark - raw: {"remains":4}
2021/01/17 14:58:35 received countdown mark - remaining: 4
2021/01/17 14:58:36 received countdown mark - raw: {"remains":3}
2021/01/17 14:58:36 received countdown mark - remaining: 3
2021/01/17 14:58:37 received countdown mark - raw: {"remains":2}
2021/01/17 14:58:37 received countdown mark - remaining: 2
2021/01/17 14:58:38 got an error
2021/01/17 14:58:38 EOF
2021/01/17 14:58:38 failure: context deadline exceeded
```

### How does it work?

#### Setting the right consumer

First and foremost, we have to realize that the "application/json" mime is not really
describing our API. Rather, the server _streams_ chunks of individual JSON bits.

The runtime does not automatically detect that fact, we need to override this, like so:
```go
import (
  ...
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-swagger/go-swagger/examples/stream-server/client"
  ...
)

	customized := httptransport.New("localhost:8000", "/", []string{"http"})
	customized.Consumers[runtime.JSONMime] = runtime.ByteStreamConsumer()
	countdowns := client.New(customized, nil)
```

This tells the runtime to use a `ByteStreamConsumer` instead of a `JSONConsumer` when consuming a response.

#### Consuming asynchronously

The runtime consumer performs a "io.Copy()" call from the body to the writer passed by the request.

If we don't want to block until this is complete, we may pass an `io.PipeWriter` as the writer for this request. Like so:

```go
	reader, writer := io.Pipe()
  ...
	_, err := countdowns.Operations.Elapse(elapsed, writer)
```

The `reader` side of this pipe may be consuming by another go routine.

#### Unmarshalling the stream

The response is just a stream of byte, so the client has to unmarshal the messages received unitarily.
In this example, the stream is separated by line feed, so we can use a `bufio.Scanner` to do the job.

Notice the use of the `cancel()` method to interrupt the ongoing request if the go routine fails.

```go
    ...
		// read response items line by line
		for scanner.Scan() {
			// each response item is JSON
			txt := scanner.Text()
			log.Printf("received countdown mark - raw: %s", txt)

			var mark models.Mark

			err := json.Unmarshal([]byte(txt), &mark)
			if err != nil {
				log.Printf("unmarshal error: %v", err)
				return
			}

			log.Printf("received countdown mark - remaining: %d", swag.Int64Value(mark.Remains))
		}
    ...
```
