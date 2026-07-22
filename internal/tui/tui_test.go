package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iamrichardd/pharos-advanced-blocking/internal/config"
)

func TestTUI_Lifecycle(t *testing.T) {
	// 1. Create a decoupled configuration structure
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
			"192.168.1.150": "IoT",
			"10.0.0.5":      "Laptops",
		},
	}

	// 2. Initialize the Bubble Tea Model
	m := New(cfg)

	// Simulate Init
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Expected Init() to return a command for blinking text input")
	}

	// Verify initial content type is "empty"
	if m.contentType != "empty" {
		t.Errorf("Expected initial contentType to be 'empty', got %q", m.contentType)
	}

	// 3. Ensure clients list is properly populated and sorted
	if len(m.clients) != 3 {
		t.Errorf("Expected 3 clients, got %d", len(m.clients))
	}

	// Ensure clients are alphabetically/lexicographically sorted by IP for stable view
	if m.clients[0].IP != "10.0.0.5" {
		t.Errorf("Expected first IP to be 10.0.0.5, got %s", m.clients[0].IP)
	}

	// 4. Simulate User Search Input
	// We'll simulate pressing the 'I', 'o', 'T' keys
	keyMsgs := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'I'}},
		{Type: tea.KeyRunes, Runes: []rune{'o'}},
		{Type: tea.KeyRunes, Runes: []rune{'T'}},
	}

	// Because Update returns a tea.Model interface, we keep track of the pointer
	var updatedModel tea.Model = m
	for _, key := range keyMsgs {
		updatedModel, _ = updatedModel.Update(key)
	}
	m2, ok := updatedModel.(*Model)
	if !ok {
		t.Fatal("Expected model type *Model")
	}

	// Verify contentType changes to "table" when search input is present
	if m2.contentType != "table" {
		t.Errorf("Expected contentType to be 'table' after search input, got %q", m2.contentType)
	}

	// 5. Verify dynamic filtering functionality
	if len(m2.filtered) != 1 {
		t.Fatalf("Expected 1 filtered result, got %d", len(m2.filtered))
	}
	if m2.filtered[0].Group != "IoT" {
		t.Errorf("Expected filtered group to be 'IoT', got '%s'", m2.filtered[0].Group)
	}
	if m2.filtered[0].IP != "192.168.1.150" {
		t.Errorf("Expected IP to be '192.168.1.150'")
	}

	// 6. Verify View rendering
	m2.ready = true
	m2.width = 100 // Simulate terminal window sizing
	m2.height = 30
	viewOutput := m2.View()

	if !strings.Contains(viewOutput, "IoT") {
		t.Errorf("Expected view to contain 'IoT', got:\n%s", viewOutput)
	}
	if !strings.Contains(viewOutput, "192.168.1.150") {
		t.Errorf("Expected view to contain '192.168.1.150'")
	}
	if !strings.Contains(viewOutput, "Pharos Advanced Blocking") {
		t.Errorf("Expected view to contain title")
	}
	if !strings.Contains(viewOutput, "ctrl+c / esc: exit | /help: commands") {
		t.Errorf("Expected view to contain footer help status line")
	}
}

func TestTUI_HelpCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Simulate typing "/help" and pressing Enter
	helpKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	_, _ = m.Update(helpKeyMsg)

	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify contentType changes to "help"
	if m.contentType != "help" {
		t.Errorf("Expected contentType to be 'help' after /help command, got %q", m.contentType)
	}

	// Verify content text contains help information
	if !strings.Contains(m.contentText, "Available Commands") {
		t.Errorf("Expected contentText to contain 'Available Commands'")
	}

	// Verify the view renders help text
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "Available Commands") {
		t.Errorf("Expected help text in view output")
	}
	if !strings.Contains(viewOutput, "/clear") {
		t.Errorf("Expected /clear command in help text")
	}
}

