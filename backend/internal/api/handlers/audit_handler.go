package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anomalyco/iptables-visualizer/internal/models"
	"github.com/anomalyco/iptables-visualizer/internal/repository"
)

type AuditHandler struct {
	auditRepo repository.AuditRepository
}

func NewAuditHandler(ar repository.AuditRepository) *AuditHandler {
	return &AuditHandler{auditRepo: ar}
}

func (h *AuditHandler) Query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := models.AuditQuery{
		Limit:  50,
		Offset: 0,
	}

	if userID := q.Get("user_id"); userID != "" {
		id, err := parseInt64(userID)
		if err == nil {
			query.UserID = id
		}
	}
	if action := q.Get("action"); action != "" {
		query.Action = models.AuditAction(action)
	}
	if resource := q.Get("resource"); resource != "" {
		query.Resource = resource
	}
	if from := q.Get("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err == nil {
			query.From = &t
		}
	}
	if to := q.Get("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err == nil {
			query.To = &t
		}
	}
	if limit := q.Get("limit"); limit != "" {
		if l, err := parseInt64(limit); err == nil && l > 0 {
			query.Limit = int(l)
		}
	}
	if offset := q.Get("offset"); offset != "" {
		if o, err := parseInt64(offset); err == nil && o >= 0 {
			query.Offset = int(o)
		}
	}

	logs, err := h.auditRepo.Query(query)
	if err != nil {
		http.Error(w, `{"error":"failed to query audit logs"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, logs)
}

func parseInt64(s string) (int64, error) {
	var i int64
	for _, b := range []byte(s) {
		if b < '0' || b > '9' {
			return 0, json.Unmarshal([]byte(s), &i)
		}
	}
	return 0, json.Unmarshal([]byte(s), &i)
}
