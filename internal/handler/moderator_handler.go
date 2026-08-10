package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/service"
	"github.com/temuka-api-service/util/rest"
)

type ModeratorHandler interface {
	GetModeratorsByCommunityID(w http.ResponseWriter, r *http.Request)
	SendModeratorRequest(w http.ResponseWriter, r *http.Request)
	RemoveModerator(w http.ResponseWriter, r *http.Request)
}

type ModeratorHandlerImpl struct {
	ModeratorService service.ModeratorService
}

func NewModeratorHandler(moderatorService service.ModeratorService) ModeratorHandler {
	return &ModeratorHandlerImpl{
		ModeratorService: moderatorService,
	}
}

func (h *ModeratorHandlerImpl) GetModeratorsByCommunityID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDStr := vars["community_id"]

	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community ID"})
		return
	}

	moderators, err := h.ModeratorService.GetModeratorsByCommunityID(r.Context(), communityID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Moderators have been retrieved",
		"data":    moderators,
	})
}

func (h *ModeratorHandlerImpl) SendModeratorRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.SendModeratorRequest

	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.ModeratorService.SendModeratorRequest(r.Context(), request); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]string{"message": "Moderator request has been sent"})
}

func (h *ModeratorHandlerImpl) RemoveModerator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moderatorIDStr := vars["id"]

	moderatorID, err := strconv.Atoi(moderatorIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid moderator ID"})
		return
	}

	if err := h.ModeratorService.RemoveModerator(r.Context(), moderatorID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]string{"message": "Moderator has been removed"})
}
