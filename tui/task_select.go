package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tormgibbs/worklogger/data"
)

// type Session struct {
// 	ID          string
// 	Description string
// }

type TaskSelectModel struct {
	Title    string
	Sessions []*data.DetailedTaskSession
	Cursor   int
	Selected *data.DetailedTaskSession
}

func RunTaskSelectUI(title string, sessions []*data.DetailedTaskSession) (*data.DetailedTaskSession, error) {
	p := tea.NewProgram(NewTaskSelectModel(title, sessions))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	selected := m.(TaskSelectModel).Selected
	if selected == nil {
		return nil, fmt.Errorf("no task selected")
	}

	return selected, nil
}

func NewTaskSelectModel(title string, sessions []*data.DetailedTaskSession) TaskSelectModel {
	return TaskSelectModel{
		Title:    title,
		Sessions: sessions,
	}
}

func (m TaskSelectModel) Init() tea.Cmd {
	return nil
}

func (m TaskSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Sessions)-1 {
				m.Cursor++
			}
		case "enter":
			m.Selected = m.Sessions[m.Cursor]
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m TaskSelectModel) View() string {
	var b strings.Builder

	b.WriteString(taskTitleStyle.Render(fmt.Sprintf("\n%s\n", m.Title)))
	b.WriteString("\n")

	for i, session := range m.Sessions {
		cursor := " "
		display := session.Task.Description

		if m.Cursor == i {
			cursor = taskCursorStyle.Render("›")
			display = taskActiveStyle.Render(session.Task.Description)
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, display))
	}

	if m.Selected == nil {
		b.WriteString("\n")
		b.WriteString(taskHintStyle.Render("↑/↓: navigate • enter: select • q: quit") + "\n")
	}

	return b.String()
}
