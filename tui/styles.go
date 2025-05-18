package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	currentIndexStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5C57")).Bold(true)
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	hintStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	tooltipStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true).PaddingTop(1)

	taskTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).Underline(true)
	taskCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	taskActiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
	taskHintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)
