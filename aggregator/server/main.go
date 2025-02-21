package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/DataExMachina-dev/demos/aggregator/server/rpcpb"
	"github.com/DataExMachina-dev/side-eye-go/sideeye"
)

var (
	grpcPort  = flag.Int("port", 6544, "Port to serve gRPC on.")
	debugPort = flag.Int("debug-port", 6545, "Port to serve debug on.")
	useLib    = flag.Bool("use-lib", false, "Use the side-eye-go library to connect to Side-Eye.")
)

func main() {
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *debugPort), nil))
	}()

	if *useLib {
		log.Println("Using the side-eye-go library.")
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stdout, os.Stdout, 100))

		logSideEyeError := func(err error) {
			log.Printf("[side-eye-go] %s", err)
		}
		logSideEyeInfo := func(format string, args ...interface{}) {
			log.Printf("[side-eye-go] "+format, args...)
		}
		if err := sideeye.Init(context.Background(),
			"aggregator-server",
			sideeye.WithErrorLogger(logSideEyeError),
			sideeye.WithInfoLogger(logSideEyeInfo),
		); err != nil {
			log.Fatalf("Failed to initialize side-eye-go: %v", err)
		} else {
			log.Println("Initialized side-eye-go.")
		}
	}

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
