package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	Accent    = lipgloss.Color("#7EE787")
	AccentDim = lipgloss.Color("#2F6F4E")
	Cyan      = lipgloss.Color("#79C0FF")
	Green     = lipgloss.Color("#8FD694")
	Red       = lipgloss.Color("#E06C75")
	Yellow    = lipgloss.Color("#D9B86C")
	White     = lipgloss.Color("#E6EDF3")
	Gray100   = lipgloss.Color("#D8DEE9")
	Gray300   = lipgloss.Color("#B6C2CF")
	Gray400   = lipgloss.Color("#8C98A7")
	Gray500   = lipgloss.Color("#6F7B89")
	Gray600   = lipgloss.Color("#505A66")
	Gray700   = lipgloss.Color("#343D49")
	Gray800   = lipgloss.Color("#202734")
	Gray850   = lipgloss.Color("#181E27")
	Gray900   = lipgloss.Color("#11161D")
)

// Layout
var (
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Background(Gray900)

	HeaderBrand = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	HeaderMode = lipgloss.NewStyle().
			Foreground(Gray900).
			Background(Accent).
			Bold(true).
			Padding(0, 1)

	HeaderDim = lipgloss.NewStyle().
			Foreground(Gray500)

	Separator = lipgloss.NewStyle().
			Foreground(Gray800)

	PanelBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Gray700)

	PanelBorderActive = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(AccentDim)

	PanelTitle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	PanelMeta = lipgloss.NewStyle().
			Foreground(Gray500)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Background(Gray900)

	FooterKey = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	FooterDesc = lipgloss.NewStyle().
			Foreground(Gray500)
)

// Session list
var (
	// Selected row — full-width inverse
	ListSelected = lipgloss.NewStyle().
			Background(Gray800).
			Foreground(White).
			Bold(true)

	// Normal row
	ListNormal = lipgloss.NewStyle().
			Foreground(Gray300)

	// Metadata on normal rows
	ListMeta = lipgloss.NewStyle().
			Foreground(Gray500)

	// Attached dot
	ListDot = lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true)

	// Cursor arrow on selected
	ListCursor = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)
)

// Preview
var (
	PreviewTitle = lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

	PreviewPath = lipgloss.NewStyle().
			Foreground(Gray500)

	PreviewContent = lipgloss.NewStyle().
			Foreground(Gray300)

	PreviewDim = lipgloss.NewStyle().
			Foreground(Gray600)
)

// Overlay
var (
	Overlay = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(AccentDim).
		Padding(1, 2)

	OverlayTitle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	OverlaySelected = lipgloss.NewStyle().
			Background(Gray800).
			Foreground(White).
			Bold(true)

	OverlayNormal = lipgloss.NewStyle().
			Foreground(Gray300)

	OverlayDim = lipgloss.NewStyle().
			Foreground(Gray500)

	OverlayPrompt = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	OverlayInput = lipgloss.NewStyle().
			Foreground(White)
)

// Status
var (
	StatusOk = lipgloss.NewStyle().
			Foreground(Green).
			Background(Gray900)

	StatusErr = lipgloss.NewStyle().
			Foreground(Red).
			Background(Gray900)

	StatusInfo = lipgloss.NewStyle().
			Foreground(Gray400).
			Background(Gray900)

	StatusWarn = lipgloss.NewStyle().
			Foreground(Yellow).
			Background(Gray900)

	ErrorText = lipgloss.NewStyle().
			Foreground(Red)

	WarnText = lipgloss.NewStyle().
			Foreground(Yellow)

	GoodText = lipgloss.NewStyle().
			Foreground(Green)

	Muted = lipgloss.NewStyle().
		Foreground(Gray500)

	Bold = lipgloss.NewStyle().
		Foreground(White).
		Bold(true)

	Subtle = lipgloss.NewStyle().
		Foreground(Gray600)
)

// Truncate returns s constrained to width terminal cells.
func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}

	ellipsis := "…"
	ellipsisW := lipgloss.Width(ellipsis)
	if width <= ellipsisW {
		return ellipsis
	}

	limit := width - ellipsisW
	var b strings.Builder
	used := 0
	for _, r := range s {
		part := string(r)
		partW := lipgloss.Width(part)
		if used+partW > limit {
			break
		}
		b.WriteRune(r)
		used += partW
	}
	return b.String() + ellipsis
}

// PadRight truncates s if necessary and pads it to width terminal cells.
func PadRight(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = Truncate(s, width)
	pad := width - lipgloss.Width(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}
