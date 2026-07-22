package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/iamrichardd/pharos-advanced-blocking/internal/config"
)

// ContentType represents the type of content being rendered
type ContentType int

const (
	ContentTypeEmpty ContentType = iota
	ContentTypeTable
	ContentTypeHelp
	ContentTypeError
	ContentTypeCommandList
	ContentTypeViewNetworkGroupMap
	ContentTypeViewGroups
	ContentTypeViewGroup
)

// String returns the string representation for debugging
func (ct ContentType) String() string {
	switch ct {
	case ContentTypeEmpty:
		return "empty"
	case ContentTypeTable:
		return "table"
	case ContentTypeHelp:
		return "help"
	case ContentTypeError:
		return "error"
	case ContentTypeCommandList:
		return "command_list"
	case ContentTypeViewNetworkGroupMap:
		return "view_networkgroupmap"
	case ContentTypeViewGroups:
		return "view_groups"
	case ContentTypeViewGroup:
		return "view_group"
	default:
		return "unknown"
	}
}

// SlashCommand represents a command with name, aliases, and description
type SlashCommand struct {
	Name        string
	Aliases     []string
	Description string
}

// Slash command registry with all available commands
var commands = []SlashCommand{
	{
		Name:        "/help",
		Aliases:     []string{"/?", "/h"},
		Description: "Show available commands",
	},
	{
		Name:        "/exit",
		Aliases:     []string{"/quit", "/q"},
		Description: "Exit the TUI",
	},
	{
		Name:        "/clear",
		Aliases:     []string{"/c"},
		Description: "Clear search and reset",
	},
	{
		Name:        "/view",
		Aliases:     []string{"/v"},
		Description: "View network mappings, groups, or group details",
	},
}

// ViewSubcommand represents a subcommand for /view
type ViewSubcommand struct {
	Name        string
	Description string
}

// View subcommands registry
var viewSubcommands = []ViewSubcommand{
	{Name: "groups", Description: "List all groups with device counts"},
	{Name: "group", Description: "Show details for a specific group (followed by group name)"},
	{Name: "networkGroupMap", Description: "Show all IP-to-group mappings"},
}

// CommandEvent represents a single command execution in the history
type CommandEvent struct {
	Timestamp time.Time
	Command   string   // e.g., "/view groups"
	Output    string   // multi-line output
	Lines     []string // split output for rendering
}

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
	clients             []ClientEntry
	filtered            []ClientEntry
	unifiedInput        textinput.Model
	width               int
	height              int
	err                 error
	ready               bool
	contentType         ContentType // Type-safe content rendering
	contentText         string      // For help/status messages
	commandMatches      []SlashCommand
	selectedCommand     int
	inTypeaheadMode     bool
	inPostTabCompletion bool           // Prevent re-entering typeahead after Tab completion
	groups              []config.Group // Groups from config
	viewGroupName       string         // Current group being viewed
	viewGroupKind       string         // "all", "blocklists", "allowed", "blocked"
	scrollOffset        int            // Current scroll position in viewport
	commandHistory      []CommandEvent // Append-only history log
	historyScroll       int            // Scroll position in history view
	firstRun            bool           // Track first-time user for welcome banner
}

// New creates and initializes a new TUI model, preparing the text input.
// Note: TUI auto-launch always loads "dnsApp.config" from current directory.
// The --config flag is not respected in auto-launch mode. Use a subcommand
// (pab map, pab deploy) with --config if you need to specify a custom config path.
func New(cfg *config.Config) *Model {
	ti := textinput.New()
	ti.Placeholder = "Search or type /help for commands"
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
		clients:             clients,
		unifiedInput:        ti,
		ready:               true,
		contentType:         ContentTypeEmpty,
		contentText:         "",
		commandMatches:      []SlashCommand{},
		selectedCommand:     0,
		inTypeaheadMode:     false,
		inPostTabCompletion: false,
		viewGroupName:       "",
		viewGroupKind:       "all",
		commandHistory:      []CommandEvent{},
		historyScroll:       0,
		firstRun:            true,
	}
	// Add groups from config
	if cfg != nil {
		m.groups = cfg.Groups
	}
	m.filterClients()
	return m
}

