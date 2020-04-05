package producer

import (
	"AKFAK/proto/metadatapb"
	"sync"

	cmap "github.com/orcaman/concurrent-map"
)

// RoundRobinPartitioner is to be filled
type RoundRobinPartitioner struct {
	// https://github.com/orcaman/concurrent-map/blob/master/concurrent_map.go
	topicCounterMap cmap.ConcurrentMap
}

// AtomicCounter is to be filled
type AtomicCounter struct {
	counter int
	mux     sync.Mutex
}

// returns the updated value
func (a *AtomicCounter) increment() {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.counter = a.counter + 1
}

func (a *AtomicCounter) set(Integer int) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.counter = Integer
}

func (a *AtomicCounter) getCount() int {
	a.mux.Lock()
	defer a.mux.Unlock()
	return a.counter
}

// getPartition return the next partition index for a topic
func (r *RoundRobinPartitioner) getPartition(topic string, partitions []*metadatapb.Partition) int {
	numPartitions := len(partitions)
	nextValue := r.getNextValue(topic)

	if numPartitions > 0 {
		part := nextValue % numPartitions
		return int(partitions[part].GetPartitionIndex())
	}
	return nextValue % numPartitions
}

func (r *RoundRobinPartitioner) getNextValue(topic string) int {
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
