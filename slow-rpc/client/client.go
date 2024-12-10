package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DataExMachina-dev/demos/slow-rpc/server/rpcpb"
)

var (
	serverAddr = flag.String("server-addr", "localhost:6543", "The server address in as host:port")
)

func main() {
	const numClients = 100
	log.Printf("Connecting to server at %s", *serverAddr)
	rpcClient, err := dialServer(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to dial server: %s", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			clientID := int32(i)
			runClient(clientID, rpcClient)
		}()
	}
	wg.Wait()
}

func runClient(clientID int32, rpcClient rpcpb.SideEyeDemoClient) {
	for {
		now := timestamppb.Now()
		_, err := rpcClient.GetInfo(context.Background(), &rpcpb.GetInfoRequest{
			ClientID:         clientID,
			RequestTimestamp: now,
		})
		if err != nil {
			log.Printf("RPC failed for clientID %d: %s", clientID, err)
			time.Sleep(time.Second)
		} else {
			duration := time.Since(now.AsTime())
			if duration > time.Second {
				log.Printf("RPC took a long time for clientID %d: %s", clientID, duration)
			}
		}
	}
}

func dialServer(addr string) (rpcpb.SideEyeDemoClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return rpcpb.NewSideEyeDemoClient(conn), nil
}