// Init initializes the Bubble Tea application and triggers configuration loading.
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// appendHistory adds a command event to the history log
func (m *Model) appendHistory(command string, output string) {
	event := CommandEvent{
		Timestamp: time.Now(),
		Command:   command,
		Output:    output,
		Lines:     strings.Split(output, "\n"),
	}
	m.commandHistory = append(m.commandHistory, event)
	m.historyScroll = 0 // Reset scroll to top of new history
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

		// Handle navigation in typeahead mode
		if m.inTypeaheadMode && len(m.commandMatches) > 0 {
			switch msg.Type {
			case tea.KeyUp:
				if m.selectedCommand > 0 {
					m.selectedCommand--
				}
				return m, nil
			case tea.KeyDown:
				if m.selectedCommand < len(m.commandMatches)-1 {
					m.selectedCommand++
				}
				return m, nil
			case tea.KeyTab:
				// Tab completion: complete to selected command name + space
				selectedCmd := m.commandMatches[m.selectedCommand]
				m.unifiedInput.SetValue(selectedCmd.Name + " ")
				m.unifiedInput.CursorEnd()
				// Exit typeahead after completion so user can type subcommand arguments without re-filtering
				m.inTypeaheadMode = false
				m.inPostTabCompletion = true // Prevent re-entering typeahead mode
				m.commandMatches = []SlashCommand{}
				m.selectedCommand = 0
				m.contentType = ContentTypeEmpty
				return m, nil
			case tea.KeyEnter:
				// Execute the selected command
				selectedCmd := m.commandMatches[m.selectedCommand]
				// Parse the command name to extract the command and any args
				parts := strings.Fields(selectedCmd.Name)
				cmd := strings.TrimPrefix(parts[0], "/")
				args := []string{}
				if len(parts) > 1 {
					args = parts[1:]
				}
				return m.executeCommand(cmd, args)
			}
		}

		// Handle scrolling in normal (non-typeahead) mode
		if !m.inTypeaheadMode {
			switch msg.Type {
			case tea.KeyUp:
				// Scroll up in history if we have history, otherwise in regular content
				if len(m.commandHistory) > 0 {
					if m.historyScroll > 0 {
						m.historyScroll--
					}
				} else if m.contentType != ContentTypeEmpty && m.scrollOffset > 0 {
					m.scrollOffset--
				}
				return m, nil
			case tea.KeyDown:
				// Scroll down in history if we have history, otherwise in regular content
				if len(m.commandHistory) > 0 {
					// Conservative limit to prevent overflow
					if m.historyScroll < 100 {
						m.historyScroll++
					}
				} else if m.contentType != ContentTypeEmpty && m.scrollOffset < 100 {
					m.scrollOffset++
				}
				return m, nil
			}
		}

		// Forward key presses to the search text input component
		var tiCmd tea.Cmd
		m.unifiedInput, tiCmd = m.unifiedInput.Update(msg)

		rawInput := m.unifiedInput.Value()
		input := strings.TrimSpace(rawInput)

		// Parse the input using unified parser
		parsed := ParseUnifiedInput(input)

		switch parsed.Type {
		case InputTypeCommand:
			// Clear post-Tab-completion flag if user starts a new command (just "/")
			if input == "/" {
				m.inPostTabCompletion = false
			}

			m.commandMatches = filterCommands(rawInput)
			// Only enter typeahead mode if:
			// 1. We have matching commands
			// 2. We're not in post-Tab-completion mode (prevents re-entry when typing subcommand args)
			if len(m.commandMatches) > 0 && !m.inPostTabCompletion {
				m.inTypeaheadMode = true
				m.contentType = ContentTypeCommandList
			} else {
				// No matches or in post-Tab-completion mode - user is typing subcommand args, stay out of typeahead
				m.inTypeaheadMode = false
				m.contentType = ContentTypeEmpty
			}
			m.selectedCommand = 0

			// Handle Enter to execute full command
			if msg.Type == tea.KeyEnter {
				if len(m.commandMatches) > 0 {
					selectedCmd := m.commandMatches[m.selectedCommand]
					return m.executeCommand(selectedCmd.Name, []string{})
				}
				// Try to execute what was typed as-is
				trimmed := strings.TrimPrefix(input, "/")
				return m.executeCommand(trimmed, parsed.Args)
			}

		case InputTypeSearch:
			// Not in slash command mode
			m.inTypeaheadMode = false
			m.inPostTabCompletion = false
			m.commandMatches = []SlashCommand{}
			m.selectedCommand = 0

			// Update filter results dynamically on every keystroke
			m.contentType = ContentTypeTable
			m.scrollOffset = 0
			m.filterClients()

		case InputTypeEmpty:
			// No-op for empty input
			m.inTypeaheadMode = false
			m.inPostTabCompletion = false
			m.commandMatches = []SlashCommand{}
			m.selectedCommand = 0
			m.contentType = ContentTypeEmpty
		}

		cmd = tea.Batch(cmd, tiCmd)
	}

	return m, cmd
}

