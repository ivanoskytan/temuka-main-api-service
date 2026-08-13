package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/service"
	rest "github.com/temuka-api-service/util/rest"
)

type CommentHandler interface {
	AddComment(w http.ResponseWriter, r *http.Request)
	ShowCommentsByPost(w http.ResponseWriter, r *http.Request)
	DeleteComment(w http.ResponseWriter, r *http.Request)
	ShowReplies(w http.ResponseWriter, r *http.Request)
	GetCommentsByUser(w http.ResponseWriter, r *http.Request)
}

type CommentHandlerImpl struct {
	CommentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) CommentHandler {
	return &CommentHandlerImpl{
		CommentService: commentService,
	}
}

func (h *CommentHandlerImpl) AddComment(w http.ResponseWriter, r *http.Request) {
	var req dto.AddCommentRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	comment, err := h.CommentService.AddComment(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Comment has been added",
		"data":    comment,
	})
}

func (h *CommentHandlerImpl) ShowCommentsByPost(w http.ResponseWriter, r *http.Request) {
	var req dto.ShowCommentsRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	comments, err := h.CommentService.ShowCommentsByPost(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Comments have been retrieved",
		"data":    comments,
	})
}

func (h *CommentHandlerImpl) GetCommentsByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	comments, err := h.CommentService.GetCommentsByUser(r.Context(), userID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Comments have been retrieved",
		"data":    comments,
	})
}

func (h *CommentHandlerImpl) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentIDStr := vars["comment_id"]

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid comment ID"})
		return
	}

	if err := h.CommentService.DeleteComment(r.Context(), commentID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]string{"message": "Comment has been deleted"})
}

func (h *CommentHandlerImpl) ShowReplies(w http.ResponseWriter, r *http.Request) {
	var req dto.ShowRepliesRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	replies, err := h.CommentService.ShowReplies(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Replies have been shown",
		"data":    replies,
	})
}
