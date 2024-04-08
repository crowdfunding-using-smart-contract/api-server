package entity

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

// ReceivedMessage represents a received websocket message
type ReceivedMessage struct {
	Action  string `json:"action"`
	Room    string `json:"room"`
	Message *any   `json:"message"`
}

// WebsocketMessage represents an emitted message
type WebsocketMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

// Encode turns the message into a byte array
func (message *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode message")
	}

	return encoding
}
