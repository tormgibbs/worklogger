package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskDescModel struct {
	Input    textinput.Model
	Quitting bool
}

func NewTaskDescModel() TaskDescModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. Fixing login bug"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 40

	return TaskDescModel{Input: ti}
}

func (m TaskDescModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TaskDescModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit

		case "enter":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m TaskDescModel) View() string {
	if m.Quitting {
		return "Cancelled.\n"
	}

	title := lipgloss.NewStyle().Bold(true).Render("\n📝 Enter a description for the new task:")
	return title + "\n\n" + m.Input.View() + "\n\n(press Enter to continue, Esc to cancel)"
}

func RunNewTaskUI() (string, error) {
	p := tea.NewProgram(NewTaskDescModel())
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	model := m.(TaskDescModel)
	if model.Quitting {
		return "", nil
	}
	return model.Input.Value(), nil
}
