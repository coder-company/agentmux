package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorPrimary   = lipgloss.Color("#7C3AED")
	ColorSecondary = lipgloss.Color("#6366F1")
	ColorMuted     = lipgloss.Color("#6B7280")
	ColorSuccess   = lipgloss.Color("#10B981")
	ColorDanger    = lipgloss.Color("#EF4444")
	ColorBorder    = lipgloss.Color("#374151")
	ColorBg        = lipgloss.Color("#111827")
	ColorFg        = lipgloss.Color("#F9FAFB")
)

// Styles
var (
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	ActivePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		MarginBottom(1)

	Selected = lipgloss.NewStyle().
			Foreground(ColorFg).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	Normal = lipgloss.NewStyle().
		Foreground(ColorFg).
		Padding(0, 1)

	Muted = lipgloss.NewStyle().
		Foreground(ColorMuted)

	Help = lipgloss.NewStyle().
		Foreground(ColorMuted).
		MarginTop(1)

	Status = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)
)
