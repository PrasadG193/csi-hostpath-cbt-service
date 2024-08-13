package driver

import (
	"context"
	"log"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type SampleDriver struct {
	csi.UnimplementedSnapshotMetadataServer
	csi.UnimplementedIdentityServer
	endpoint string
}

func NewSampleDriver(endpoint string) *SampleDriver {
	return &SampleDriver{endpoint: endpoint}
}

func (sd *SampleDriver) GetMetadataDelta(req *csi.GetMetadataDeltaRequest, stream csi.SnapshotMetadata_GetMetadataDeltaServer) error {
	log.Println("Received request::", req.String())
	// Generate and send fake data
	resp := csi.GetMetadataDeltaResponse{
		BlockMetadataType:   csi.BlockMetadataType_FIXED_LENGTH,
		VolumeCapacityBytes: 1024 * 1024 * 1024,
		BlockMetadata: []*csi.BlockMetadata{
			&csi.BlockMetadata{
				SizeBytes: 1024 * 1024,
			},
			&csi.BlockMetadata{
				SizeBytes: 1024 * 1024,
			},
		},
	}
	for i := 1; i <= 20; i++ {
		resp.BlockMetadata[0].ByteOffset = int64(i)
		resp.BlockMetadata[1].ByteOffset = int64(i + 1)
		i++
		log.Println("Sending response to external-snap-session-svc")
		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	log.Println("End of the session")
	return nil
}

func (sd *SampleDriver) GetMetadataAllocated(req *csi.GetMetadataAllocatedRequest, stream csi.SnapshotMetadata_GetMetadataAllocatedServer) error {
	return nil
}

func (sd *SampleDriver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	log.Println("Using default capabilities")
	caps := []*csi.PluginCapability{
		{
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: csi.PluginCapability_Service_SNAPSHOT_METADATA_SERVICE,
				},
			},
		},
	}
	return &csi.GetPluginCapabilitiesResponse{Capabilities: caps}, nil
}

func (sd *SampleDriver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	log.Println("Using default GetPluginInfo")
	return &csi.GetPluginInfoResponse{
		Name: os.Getenv("DRIVER_NAME"),
	}, nil
}

func (sd *SampleDriver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}

func (sd *SampleDriver) Run() error {
	s := NewServer()
	s.Start(sd.endpoint, sd, sd)
	defer s.Stop()
	return nil
}
