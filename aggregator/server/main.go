package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/DataExMachina-dev/demos/aggregator/server/rpcpb"

	"net/http"
	_ "net/http/pprof"
)

var (
	grpcPort  = flag.Int("port", 6544, "Port to serve gRPC on.")
	debugPort = flag.Int("debug-port", 6545, "Port to serve debug on.")
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *debugPort), nil))
	}()

	s := startServer()
	defer s.stopServer()
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&rpcpb.Aggregator_ServiceDesc, s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", *grpcPort, err)
	}
	fmt.Printf("Serving gRPC on port %d.\n", *grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("grpc service exited: %v", err)
	}
}
