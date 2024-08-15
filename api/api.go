package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pe-Gomes/memory-crud-go/infra"
)

type apiHandler struct {
	r     *chi.Mux
	q     *infra.AppDB
	v     *validator.Validate
	db    *infra.AppDB
	mutex *sync.Mutex
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(db *infra.AppDB) http.Handler {
	v := validator.New(validator.WithRequiredStructEnabled())
	a := &apiHandler{v: v, db: db, mutex: &sync.Mutex{}}

	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.RequestID, middleware.Logger)
	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", a.handleListUsers)
			r.Post("/", a.handleCreateUser)
			r.Get("/{userID}", a.handleGetUser)
			r.Delete("/{userID}", a.handleDeleteUser)
			r.Put("/{userID}", a.handleUpdateUser)
		})
	})

	a.r = r
	return a
}

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func handleJSON(w http.ResponseWriter, r Response, code int) {
	data, err := json.Marshal(r)
	if err != nil {
		slog.Error("error marshal", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

func (h apiHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.db.ListUsers()

	res := make([]UserResponse, 0, len(users))

	for _, user := range users {
		res = append(res, UserResponse{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Biography: user.Biography,
		})
	}

	handleJSON(w, Response{Data: res}, http.StatusOK)
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=20"`
	LastName  string `json:"last_name" validate:"required,min=2,max=20"`
	Biography string `json:"biography" validate:"required,min=2,max=200"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

func (h apiHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handleJSON(w, Response{Message: "invalid request"}, http.StatusUnprocessableEntity)
		return
	}

	if err := h.v.Struct(body); err != nil {
		handleJSON(w, Response{Message: fmt.Sprintf("Invalid input:%s", err.Error())}, http.StatusUnprocessableEntity)
		return
	}

	userID := h.db.CreateUser(infra.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Biography: body.Biography,
	})

	fmt.Println(*h.db)
	handleJSON(w, Response{Data: CreateUserResponse{ID: userID.String()}}, http.StatusCreated)
	return
}

func (h apiHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		handleJSON(w, Response{Message: "invalid user id"}, http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUser(infra.ID(userID))
	if err != nil {
		handleJSON(w, Response{Message: "could not find user"}, http.StatusNotFound)
		return
	}

	handleJSON(w, Response{Data: UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Biography: user.Biography,
	}},
		http.StatusOK,
	)
}

func (h apiHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		handleJSON(w, Response{Message: "invalid uuid"}, http.StatusBadRequest)
		return
	}

	err = h.db.DeleteUser(infra.ID(userID))
	if err != nil {
		handleJSON(w, Response{Message: "could not delete user"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=20"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=20"`
	Biography string `json:"biography" validate:"omitempty,min=2,max=200"`
}

func (h apiHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	rawUserID := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		handleJSON(w, Response{Message: "invalid uuid"}, http.StatusBadRequest)
		return
	}

	var body UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		handleJSON(w, Response{Message: "invalid request"}, http.StatusUnprocessableEntity)
		return
	}

	if err := h.v.Struct(body); err != nil {
		handleJSON(w, Response{Message: fmt.Sprintf("Invalid input: %s", err.Error())}, http.StatusUnprocessableEntity)
		return
	}

	err = h.db.UpdateUser(infra.ID(userID), infra.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Biography: body.Biography,
	})
	if err != nil {
		handleJSON(w, Response{Message: "could not update user"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
