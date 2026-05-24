package views

import (
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
	w := p.Width * 50 / 100
	if w < 44 {
		w = 44
	}
	if w > 72 {
		w = 72
	}

	maxItems := p.Height - 10
	if maxItems < 3 {
		maxItems = 3
	}

	title := styles.OverlayTitle.Render("Commands")
	prompt := styles.OverlayPrompt.Render("❯ ") +
		styles.OverlayInput.Render(p.Query) +
		styles.OverlayPrompt.Render("│")

	var items string
	visible := p.Filtered
	if len(visible) > maxItems {
		visible = visible[:maxItems]
	}

	for i, a := range visible {
		name := a.Name
		if a.Key != "" {
			name += "  " + styles.OverlayDim.Render(a.Key)
		}

		if i == p.Cursor {
			items += styles.OverlaySelected.Width(w-6).Render(name) + "\n"
		} else {
			items += styles.OverlayNormal.Render(name) + "\n"
		}
	}

	if len(p.Filtered) == 0 {
		items = styles.Muted.Render("  no matches")
	}

	footer := styles.HeaderDim.Render("↑↓ move · ⏎ select · esc close")
	content := title + "\n" + prompt + "\n\n" + items + "\n" + footer
	box := styles.Overlay.Width(w).Render(content)

	return lipgloss.Place(p.Width, p.Height,
		lipgloss.Center, lipgloss.Center, box)
}
