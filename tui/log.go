package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tormgibbs/worklogger/data"
)

type LogModel struct {
	Logs []data.Log
}

func RunLogUI(logs []data.Log) (LogModel, error) {
	p := tea.NewProgram(NewLogModel(logs))
	m, err := p.Run()
	if err != nil {
		return LogModel{}, err
	}

	return m.(LogModel), nil
}

func NewLogModel(logs []data.Log) LogModel {
	return LogModel{logs}
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

	for _, log := range m.Logs {
		b.WriteString(fmt.Sprintf("[%s]\n", log.Date))

		for _, s := range log.Sessions {
			start := s.StartedAt.Format("15:04")

			end := "ongoing"
			if !s.EndedAt.IsZero() {
				end = s.EndedAt.Format("15:04")
			}

			duration := fmtDuration(s.TotalTime)

			b.WriteString(fmt.Sprintf("🕒 %s - %s | Task: \"%s\" | ⏱ %s\n",
				start, end, s.Task, duration))

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

func fmtDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", h, m)
}
