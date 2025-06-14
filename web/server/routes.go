package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes(h *Handler) http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/summary", h.getSummary)
	router.HandlerFunc(http.MethodGet, "/api/stats/daily", h.getDailyStats)
	router.HandlerFunc(http.MethodGet, "/api/stats/weekly", h.getWeeklyStats)
	router.HandlerFunc(http.MethodGet, "/api/stats/monthly", h.getMonthlyStats)
	router.HandlerFunc(http.MethodGet, "/api/sessions", h.getSessions)
	router.HandlerFunc(http.MethodGet, "/api/export.csv", h.exportAllDataCSV)

	return enableCORS(router)
}
