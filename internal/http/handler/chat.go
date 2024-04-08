package handler

import (
	"errors"
	"fmt"
	"fund-o/api-server/cmd/ws"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ChatHandler struct {
	channelUsecase usecase.ChannelUsecase
	messageUsecase usecase.MessageUsecase
	socketService  ws.SocketService
}

type ChatHandlerOptions struct {
	usecase.ChannelUsecase
	usecase.MessageUsecase
	ws.SocketService
}

func NewChatHandler(options *ChatHandlerOptions) *ChatHandler {
	return &ChatHandler{
		channelUsecase: options.ChannelUsecase,
		messageUsecase: options.MessageUsecase,
		socketService:  options.SocketService,
	}
}

// GetOrCreateChannel Godoc
// @summary Get or create a channel
// @description Get or create a channel
// @tags chat
// @accept json
// @produce json
// @security ApiKeyAuth
// @param id path string true "Recipient ID"
// @success 200 {object} handler.ResultResponse[entity.ChannelDto] "OK"
// @failure 401 {object} handler.ErrorResponse "Unauthorized"
// @failure 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /channels/{id} [get]
func (h *ChatHandler) GetOrCreateChannel(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID
	channelID := c.Param("id")

	// Check if channel already exists
	existingChannel, err := h.channelUsecase.GetExistingChannel(userID, channelID)
	if err == nil {
		fmt.Println("Existing Channel: ", existingChannel)
		c.JSON(makeHttpResponse(http.StatusOK, existingChannel))
		return
	}

	var payload entity.ChannelCreatePayload
	payload.Name = fmt.Sprintf("%s_%s", userID, channelID)
	payload.Members = []string{userID, channelID}

	channel, err := h.channelUsecase.CreateChannel(&payload)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, channel))
}

func (h *ChatHandler) GetOwnChannels(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID

	channels, err := h.channelUsecase.GetChannelByUserID(userID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	response := make([]entity.GetOwnChannelsResponse, len(channels))
	for i, c := range channels {
		var receiver entity.UserDto
		for _, m := range c.Members {
			if m.ID != userID {
				receiver = m
			}
		}
		response[i] = entity.GetOwnChannelsResponse{
			Receiver:    receiver,
			LastMessage: c.Messages[len(c.Messages)-1],
		}
	}

	c.JSON(makeHttpResponse(http.StatusOK, response))
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID
	channelID := c.Param("id")

	var payload entity.MessageCreatePayload
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	channel, err := h.channelUsecase.GetExistingChannel(userID, channelID)
	if err != nil {
		err := errors.New("channel not found")
		c.JSON(makeHttpErrorResponse(http.StatusNotFound, err.Error()))
		return
	}

	authorID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	parsedChannelID, err := uuid.Parse(channel.ID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	message, err := h.messageUsecase.CreateChannelMessage(parsedChannelID, &entity.MessageCreatePayload{
		Text:       payload.Text,
		Attachment: payload.Attachment,
		AuthorID:   authorID,
	})
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	h.socketService.EmitNewMessage(channel.Name, message)

	c.JSON(makeHttpResponse(http.StatusOK, message))
}
