module github.com/DataExMachina-dev/demos/aggregator/server

go 1.23

require (
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/DataExMachina-dev/side-eye-go v0.0.0-20250129201350-c4ff5931163d // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0 // indirect
)

replace github.com/DataExMachina-dev/side-eye-go => ../../../side-eye/side-eye-go
