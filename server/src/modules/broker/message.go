package broker

import (
	"fmt"
	"strings"

	SliceUtils "github.com/matheusm25/thoth/src/utils/slices"
)

type Message struct {
	ID          string
	Topic       string
	MessageType string
	Payload     string
}

func (m *Message) Validate() (bool, error) {
	validMessageTypes := []string{
		"PUBLISH",
		"SUBSCRIBE",
		"UNSUBSCRIBE",
		"ACKNOWLEDGE",
		"UNACKNOWLEDGE",
	}

	if !SliceUtils.Contains(validMessageTypes, strings.ToUpper(m.MessageType)) {
		return false, fmt.Errorf("invalid message type: %v", m.MessageType)
	}

	if m.Topic == "" {
		return false, fmt.Errorf("topic cannot be empty")
	}

	if m.MessageType == "PUBLISH" && m.Payload == "" {
		return false, fmt.Errorf("payload cannot be empty")
	}

	if m.ID == "" {
		return false, fmt.Errorf("ID cannot be empty")
	}

	return true, nil
}