func TestTUI_ClearCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
			"192.168.1.150": "IoT",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// First, do a search
	for _, r := range "192" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if m.contentType != "table" {
		t.Errorf("Expected contentType to be 'table' after search")
	}
	if len(m.filtered) != 2 {
		t.Errorf("Expected 2 filtered results, got %d", len(m.filtered))
	}

	// Clear the search input and type /clear command
	m.searchInput.SetValue("")
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "clear" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify contentType changes back to "empty"
	if m.contentType != "empty" {
		t.Errorf("Expected contentType to be 'empty' after /clear command, got %q", m.contentType)
	}

	// Verify search box is cleared
	if m.searchInput.Value() != "" {
		t.Errorf("Expected search input to be cleared, got %q", m.searchInput.Value())
	}

	// Verify the view renders empty state hint
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "Start typing to search") {
		t.Errorf("Expected empty state hint in view output")
	}
}

func TestTUI_SlashCommandTypeahead(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to enter typeahead mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Verify we're in typeahead mode
	if !m.inTypeaheadMode {
		t.Errorf("Expected to be in typeahead mode after typing /")
	}

	// Verify contentType is "command_list"
	if m.contentType != "command_list" {
		t.Errorf("Expected contentType to be 'command_list', got %q", m.contentType)
	}

	// Verify all commands are shown for just "/"
	if len(m.commandMatches) != 4 {
		t.Errorf("Expected 4 command matches for '/', got %d", len(m.commandMatches))
	}

	// Verify view renders command list
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "Available Slash Commands") {
		t.Errorf("Expected command list header in view output")
	}
	if !strings.Contains(viewOutput, "/help") {
		t.Errorf("Expected /help command in view output")
	}
}

func TestTUI_CommandFiltering_Help(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/h" to filter help commands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	// Verify filtered results
	if len(m.commandMatches) != 1 {
		t.Errorf("Expected 1 command match for '/h', got %d", len(m.commandMatches))
	}
	if len(m.commandMatches) > 0 && !strings.HasPrefix(m.commandMatches[0].Name, "/help") {
		t.Errorf("Expected /help to match /h")
	}

	// Verify view renders only matching command
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "/help") {
		t.Errorf("Expected /help in filtered results")
	}
}

func TestTUI_CommandFiltering_Exit(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/e" to filter exit commands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	// Verify filtered results
	if len(m.commandMatches) != 1 {
		t.Errorf("Expected 1 command match for '/e', got %d", len(m.commandMatches))
	}
	if len(m.commandMatches) > 0 && !strings.HasPrefix(m.commandMatches[0].Name, "/exit") {
		t.Errorf("Expected /exit to match /e")
	}
}

func TestTUI_CommandFiltering_Alias(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/q" which is an alias for /exit
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Verify /exit command is matched by alias
	if len(m.commandMatches) != 1 {
		t.Errorf("Expected 1 command match for '/q' alias, got %d", len(m.commandMatches))
	}
	if len(m.commandMatches) > 0 && m.commandMatches[0].Name != "/exit" {
		t.Errorf("Expected /exit to match /q alias")
	}
}

func TestTUI_ArrowNavigation_Down(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to show all commands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Verify initial selection is 0
	if m.selectedCommand != 0 {
		t.Errorf("Expected initial selection to be 0, got %d", m.selectedCommand)
	}

	// Press down arrow
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Verify selection moved to 1
	if m.selectedCommand != 1 {
		t.Errorf("Expected selection to be 1 after down arrow, got %d", m.selectedCommand)
	}

	// Press down arrow again
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Verify selection moved to 2
	if m.selectedCommand != 2 {
		t.Errorf("Expected selection to be 2 after second down arrow, got %d", m.selectedCommand)
	}
}

func TestTUI_ArrowNavigation_Up(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to show all commands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Move down twice
	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if m.selectedCommand != 2 {
		t.Errorf("Expected selection to be 2, got %d", m.selectedCommand)
	}

	// Press up arrow
	m.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Verify selection moved to 1
	if m.selectedCommand != 1 {
		t.Errorf("Expected selection to be 1 after up arrow, got %d", m.selectedCommand)
	}
}

