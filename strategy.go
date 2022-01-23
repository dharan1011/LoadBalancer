package main

import "sync/atomic"

type LoadBalancingStrategy interface {
	GetBackend(backends []*Backend) *Backend
}

// Round Robing Algorithm
type RoundRobin struct {
	currentIdx uint64
}

func (rr *RoundRobin) NextIndex(backends []*Backend) int {
	return int(atomic.AddUint64(&rr.currentIdx, 1) % uint64(len(backends)))
}

func (rr *RoundRobin) GetBackend(backends []*Backend) *Backend {
	next := rr.NextIndex(backends)
	length := next + len(backends)
	for i := next; i < length; i++ {
		idx := i % len(backends)
		if backends[idx].IsAlive() {
			atomic.StoreUint64(&rr.currentIdx, uint64(idx))
			return backends[i]
		}
	}
	return nil
}

// END Round Robin Algorithm
