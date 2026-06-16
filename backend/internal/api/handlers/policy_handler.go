package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/anomalyco/iptables-visualizer/internal/engine"
	"github.com/anomalyco/iptables-visualizer/internal/models"
	"github.com/anomalyco/iptables-visualizer/internal/repository"
)

type PolicyHandler struct {
	policyRepo repository.PolicyRepository
	auditRepo  repository.AuditRepository
	compiler   *engine.Compiler
	validator  *engine.Validator
}

func NewPolicyHandler(pr repository.PolicyRepository, ar repository.AuditRepository, comp *engine.Compiler, val *engine.Validator) *PolicyHandler {
	return &PolicyHandler{policyRepo: pr, auditRepo: ar, compiler: comp, validator: val}
}

func (h *PolicyHandler) List(w http.ResponseWriter, r *http.Request) {
	policies, err := h.policyRepo.List(100, 0)
	if err != nil {
		http.Error(w, `{"error":"failed to list policies"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, policies)
}

func (h *PolicyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, policy)
}

func (h *PolicyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"policy name is required"}`, http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	policy := &models.Policy{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		Status:      models.PolicyDraft,
		CreatedBy:   userID,
		Tags:        req.Tags,
	}

	if err := h.policyRepo.Create(policy); err != nil {
		http.Error(w, `{"error":"failed to create policy"}`, http.StatusInternalServerError)
		return
	}

	h.logAudit(r, models.AuditCreatePolicy, "policy:"+policy.Name, "created")
	writeJSON(w, http.StatusCreated, policy)
}

func (h *PolicyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}

	var req models.UpdatePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Name != "" {
		policy.Name = req.Name
	}
	if req.Description != "" {
		policy.Description = req.Description
	}
	if req.Rules != nil {
		policy.Rules = req.Rules
	}
	if req.Tags != nil {
		policy.Tags = req.Tags
	}

	policy.Status = models.PolicyDraft

	if err := h.policyRepo.Update(policy); err != nil {
		http.Error(w, `{"error":"failed to update policy"}`, http.StatusInternalServerError)
		return
	}

	h.logAudit(r, models.AuditUpdatePolicy, "policy:"+policy.Name, "updated")
	writeJSON(w, http.StatusOK, policy)
}

func (h *PolicyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}

	if err := h.policyRepo.Delete(id); err != nil {
		http.Error(w, `{"error":"failed to delete policy"}`, http.StatusInternalServerError)
		return
	}

	h.logAudit(r, models.AuditDeletePolicy, "policy:"+policy.Name, "deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (h *PolicyHandler) Validate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}

	results := h.validator.Validate(policy)
	h.logAudit(r, models.AuditValidatePolicy, "policy:"+policy.Name, "validated")
	writeJSON(w, http.StatusOK, results)
}

func (h *PolicyHandler) DryRun(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}

	driverName := r.URL.Query().Get("driver")
	if driverName == "" {
		driverName = "iptables"
	}

	rules, err := h.compiler.Compile(policy, driverName)
	if err != nil {
		http.Error(w, `{"error":"compilation failed: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	h.logAudit(r, models.AuditDryRunPolicy, "policy:"+policy.Name, "dry-run on "+driverName)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"policy_id": id,
		"driver":    driverName,
		"rules":     rules,
		"count":     len(rules),
	})
}

func (h *PolicyHandler) Apply(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid policy id"}`, http.StatusBadRequest)
		return
	}

	policy, err := h.policyRepo.FindByID(id)
	if err != nil {
		http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
		return
	}

	results := h.validator.Validate(policy)
	if !results.Valid {
		http.Error(w, `{"error":"policy validation failed","issues":`+toJSON(results.Issues)+`}`, http.StatusBadRequest)
		return
	}

	driverName := r.URL.Query().Get("driver")
	if driverName == "" {
		driverName = "iptables"
	}

	rules, err := h.compiler.Compile(policy, driverName)
	if err != nil {
		http.Error(w, `{"error":"compilation failed"}`, http.StatusBadRequest)
		return
	}

	if err := h.policyRepo.UpdateStatus(id, models.PolicyActive); err != nil {
		http.Error(w, `{"error":"failed to update policy status"}`, http.StatusInternalServerError)
		return
	}

	h.logAudit(r, models.AuditApplyPolicy, "policy:"+policy.Name, "applied on "+driverName)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"policy_id": id,
		"driver":    driverName,
		"rules":     rules,
		"count":     len(rules),
		"status":    models.PolicyActive,
	})
}

func (h *PolicyHandler) logAudit(r *http.Request, action models.AuditAction, resource, details string) {
	userID, _ := r.Context().Value("user_id").(int64)
	username, _ := r.Context().Value("username").(string)
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	go h.auditRepo.Create(&models.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Details:   details,
		IPAddress: ip,
		Success:   true,
	})
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
