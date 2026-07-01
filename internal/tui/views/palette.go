package views

import (
	"fmt"
	"strings"

	"agentmux/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// Action is a searchable command palette entry.
type Action struct {
	Name string
	Desc string
	Key  string
	Do   func()
}

// PaletteView is the fuzzy command palette.
type PaletteView struct {
	Actions  []Action
	Filtered []Action
	Query    string
	Cursor   int
	Width    int
	Height   int
}

// NewPalette creates a command palette with the given actions.
func NewPalette(actions []Action) *PaletteView {
	return &PaletteView{
		Actions:  actions,
		Filtered: actions,
	}
}

// UpdateActions replaces the action list.
func (p *PaletteView) UpdateActions(actions []Action) {
	p.Actions = actions
	p.SetQuery(p.Query)
}

// SetQuery updates the search filter.
func (p *PaletteView) SetQuery(q string) {
	p.Query = q
	p.Cursor = 0
	if q == "" {
		p.Filtered = p.Actions
		return
	}
	lower := strings.ToLower(q)
	var filtered []Action
	for _, a := range p.Actions {
		if fuzzyContains(strings.ToLower(a.Name), lower) ||
			strings.Contains(strings.ToLower(a.Desc), lower) {
			filtered = append(filtered, a)
		}
	}
	p.Filtered = filtered
}

func fuzzyContains(haystack, needle string) bool {
	hi := 0
	for _, ch := range needle {
		found := false
		for hi < len(haystack) {
			if rune(haystack[hi]) == ch {
				hi++
				found = true
				break
			}
			hi++
		}
		if !found {
			return false
		}
	}
	return true
}

// TypeChar appends a character to the query.
func (p *PaletteView) TypeChar(ch rune) {
	p.SetQuery(p.Query + string(ch))
}

// Backspace removes the last character.
func (p *PaletteView) Backspace() {
	if len(p.Query) > 0 {
		p.SetQuery(p.Query[:len(p.Query)-1])
	}
}

// Selected returns the highlighted action, or nil.
func (p *PaletteView) Selected() *Action {
	if len(p.Filtered) == 0 {
		return nil
	}
	if p.Cursor >= len(p.Filtered) {
		p.Cursor = len(p.Filtered) - 1
	}
	return &p.Filtered[p.Cursor]
}

// MoveUp moves the cursor up.
func (p *PaletteView) MoveUp() {
	if p.Cursor > 0 {
		p.Cursor--
	}
}

// MoveDown moves the cursor down.
func (p *PaletteView) MoveDown() {
	if p.Cursor < len(p.Filtered)-1 {
		p.Cursor++
	}
}

// Render returns the palette as a centered overlay.
func (p *PaletteView) Render() string {
	w := p.Width * 60 / 100
	if w < 48 {
		w = 48
	}
	if w > 84 {
		w = 84
	}
	if p.Width > 0 && w > p.Width-4 {
		w = p.Width - 4
	}
	if w < 24 {
		w = 24
	}

	innerW := w - 6
	if innerW < 12 {
		innerW = 12
	}

	maxItems := p.Height - 11
	if maxItems < 3 {
		maxItems = 3
	}
	if maxItems > 10 {
		maxItems = 10
	}

	title := styles.OverlayTitle.Render("Commands")
	count := styles.OverlayDim.Render(resultCount(len(p.Filtered), len(p.Actions)))
	titleRow := title
	gap := innerW - lipgloss.Width(title) - lipgloss.Width(count)
	if gap > 0 {
		titleRow = title + lipgloss.NewStyle().Width(gap).Render("") + count
	}

	query := p.Query
	if query == "" {
		query = styles.OverlayDim.Render("type to filter")
	} else {
		query = styles.OverlayInput.Render(styles.Truncate(query, innerW-4))
	}
	prompt := styles.OverlayPrompt.Render("› ") + query + styles.OverlayPrompt.Render("█")

	var items string
	start, visible := p.visibleActions(maxItems)

	for i, a := range visible {
		selected := start+i == p.Cursor
		items += renderActionRow(a, innerW, selected) + "\n"
	}

	if len(p.Filtered) == 0 {
		items = styles.Muted.Render(styles.Truncate("  No commands match \""+p.Query+"\"", innerW))
	}

	footer := styles.HeaderDim.Render(styles.Truncate("type to filter · ↑↓ move · enter run · esc close", innerW))
	content := titleRow + "\n" + prompt + "\n\n" + items + "\n" + footer
	box := styles.Overlay.Width(w).Render(content)

	return lipgloss.Place(p.Width, p.Height,
		lipgloss.Center, lipgloss.Center, box)
}

func (p *PaletteView) visibleActions(maxItems int) (int, []Action) {
	if len(p.Filtered) <= maxItems {
		return 0, p.Filtered
	}

	start := p.Cursor - maxItems/2
	if start < 0 {
		start = 0
	}
	if start+maxItems > len(p.Filtered) {
		start = len(p.Filtered) - maxItems
	}
	return start, p.Filtered[start : start+maxItems]
}

func renderActionRow(a Action, width int, selected bool) string {
	key := ""
	if a.Key != "" {
		key = "[" + a.Key + "]"
	}

	keyW := lipgloss.Width(key)
	nameW := width
	descW := 0
	if width >= 54 && a.Desc != "" {
		nameW = 24
		descW = width - nameW - keyW - 4
		if descW < 12 {
			descW = 0
			nameW = width - keyW - 2
		}
	} else if keyW > 0 {
		nameW = width - keyW - 2
	}
	if nameW < 1 {
		nameW = 1
	}

	name := styles.PadRight(a.Name, nameW)
	line := name
	if descW > 0 {
		line += "  " + styles.PadRight(a.Desc, descW)
	}
	if key != "" {
		gap := width - lipgloss.Width(line) - keyW
		if gap < 1 {
			gap = 1
		}
		line += strings.Repeat(" ", gap) + key
	}
	line = styles.Truncate(line, width)

	if selected {
		return styles.OverlaySelected.Width(width).Render(line)
	}

	if descW > 0 {
		namePart := styles.OverlayNormal.Render(styles.PadRight(a.Name, nameW))
		descPart := styles.OverlayDim.Render(styles.PadRight(a.Desc, descW))
		line = namePart + "  " + descPart
		if key != "" {
			gap := width - lipgloss.Width(line) - keyW
			if gap < 1 {
				gap = 1
			}
			line += strings.Repeat(" ", gap) + styles.OverlayDim.Render(key)
		}
		return line
	}

	return styles.OverlayNormal.Render(line)
}

func resultCount(filtered, total int) string {
	if total == 0 {
		return "no actions"
	}
	if filtered == total {
		return fmt.Sprintf("%d actions", total)
	}
	return fmt.Sprintf("%d of %d", filtered, total)
}
