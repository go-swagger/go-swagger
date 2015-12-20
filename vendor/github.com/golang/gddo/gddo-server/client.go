// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// This file implements an http.Client with request timeouts set by command
// line flags.

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	dialTimeout    = flag.Duration("dial_timeout", 5*time.Second, "Timeout for dialing an HTTP connection.")
	requestTimeout = flag.Duration("request_timeout", 20*time.Second, "Time out for roundtripping an HTTP request.")
)

type timeoutConn struct {
	net.Conn
}

func (c timeoutConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	c.Conn.SetReadDeadline(time.Time{})
	return n, err
}

func timeoutDial(network, addr string) (net.Conn, error) {
	c, err := net.DialTimeout(network, addr, *dialTimeout)
	if err != nil {
		return c, err
	}
	// The net/http transport CancelRequest feature does not work until after
	// the TLS handshake is complete. To help catch hangs during the TLS
	// handshake, we set a deadline on the connection here and clear the
	// deadline when the first read on the connection completes. This is not
	// perfect, but it does catch the case where the server accepts and ignores
	// a connection.
	c.SetDeadline(time.Now().Add(*requestTimeout))
	return timeoutConn{c}, nil
}

type transport struct {
	t http.Transport
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	timer := time.AfterFunc(*requestTimeout, func() {
		t.t.CancelRequest(req)
		log.Printf("Canceled request for %s", req.URL)
	})
	defer timer.Stop()
	if req.URL.Host == "api.github.com" && gitHubCredentials != "" {
		if req.URL.RawQuery == "" {
			req.URL.RawQuery = gitHubCredentials
		} else {
			req.URL.RawQuery += "&" + gitHubCredentials
		}
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	return t.t.RoundTrip(req)
}

var httpClient = &http.Client{Transport: &transport{
	t: http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  timeoutDial,
		ResponseHeaderTimeout: *requestTimeout / 2,
	}}}
