package cluster

import "fmt"

// ClusterResource encapsulates metadata for a Kafka cluster
type ClusterResource struct {
	clusterID string
}

// InitClusterResource create new instance of ClusterResource and return its pointer type
func InitClusterResource(clusterID string) *ClusterResource {
	return &ClusterResource{
		clusterID: clusterID,
	}
}

// GetClusterID return ClusterResource clusterID
func (clstRes *ClusterResource) GetClusterID() string {
	return clstRes.clusterID
}

// String return the string representation of ClusterResource
func (clstRes *ClusterResource) String() string {
	return fmt.Sprintf("ClusterResource(clusterID=%v)", clstRes.clusterID)
}
