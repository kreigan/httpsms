package events

import (
	"time"

	"github.com/google/uuid"
)

// EventTypeMessagePhoneSent is emitted when the phone sends a message
const EventTypeMessagePhoneSent = "message.phone.sent"

// MessagePhoneSentPayload is the payload of the EventTypeMessagePhoneSent event
type MessagePhoneSentPayload struct {
	ID        uuid.UUID `json:"id"`
	Owner     string    `json:"owner"`
	Contact   string    `json:"contact"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}