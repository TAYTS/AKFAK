package consumermetadata

import (
	"AKFAK/proto/consumepb"
	"AKFAK/proto/consumermetadatapb"
	"sync"
)

// ConsumerMetadata is a wrapper to consumermetadatapb.MetadataConsumerState that provide
// thread safe for managing the consumer state metadata
type ConsumerMetadata struct {
	*consumermetadatapb.MetadataConsumerState
	brokerAssignmentsMap map[int][]*consumepb.MetadataAssignment
	mux                  sync.RWMutex
}

// InitConsumerMetadata create a thread safe wrapper for handling the consumer metadata
func InitConsumerMetadata(cmetadata *consumermetadatapb.MetadataConsumerState) *ConsumerMetadata {
	cstate := &ConsumerMetadata{MetadataConsumerState: cmetadata}
	cstate.populateMetadata()
	return cstate
}

// UpdateConsumerMetadata replace the current MetadataConsumerState local cache
func (cstate *ConsumerMetadata) UpdateConsumerMetadata(cnewstateMeta *consumermetadatapb.MetadataConsumerState) {
	cstate.mux.Lock()
	cstate.MetadataConsumerState = cnewstateMeta
	cstate.mux.Unlock()
	cstate.populateMetadata()
}


// GetOffset returns the offset of assignment
func (cstate *ConsumerMetadata) GetOffset(assignment *consumepb.MetadataAssignment, cgID int32) int32 {
	topicName := assignment.GetTopicName()
	partitionIndex := assignment.GetPartitionIndex()
	cstate.mux.RLock()
	for _, group := range cstate.MetadataConsumerState.GetConsumerGroups() {
		if group.GetID() == cgID {
			for _, a := range group.GetAssignments() {
				if a.GetTopicName() == topicName && a.GetPartitionIndex() == partitionIndex {
					return a.Offset
				}
			}
		}
	}
	cstate.mux.RUnlock()
	return 0
}

// UpdateOffset increases offset of assignment in metadata by one
func (cstate *ConsumerMetadata) UpdateOffset(assignment *consumepb.MetadataAssignment, cgID int32) {
	cstate.mux.Lock()
	newOffset := assignment.GetOffset() + 1

	topicName := assignment.GetTopicName()
	partitionIndex := assignment.GetPartitionIndex()
	for _, group := range cstate.MetadataConsumerState.GetConsumerGroups() {
		if group.GetID() == cgID {
			for _, a := range group.GetAssignments() {
				if a.GetTopicName() == topicName && a.GetPartitionIndex() == partitionIndex {
					a.Offset = newOffset
					break
				}
			}
		}
	}
	cstate.mux.Unlock()
}

// UpdateAssignments update assignment of consumer group
func (cstate *ConsumerMetadata) UpdateAssignments(cg int, assignments []*consumepb.MetadataAssignment) {
	cstate.mux.Lock()
	for _, grp := range cstate.MetadataConsumerState.GetConsumerGroups() {
		if int(grp.GetID()) == cg {
			grp.Assignments = append(grp.Assignments, assignments...)
			break
		}
	}
	cstate.mux.Unlock()
	cstate.populateMetadata()
}

func (cstate *ConsumerMetadata) populateMetadata() {
	cstate.mux.RLock()
	brokerAssignmentsMap := make(map[int][]*consumepb.MetadataAssignment)

	for _, grp := range cstate.MetadataConsumerState.GetConsumerGroups() {
		for _, assignment := range grp.GetAssignments() {
			brokerAssignmentsMap[int(assignment.GetBroker())] = append(brokerAssignmentsMap[int(assignment.GetBroker())],
				assignment)
		}
	}

	cstate.mux.RUnlock()

	cstate.mux.Lock()
	cstate.brokerAssignmentsMap = brokerAssignmentsMap
	cstate.mux.Unlock()
}

func (cstate *ConsumerMetadata) GetBrokerAssignmentsMap() map[int][]*consumepb.MetadataAssignment {
	cstate.mux.Lock()
	defer cstate.mux.Unlock()
	return cstate.brokerAssignmentsMap
}
