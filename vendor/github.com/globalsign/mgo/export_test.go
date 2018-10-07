package mgo

import (
	"net"
	"time"
)

func HackPingDelay(newDelay time.Duration) (restore func()) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	oldDelay := pingDelay
	restore = func() {
		globalMutex.Lock()
		pingDelay = oldDelay
		globalMutex.Unlock()
	}
	pingDelay = newDelay
	return
}

func (s *Session) Cluster() *mongoCluster {
	return s.cluster()
}

func (cluster *mongoCluster) Server(addr string) *mongoServer {
	tcpaddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	return cluster.server(addr, tcpaddr)
}
