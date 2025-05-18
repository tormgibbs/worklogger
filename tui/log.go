package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Commit struct {
	Message string
}

type WorkSession struct {
	StartTime string
	EndTime   string
	Task      string
	Duration  string
	Commits   []Commit
}

type DayLog struct {
	Date     string
	Sessions []WorkSession
}

type LogModel struct {
	DayLogs []DayLog
}

// func NewLogModel() LogModel {
// 	sessions := db.GetSessionsWithCommits() // your custom function
// 	return LogModel{DayLogs: formatIntoDayLogs(sessions)}
// }

func NewLogModel() LogModel {
	// Load your real data here instead
	return LogModel{
		DayLogs: []DayLog{
			{
				Date: "May 15",
				Sessions: []WorkSession{
					{
						StartTime: "10:32",
						EndTime:   "12:03",
						Task:      "Implement user login",
						Duration:  "1h 31m",
						Commits: []Commit{
							{Message: "feat: login page UI + validation"},
							{Message: "fix: auth token refresh"},
						},
					},
					{
						StartTime: "14:15",
						EndTime:   "15:45",
						Task:      "Refactor session API",
						Duration:  "1h 30m",
						Commits: []Commit{
							{Message: "refactor: cleaned up controller logic"},
						},
					},
				},
			},
			{
				Date: "May 14",
				Sessions: []WorkSession{
					{
						StartTime: "11:00",
						EndTime:   "12:00",
						Task:      "Write README docs",
						Duration:  "1h 0m",
						Commits: []Commit{
							{Message: "docs: initial usage guide"},
						},
					},
				},
			},
		},
	}
}

func (m LogModel) Init() tea.Cmd {
	return nil
}

func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m LogModel) View() string {
	var b strings.Builder

	b.WriteString("📅 Work Sessions\n\n")

	for _, log := range m.DayLogs {
		b.WriteString(fmt.Sprintf("[%s]\n", log.Date))

		for _, s := range log.Sessions {
			b.WriteString(fmt.Sprintf("🕒 %s - %s | Task: \"%s\" | ⏱ %s\n",
				s.StartTime, s.EndTime, s.Task, s.Duration))

			if len(s.Commits) > 0 {
				b.WriteString("  - Commits:\n")
				for _, c := range s.Commits {
					b.WriteString(fmt.Sprintf("    ✔ %s\n", c.Message))
				}
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("(press 'q' to quit)\n")
	return b.String()
}

// func formatIntoDayLogs(sessions []WorkSession) []DayLog {
// 	grouped := make(map[string][]WorkSession)

// 	// group by date
// 	for _, s := range sessions {
// 		grouped[s.Date] = append(grouped[s.Date], s)
// 	}

// 	// convert map into sorted slice
// 	var dayLogs []DayLog
// 	for date, sess := range grouped {
// 		dayLogs = append(dayLogs, DayLog{
// 			Date:     date,
// 			Sessions: sess,
// 		})
// 	}

// 	// optional: sort by date if you want
// 	// sort.Slice(dayLogs, func(i, j int) bool {
// 	//     return dayLogs[i].Date > dayLogs[j].Date // newest first
// 	// })

// 	return dayLogs
// }
