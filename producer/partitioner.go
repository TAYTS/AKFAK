package main

import (
	"github.com/orcaman/concurrent-map"
	"sync"
)

type Partitioner interface {
	getPartition() int
	nextValue() int
}

type RoundRobinPartitioner struct {
	partitioner Partitioner
	// https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go
	topicCounterMap cmap.ConcurrentMap
}

type AtomicCounter struct {
	counter int
	mux sync.Mutex
}

// returns the updated value
func (a AtomicCounter) increment() {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.counter = a.counter + 1
}

func (a AtomicCounter) set(Integer int) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.counter = Integer
}

func (a AtomicCounter) getCount() int {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.counter
}

// where cluster.partitionsForTopic(topic) gives an empty list or a list from the Map of the topic
func (r RoundRobinPartitioner) getPartition (topic string, key []byte, value []byte, cluster Cluster) int {
	var partitions []PartitionInfo = cluster.partitionsForTopic(topic)
	var availablePartitions []PartitionInfo = cluster.availablePartitionsForTopic(topic)
	numPartitions := len(partitions)
	numAvailPartitions := len(partitions)

	nextValue := r.getNextValue(topic)

	if len(availablePartitions) == 0 {
		part := nextValue % numAvailPartitions
		return availablePartitions[part].partition()
	} else {
		return nextValue % numPartitions
	}
}

func (r RoundRobinPartitioner) getNextValue(topic string) int {
	// check if topic in topicCounterMap. If yes, getCount the current counter value and increment it.
	// Otherwise, add it to the concurrent map and return the initialised value of 0.
	var counter AtomicCounter
	c, ok := r.topicCounterMap.Get(topic)
	if ok {
		// set counter to getCount the counter of topicCounterMap
		counter.set(c.(int))
		// increase the counter and getCount the updated value - thread safe
		counter.increment()
		// update the topicCounterMap with the updated value
		r.topicCounterMap.Set(topic, counter.getCount())
	} else {
		counter.set(0)
		r.topicCounterMap.Set(topic, 0)
	}
	return counter.getCount()
}

func main() {
	//var r RoundRobinPartitioner
}

