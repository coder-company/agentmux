package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	Purple    = lipgloss.Color("#A78BFA")
	PurpleDim = lipgloss.Color("#7C3AED")
	Cyan      = lipgloss.Color("#67E8F9")
	Green     = lipgloss.Color("#6EE7B7")
	Red       = lipgloss.Color("#FCA5A5")
	Yellow    = lipgloss.Color("#FDE68A")
	White     = lipgloss.Color("#F9FAFB")
	Gray100   = lipgloss.Color("#F3F4F6")
	Gray300   = lipgloss.Color("#D1D5DB")
	Gray400   = lipgloss.Color("#9CA3AF")
	Gray500   = lipgloss.Color("#6B7280")
	Gray600   = lipgloss.Color("#4B5563")
	Gray700   = lipgloss.Color("#374151")
	Gray800   = lipgloss.Color("#1F2937")
	Gray900   = lipgloss.Color("#111827")
)

// Layout
var (
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Gray400).
			Padding(0, 1)

	HeaderBrand = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	HeaderDim = lipgloss.NewStyle().
			Foreground(Gray600)

	Separator = lipgloss.NewStyle().
			Foreground(Gray700)

	PanelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gray700)

	PanelBorderActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Purple)

	PanelTitle = lipgloss.NewStyle().
			Foreground(Gray400).
			Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Gray600).
			Padding(0, 1)

	FooterKey = lipgloss.NewStyle().
			Foreground(Gray300).
			Bold(true)

	FooterDesc = lipgloss.NewStyle().
			Foreground(Gray600)
)

// Session list
var (
	// Selected row — full-width inverse
	ListSelected = lipgloss.NewStyle().
			Background(PurpleDim).
			Foreground(White).
			Bold(true).
			Padding(0, 1)

	// Normal row
	ListNormal = lipgloss.NewStyle().
			Foreground(Gray300).
			Padding(0, 1)

	// Metadata on normal rows
	ListMeta = lipgloss.NewStyle().
			Foreground(Gray500)

	// Attached dot
	ListDot = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	// Cursor arrow on selected
	ListCursor = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)
)

// Preview
var (
	PreviewTitle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	PreviewPath = lipgloss.NewStyle().
			Foreground(Gray500)

	PreviewContent = lipgloss.NewStyle().
			Foreground(Gray400)

	PreviewDim = lipgloss.NewStyle().
			Foreground(Gray600)
)

// Overlay
var (
	Overlay = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Purple).
		Padding(1, 2)

	OverlayTitle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	OverlaySelected = lipgloss.NewStyle().
			Background(PurpleDim).
			Foreground(White).
			Bold(true).
			Padding(0, 1)

	OverlayNormal = lipgloss.NewStyle().
			Foreground(Gray300).
			Padding(0, 1)

	OverlayDim = lipgloss.NewStyle().
			Foreground(Gray500)

	OverlayPrompt = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	OverlayInput = lipgloss.NewStyle().
			Foreground(White)
)

// Status
var (
	StatusOk = lipgloss.NewStyle().
			Foreground(Green).
			Padding(0, 1)

	StatusErr = lipgloss.NewStyle().
			Foreground(Red).
			Padding(0, 1)

	StatusInfo = lipgloss.NewStyle().
			Foreground(Gray400).
			Padding(0, 1)

	Muted = lipgloss.NewStyle().
		Foreground(Gray500)

	Bold = lipgloss.NewStyle().
		Foreground(White).
		Bold(true)
)
