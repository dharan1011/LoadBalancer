package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/dharan1011/LoadBalancer/lb"
)

func main() {
	url, err := url.Parse(fmt.Sprintf("%s:%d", "localhost", 8080))
	if err != nil {
		panic(err)
	}
	serverPool := lb.ServerPoolManager{}
	serverPool.AddBackend(&lb.Backend{
		URL:         url,
		ReveseProxy: httputil.NewSingleHostReverseProxy(url),
	})
	strategy := lb.RoundRobin{}
	lb := lb.LoadBalancer{
		Port:                  9000,
		ServerPool:            &serverPool,
		LoadBalancingStrategy: &strategy,
	}
	lb.Run()
}
