package driver

import (
	"log"

	csigrpc "github.com/PrasadG193/external-snapshot-metadata/pkg/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	cleanup func()
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Stop() {
	s.server.GracefulStop()
	s.cleanup()
}

func (s *Server) Start(endpoint string, md csigrpc.SnapshotMetadataServer) {
	listener, cleanup, err := Listen(endpoint)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	s.server = server

	s.cleanup = cleanup

	csigrpc.RegisterSnapshotMetadataServer(server, md)

	log.Printf("Listening for connections on address: %#v", listener.Addr())
	server.Serve(listener)
}
