package dto

import (
	"task-management-backend/internal/model/db"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
	)
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}

func ValidateTaskFields(title, status string) error {
	payload := struct {
		Title  string
		Status string
	}{Title: title, Status: status}

	return validation.ValidateStruct(&payload,
		validation.Field(&payload.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&payload.Status, validation.Required, validation.In(
			db.TaskStatusPending,
			db.TaskStatusInProgress,
			db.TaskStatusCompleted,
		)),
	)
}
