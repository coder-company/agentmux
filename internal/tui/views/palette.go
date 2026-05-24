package views

import (
	"strings"

	"agentmux/internal/tui/styles"
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
	Active   bool
}

// NewPalette creates a command palette with the given actions.
func NewPalette(actions []Action) *PaletteView {
	return &PaletteView{
		Actions:  actions,
		Filtered: actions,
	}
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
		if strings.Contains(strings.ToLower(a.Name), lower) ||
			strings.Contains(strings.ToLower(a.Desc), lower) {
			filtered = append(filtered, a)
		}
	}
	p.Filtered = filtered
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

// Render returns the palette overlay.
func (p *PaletteView) Render() string {
	w := p.Width / 2
	if w < 40 {
		w = 40
	}

	prompt := styles.Status.Render("> ") + p.Query + "█"
	var items string
	for i, a := range p.Filtered {
		label := a.Name
		if a.Key != "" {
			label += "  " + styles.Muted.Render("["+a.Key+"]")
		}
		if i == p.Cursor {
			items += styles.Selected.Width(w-4).Render(label) + "\n"
		} else {
			items += styles.Normal.Width(w-4).Render(label) + "\n"
		}
	}
	if len(p.Filtered) == 0 {
		items = styles.Muted.Render("  No matches")
	}

	content := styles.Title.Render("Command Palette") + "\n" + prompt + "\n\n" + items
	return styles.ActivePanel.Width(w).Render(content)
}
