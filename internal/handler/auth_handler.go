package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/dto"
	"github.com/temuka-api-service/internal/model"
	"github.com/temuka-api-service/internal/service"
	rest "github.com/temuka-api-service/util/rest"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

type AuthHandlerImpl struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
	}
}

func (c *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest

	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	newUser, err := c.AuthService.Register(r.Context(), request)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := struct {
		Message string      `json:"message"`
		Data    *model.User `json:"data"`
	}{
		Message: "New user has been registered",
		Data:    newUser,
	}

	rest.WriteResponse(w, http.StatusOK, response)
}

func (c *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest

	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	data, err := c.AuthService.Login(r.Context(), request)
	if err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "User has login successfully",
		Token:   data["token"].(string),
	}

	rest.WriteResponse(w, http.StatusOK, response)
}

func (c *AuthHandlerImpl) ResetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["id"]

	var request dto.ResetPasswordRequest
	if err := rest.ReadRequest(r, &request); err != nil {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if userIDstr == "" {
		rest.WriteResponse(w, http.StatusBadRequest, map[string]string{"error": "User id is required"})
		return
	}

	request.UserID = userIDstr

	if err := c.AuthService.ResetPassword(r.Context(), request); err != nil {
		rest.WriteResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Password was reset successfully",
	}
	rest.WriteResponse(w, http.StatusOK, response)
}
