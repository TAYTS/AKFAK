package main

import (
	"AKFAK/consumermetadata"
	"AKFAK/proto/consumepb"
	"AKFAK/proto/consumermetadatapb"
	"reflect"
	"testing"
)

func TestOffset(t *testing.T) {
	var assignments [] *consumepb.MetadataAssignment
	assignment := consumepb.MetadataAssignment{
		TopicName:            "DSC",
		PartitionIndex:       0,
		Offset:               0,
		Broker:               0,
	}
	assignments = append(assignments, &assignment)
	cg := consumermetadatapb.ConsumerGroup{
		ID:                   1,
		Assignments:          assignments,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	var consumerGroups []*consumermetadatapb.ConsumerGroup
	consumerGroups = append(consumerGroups, &cg)

	metaDataConsumerState := consumermetadatapb.MetadataConsumerState{
		ConsumerGroups:       consumerGroups,
	}
	c := consumermetadata.InitConsumerMetadata(&metaDataConsumerState)
	c.UpdateOffset(&assignment, cg.GetID())

	var updatedOffset int32
	for _, v := range c.ConsumerGroups {
		if v.GetID() == cg.GetID() {
			for _, a := range v.Assignments {
				updatedOffset = a.GetOffset()
			}
		}
	}
	if !reflect.DeepEqual(updatedOffset, int32(1)) {
		t.Errorf("Evaluated Offset: %+v. Expected Offset: %+v", updatedOffset, int32(1))
	}
}
