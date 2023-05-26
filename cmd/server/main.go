package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pgrpc "github.com/PrasadG193/external-snapshot-session-service/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// unix socket
	PROTOCOL = "unix"
	SOCKET   = "/csi/csi.sock"
)

func main() {
	os.Remove(SOCKET)
	listener, err := net.Listen(PROTOCOL, SOCKET)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	s := grpc.NewServer()
	reflection.Register(s)
	pgrpc.RegisterSnapshotMetadataServer(s, &Server{})
	fmt.Println("SERVER STARTED!")
	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	pgrpc.UnimplementedSnapshotMetadataServer
}

func (s *Server) GetDelta(req *pgrpc.GetDeltaRequest, stream pgrpc.SnapshotMetadata_GetDeltaServer) error {
	fmt.Println("Received request::", req.String())
	resp := pgrpc.GetDeltaResponse{
		BlockMetadataType: pgrpc.BlockMetadataType_FIXED_LENGTH,
		VolumeSizeBytes:   1024 * 1024 * 1024,
		BlockMetadata: []*pgrpc.BlockMetadata{
			&pgrpc.BlockMetadata{
				ByteOffset: 1,
				SizeBytes:  1024 * 1024,
			},
		},
	}
	for i := 1; i <= 10; i++ {
		resp.BlockMetadata[0].ByteOffset = uint64(i)
		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	return nil
}