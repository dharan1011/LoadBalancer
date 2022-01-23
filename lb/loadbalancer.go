package lb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LoadBalancer struct {
	port                  int
	id                    uint64
	serverPool            *ServerPoolManager
	loadBalancingStrategy LoadBalancingStrategy
}

func (lb *LoadBalancer) getNextBackendNode() *Backend {
	return lb.loadBalancingStrategy.GetBackend(lb.serverPool.GetServerPool())
}

// For every request handler run in go routine
func (lb *LoadBalancer) handler(w http.ResponseWriter, r *http.Request) {
	nextNode := lb.getNextBackendNode()
	if nextNode != nil {
		nextNode.ReveseProxy.ServeHTTP(w, r)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "service unavailable",
	})
}

func (lb *LoadBalancer) Run() {
	loadBalancer := http.Server{
		Addr:    fmt.Sprintf(":%d", lb.port),
		Handler: http.HandlerFunc(lb.handler),
	}
	go lb.serverPool.RunPoolHealthCheck()
	if err := loadBalancer.ListenAndServe(); err != nil {
		panic(err)
	} else {
		log.Printf("Load Balancer [ID=%d] started. Listening on port %d", lb.id, lb.port)
	}
}
