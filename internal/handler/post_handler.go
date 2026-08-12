package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/service"
	"github.com/temuka-api-service/util/rest"
)

type PostHandler interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetPostDetail(w http.ResponseWriter, r *http.Request)
	GetUserPosts(w http.ResponseWriter, r *http.Request)
	UpdatePost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
	GetTimelinePosts(w http.ResponseWriter, r *http.Request)
	LikePost(w http.ResponseWriter, r *http.Request)
	UnlikePost(w http.ResponseWriter, r *http.Request)
	SavePost(w http.ResponseWriter, r *http.Request)
	UnsavePost(w http.ResponseWriter, r *http.Request)
	GetSavedPostsByUser(w http.ResponseWriter, r *http.Request)
}

type PostHandlerImpl struct {
	postService service.PostService
}

func NewPostHandler(s service.PostService) PostHandler {
	return &PostHandlerImpl{postService: s}
}

func (h *PostHandlerImpl) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	post, err := h.postService.CreatePost(r.Context(), &req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Post created", Data: post}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) GetPostDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var currentUserID int
	if userID, ok := r.Context().Value("user_id").(int); ok {
		currentUserID = userID
	}
	post, err := h.postService.GetPostDetail(r.Context(), id, currentUserID)
	if err != nil {
		rest.WriteResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}
	resp := dto.MessageResponse{Message: "Post detail retrieved", Data: post}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
	posts, err := h.postService.GetUserPosts(r.Context(), userID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "User posts retrieved", Data: posts}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdatePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), id, &req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := dto.MessageResponse{Message: "Post updated", Data: post}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.postService.DeletePost(r.Context(), id); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "Post deleted"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) GetTimelinePosts(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
	posts, err := h.postService.GetTimelinePosts(r.Context(), userID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "Timeline posts retrieved", Data: posts}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) LikePost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.LikePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.postService.LikePost(r.Context(), postID, req.UserID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "You have liked this post"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) UnlikePost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UnlikePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.postService.UnlikePost(r.Context(), postID, req.UserID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "You have unliked this post"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) SavePost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.SavePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.postService.SavePost(r.Context(), postID, req.UserID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "You have saved this post"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) UnsavePost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UnsavePostRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	if err := h.postService.UnsavePost(r.Context(), postID, req.UserID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "You have unsaved this post"}
	rest.WriteResponse(w, http.StatusOK, resp)
}

func (h *PostHandlerImpl) GetSavedPostsByUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
	posts, err := h.postService.GetSavedPostsByUser(r.Context(), userID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := dto.MessageResponse{Message: "Saved posts retrieved", Data: posts}
	rest.WriteResponse(w, http.StatusOK, resp)
}