func TestTUI_ExitTypeaheadWithRegularText(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
			"192.168.1.150": "IoT",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Start in typeahead mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	if !m.inTypeaheadMode {
		t.Errorf("Expected to be in typeahead mode")
	}

	// Clear and type regular text (backspace then regular char)
	m.searchInput.SetValue("")
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})

	// Verify we're NOT in typeahead mode anymore
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode when typing regular text")
	}

	// Verify contentType is "table" for regular search
	if m.contentType != "table" {
		t.Errorf("Expected contentType to be 'table' for regular search, got %q", m.contentType)
	}
}

func TestTUI_CommandListHighlight(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" and navigate
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Verify selectedCommand is set
	if m.selectedCommand != 1 {
		t.Errorf("Expected selectedCommand to be 1")
	}

	// Render the view and verify it contains the selected command indicator
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "/exit") {
		t.Errorf("Expected /exit in view output")
	}
}

func TestTUI_ClearBackToNormalSearch(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Start typeahead
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	if !m.inTypeaheadMode {
		t.Errorf("Expected typeahead mode")
	}

	// Type "/help"
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify in typeahead mode
	if !m.inTypeaheadMode {
		t.Errorf("Expected to be in typeahead mode")
	}

	// Clear the input
	m.searchInput.SetValue("")

	// Verify empty state
	if m.searchInput.Value() != "" {
		t.Errorf("Expected search input to be cleared")
	}
}

func TestTUI_TabCompletion(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to enter typeahead mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Type "hel" to filter to /help command
	for _, r := range "hel" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify we have a single match for /help
	if len(m.commandMatches) != 1 {
		t.Errorf("Expected 1 command match for '/hel', got %d", len(m.commandMatches))
	}
	if m.commandMatches[0].Name != "/help" {
		t.Errorf("Expected /help command, got %q", m.commandMatches[0].Name)
	}

	// Press Tab to complete
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify input is now "/help " (with trailing space)
	if m.searchInput.Value() != "/help " {
		t.Errorf("Expected '/help ' after Tab completion, got %q", m.searchInput.Value())
	}

	// Verify exited typeahead mode (user can now type subcommand arguments)
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode after Tab completion")
	}

	// Verify contentType is empty (ready to accept command arguments)
	if m.contentType != "empty" {
		t.Errorf("Expected contentType 'empty' after Tab completion, got %q", m.contentType)
	}
}

func TestTUI_TabCompletion_MultiMatch(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to enter typeahead mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Verify all 4 commands are shown
	if len(m.commandMatches) != 4 {
		t.Errorf("Expected 4 command matches for '/', got %d", len(m.commandMatches))
	}

	// Press Tab with multiple matches (completes to selected command)
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify Tab completes to the selected command (/help is default selected)
	if m.searchInput.Value() != "/help " {
		t.Errorf("Expected '/help ' (default selected), got %q", m.searchInput.Value())
	}

	// Verify exited typeahead after completion
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode after Tab completion")
	}
}

func TestTUI_TabCompletion_OutsideTypeahead(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "abc" (normal search, no slash command)
	for _, r := range "abc" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify not in typeahead mode
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to be in typeahead mode for normal search")
	}

	// Press Tab (should be a no-op, not crash)
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify Tab had no effect outside typeahead
	if m.searchInput.Value() != "abc" {
		t.Errorf("Expected Tab to be no-op outside typeahead, got input %q", m.searchInput.Value())
	}
}

func TestTUI_TabCompletion_AfterArrowNav(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/" to enter typeahead mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	// Verify /help is selected (default selectedCommand = 0)
	if m.commandMatches[0].Name != "/help" {
		t.Errorf("Expected first match to be /help")
	}

	// Navigate Down to /exit (assuming order: /help, /exit, /clear, /view)
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Verify /exit is now selected
	if m.selectedCommand != 1 {
		t.Errorf("Expected selectedCommand to be 1 after Down, got %d", m.selectedCommand)
	}

	// Press Tab to complete the currently selected command (/exit)
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify Tab completed to /exit, not /help
	if m.searchInput.Value() != "/exit " {
		t.Errorf("Expected '/exit ' after Tab with Down navigation, got %q", m.searchInput.Value())
	}
}

