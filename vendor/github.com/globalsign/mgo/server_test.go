// mgo - MongoDB driver for Go
//
// Copyright (c) 2018 Canonical Ltd
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package mgo_test

import (
	"time"

	"github.com/globalsign/mgo"
	. "gopkg.in/check.v1"
)

func (s *S) TestServerRecoversFromAbend(c *C) {
	session, err := mgo.Dial("localhost:40001")
	c.Assert(err, IsNil)
	defer session.Close()
	// Peek behind the scenes
	cluster := session.Cluster()
	server := cluster.Server("127.0.0.1:40001")

	info := &mgo.DialInfo{
		Timeout: time.Second,
		PoolLimit: 100,
	}
	
	sock, abended, err := server.AcquireSocket(info)
	c.Assert(err, IsNil)
	c.Assert(sock, NotNil)
	sock.Release()
	c.Check(abended, Equals, false)
	// Forcefully abend this socket
	sock.Close()
	server.AbendSocket(sock)
	// Next acquire notices the connection was abnormally ended
	sock, abended, err = server.AcquireSocket(info)
	c.Assert(err, IsNil)
	sock.Release()
	c.Check(abended, Equals, true)
	// cluster.AcquireSocketWithPoolTimeout should fix the abended problems
	sock, err = cluster.AcquireSocketWithPoolTimeout(mgo.Primary, false, time.Minute, nil, info)
	c.Assert(err, IsNil)
	sock.Release()
	sock, abended, err = server.AcquireSocket(info)
	c.Assert(err, IsNil)
	c.Check(abended, Equals, false)
	sock.Release()
}
