package ui

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var Theme = huh.ThemeCharm(true)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2"))

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))
)
