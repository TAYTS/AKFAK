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

	// add the requesting broker to live brokers
	exist := false
	reqBrk := req.GetBroker()
	for _, brk := range zk.clusterMetadata.GetLiveBrokers() {
		if brk.ID == reqBrk.GetID() && brk.Host == reqBrk.GetHost() && brk.Port == reqBrk.GetPort() {
			exist = true
			break
		}
	}
	if !exist {
		zk.clusterMetadata.LiveBrokers = append(zk.clusterMetadata.LiveBrokers, reqBrk)
	}
	zk.mux.Unlock()

	// return response
	return &zkmessagepb.GetClusterMetadataResponse{
		ClusterInfo: &zk.clusterMetadata,
	}, nil
}