// executeCommand executes a slash command and updates the model state
// cmd should be the command name (without leading slash), and args are the parsed arguments
func (m *Model) executeCommand(cmd string, args []string) (tea.Model, tea.Cmd) {
	// cmd and args are already parsed by ParseUnifiedInput
	cmdLower := strings.ToLower(cmd)

	// Check if it matches any command or alias
	for _, c := range commands {
		if strings.EqualFold(c.Name, "/"+cmd) {
			cmdLower = strings.ToLower(cmd)
			break
		}
		for _, alias := range c.Aliases {
			if strings.EqualFold(alias, "/"+cmd) {
				cmdLower = strings.TrimPrefix(strings.ToLower(c.Name), "/")
				break
			}
		}
	}

	// Reconstruct the full command string for history
	input := "/" + cmd
	if len(args) > 0 {
		input += " " + strings.Join(args, " ")
	}

	switch cmdLower {
	case "exit", "quit":
		return m, tea.Quit
	case "help", "?":
		// Show help text in content area
		helpOutput := helpText()
		m.appendHistory(input, helpOutput)
		m.contentType = ContentTypeHelp
		m.contentText = helpOutput
		m.scrollOffset = 0
		m.unifiedInput.SetValue("")
		m.inTypeaheadMode = false
		m.inPostTabCompletion = false
		m.firstRun = false // Dismiss first-run banner when a command is executed
		return m, nil
	case "clear", "c":
		m.commandHistory = []CommandEvent{}
		m.historyScroll = 0
		m.unifiedInput.SetValue("")
		m.contentType = ContentTypeEmpty
		m.contentText = ""
		m.scrollOffset = 0
		m.commandMatches = []SlashCommand{}
		m.selectedCommand = 0
		m.inTypeaheadMode = false
		m.inPostTabCompletion = false
		m.filterClients()
		m.firstRun = false // Dismiss first-run banner when a command is executed
		return m, nil
	case "view", "v":
		viewOutput := m.handleViewWithOutput(args)
		m.appendHistory(input, viewOutput)
		m.scrollOffset = 0
		m.unifiedInput.SetValue("")
		m.inTypeaheadMode = false
		m.inPostTabCompletion = false
		m.firstRun = false // Dismiss first-run banner when a command is executed
		return m, nil
	default:
		// Unknown command, clear it
		m.unifiedInput.SetValue("")
		m.inTypeaheadMode = false
		return m, nil
	}
}

