// +build ignore

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-swagger/go-swagger/fixtures/bugs/883/gen-fixture-883/client"
	"github.com/go-swagger/go-swagger/fixtures/bugs/883/gen-fixture-883/client/operations"
)

// Buffer knows how to UnmarshalText
type Buffer struct {
	*bytes.Buffer
}

// UnmarshalText handles text/plain
func (b *Buffer) UnmarshalText(text []byte) error {
	_, err := b.Write(text)
	return err
}

// NewBuffer creates a new buffer that knows how to unmarshal text/plain
func NewBuffer() *Buffer {
	return &Buffer{
		Buffer: bytes.NewBuffer(nil),
	}
}

func main() {

	blockingMode := true
	chunkingMode := true

	if len(os.Args) > 1 {
		for _, arg := range os.Args {
			switch arg {
			case "nonblocking":
				blockingMode = false
			case "nochunking":
				chunkingMode = false
			}
		}
	}

	log.Printf("asking jigsaw server with blockingMode: %t, chunkingMode: %t", blockingMode, chunkingMode)

	if blockingMode {
		if err := chunkedBlocking(chunkingMode); err != nil {
			log.Fatalf("error: %v", err)
		}
		return
	}

	if err := chunkedNonBlocking(chunkingMode); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func customTransport(withChunks bool) *httptransport.Runtime {
	// calling htts://jigsaw.w3.org/
	transport := httptransport.New("jigsaw.w3.org", "/", []string{"https"})

	if withChunks {
		// the jigsaw API enables chunks for http 1.1 clients.
		// No chunking takes place with http 2.0
		http1Only := http.DefaultTransport.(*http.Transport)
		// this disables http 2.0
		http1Only.TLSNextProto = map[string]func(authority string, c *tls.Conn) http.RoundTripper{}
		transport.Transport = http1Only
	}

	return transport
}

// chunkedBlocking consumes some text/plain resource, blocking for the response to be completely sent
func chunkedBlocking(withChunks bool) error {

	c := client.New(customTransport(withChunks), nil).Operations

	// we just need to specify a buffer that knows how to UnmarshalText()
	buf := NewBuffer()
	_, err := c.Chunked(operations.NewChunkedParams(), buf)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}

	log.Printf("result: %v", string(data))
	return nil
}

// chunkedNonBlocking consumes some text/plain resource asynchronously
func chunkedNonBlocking(withChunks bool) error {
	transport := customTransport(withChunks)

	// override text/plain consumer to consume this as a stream
	transport.Consumers[runtime.TextMime] = runtime.ByteStreamConsumer()

	c := client.New(transport, nil).Operations

	reader, writer := io.Pipe()

	scanner := bufio.NewScanner(reader)

	ctx, cancel := context.WithCancel(context.Background())

	// consumes asynchronously the response buffer
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer cancel()

		line := 1
		// read response items line by line
		for scanner.Scan() {
			// each response item is JSON
			txt := scanner.Text()
			log.Printf("line[%d]: %s", line, txt)
			line++
		}

		if err := scanner.Err(); err != nil {
			log.Printf("scanner err: %v", err)
		}

		log.Println("EOF")
	}(&wg)

	_, err := c.Chunked(operations.NewChunkedParamsWithContext(ctx), writer)

	if err == nil {
		log.Printf("response complete")
	} else {
		log.Printf("got an error")
	}

	_ = writer.Close()

	wg.Wait()
	return err
}
