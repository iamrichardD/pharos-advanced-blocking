package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/iamrichardd/pharos-advanced-blocking/internal/config"
)

// Brand Colors aligned with Pharos aesthetics
var (
	pharosBlue  = lipgloss.Color("#005f87") // Sleek blue
	textColor   = lipgloss.Color("#d0d0d0")
	accentColor = lipgloss.Color("#5fafd7")
	borderBlue  = lipgloss.Color("#0087af")
	errorColor  = lipgloss.Color("#d70000")
)

// Visual Styles
var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderBlue).
		Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(pharosBlue).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		Underline(true)

	rowStyle = lipgloss.NewStyle().
		Foreground(textColor)

	altRowStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9e9e9e"))

	errorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)
		
	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#767676")).
		MarginTop(1)
)

// ClientEntry represents a single client mapping to a filtering group.
type ClientEntry struct {
	IP    string
	Group string
}

// Model defines the state for the Bubble Tea TUI
type Model struct {
	clients     []ClientEntry
	filtered    []ClientEntry
	searchInput textinput.Model
	width       int
	height      int
	err         error
	ready       bool
}

// New creates and initializes a new TUI model, preparing the text input.
func New(cfg *config.Config) *Model {
	ti := textinput.New()
	ti.Placeholder = "Search IP or Group..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	var clients []ClientEntry
	if cfg != nil {
		for ip, group := range cfg.NetworkGroupMap {
			clients = append(clients, ClientEntry{IP: ip, Group: group})
		}
	}

	// Sort clients by IP for consistent display
	slices.SortFunc(clients, func(a, b ClientEntry) int {
		return cmp.Compare(a.IP, b.IP)
	})

	m := &Model{
		clients:     clients,
		searchInput: ti,
		ready:       true,
	}
	m.filterClients()
	return m
}

// Init initializes the Bubble Tea application and triggers configuration loading.
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming events (key presses, window resizing) and state changes.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle standard viewport size updates to gracefully scale the TUI
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		// Forward key presses to the search text input component
		var tiCmd tea.Cmd
		m.searchInput, tiCmd = m.searchInput.Update(msg)
		
		// Update filter results dynamically on every keystroke
		m.filterClients()
		cmd = tea.Batch(cmd, tiCmd)
	}

	return m, cmd
}

// filterClients updates the filtered client list based on the search input query.
func (m *Model) filterClients() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filtered = make([]ClientEntry, len(m.clients))
		copy(m.filtered, m.clients)
		return
	}

	var filtered []ClientEntry
	for _, c := range m.clients {
		if strings.Contains(strings.ToLower(c.IP), query) || strings.Contains(strings.ToLower(c.Group), query) {
			filtered = append(filtered, c)
		}
	}
	m.filtered = filtered
}

// View renders the TUI into a beautifully formatted string using Lipgloss.
func (m *Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error loading config: %v", m.err))
	}

	if !m.ready {
		return "Loading Pharos Config..."
	}

	title := titleStyle.Render("Pharos Advanced Blocking TUI")

	// Render the dynamic search bar
	searchBox := lipgloss.NewStyle().MarginBottom(1).Render(m.searchInput.View())

	// Render the ASCII table
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf("%-30s %s", "Client IP", "Group")) + "\n")
	
	for i, c := range m.filtered {
		row := fmt.Sprintf("%-30s %s", c.IP, c.Group)
		// Alternate row coloring for sleek UI readability
		if i%2 == 0 {
			b.WriteString(rowStyle.Render(row) + "\n")
		} else {
			b.WriteString(altRowStyle.Render(row) + "\n")
		}
	}

	if len(m.filtered) == 0 {
		b.WriteString(altRowStyle.Render("No clients found matching the search.") + "\n")
	}

	// Add footer help status line
	footer := footerStyle.Render("ctrl+c / esc: exit | start typing to search")

	content := lipgloss.JoinVertical(lipgloss.Left, title, searchBox, b.String(), footer)

	// Apply responsive padding and borders to wrap the layout
	if m.width > 0 {
		return baseStyle.Width(m.width - 4).Render(content)
	}

	return baseStyle.Render(content)
}
