package ws

import (
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"github.com/rs/zerolog/log"
)

type SocketService interface {
	EmitNewMessage(room string, message *entity.MessageDto)
}

type socketService struct {
	hub *Hub
}

type SocketServiceConfig struct {
	*Hub
}

func NewSocketService(config *SocketServiceConfig) SocketService {
	return &socketService{
		hub: config.Hub,
	}
}

func (s *socketService) EmitNewMessage(room string, message *entity.MessageDto) {
	data, err := json.Marshal(entity.WebsocketMessage{
		Action: NewMessageAction,
		Data:   message,
	})

	if err != nil {
		log.Error().Err(err).Msg("error marshalling response")
	}

	fmt.Println("Broadcasting to room: ", room, string(data))
	s.hub.BroadcastToRoom(data, room)
}