// handleView processes /view subcommands
func (m *Model) handleView(args []string) {
	if len(args) == 0 {
		m.contentType = ContentTypeHelp
		return
	}

	switch args[0] {
	case "networkGroupMap", "map":
		m.contentType = ContentTypeViewNetworkGroupMap
	case "groups":
		m.contentType = ContentTypeViewGroups
	case "group":
		if len(args) < 2 {
			// /view group needs a group name
			m.contentType = ContentTypeHelp
			return
		}
		groupName := args[1]
		kind := "all"
		if len(args) > 2 {
			kind = args[2] // "blocklists", "allowed", "blocked"
		}
		m.viewGroupName = groupName
		m.viewGroupKind = kind
		m.contentType = ContentTypeViewGroup
	default:
		m.contentType = ContentTypeHelp
	}
}

// handleViewWithOutput processes /view subcommands and returns the output
func (m *Model) handleViewWithOutput(args []string) string {
	if len(args) == 0 {
		m.contentType = ContentTypeHelp
		return ""
	}

	switch args[0] {
	case "networkGroupMap", "map":
		m.contentType = ContentTypeViewNetworkGroupMap
		return m.renderNetworkGroupMap()
	case "groups":
		m.contentType = ContentTypeViewGroups
		return m.renderGroupsList()
	case "group":
		if len(args) < 2 {
			// /view group needs a group name
			m.contentType = ContentTypeHelp
			return "group name required\nUsage: /view group <name>"
		}
		groupName := args[1]
		kind := "all"
		if len(args) > 2 {
			kind = args[2] // "blocklists", "allowed", "blocked"
		}
		m.viewGroupName = groupName
		m.viewGroupKind = kind
		m.contentType = ContentTypeViewGroup
		return m.renderGroupDetail()
	default:
		m.contentType = ContentTypeHelp
		return "Unknown view subcommand\n" + viewSubcommandHelp()
	}
}

// viewSubcommandHelp returns help text for view subcommands
func viewSubcommandHelp() string {
	return `Available View Subcommands:
  /view networkGroupMap    Show all IP to Group mappings
  /view groups             List all configured groups
  /view group <name>       Show group details (all domains)`
}

// findGroup helper finds a group by case-insensitive name match
func (m *Model) findGroup(name string) *config.Group {
	for i := range m.groups {
		if strings.ToLower(m.groups[i].Name) == strings.ToLower(name) {
			return &m.groups[i]
		}
	}
	return nil
}

// renderNetworkGroupMap renders all IP to Group mappings
func (m *Model) renderNetworkGroupMap() string {
	if len(m.clients) == 0 {
		return "No client mappings found."
	}
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf("%-30s %s", "Client IP", "Group")) + "\n")
	for i, c := range m.clients {
		row := fmt.Sprintf("%-30s %s", c.IP, c.Group)
		if i%2 == 0 {
			b.WriteString(rowStyle.Render(row) + "\n")
		} else {
			b.WriteString(altRowStyle.Render(row) + "\n")
		}
	}
	return b.String()
}

// renderGroupsList renders a list of all groups
func (m *Model) deviceCount(groupName string) int {
	count := 0
	for _, c := range m.clients {
		if c.Group == groupName {
			count++
		}
	}
	return count
}

func (m *Model) renderGroupsList() string {
	if len(m.groups) == 0 {
		return "No groups configured."
	}
	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf("%-30s %s", "Group Name", "Devices")) + "\n")
	for i, g := range m.groups {
		count := m.deviceCount(g.Name)
		row := fmt.Sprintf("%-30s %d", g.Name, count)
		if i%2 == 0 {
			b.WriteString(rowStyle.Render(row) + "\n")
		} else {
			b.WriteString(altRowStyle.Render(row) + "\n")
		}
	}
	return b.String()
}

