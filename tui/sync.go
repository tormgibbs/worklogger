package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	SyncOptionExisting     = "Associate with existing task"
	SyncOptionUnassociated = "Leave commits unassociated"
	SyncOptionNew          = "Create a new task"
	SyncOptionCancel       = "Cancel sync"
)

var SyncChoices = []string{
	SyncOptionExisting,
	SyncOptionNew,
	SyncOptionUnassociated,
	SyncOptionCancel,
}

var syncTooltips = map[string]string{
	"Associate with existing task": "Pick a task you’ve already started to sync the commits with it.",
	"Create a new task":            "Start fresh with a new task and sync the commits to it.",
	"Leave commits unassociated":   "Skip attaching commits to any task — just log them raw.",
	"Cancel sync":                  "Abort the sync process. No changes will be made.",
}

type SyncModel struct {
	Title    string
	Choices  []string
	Cursor   int
	Selected string
}

func RunSyncUI(title string) (SyncModel, error) {
	p := tea.NewProgram(NewSyncModel(title))
	m, err := p.Run()
	if err != nil {
		return SyncModel{}, err
	}

	return m.(SyncModel), nil
}

func NewSyncModel(title string) SyncModel {
	return SyncModel{
		Title:   title,
		Choices: SyncChoices,
	}
}

func (m SyncModel) Init() tea.Cmd {
	return nil
}

func (m SyncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter":
			m.Selected = m.Choices[m.Cursor]
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SyncModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("\n%s\n", m.Title)))
	b.WriteString("\n")

	for i, choice := range m.Choices {
		cursor := " "
		displayChoice := choice

		if m.Cursor == i {
			cursor = cursorStyle.Render("›")
		}

		if m.Selected == choice {
			displayChoice = selectedStyle.Render(choice)
		} else if m.Cursor == i {
			displayChoice = currentIndexStyle.Render(choice)
		}

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, displayChoice))
	}

	if m.Selected == "" {
		currentChoice := m.Choices[m.Cursor]
		tooltip := syncTooltips[currentChoice]

		b.WriteString("\n")
		b.WriteString(hintStyle.Render("↑/↓: navigate • enter: select • q: quit") + "\n")
		b.WriteString(tooltipStyle.Render("💡 "+tooltip) + "\n")
	}

	return b.String()
}
