package rest

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strconv"

	apperrors "task-management-backend/internal/constant/errors"
	"task-management-backend/internal/handler/middleware"
	"task-management-backend/internal/model/dto"
	"task-management-backend/internal/model/response"
	"task-management-backend/internal/module"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	module *module.Module
}

func New(module *module.Module) *Handler {
	return &Handler{module: module}
}

func (h *Handler) Module() *module.Module {
	return h.module
}

func (h *Handler) Register(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid payload")
		return
	}

	token, user, err := h.Module().User.Register(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusCreated, response.AuthResponse{Token: token, User: user})
}

func (h *Handler) Login(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid payload")
		return
	}

	token, user, err := h.Module().User.Login(r.Context(), req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.AuthResponse{Token: token, User: user})
}

func (h *Handler) Me(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	user, err := h.Module().User.GetProfile(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.ProfileResponse{User: user})
}

func (h *Handler) ListTasks(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	tasks, err := h.Module().User.ListTasks(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.TaskListResponse{Tasks: tasks})
}

func (h *Handler) CreateTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid payload")
		return
	}

	task, err := h.Module().User.CreateTask(r.Context(), userID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusCreated, response.TaskResponse{Task: task})
}

func (h *Handler) UpdateTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	taskID, err := parseTaskID(r)
	if err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid payload")
		return
	}

	task, err := h.Module().User.UpdateTask(r.Context(), userID, taskID, req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.TaskResponse{Task: task})
}

func (h *Handler) UpdateTaskStatus(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	taskID, err := parseTaskID(r)
	if err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid task id")
		return
	}

	var req dto.UpdateTaskStatusRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid payload")
		return
	}

	task, err := h.Module().User.UpdateTaskStatus(r.Context(), userID, taskID, req.Status)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.TaskResponse{Task: task})
}

func (h *Handler) DeleteTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, nethttp.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
		return
	}

	taskID, err := parseTaskID(r)
	if err != nil {
		response.Error(w, nethttp.StatusBadRequest, "invalid task id")
		return
	}

	if err = h.Module().User.DeleteTask(r.Context(), userID, taskID); err != nil {
		h.writeError(w, err)
		return
	}

	response.JSON(w, nethttp.StatusOK, response.MessageResponse{Message: "Task deleted"})
}

func parseTaskID(r *nethttp.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
}

func (h *Handler) writeError(w nethttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		response.Error(w, nethttp.StatusConflict, err.Error())
	case errors.Is(err, apperrors.ErrInvalidCredentials), errors.Is(err, apperrors.ErrUnauthorized):
		response.Error(w, nethttp.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrTaskNotFound):
		response.Error(w, nethttp.StatusNotFound, err.Error())
	default:
		response.Error(w, nethttp.StatusBadRequest, err.Error())
	}
}