// renderGroupDetail renders details for a specific group
func (m *Model) renderGroupDetail() string {
	g := m.findGroup(m.viewGroupName)
	if g == nil {
		return fmt.Sprintf("Group '%s' not found.", m.viewGroupName)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Group: %s\n\n", g.Name))

	if m.viewGroupKind == "all" || m.viewGroupKind == "blocklists" {
		b.WriteString("Blocked Domains:\n")
		for _, domain := range sanitize(g.Blocked) {
			b.WriteString(fmt.Sprintf("  • %s\n", domain))
		}
		b.WriteString("\n")
	}

	if m.viewGroupKind == "all" || m.viewGroupKind == "allowed" {
		b.WriteString("Allowed Domains:\n")
		for _, domain := range sanitize(g.Allowed) {
			b.WriteString(fmt.Sprintf("  • %s\n", domain))
		}
	}

	return b.String()
}

// sanitize removes control/ANSI characters from strings
func sanitize(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		// Remove control characters and ANSI escapes
		cleaned := strings.Map(func(r rune) rune {
			if r < 32 || r == 127 {
				return -1 // remove control chars
			}
			return r
		}, item)
		result = append(result, cleaned)
	}
	return result
}

// filterCommands filters the slash commands based on the input text
// Handles both top-level commands and view subcommands
func filterCommands(input string) []SlashCommand {
	trimmed := strings.TrimSpace(input)

	if trimmed == "/" {
		// Show all commands
		return commands
	}

	// Check if we're looking for view subcommands
	// Handle /view with or without space, and /v with space
	// Examples: "/view", "/view ", "/view g", "/v ", "/v g"
	// This allows prefix filtering while preventing re-entering typeahead after Tab completion
	if input == "/view" || strings.HasPrefix(input, "/view ") || strings.HasPrefix(input, "/v ") {
		return filterViewSubcommands(input)
	}

	// Filter top-level commands
	var matches []SlashCommand
	trimmedLower := strings.ToLower(trimmed)

	for _, cmd := range commands {
		// Check if command name matches prefix
		if strings.HasPrefix(strings.ToLower(cmd.Name), trimmedLower) {
			matches = append(matches, cmd)
			continue
		}

		// Check if any alias matches prefix
		for _, alias := range cmd.Aliases {
			if strings.HasPrefix(strings.ToLower(alias), trimmedLower) {
				matches = append(matches, cmd)
				break
			}
		}
	}

	return matches
}

// filterViewSubcommands filters view subcommands based on input
// Input should start with "/view " or "/v " followed by optional prefix text
// For example: "/view g" matches "/view group" and "/view groups"
func filterViewSubcommands(input string) []SlashCommand {
	// Normalize input: handle /view with or without space, and /v with space
	var prefix string
	if strings.HasPrefix(input, "/view ") {
		prefix = strings.TrimPrefix(input, "/view ")
	} else if strings.HasPrefix(input, "/view") && input != "/view" {
		// Handle "/view<something>" like "/viewg" or "/view" without space after
		prefix = strings.TrimPrefix(input, "/view")
	} else if input == "/view" {
		// Bare "/view" shows all subcommands
		prefix = ""
	} else if strings.HasPrefix(input, "/v ") {
		prefix = strings.TrimPrefix(input, "/v ")
	} else if input == "/v" {
		// Bare "/v" shows all subcommands
		prefix = ""
	} else {
		return []SlashCommand{}
	}

	// Filter subcommands by prefix
	var matches []SlashCommand
	prefixLower := strings.ToLower(prefix)
	for _, sub := range viewSubcommands {
		if strings.HasPrefix(strings.ToLower(sub.Name), prefixLower) {
			fullCmd := "/view " + sub.Name
			matches = append(matches, SlashCommand{
				Name:        fullCmd,
				Aliases:     []string{},
				Description: sub.Description,
			})
		}
	}
	return matches
}

