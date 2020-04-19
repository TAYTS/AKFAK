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

// UpdateOffset increases offset of assignment in metadata by one
func (cstate *ConsumerMetadata) UpdateOffset(assignment *consumepb.MetadataAssignment) {
	// TODO
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
