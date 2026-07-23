package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/service"
	rest "github.com/temuka-api-service/util/rest"
	"github.com/temuka-api-service/util/websocket"
	ws "github.com/temuka-api-service/util/websocket"
)

type ChatHandler interface {
	ServeWebsocket(w http.ResponseWriter, r *http.Request)
	AddConversation(w http.ResponseWriter, r *http.Request)
	AddMessage(w http.ResponseWriter, r *http.Request)
	AddParticipant(w http.ResponseWriter, r *http.Request)
	GetConversationsByUserID(w http.ResponseWriter, r *http.Request)
	GetConversationDetail(w http.ResponseWriter, r *http.Request)
	DeleteConversation(w http.ResponseWriter, r *http.Request)
	RetrieveMessages(w http.ResponseWriter, r *http.Request)
}

type ChatHandlerImpl struct {
	Hub                 *ws.Hub
	ConversationService service.ConversationService
}

func NewChatHandler(hub *websocket.Hub, conversationService service.ConversationService) ChatHandler {
	return &ChatHandlerImpl{Hub: hub, ConversationService: conversationService}
}

func (h *ChatHandlerImpl) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("[WS Check] Incoming connection request from:", r.RemoteAddr)
	query := r.URL.Query()
	conversationID, err1 := strconv.Atoi(query.Get("conversation_id"))
	userID, err2 := strconv.Atoi(query.Get("user_id"))

	if err1 != nil || err2 != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation_id or user_id parameter"})
		return
	}

	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS Upgrade Error]: %v", err)
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to upgrade to WebSocket"})
		return
	}

	client := &ws.Client{
		Hub:            h.Hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		ConversationID: conversationID,
	}

	client.Hub.Register <- client

	go client.WritePump()

	go client.ReadPump(func(msg ws.WSMessage) {
		req := dto.AddMessageRequest{
			ConversationID: msg.ConversationID,
			UserID:         msg.SenderID,
			ParticipantID:  msg.ParticipantID,
			Text:           msg.Text,
		}

		ctx := context.Background()
		if _, err := h.ConversationService.AddMessage(ctx, req); err != nil {
			return
		}

	})
}

func (h *ChatHandlerImpl) AddConversation(w http.ResponseWriter, r *http.Request) {
	var req dto.AddConversationRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.UserID == 0 {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "User id is required"})
		return
	}

	if len(req.ParticipantIDs) == 0 {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "At least one participant id is required"})
		return
	}

	if len(req.ParticipantIDs) > 1 && req.Title == "" {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Title is required for group conversations"})
		return
	}

	conversation, err := h.ConversationService.AddConversation(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Conversation has been created", Data: conversation}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) AddMessage(w http.ResponseWriter, r *http.Request) {
	var req dto.AddMessageRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	message, err := h.ConversationService.AddMessage(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Message has been created", Data: message}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) AddParticipant(w http.ResponseWriter, r *http.Request) {
	var req dto.AddParticipantRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.ConversationService.AddParticipant(r.Context(), req); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Participant has been added"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) GetConversationsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDstr := mux.Vars(r)["user_id"]
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user id"})
		return
	}

	conversations, err := h.ConversationService.GetConversationsByUserID(r.Context(), userID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Conversations have been retrieved", Data: conversations}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) GetConversationDetail(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation id"})
		return
	}

	conversation, err := h.ConversationService.GetConversationDetail(r.Context(), id)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Conversation detail has been retrieved", Data: conversation}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation id"})
		return
	}

	if err := h.ConversationService.DeleteConversation(r.Context(), id); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Conversation has been deleted"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *ChatHandlerImpl) RetrieveMessages(w http.ResponseWriter, r *http.Request) {
	conversationIDstr := mux.Vars(r)["conversation_id"]
	conversationID, err := strconv.Atoi(conversationIDstr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation id"})
		return
	}

	messages, err := h.ConversationService.RetrieveMessages(r.Context(), conversationID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Messages have been retrieved", Data: messages}
	rest.WriteResponse(w, http.StatusOK, resp)
}
