package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	GitHubOAuth = "GitHub"
	LocalAuth   = "Local"
)

var AuthChoices = []string{GitHubOAuth, LocalAuth}

type AuthModel struct {
	Title    string
	Choices  []string
	Cursor   int
	Selected string
}

func RunAuthUI(title string) (AuthModel, error) {
	p := tea.NewProgram(NewAuthModel(title))

	m, err := p.Run()
	if err != nil {
		return AuthModel{}, err
	}

	return m.(AuthModel), nil
}

func NewAuthModel(title string) AuthModel {
	return AuthModel{
		Choices: AuthChoices,
		Title: title,
	}
}

func (m AuthModel) Init() tea.Cmd {
	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m AuthModel) View() string {

	s := strings.Builder{}

	// s.WriteString(titleStyle.Render("\n🔐 Choose a login option\n"))
	s.WriteString(titleStyle.Render(fmt.Sprintf("\n%s\n", m.Title)))
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
