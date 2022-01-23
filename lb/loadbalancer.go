package lb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LoadBalancer struct {
	Port                  int
	id                    uint64
	ServerPool            *ServerPoolManager
	LoadBalancingStrategy LoadBalancingStrategy
}

func (lb *LoadBalancer) getNextBackendNode() *Backend {
	return lb.LoadBalancingStrategy.GetBackend(lb.ServerPool.GetServerPool())
}

// For every request handler run in go routine
func (lb *LoadBalancer) handler(w http.ResponseWriter, r *http.Request) {
	nextNode := lb.getNextBackendNode()
	if nextNode != nil {
		nextNode.ReveseProxy.ServeHTTP(w, r)
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "service unavailable",
		})
	}
}

func (lb *LoadBalancer) Run() {
	loadBalancer := http.Server{
		Addr:    fmt.Sprintf(":%d", lb.Port),
		Handler: http.HandlerFunc(lb.handler),
	}
	go lb.ServerPool.RunPoolHealthCheck()
	go lb.ServerPool.SchedulePoolHealthCheck()
	log.Printf("Starting Load Balancer [ID=%d] started. Listening on port %d", lb.id, lb.Port)
	if err := loadBalancer.ListenAndServe(); err != nil {
		panic(err)
	}
}
