package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GET /api/report/hari-ini
func (h *ReportHandler) HandleTodayReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.service.GetTodaySummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET /api/report?start_date=2026-01-01&end_date=2026-02-01
func (h *ReportHandler) HandleReportByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")

	if start == "" || end == "" {
		http.Error(w, "start_date and end_date required", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetSummaryByDate(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
