package handlers

import (
	"net/http"

	"github.com/anomalyco/iptables-visualizer/internal/capture"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) GetFirewall(w http.ResponseWriter, r *http.Request) {
	state, err := capture.Capture()
	if err != nil {
		http.Error(w, `{"error":"failed to capture firewall state: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, state)
}
