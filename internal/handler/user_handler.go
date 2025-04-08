package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"users/internal/model"
	"users/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate model.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.userService.CreateUser(r.Context(), userCreate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var userUpdate model.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := h.userService.UpdateUser(r.Context(), id, userUpdate)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
