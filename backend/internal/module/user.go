package module

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apperrors "task-management-backend/internal/constant/errors"
	"task-management-backend/internal/model/db"
	"task-management-backend/internal/model/dto"
	"task-management-backend/internal/storage/persistence"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserModule struct {
	store     *persistence.Store
	jwtSecret []byte
}

type AuthClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func NewUserModule(store *persistence.Store, jwtSecret string) *UserModule {
	return &UserModule{store: store, jwtSecret: []byte(jwtSecret)}
}

func (m *UserModule) Register(ctx context.Context, req dto.RegisterRequest) (string, db.User, error) {
	if err := req.Validate(); err != nil {
		return "", db.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", db.User{}, err
	}

	user, err := m.store.CreateUser(ctx, req.Name, req.Email, string(hash))
	if err != nil {
		if persistence.IsUniqueViolation(err) {
			return "", db.User{}, apperrors.ErrUserAlreadyExists
		}
		return "", db.User{}, err
	}

	token, err := m.generateToken(user.ID)
	if err != nil {
		return "", db.User{}, err
	}

	user.PasswordHash = ""
	return token, user, nil
}

func (m *UserModule) Login(ctx context.Context, req dto.LoginRequest) (string, db.User, error) {
	if err := req.Validate(); err != nil {
		return "", db.User{}, err
	}

	user, err := m.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", db.User{}, apperrors.ErrInvalidCredentials
		}
		return "", db.User{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", db.User{}, apperrors.ErrInvalidCredentials
	}

	token, err := m.generateToken(user.ID)
	if err != nil {
		return "", db.User{}, err
	}

	user.PasswordHash = ""
	return token, user, nil
}

func (m *UserModule) GetProfile(ctx context.Context, userID int64) (db.User, error) {
	user, err := m.store.GetUserByID(ctx, userID)
	if err != nil {
		return db.User{}, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (m *UserModule) ListTasks(ctx context.Context, userID int64) ([]db.Task, error) {
	return m.store.ListTasksByUserID(ctx, userID)
}

func (m *UserModule) CreateTask(ctx context.Context, userID int64, req dto.CreateTaskRequest) (db.Task, error) {
	status := req.Status
	if status == "" {
		status = db.TaskStatusPending
	}

	if err := dto.ValidateTaskFields(req.Title, status); err != nil {
		return db.Task{}, err
	}

	return m.store.CreateTask(ctx, userID, req.Title, req.Description, status)
}

func (m *UserModule) UpdateTask(ctx context.Context, userID, taskID int64, req dto.UpdateTaskRequest) (db.Task, error) {
	status := req.Status
	if status == "" {
		status = db.TaskStatusPending
	}

	if err := dto.ValidateTaskFields(req.Title, status); err != nil {
		return db.Task{}, err
	}

	task, err := m.store.UpdateTask(ctx, userID, taskID, req.Title, req.Description, status)
	if persistence.IsNotFound(err) {
		return db.Task{}, apperrors.ErrTaskNotFound
	}

	return task, err
}

func (m *UserModule) UpdateTaskStatus(ctx context.Context, userID, taskID int64, status string) (db.Task, error) {
	if err := dto.ValidateTaskFields("x", status); err != nil {
		return db.Task{}, err
	}

	task, err := m.store.UpdateTaskStatus(ctx, userID, taskID, status)
	if persistence.IsNotFound(err) {
		return db.Task{}, apperrors.ErrTaskNotFound
	}
	return task, err
}

func (m *UserModule) DeleteTask(ctx context.Context, userID, taskID int64) error {
	err := m.store.DeleteTask(ctx, userID, taskID)
	if persistence.IsNotFound(err) {
		return apperrors.ErrTaskNotFound
	}
	return err
}

func (m *UserModule) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		return m.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return 0, apperrors.ErrUnauthorized
	}

	return claims.UserID, nil
}

func (m *UserModule) generateToken(userID int64) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtSecret)
}
