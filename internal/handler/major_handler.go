package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/service"
	"github.com/temuka-api-service/util/rest"
)

type MajorHandler interface {
	AddMajor(w http.ResponseWriter, r *http.Request)
	GetMajors(w http.ResponseWriter, r *http.Request)
	GetMajorDetail(w http.ResponseWriter, r *http.Request)
	GetMajorsByUniversity(w http.ResponseWriter, r *http.Request)
	AddMajorReview(w http.ResponseWriter, r *http.Request)
	GetMajorReviews(w http.ResponseWriter, r *http.Request)
}

type MajorHandlerImpl struct {
	MajorService service.MajorService
}

func NewMajorHandler(service service.MajorService) MajorHandler {
	return &MajorHandlerImpl{
		MajorService: service,
	}
}

func (h *MajorHandlerImpl) AddMajor(w http.ResponseWriter, r *http.Request) {
	var req dto.AddMajorRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	major, err := h.MajorService.AddMajor(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Major has been added successfully",
		"data":    major,
	})
}

func (h *MajorHandlerImpl) GetMajors(w http.ResponseWriter, r *http.Request) {
	majors, err := h.MajorService.GetMajors(r.Context())
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Majors list retrieved successfully",
		"data":    majors,
	})
}

func (h *MajorHandlerImpl) GetMajorDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid major ID key format"})
		return
	}

	major, err := h.MajorService.GetMajorDetail(r.Context(), id)
	if err != nil {
		rest.WriteResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Major detail tracking complete",
		"data":    major,
	})
}

func (h *MajorHandlerImpl) GetMajorsByUniversity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uniID, err := strconv.Atoi(vars["university_id"])
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid university ID format"})
		return
	}

	majors, err := h.MajorService.GetMajorsByUniversity(r.Context(), uniID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Majors for specified university retrieved",
		"data":    majors,
	})
}

func (h *MajorHandlerImpl) AddMajorReview(w http.ResponseWriter, r *http.Request) {
	var req dto.AddMajorReviewRequest
	if err := rest.ReadRequest(r, &req); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body payload"})
		return
	}

	review, err := h.MajorService.AddMajorReview(r.Context(), req)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Major review submitted, calculations saved",
		"data":    review,
	})
}

func (h *MajorHandlerImpl) GetMajorReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	majorID, err := strconv.Atoi(vars["major_id"])
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid major ID metric format"})
		return
	}

	reviews, err := h.MajorService.GetMajorReviews(r.Context(), majorID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rest.WriteResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Major comprehensive review reports tracked",
		"data":    reviews,
	})
}
