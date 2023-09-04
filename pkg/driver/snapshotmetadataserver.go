package driver

import (
	"log"

	csigrpc "github.com/PrasadG193/external-snapshot-metadata/pkg/grpc"
)

type SampleDriver struct {
	endpoint string
}

func NewSampleDriver(endpoint string) *SampleDriver {
	return &SampleDriver{endpoint: endpoint}
}

func (sd *SampleDriver) GetDelta(req *csigrpc.GetDeltaRequest, stream csigrpc.SnapshotMetadata_GetDeltaServer) error {
	log.Println("Received request::", req.String())
	// Generate and send fake data
	resp := csigrpc.GetDeltaResponse{
		BlockMetadataType: csigrpc.BlockMetadataType_FIXED_LENGTH,
		VolumeSizeBytes:   1024 * 1024 * 1024,
		BlockMetadata: []*csigrpc.BlockMetadata{
			&csigrpc.BlockMetadata{
				SizeBytes: 1024 * 1024,
			},
			&csigrpc.BlockMetadata{
				SizeBytes: 1024 * 1024,
			},
		},
	}
	for i := 1; i <= 20; i++ {
		resp.BlockMetadata[0].ByteOffset = uint64(i)
		resp.BlockMetadata[1].ByteOffset = uint64(i + 1)
		i++
		log.Println("Sending response to external-snap-session-svc")
		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	log.Println("End of the session")
	return nil
}

func (sd *SampleDriver) GetAllocated(req *csigrpc.GetAllocatedRequest, stream csigrpc.SnapshotMetadata_GetAllocatedServer) error {
	return nil
}

func (sd *SampleDriver) Run() error {
	s := NewServer()
	s.Start(sd.endpoint, sd)
	defer s.Stop()
	return nil
}
