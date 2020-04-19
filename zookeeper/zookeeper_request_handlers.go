package zookeeper

import (
	"AKFAK/proto/commonpb"
	"AKFAK/proto/zkmessagepb"
	"context"
	"errors"
	"log"
)

// GetClusterMetadata return the current cluster state stored in the ZK
func (zk *Zookeeper) GetClusterMetadata(ctx context.Context, req *zkmessagepb.GetClusterMetadataRequest) (*zkmessagepb.GetClusterMetadataResponse, error) {
	zk.mux.Lock()
	// check if the controller has set
	controllerSet := zk.clusterMetadata.GetController().GetID() != -1
	if !controllerSet {
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

	if controllerSet {
		// update controller
		log.Printf("ZK update controller for new Broker %v\n", reqBrk.GetID())
		go zk.updateControllerMetadata()
	}

	// return response
	return &zkmessagepb.GetClusterMetadataResponse{
		ClusterInfo: zk.clusterMetadata.MetadataCluster,
	}, nil
}

// UpdateClusterMetadata update the ZK local cluster cache and flush to disk to persist the state
func (zk *Zookeeper) UpdateClusterMetadata(ctx context.Context, req *zkmessagepb.UpdateClusterMetadataRequest) (*zkmessagepb.UpdateClusterMetadataResponse, error) {
	newClsInfo := req.GetNewClusterInfo()

	// flush new cluster state into disk
	err := WriteClusterStateToFile(zk.config.DataDir, *newClsInfo)
	if err != nil {
		return &zkmessagepb.UpdateClusterMetadataResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_FAIL}}, errors.New("Fail to flush data to disk")
	}

	// update local cluster metadata cache
	zk.clusterMetadata.UpdateClusterMetadata(newClsInfo)

	return &zkmessagepb.UpdateClusterMetadataResponse{Response: &commonpb.Response{Status: commonpb.ResponseStatus_SUCCESS, Message: "Successfully updated the cluster metadata to ZK"}}, nil
}
