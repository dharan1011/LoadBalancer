package lb

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL         *url.URL
	alive       bool
	mutex       sync.RWMutex
	ReveseProxy *httputil.ReverseProxy
}

func (b *Backend) SetAlive(alive bool) {
	b.mutex.Lock()
	b.alive = alive
	b.mutex.Unlock()
}

func (b *Backend) IsAlive() (alive bool) {
	b.mutex.RLock()
	alive = b.alive
	b.mutex.RUnlock()
	return
}

func (b *Backend) HealthCheck() bool {
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		log.Printf("Host unreachable : %s\n", b.URL.Host)
		return false
	}
	defer conn.Close()
	return true
}
