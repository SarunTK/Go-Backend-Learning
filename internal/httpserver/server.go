package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go101/internal/auth"
	"go101/internal/model"
	"go101/internal/service"

	"github.com/gorilla/mux"
)

type Server struct {
	service   *service.UserService
	jwtSecret string
}

func NewServer(svc *service.UserService, jwtSecret string) *Server {
	return &Server{service: svc, jwtSecret: jwtSecret}
}

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", s.healthzHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/login", s.loginHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users", s.authMiddleware(s.createUserHandler)).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users", s.authMiddleware(s.listUsersHandler)).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id:[0-9]+}", s.authMiddleware(s.getUserHandler)).Methods(http.MethodGet)
	return r
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(s.jwtSecret, parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		_ = claims
		next(w, r)
	}
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJWT(s.jwtSecret, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token, "email": req.Email})
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.service.CreateUser(context.Background(), req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (s *Server) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.service.ListUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string][]model.User{"users": users})
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := s.service.GetUser(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
