package topicutils

type TopicOperation string

const (
	// LIST_TOPIC string constant for CLI command
	LIST_TOPIC TopicOperation = "list"
	// CREATE_TOPIC string constant for CLI command
	CREATE_TOPIC TopicOperation = "create"
	// DELETE_TOPIC string constant for CLI command
	DELETE_TOPIC TopicOperation = "delete"
)
