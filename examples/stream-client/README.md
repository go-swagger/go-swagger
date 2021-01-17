# Streaming client

This example demonstrates how to build a client from the generated client
to work against a streaming server.

Thanks to @nelz9999, who raised this in issue #883.


This example calls `https://jigsaw.w3c.org/` so as to contentrate on the client
side. For an example mixing the generated server and a client, see [here](../streaming-server/README.md).

The URL called by the API returns a blob of ASCII digits (1000 lines).

## Try it

### Blocking client

The client sends a request and blocks until the server has finished sending the response.

```
SWAGGER_DEBUG=1 go run jigsaw.go

2021/01/17 17:09:59 asking jigsaw server with blockingMode: true, chunkingMode: true
GET /HTTP/ChunkedScript HTTP/1.1
Host: jigsaw.w3.org
User-Agent: Go-http-client/1.1
Accept: application/json
Accept-Encoding: gzip


HTTP/1.1 200 OK
Transfer-Encoding: chunked
Cache-Control: max-age=0
Content-Type: text/plain
...
Server: Jigsaw/2.3.0-beta2
Strict-Transport-Security: max-age=15552015; includeSubDomains; preload
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block

2021/01/17 17:10:00 result: This output will be chunked encoded by the server, if your client is HTTP/1.1
Below this line, is 1000 repeated lines of 0-9.
-------------------------------------------------------------------------
01234567890123456789012345678901234567890123456789012345678901234567890
...
01234567890123456789012345678901234567890123456789012345678901234567890
```

```
# http2 client is detected by jigsaw: response is not chunked
SWAGGER_DEBUG=1 go run jigsaw.go nochunking

2021/01/17 17:10:12 asking jigsaw server with blockingMode: true, chunkingMode: false
GET /HTTP/ChunkedScript HTTP/1.1
Host: jigsaw.w3.org
User-Agent: Go-http-client/1.1
Accept: application/json
Accept-Encoding: gzip


HTTP/2.0 200 OK
Connection: close
Cache-Control: max-age=0
Content-Type: text/plain
...
Server: Jigsaw/2.3.0-beta3
Strict-Transport-Security: max-age=15552015; includeSubDomains; preload
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block
...

2021/01/17 17:10:12 result: This output will be chunked encoded by the server, if your client is HTTP/1.1
Below this line, is 1000 repeated lines of 0-9.
-------------------------------------------------------------------------
01234567890123456789012345678901234567890123456789012345678901234567890
...
01234567890123456789012345678901234567890123456789012345678901234567890
```

### Non-blocking client

The client sends a request and consumes the chunks as they are produced.

```
SWAGGER_DEBUG=1 go run jigsaw.go nonblocking

2021/01/17 17:09:23 asking jigsaw server with blockingMode: false, chunkingMode: true
GET /HTTP/ChunkedScript HTTP/1.1
Host: jigsaw.w3.org
User-Agent: Go-http-client/1.1
Accept: application/json
Accept-Encoding: gzip

HTTP/1.1 200 OK
Transfer-Encoding: chunked
Cache-Control: max-age=0
Content-Type: text/plain
...
Server: Jigsaw/2.3.0-beta4
Strict-Transport-Security: max-age=15552015; includeSubDomains; preload
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block

2021/01/17 17:09:23 line[1]: This output will be chunked encoded by the server, if your client is HTTP/1.1
2021/01/17 17:09:23 line[2]: Below this line, is 1000 repeated lines of 0-9.
2021/01/17 17:09:23 line[3]: -------------------------------------------------------------------------
2021/01/17 17:09:23 line[4]: 01234567890123456789012345678901234567890123456789012345678901234567890
2021/01/17 17:09:23 line[5]: 01234567890123456789012345678901234567890123456789012345678901234567890
...
2021/01/17 17:09:23 line[1003]: 01234567890123456789012345678901234567890123456789012345678901234567890
2021/01/17 17:09:23 EOF
```

```
SWAGGER_DEBUG=1 go run jigsaw.go nonblocking nochunking

2021/01/17 17:09:39 asking jigsaw server with blockingMode: false, chunkingMode: false
GET /HTTP/ChunkedScript HTTP/1.1
Host: jigsaw.w3.org
User-Agent: Go-http-client/1.1
Accept: application/json
Accept-Encoding: gzip

HTTP/2.0 200 OK
Connection: close
Cache-Control: max-age=0
Content-Type: text/plain
...
Server: Jigsaw/2.3.0-beta3
Strict-Transport-Security: max-age=15552015; includeSubDomains; preload
X-Frame-Options: deny
X-Xss-Protection: 1; mode=block

2021/01/17 17:09:39 line[1]: This output will be chunked encoded by the server, if your client is HTTP/1.1
2021/01/17 17:09:39 line[2]: Below this line, is 1000 repeated lines of 0-9.
2021/01/17 17:09:39 line[3]: -------------------------------------------------------------------------
2021/01/17 17:09:39 line[4]: 01234567890123456789012345678901234567890123456789012345678901234567890
2021/01/17 17:09:39 line[5]: 01234567890123456789012345678901234567890123456789012345678901234567890
...
2021/01/17 17:09:39 line[1002]: 01234567890123456789012345678901234567890123456789012345678901234567890
2021/01/17 17:09:39 line[1003]: 01234567890123456789012345678901234567890123456789012345678901234567890
2021/01/17 17:09:39 EOF
```

## How does it work?

### Blocking: customizing the `io.Writer`

We have specified a `binary` format for the response (instead of, say, `[]string`).

We need to provide the runtime consumer with some means to unmarshal `text/plain` in this buffer.

The simplest way is to equip the destination buffer with a `UnmarshalText()` method. Like so:
```go
type Buffer struct {
	*bytes.Buffer
}

// UnmarshalText handles text/plain
func (b *Buffer) UnmarshalText(text []byte) error {
	_, err := b.Write(text)
	return err
}

buf := NewBuffer()
_, err := c.Chunked(operations.NewChunkedParams(), buf)
```

### Non-blocking: customizing the consumer

We need to instruct the runtime client to use a `ByteStreamConsumer` when the response is `text/plain`.
Like so:

```go
	transport.Consumers[runtime.TextMime] = runtime.ByteStreamConsumer()
```

#### Non-blocking: unmarshalling the stream

The response is just a stream of byte, so the client has to unmarshal the messages received unitarily.
In this example, the stream is separated by line feed, so we can use a `bufio.Scanner` to do the job.

Notice the use of the `cancel()` method to interrupt the ongoing request if the go routine fails.
