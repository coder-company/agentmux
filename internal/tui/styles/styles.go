package styles

import "github.com/charmbracelet/lipgloss"

// Colors - a refined dark theme with purple accents.
var (
	ColorPrimary    = lipgloss.Color("#A78BFA") // violet-400
	ColorAccent     = lipgloss.Color("#7C3AED") // violet-600
	ColorMuted      = lipgloss.Color("#6B7280") // gray-500
	ColorSubtle     = lipgloss.Color("#4B5563") // gray-600
	ColorSuccess    = lipgloss.Color("#34D399") // emerald-400
	ColorDanger     = lipgloss.Color("#F87171") // red-400
	ColorWarning    = lipgloss.Color("#FBBF24") // amber-400
	ColorBorder     = lipgloss.Color("#374151") // gray-700
	ColorBorderDim  = lipgloss.Color("#1F2937") // gray-800
	ColorFg         = lipgloss.Color("#F9FAFB") // gray-50
	ColorFgDim      = lipgloss.Color("#D1D5DB") // gray-300
	ColorBg         = lipgloss.Color("#111827") // gray-900
	ColorBgElevated = lipgloss.Color("#1F2937") // gray-800
)

// Layout styles
var (
	// Header bar at the top
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorFg).
		Padding(0, 1)

	HeaderBrand = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	HeaderMeta = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Panel containers
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	PanelActive = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1)

	PanelTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	PanelTitleDim = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMuted)

	// Footer help bar
	Footer = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1)

	FooterKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgDim)

	FooterDesc = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	FooterSep = lipgloss.NewStyle().
			Foreground(ColorBorderDim)
)

// Content styles
var (
	// List items
	Selected = lipgloss.NewStyle().
			Foreground(ColorFg).
			Background(ColorAccent).
			Bold(true).
			Padding(0, 1)

	Normal = lipgloss.NewStyle().
		Foreground(ColorFgDim).
		Padding(0, 1)

	NormalDim = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	// Indicators
	Attached = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	Detached = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	// Text
	Muted = lipgloss.NewStyle().
		Foreground(ColorMuted)

	Subtle = lipgloss.NewStyle().
		Foreground(ColorSubtle)

	Success = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	Danger = lipgloss.NewStyle().
		Foreground(ColorDanger)

	Warning = lipgloss.NewStyle().
		Foreground(ColorWarning)

	Bold = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorFg)

	// Status bar
	StatusBar = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	StatusSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Padding(0, 1)

	StatusError = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Padding(0, 1)

	// Overlay/modal
	Overlay = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2)

	// Input prompt
	Prompt = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	Input = lipgloss.NewStyle().
		Foreground(ColorFg)

	Cursor = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)
)
