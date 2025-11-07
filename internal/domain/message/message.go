package message

import (
	"errors"
	"fmt"
)

const (
	FailedStatus = "Failed"

	FailedMessage     = "⚠️ failed ping. Resource: %q"
	ReadyAgainMessage = "✅ ping ready again from resource %q"
)

var (
	ErrInvalidResource = errors.New("message resource is required")
	ErrInvalidStatus   = errors.New("message status is required")
)

type Message struct {
	resource, status string
}

func (m Message) Resource() string {
	return m.resource
}

func (m Message) Status() string {
	return m.status
}

func (m Message) IsFailed() bool {
	return m.status == FailedStatus
}

func (m Message) Message() string {
	msg := ReadyAgainMessage
	if m.IsFailed() {
		msg = FailedMessage
	}

	return fmt.Sprintf(msg, m.resource)
}

func (m Message) validate() error {
	switch {
	case m.resource == "":
		return ErrInvalidResource
	case m.status == "":
		return ErrInvalidStatus
	default:
		return nil
	}
}

func NewMessage(resource string, status string) (*Message, error) {
	m := Message{resource: resource, status: status}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}