func TestTUI_TabCompletion_SubcommandEntry(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/v" to filter to /view
	for _, r := range "/v" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify /view is in matches
	if len(m.commandMatches) < 1 {
		t.Errorf("Expected at least 1 match for '/v'")
	}

	// Press Tab to complete to "/view "
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.searchInput.Value() != "/view " {
		t.Errorf("Expected '/view ' after Tab, got %q", m.searchInput.Value())
	}

	// Now type subcommand "groups" without re-filtering
	for _, r := range "groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify user can type subcommand after Tab completion
	if m.searchInput.Value() != "/view groups" {
		t.Errorf("Expected '/view groups' after typing subcommand, got %q", m.searchInput.Value())
	}

	// Verify not back in typeahead mode (so Enter will use the raw typed string)
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to be in typeahead mode after typing subcommand")
	}
}

func TestTUI_ViewSubcommandFiltering_Prefix(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/view g" to filter subcommands starting with "g"
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view g" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify we're in typeahead mode
	if !m.inTypeaheadMode {
		t.Errorf("Expected to be in typeahead mode after typing /view g")
	}

	// Verify we have 2 matches: /view group and /view groups
	if len(m.commandMatches) != 2 {
		t.Errorf("Expected 2 command matches for '/view g', got %d", len(m.commandMatches))
	}

	// Verify matches are the group-related subcommands
	for _, cmd := range m.commandMatches {
		if !strings.HasPrefix(cmd.Name, "/view group") {
			t.Errorf("Expected match to start with '/view group', got %q", cmd.Name)
		}
	}

	// Verify view renders the filtered subcommands
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "/view group") {
		t.Errorf("Expected /view group in filtered results")
	}
}

func TestTUI_ViewSubcommandFiltering_NoMatch(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/view xyz" to filter with no matching subcommands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view xyz" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify we're NOT in typeahead mode (no matches)
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to be in typeahead mode when no subcommands match")
	}

	// Verify we have 0 matches
	if len(m.commandMatches) != 0 {
		t.Errorf("Expected 0 command matches for '/view xyz', got %d", len(m.commandMatches))
	}

	// Verify contentType is "empty" (no matches shown)
	if m.contentType != "empty" {
		t.Errorf("Expected contentType 'empty' when no matches found, got %q", m.contentType)
	}
}

func TestTUI_ViewSubcommandAfterTabNoReentry(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
		Groups: []config.Group{
			{
				Name:    "group1",
				Blocked: []string{"blocked.com"},
				Allowed: []string{"allowed.com"},
			},
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Type "/view " to show all subcommands
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view " {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify in typeahead mode with 3 matches
	if !m.inTypeaheadMode {
		t.Errorf("Expected to be in typeahead mode after /view ")
	}
	if len(m.commandMatches) != 3 {
		t.Errorf("Expected 3 subcommand matches for '/view ', got %d", len(m.commandMatches))
	}

	// Press Tab to complete to first match ("/view groups")
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify we exited typeahead mode
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode after Tab completion")
	}

	// Verify search input is now a full command with trailing space
	currentInput := m.searchInput.Value()
	if !strings.HasPrefix(currentInput, "/view ") {
		t.Errorf("Expected input to start with '/view ', got %q", currentInput)
	}

	// Now type a space (simulating user continuing to type)
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

	// Verify we're still NOT in typeahead mode (no re-entry)
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to re-enter typeahead mode after typing space following Tab completion")
	}

	// Verify contentType remains "empty" (not command_list)
	if m.contentType != "empty" {
		t.Errorf("Expected contentType 'empty' after typing space, got %q", m.contentType)
	}
}
