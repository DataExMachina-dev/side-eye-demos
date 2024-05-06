package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/DataExMachina-dev/demos/slow-rpc/server/rpcpb"
)

var (
	grpcPort = flag.Int("port", 6543, "Port to serve gRPC on.")
)

func main() {
	s := startServer()
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&rpcpb.SideEyeDemo_ServiceDesc, s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %d for api_auth server: %v", *grpcPort, err)
	}
	fmt.Printf("Serving gRPC on port %d.\n", *grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("grpc service exited: %v", err)
	}
}
