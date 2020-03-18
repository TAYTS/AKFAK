package cluster

import "fmt"

// Node store the information about the Kafka node
type Node struct {
	id       int
	idString string
	host     string
	port     int
}

///////////////////////////////////
// 		   Public Methods		 //
///////////////////////////////////

// InitNode create new Node with the id, host and port and return the Node pointer type
func InitNode(id int, host string, port int) *Node {
	node := &Node{
		id:       id,
		idString: fmt.Sprintf("%v", id),
		host:     host,
		port:     port,
	}

	return node
}

// NoNode return empty/invalid node
func NoNode() *Node {
	return InitNode(-1, "", -1)
}

// IsEmpty check if the current node is empty/invalid node
func (node *Node) IsEmpty() bool {
	return node.host == "" || node.port < 0
}

// GetID return the node id
func (node *Node) GetID() int {
	return node.id
}

// GetIDString return the node idString
func (node *Node) GetIDString() string {
	return node.idString
}

// GetHost return the node host
func (node *Node) GetHost() string {
	return node.host
}

// GetPort return the node port
func (node *Node) GetPort() int {
	return node.port
}

// String return the string representation of Node
func (node *Node) String() string {
	return fmt.Sprintf("%v:%v (id: %v)", node.host, node.port, node.idString)
}
