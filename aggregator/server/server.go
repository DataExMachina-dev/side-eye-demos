package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DataExMachina-dev/demos/aggregator/server/rpcpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type clientRecords struct {
	mu struct {
		sync.Mutex
		records []int64
	}
}

type server struct {
	rpcpb.UnimplementedAggregatorServer

	clients sync.Map

	cancel atomic.Bool
	done   chan bool
}

// server implements the gRPC server interface.
var _ rpcpb.AggregatorServer = &server{}

func startServer() *server {
	s := &server{
		done: make(chan bool),
	}
	go s.compact()
	return s
}

func (s *server) stopServer() {
	s.cancel.Store(true)
	<-s.done
}

func (s *server) Record(
	ctx context.Context, request *rpcpb.RecordRequest,
) (*rpcpb.RecordResponse, error) {
	v, _ := s.clients.LoadOrStore(request.ClientId, &clientRecords{})
	cr := v.(*clientRecords)
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.mu.records = append(cr.mu.records, request.Records...)
	if len(cr.mu.records) > 40000 {
		cr.mu.records = cr.mu.records[:40000]
	}
	return &rpcpb.RecordResponse{}, nil
}

// Swaps median value to the first
func median(a, b, c *int64) {
	i0 := *a < *b
	i1 := *a < *c
	if i0 != i1 {
		return
	}
	i2 := *b < *c
	if i0 == i2 {
		*a, *b = *b, *a
	} else {
		*a, *c = *c, *a
	}
}

// Partitions using first number as a pivot
func partition(numbers []int64) (i int) {
	i = 1
	for j := 1; j < len(numbers); j++ {
		if numbers[j] < numbers[0] {
			numbers[i], numbers[j] = numbers[j], numbers[i]
			i++
		}
	}
	return
}

func sort(numbers []int64) {
	if len(numbers) < 2 {
		return
	}
	median(&numbers[0], &numbers[len(numbers)/2], &numbers[len(numbers)-1])
	p := partition(numbers)
	sort(numbers[:p])
	sort(numbers[p:])
}

func compactRecords(records []int64) {
	fmt.Printf("compacting %d records\n", len(records))
	now := timestamppb.Now()
	sort(records)
	duration := time.Since(now.AsTime())
	fmt.Printf("compaction took %s\n", duration)
}

func (s *server) compact() {
	fmt.Printf("compactor started\n")
	defer close(s.done)
	for {
		if s.cancel.Load() {
			return
		}
		s.clients.Range(func(k, v any) bool {
			if s.cancel.Load() {
				return false
			}
			cr := v.(*clientRecords)
			cr.mu.Lock()
			defer cr.mu.Unlock()
			if len(cr.mu.records) > 0 {
				compactRecords(cr.mu.records)
				cr.mu.records = nil
			}
			return true
		})
		time.Sleep(200 * time.Millisecond)
	}
}
