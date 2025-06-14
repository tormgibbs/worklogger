package server

import (
	"database/sql"
	"encoding/csv"
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

func (h *Handler) exportAllDataCSV(w http.ResponseWriter, r *http.Request) {
	summary, err := data.GetSummaryStats(h.DB)
	if err != nil {
		http.Error(w, "Failed to get summary", http.StatusInternalServerError)
		return
	}

	daily, err := data.GetDailyStats(h.DB)
	if err != nil {
		http.Error(w, "Failed to get daily stats", http.StatusInternalServerError)
		return
	}

	weekly, err := data.GetWeeklyStats(h.DB)
	if err != nil {
		http.Error(w, "Failed to get weekly stats", http.StatusInternalServerError)
		return
	}

	monthly, err := data.GetMonthlyStats(h.DB)
	if err != nil {
		http.Error(w, "Failed to get monthly stats", http.StatusInternalServerError)
		return
	}

	sessions, err := data.GetSessions(h.DB)
	if err != nil {
		http.Error(w, "Failed to get sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
	w.Header().Set("Content-Type", "text/csv")
	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Section", "Metric", "Value", "Change"})
	writer.Write([]string{"Summary", "Today's Hours", fmt.Sprintf("%v", summary.TodayHours.Value), fmt.Sprintf("%v", summary.TodayHours.Change)})
	writer.Write([]string{"Summary", "Week Hours", fmt.Sprintf("%v", summary.WeekHours.Value), fmt.Sprintf("%v", summary.WeekHours.Change)})
	writer.Write([]string{"Summary", "Sessions Today", fmt.Sprintf("%v", summary.SessionsToday.Value), fmt.Sprintf("%v", summary.SessionsToday.Change)})
	writer.Write([]string{"Summary", "Productivity Score", fmt.Sprintf("%v", summary.ProductivityScore.Value), fmt.Sprintf("%v", summary.ProductivityScore.Change)})

	writer.Write([]string{}) // empty row
	writer.Write([]string{"Section", "Date", "Hours", "Sessions"})
	for _, d := range daily {
		writer.Write([]string{"Daily", d.Date, fmt.Sprintf("%.2f", d.Hours), fmt.Sprintf("%d", d.Sessions)})
	}

	writer.Write([]string{})
	writer.Write([]string{"Section", "Week Start", "Hours", "Sessions"})
	for _, wStat := range weekly {
		writer.Write([]string{"Weekly", wStat.Start, fmt.Sprintf("%.2f", wStat.Hours), fmt.Sprintf("%d", wStat.Sessions)})
	}

	writer.Write([]string{})
	writer.Write([]string{"Section", "Month", "Hours", "Sessions"})
	for _, m := range monthly {
		writer.Write([]string{"Monthly", m.Month, fmt.Sprintf("%.2f", m.Hours), fmt.Sprintf("%d", m.Sessions)})
	}

	writer.Write([]string{})
	writer.Write([]string{"Section", "Task", "Start Time", "End Time", "Duration", "Status"})
	for _, s := range sessions {
		endTime := ""
		if s.EndTime != nil {
			endTime = *s.EndTime
		}

		writer.Write([]string{
			"Session",
			s.Task,
			s.StartTime,
			endTime,
			s.Duration,
			s.Status,
		})
	}
}
