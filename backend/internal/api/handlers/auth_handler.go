package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/anomalyco/iptables-visualizer/internal/config"
	"github.com/anomalyco/iptables-visualizer/internal/models"
	"github.com/anomalyco/iptables-visualizer/internal/repository"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	auditRepo repository.AuditRepository
	cfg      config.JWTConfig
}

func NewAuthHandler(ur repository.UserRepository, ar repository.AuditRepository, cfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{userRepo: ur, auditRepo: ar, cfg: cfg}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		h.logAudit(0, req.Username, models.AuditLogin, "login", "failed: user not found", r.RemoteAddr, false)
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.logAudit(user.ID, user.Username, models.AuditLogin, "login", "failed: wrong password", r.RemoteAddr, false)
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		h.logAudit(user.ID, user.Username, models.AuditLogin, "login", "failed: inactive account", r.RemoteAddr, false)
		http.Error(w, `{"error":"account is disabled"}`, http.StatusForbidden)
		return
	}

	expiresAt := time.Now().Add(h.cfg.Expiration)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     string(user.Role),
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.cfg.Secret))
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	h.logAudit(user.ID, user.Username, models.AuditLogin, "login", "success", r.RemoteAddr, true)

	resp := models.LoginResponse{
		Token:     tokenStr,
		User:      *user,
		ExpiresAt: expiresAt.Unix(),
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		Role:         req.Role,
		Email:        req.Email,
		IsActive:     true,
	}

	if user.Role == "" {
		user.Role = models.RoleViewer
	}

	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, `{"error":"failed to create user: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	adminUser := r.Context().Value("username").(string)
	adminID := r.Context().Value("user_id").(int64)
	h.logAudit(adminID, adminUser, models.AuditCreateUser, "user:"+user.Username, "created", r.RemoteAddr, true)

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.List(100, 0)
	if err != nil {
		http.Error(w, `{"error":"failed to list users"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *AuthHandler) logAudit(userID int64, username string, action models.AuditAction, resource, details, ip string, success bool) {
	go func() {
		h.auditRepo.Create(&models.AuditLog{
			UserID:    userID,
			Username:  username,
			Action:    action,
			Resource:  resource,
			Details:   details,
			IPAddress: ip,
			Success:   success,
		})
	}()
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
