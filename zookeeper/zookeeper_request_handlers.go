package zookeeper

import (
	"AKFAK/proto/zkmessagepb"
	"context"
)

// GetClusterMetadata return the current cluster state stored in the ZK
func (zk *Zookeeper) GetClusterMetadata(ctx context.Context, req *zkmessagepb.GetClusterMetadataRequest) (*zkmessagepb.GetClusterMetadataResponse, error) {
	zk.mux.Lock()
	// check if the controller has set
	if zk.clusterMetadata.GetController().GetID() == -1 {
		// set the requesting broker as the controller
		zk.clusterMetadata.Controller = req.GetBroker()
	}
	zk.mux.Unlock()

	// return response
	return &zkmessagepb.GetClusterMetadataResponse{
		ClusterInfo: &zk.clusterMetadata,
	}, nil
}
