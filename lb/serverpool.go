package lb

import (
	"log"
	"net/url"
	"time"
)

type ServerPoolManager struct {
	backends []*Backend
}

func (sp *ServerPoolManager) AddBackend(backend *Backend) {
	sp.backends = append(sp.backends, backend)
}

func (sp *ServerPoolManager) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range sp.backends {
		if b.URL == backendUrl {
			b.SetAlive(alive)
			break
		}
	}
}

func (sp *ServerPoolManager) GetServerPool() []*Backend {
	return sp.backends
}

func (sp *ServerPoolManager) poolHealthCheck() {
	log.Println("Starting Backend Pool check")
	for _, backend := range sp.backends {
		if backend == nil {
			continue
		}
		isBacknedAlive := backend.IsAlive()
		backend.SetAlive(isBacknedAlive)
		var status string
		if isBacknedAlive {
			status = "UP"
		} else {
			status = "DOWN"
		}
		log.Printf("%s [%s]\n", backend.URL, status)
	}
	log.Println("Backend Pool Check Completed")
}

func (sp *ServerPoolManager) RunPoolHealthCheck() {
	t := time.NewTicker(time.Minute * 2)
	for {
		select {
		case <-t.C:
			sp.poolHealthCheck()
		}
	}
}
