package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	currentIndexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5C57")).Bold(true)
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)



type LoginModel struct {
	Choices  []string
	Cursor   int
	Selected string
}

func RunLoginUI() (LoginModel, error) {
	p := tea.NewProgram(NewLoginModel())

	m, err := p.Run()
	if err != nil {
		return LoginModel{}, err
	}

	lm := m.(LoginModel)
	return lm, nil
}

func NewLoginModel() LoginModel {
	return LoginModel{
		Choices: []string{"Github", "Local Authentication"},
	}
}

func (m LoginModel) Init() tea.Cmd {
	return nil
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m LoginModel) View() string {

	s := strings.Builder{}

	s.WriteString(titleStyle.Render("\n🔐 Choose a login option\n"))
	s.WriteString("\n")

	for i, choice := range m.Choices {
		cursor := " "
		displayChoice := choice

		if m.Cursor == i {
			cursor = cursorStyle.Render(">")
		}

		if m.Selected != "" && m.Selected == choice {
			displayChoice = selectedStyle.Render(choice)
		} else if m.Cursor == i {
			displayChoice = currentIndexStyle.Render(choice)
		}

		s.WriteString(fmt.Sprintf("%s %s\n", cursor, displayChoice))
	}

	if m.Selected == "" {
		s.WriteString("\n↑/↓: navigate • enter: select • q: quit")
	}

	return s.String()
}
