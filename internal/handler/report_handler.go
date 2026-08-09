package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/service"
	"github.com/temuka-api-service/util/rest"
)

type ReportHandler interface {
	CreateReport(w http.ResponseWriter, r *http.Request)
	UpdateReport(w http.ResponseWriter, r *http.Request)
	GetReportsByCommunity(w http.ResponseWriter, r *http.Request)
	DeleteReport(w http.ResponseWriter, r *http.Request)
}

type ReportHandlerImpl struct {
	ReportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) ReportHandler {
	return &ReportHandlerImpl{
		ReportService: reportService,
	}
}

func (h *ReportHandlerImpl) CreateReport(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateReportRequest

	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	report := &model.Report{
		CommunityID:     request.CommunityID,
		ReportedUserID:  request.ReportedUserID,
		CommunityRuleID: request.CommunityRuleID,
		TargetType:      request.TargetType,
		TargetID:        request.TargetID,
		Reason:          request.Reason,
		Status:          "pending",
	}

	if err := h.ReportService.CreateReport(r.Context(), report); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := dto.MessageResponse{
		Message: "Report has been created successfully",
	}
	rest.WriteResponse(w, http.StatusOK, response)
}

func (h *ReportHandlerImpl) UpdateReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportIDStr := vars["id"]

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
		return
	}

	var request dto.UpdateReportRequest
	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	report := &model.Report{
		ID:              reportID,
		CommunityID:     request.CommunityID,
		ReportedUserID:  request.ReportedUserID,
		CommunityRuleID: request.CommunityRuleID,
		TargetType:      request.TargetType,
		TargetID:        request.TargetID,
		Reason:          request.Reason,
	}

	if err := h.ReportService.UpdateReport(r.Context(), report); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := dto.MessageResponse{
		Message: "Report has been updated successfully",
	}
	rest.WriteResponse(w, http.StatusOK, response)
}

func (h *ReportHandlerImpl) GetReportsByCommunity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	communityIDStr := vars["community_id"]

	communityID, err := strconv.Atoi(communityIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid community ID"})
		return
	}

	reports, err := h.ReportService.GetReportsByCommunity(r.Context(), communityID)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := dto.MessageResponse{
		Message: "Reports have been retrieved successfully",
		Data:    reports,
	}
	rest.WriteResponse(w, http.StatusOK, response)
}

func (h *ReportHandlerImpl) DeleteReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportIDStr := vars["id"]

	reportID, err := strconv.Atoi(reportIDStr)
	if err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid report ID"})
		return
	}

	if err := h.ReportService.DeleteReport(r.Context(), reportID); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := dto.MessageResponse{
		Message: "Report has been deleted",
	}
	rest.WriteResponse(w, http.StatusOK, response)
}
