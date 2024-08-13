package driver

import (
	"log"

	"github.com/container-storage-interface/spec/lib/go/csi"
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

func (s *Server) Start(endpoint string, md csi.SnapshotMetadataServer, ids csi.IdentityServer) {
	listener, cleanup, err := Listen(endpoint)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	s.server = server

	s.cleanup = cleanup

	csi.RegisterSnapshotMetadataServer(server, md)
	csi.RegisterIdentityServer(server, ids)

	log.Printf("Listening for connections on address: %#v", listener.Addr())
	server.Serve(listener)
}
