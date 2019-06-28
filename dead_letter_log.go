package psdll

import (
	"time"
)

// DeadLetterLog is the log format for the mispublished message.
type DeadLetterLog struct {
	Message   `json:"message"`
	Project   string    `json:"project"`
	Topic     string    `json:"topic"`
	Publisher string    `json:"publisher"`
	PodName   string    `json:"pod_name"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
}

// Message has a data and attributes.
type Message struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}
