package main

import (
	"context"
	"flag"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/DataExMachina-dev/demos/aggregator/server/rpcpb"
	_ "net/http/pprof"
)

var (
	serverAddr = flag.String("server-addr", "localhost:6544", "The server address in as host:port")
)

func main() {
	const numClients = 5
	log.Printf("Connecting to server at %s", *serverAddr)
	rpcClient, err := dialServer(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to dial server: %s", err)
	}

	go http.ListenAndServe(":8081", nil)

	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			clientID := int64(i)
			runClient(clientID, rpcClient)
		}()
	}
	wg.Wait()
}

func runClient(clientID int64, rpcClient rpcpb.AggregatorClient) {
	records := make([]int64, 10000)
	for {
		if clientID != 3 {
			for i, _ := range records {
				records[i] = rand.Int64()
			}
		}
		_, err := rpcClient.Record(context.Background(), &rpcpb.RecordRequest{
			ClientId: clientID,
			Records:  records,
		})
		if err != nil {
			log.Printf("RPC failed: %s", err)
		}
		time.Sleep(40 * time.Millisecond)
	}
}

func dialServer(addr string) (rpcpb.AggregatorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return rpcpb.NewAggregatorClient(conn), nil
}
