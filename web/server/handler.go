package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/tormgibbs/worklogger/data"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) getSummary(w http.ResponseWriter, r *http.Request) {
	stats, err := data.GetSummaryStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get summary stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) getDailyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := data.GetDailyStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get daily stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) getWeeklyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := data.GetWeeklyStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get weekly stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) getMonthlyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := data.GetMonthlyStats(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get monthly stats", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (h *Handler) getSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := data.GetSessions(h.DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}
