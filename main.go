package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/dharan1011/LoadBalancer/lb"
)

func main() {
	url, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("URL : ", url.Host)
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
