package response

import "task-management-backend/internal/model/db"

type AuthResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

type ProfileResponse struct {
	User db.User `json:"user"`
}

type TaskListResponse struct {
	Tasks []db.Task `json:"tasks"`
}

type TaskResponse struct {
	Task db.Task `json:"task"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