// renderCommandList renders the filtered command list with descriptions
func renderCommandList(commands []SlashCommand, selected int) string {
	if len(commands) == 0 {
		return "No matching commands found."
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render("Available Slash Commands:") + "\n\n")

	for i, cmd := range commands {
		// Build aliases string
		var aliasStr string
		if len(cmd.Aliases) > 0 {
			aliasStr = strings.Join(cmd.Aliases, ", ")
		} else {
			aliasStr = ""
		}

		// Format the command line
		var cmdLine string
		if aliasStr != "" {
			cmdLine = fmt.Sprintf("  %s, %-10s - %s", cmd.Name, aliasStr, cmd.Description)
		} else {
			cmdLine = fmt.Sprintf("  %s - %s", cmd.Name, cmd.Description)
		}

		// Highlight selected command
		if i == selected {
			b.WriteString(rowStyle.Foreground(accentColor).Bold(true).Render("> "+cmdLine) + "\n")
		} else {
			b.WriteString(rowStyle.Render(cmdLine) + "\n")
		}
	}

	b.WriteString("\nStart typing to narrow down...")
	return b.String()
}

// filterClients updates the filtered client list based on the search input query.
func (m *Model) filterClients() {
	query := strings.ToLower(m.unifiedInput.Value())
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

// helpText returns the help message for available commands
func helpText() string {
	return `Available Commands:
  /help or /?              Show this help message
  /clear                   Clear search and show all clients
  /exit or /quit           Exit the application

View Commands:
  /view networkGroupMap    Show all IP to Group mappings
  /view groups             List all configured groups
  /view group <name>       Show group details (all domains)
  /view group <name> blocklists  Show blocked domains
  /view group <name> allowed     Show allowed domains

Tips:
  • Start typing to search by IP or Group
  • Results update as you type
  • Use arrow keys to navigate if needed
  • Press Enter to confirm search queries`
}

// renderTable renders the client table content
func (m *Model) renderTable() string {
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

	return b.String()
}

// renderHistory renders the full command history from oldest to newest
func (m *Model) renderHistory() string {
	if len(m.commandHistory) == 0 {
		return "No command history yet. Type a command to get started.\n(Use / to see available commands)"
	}

	var output strings.Builder
	contentHeight := m.height - 10 // Account for: title(2) + search(3) + footer(1) + border+padding(4)
	if contentHeight < 3 {
		contentHeight = 3
	}

	// Calculate total lines needed
	totalLines := 0
	for _, event := range m.commandHistory {
		totalLines += len(event.Lines) + 3 // +3 for timestamp line, separator, spacing
	}

	// Build the full history text first
	var fullHistoryLines []string
	for _, event := range m.commandHistory {
		timestamp := event.Timestamp.Format("15:04:05")
		fullHistoryLines = append(fullHistoryLines, fmt.Sprintf("%s | %s", timestamp, event.Command))

		// Add output lines
		for _, line := range event.Lines {
			fullHistoryLines = append(fullHistoryLines, line)
		}

		// Add separator
		fullHistoryLines = append(fullHistoryLines, "---")
	}

	// Apply scroll offset to show scrolled view
	startIdx := m.historyScroll
	endIdx := startIdx + contentHeight
	if endIdx > len(fullHistoryLines) {
		endIdx = len(fullHistoryLines)
	}
	if startIdx >= len(fullHistoryLines) {
		startIdx = len(fullHistoryLines) - contentHeight
		if startIdx < 0 {
			startIdx = 0
		}
	}

	var lines []string
	if startIdx > 0 || endIdx < len(fullHistoryLines) {
		lines = fullHistoryLines[startIdx:endIdx]
	} else {
		lines = fullHistoryLines
	}

	// Pad with blank lines to fill available height
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	output.WriteString(strings.Join(lines, "\n"))
	return output.String()
}

// renderContent renders content based on the current content type
func (m *Model) renderContent() string {
	contentHeight := m.height - 10 // Account for: title(2) + search(3) + footer(1) + border+padding(4)
	if contentHeight < 3 {
		contentHeight = 3
	}

	var content string
	switch m.contentType {
	case ContentTypeHelp:
		content = m.contentText
	case ContentTypeCommandList:
		content = renderCommandList(m.commandMatches, m.selectedCommand)
	case ContentTypeTable:
		content = m.renderTable()
	case ContentTypeError:
		content = m.contentText
	case ContentTypeViewNetworkGroupMap:
		content = m.renderNetworkGroupMap()
	case ContentTypeViewGroups:
		content = m.renderGroupsList()
	case ContentTypeViewGroup:
		content = m.renderGroupDetail()
	default: // ContentTypeEmpty
		content = "Start typing to search by IP or Group, or type /help for commands"
	}

	// Split content by lines and apply scroll offset
	lines := strings.Split(content, "\n")

	// Apply scroll offset to show scrolled view
	startIdx := m.scrollOffset
	endIdx := startIdx + contentHeight
	if endIdx > len(lines) {
		endIdx = len(lines)
	}
	if startIdx >= len(lines) {
		startIdx = len(lines) - contentHeight
		if startIdx < 0 {
			startIdx = 0
		}
	}
	if startIdx > 0 || endIdx < len(lines) {
		lines = lines[startIdx:endIdx]
	}

	// Pad with blank lines to fill available height
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// View renders the TUI into a beautifully formatted string using Lipgloss.
func (m *Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error loading config: %v", m.err))
	}

	if !m.ready {
		return "Loading Pharos Config..."
	}

	// Fixed title at top
	title := titleStyle.Render("Pharos Advanced Blocking")

	// First-run welcome banner
	if m.firstRun && m.unifiedInput.Value() == "" {
		welcome := lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Bold(true).
			Render("Welcome to Pharos Advanced Blocking!")

		subtitle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render("Quick start: Type to search by IP/Group, or /help for commands")

		welcomeBox := lipgloss.JoinVertical(
			lipgloss.Top,
			welcome,
			subtitle,
			"",
			"Tip: Press / to see all available commands",
		)

		searchBox := lipgloss.NewStyle().
			Padding(0, 1).
			MarginTop(1).
			Render(m.unifiedInput.View())

		footer := footerStyle.Render("ctrl+c / esc: exit | /help: commands | /clear: reset")

		layout := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			welcomeBox,
			searchBox,
			footer,
		)

		if m.width > 0 && m.height > 0 {
			return baseStyle.Width(m.width - 4).Height(m.height - 2).Render(layout)
		}
		return baseStyle.Render(layout)
	}

	// Dismiss banner on first keystroke
	if m.firstRun && m.unifiedInput.Value() != "" {
		m.firstRun = false
	}

	// Dynamic content area in the middle
	// If we have history, show it; otherwise show regular content
	var renderedContent string
	if len(m.commandHistory) > 0 {
		renderedContent = m.renderHistory()
	} else {
		renderedContent = m.renderContent()
	}

	contentBox := lipgloss.NewStyle().
		Padding(0, 1).
		Render(renderedContent)

	// Search box (fixed above footer)
	searchBox := lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1).
		Render(m.unifiedInput.View())

	// Footer help status line (fixed at bottom)
	footerText := "ctrl+c / esc: exit | /help: commands | /clear: reset"
	// Add scroll hint if we have history
	if len(m.commandHistory) > 0 && !m.inTypeaheadMode {
		footerText += " | ↑↓: scroll through history"
		footerText += fmt.Sprintf(" | (%d commands in history)", len(m.commandHistory))
	} else if m.contentType != ContentTypeEmpty && !m.inTypeaheadMode {
		footerText += " | ↑↓: scroll"
	}
	footer := footerStyle.Render(footerText)

	// Assemble the layout vertically: title, content, search, footer
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		contentBox,
		searchBox,
		footer,
	)

	// Apply responsive padding and borders to wrap the layout
	if m.width > 0 && m.height > 0 {
		return baseStyle.Width(m.width - 4).Height(m.height - 2).Render(layout)
	}

	return baseStyle.Render(layout)
}
