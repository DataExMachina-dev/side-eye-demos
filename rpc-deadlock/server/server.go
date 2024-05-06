package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/durationpb"
	"log"
	"maps"
	"sync"
	"time"

	"github.com/DataExMachina-dev/demos/rpc-deadlock/server/rpcpb"
)

type clientState struct {
	clientID int32
	numCalls int
	mu       struct {
		sync.Mutex
		processing                     bool
		lastBackgroundServiceStartTime time.Time
	}
}

type server struct {
	rpcpb.UnimplementedSideEyeDemoServer
	mu struct {
		sync.Mutex
		clientStates map[int32]*clientState
	}
}

// server implements the gRPC server interface.
var _ rpcpb.SideEyeDemoServer = &server{}

func startServer() *server {
	s := &server{}
	s.mu.clientStates = make(map[int32]*clientState)
	go s.backgroundProcessor()
	return s
}

func (s *server) GetInfo(
	ctx context.Context, request *rpcpb.GetInfoRequest,
) (*rpcpb.GetInfoResponse, error) {
	cs := s.ensureClientState(request.ClientID)
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.numCalls++
	// Simulate a short processing time.
	time.Sleep(10 * time.Millisecond)
	return &rpcpb.GetInfoResponse{
		Result:         fmt.Sprintf("num_calls: %d", cs.numCalls),
		ProcessingTime: durationpb.New(time.Since(request.RequestTimestamp.AsTime())),
	}, nil
}

func (s *server) ensureClientState(clientID int32) *clientState {
	s.mu.Lock()
	defer s.mu.Unlock()
	cs, ok := s.mu.clientStates[clientID]
	if ok {
		return cs
	}
	cs = &clientState{clientID: clientID}
	s.mu.clientStates[clientID] = cs
	return cs
}

func (s *server) backgroundProcessor() {
	// Process the entries in s.mu.clientStates over and over. Each iteration goes
	// through every entry.
	for {
		// Take a snapshot of all clients currently registered while holding the
		// server's lock.
		s.mu.Lock()
		clientStates := maps.Clone(s.mu.clientStates)
		s.mu.Unlock()

		// Go through each client and process it. The client's lock is held while
		// processing.
		for _, c := range clientStates {
			s.processClient(c)
		}
		// Don't spin too hot case there are no clients.
		if len(clientStates) == 0 {
			time.Sleep(time.Second)
		}
	}
}

func (s *server) processClient(cs *clientState) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	log.Printf("processing client %d", cs.clientID)
	processingStart := time.Now()
	cs.mu.processing = true
	cs.mu.lastBackgroundServiceStartTime = time.Now()

	// Simulate a lengthy operation.
	time.Sleep(time.Minute)

	cs.mu.processing = false
	log.Printf("finished processing client %d; processing took: %s",
		cs.clientID, time.Since(processingStart))
}
