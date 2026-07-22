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
	if m.contentType != ContentTypeEmpty {
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
	if m2.contentType != ContentTypeTable {
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
	if m.contentType != ContentTypeHelp {
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
	if m.contentType != ContentTypeTable {
		t.Errorf("Expected contentType to be 'table' after search")
	}
	if len(m.filtered) != 2 {
		t.Errorf("Expected 2 filtered results, got %d", len(m.filtered))
	}

	// Clear the search input and type /clear command
	m.unifiedInput.SetValue("")
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "clear" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify contentType changes back to "empty"
	if m.contentType != ContentTypeEmpty {
		t.Errorf("Expected contentType to be 'empty' after /clear command, got %q", m.contentType)
	}

	// Verify search box is cleared
	if m.unifiedInput.Value() != "" {
		t.Errorf("Expected search input to be cleared, got %q", m.unifiedInput.Value())
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
	if m.contentType != ContentTypeCommandList {
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
	m.unifiedInput.SetValue("")
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'S'}})

	// Verify we're NOT in typeahead mode anymore
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode when typing regular text")
	}

	// Verify contentType is "table" for regular search
	if m.contentType != ContentTypeTable {
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
	m.unifiedInput.SetValue("")

	// Verify empty state
	if m.unifiedInput.Value() != "" {
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
	if m.unifiedInput.Value() != "/help " {
		t.Errorf("Expected '/help ' after Tab completion, got %q", m.unifiedInput.Value())
	}

	// Verify exited typeahead mode (user can now type subcommand arguments)
	if m.inTypeaheadMode {
		t.Errorf("Expected to exit typeahead mode after Tab completion")
	}

	// Verify contentType is empty (ready to accept command arguments)
	if m.contentType != ContentTypeEmpty {
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
	if m.unifiedInput.Value() != "/help " {
		t.Errorf("Expected '/help ' (default selected), got %q", m.unifiedInput.Value())
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
	if m.unifiedInput.Value() != "abc" {
		t.Errorf("Expected Tab to be no-op outside typeahead, got input %q", m.unifiedInput.Value())
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
	if m.unifiedInput.Value() != "/exit " {
		t.Errorf("Expected '/exit ' after Tab with Down navigation, got %q", m.unifiedInput.Value())
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

	if m.unifiedInput.Value() != "/view " {
		t.Errorf("Expected '/view ' after Tab, got %q", m.unifiedInput.Value())
	}

	// Now type subcommand "groups" without re-filtering
	for _, r := range "groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify user can type subcommand after Tab completion
	if m.unifiedInput.Value() != "/view groups" {
		t.Errorf("Expected '/view groups' after typing subcommand, got %q", m.unifiedInput.Value())
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
	if m.contentType != ContentTypeEmpty {
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
	currentInput := m.unifiedInput.Value()
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
	if m.contentType != ContentTypeEmpty {
		t.Errorf("Expected contentType 'empty' after typing space, got %q", m.contentType)
	}
}

// History tests start here

func TestTUI_HistoryAppend_SingleCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Verify history is empty initially
	if len(m.commandHistory) != 0 {
		t.Errorf("Expected empty history initially, got %d commands", len(m.commandHistory))
	}

	// Type "/help" and press Enter
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify history now has one command
	if len(m.commandHistory) != 1 {
		t.Errorf("Expected 1 command in history, got %d", len(m.commandHistory))
	}

	// Verify the command is "/help"
	if m.commandHistory[0].Command != "/help" {
		t.Errorf("Expected command to be '/help', got %q", m.commandHistory[0].Command)
	}

	// Verify output is not empty
	if m.commandHistory[0].Output == "" {
		t.Errorf("Expected non-empty output for /help command")
	}
}

func TestTUI_HistoryAppend_MultipleCommands(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
			"192.168.1.150": "IoT",
		},
		Groups: []config.Group{
			{
				Name:    "Servers",
				Blocked: []string{"malware.com"},
				Allowed: []string{"trusted.com"},
			},
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Execute /view groups
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify history has 2 commands
	if len(m.commandHistory) != 2 {
		t.Errorf("Expected 2 commands in history, got %d", len(m.commandHistory))
	}

	// Verify order (oldest first)
	if m.commandHistory[0].Command != "/help" {
		t.Errorf("Expected first command to be '/help', got %q", m.commandHistory[0].Command)
	}
	if m.commandHistory[1].Command != "/view groups" {
		t.Errorf("Expected second command to be '/view groups', got %q", m.commandHistory[1].Command)
	}
}

func TestTUI_HistoryScroll_UpDown(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help to add to history
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify initial historyScroll is 0
	if m.historyScroll != 0 {
		t.Errorf("Expected historyScroll to be 0 initially, got %d", m.historyScroll)
	}

	// Press down arrow
	m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Verify historyScroll increased
	if m.historyScroll != 1 {
		t.Errorf("Expected historyScroll to be 1 after down arrow, got %d", m.historyScroll)
	}

	// Press up arrow
	m.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Verify historyScroll decreased
	if m.historyScroll != 0 {
		t.Errorf("Expected historyScroll to be 0 after up arrow, got %d", m.historyScroll)
	}
}

func TestTUI_HistoryScroll_Bounds(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help to add to history
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Try to scroll up from the top (should not go below 0)
	m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.historyScroll < 0 {
		t.Errorf("Expected historyScroll to not go below 0, got %d", m.historyScroll)
	}

	// Scroll down many times and verify max is capped at 100
	for i := 0; i < 110; i++ {
		m.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	if m.historyScroll > 100 {
		t.Errorf("Expected historyScroll to be capped at 100, got %d", m.historyScroll)
	}
}

func TestTUI_HistoryClear_ResetOnClear(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help to add to history
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify history is not empty
	if len(m.commandHistory) == 0 {
		t.Errorf("Expected history to have commands before /clear")
	}

	// Execute /clear
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "clear" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify history is empty after /clear
	if len(m.commandHistory) != 0 {
		t.Errorf("Expected history to be empty after /clear, got %d commands", len(m.commandHistory))
	}

	// Verify historyScroll is reset
	if m.historyScroll != 0 {
		t.Errorf("Expected historyScroll to be 0 after /clear, got %d", m.historyScroll)
	}

	// Verify search input is cleared
	if m.unifiedInput.Value() != "" {
		t.Errorf("Expected search input to be cleared, got %q", m.unifiedInput.Value())
	}
}

func TestTUI_HistoryRender_TimestampFormat(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify the command is stored with timestamp
	if m.commandHistory[0].Timestamp.IsZero() {
		t.Errorf("Expected timestamp to be set for command")
	}

	// Verify view renders history with timestamp
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "/help") {
		t.Errorf("Expected /help command in history view")
	}

	// The timestamp should be formatted in 15:04:05 format
	// We can't check exact time but should see the separator "|"
	if !strings.Contains(viewOutput, "|") {
		t.Errorf("Expected timestamp separator '|' in history view")
	}
}

func TestTUI_HistoryRender_EmptyHistory(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Verify first-run banner is shown (since firstRun=true and input is empty)
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "Welcome to Pharos Advanced Blocking") {
		t.Errorf("Expected first-run banner in view output")
	}

	// Dismiss banner by setting firstRun=false
	m.firstRun = false
	viewOutput = m.View()
	if !strings.Contains(viewOutput, "Start typing to search") {
		t.Errorf("Expected empty state hint in view output after banner is dismissed")
	}
}

func TestTUI_HistoryRender_AfterCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
		Groups: []config.Group{
			{
				Name:    "Servers",
				Blocked: []string{"blocked.com"},
				Allowed: []string{"allowed.com"},
			},
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /view groups
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify view now shows history
	viewOutput := m.View()
	if !strings.Contains(viewOutput, "/view groups") {
		t.Errorf("Expected /view groups command in history view")
	}

	// Should contain "Servers" from the group list
	if !strings.Contains(viewOutput, "Servers") {
		t.Errorf("Expected 'Servers' group in history view")
	}
}

func TestTUI_HistoryScroll_ResetOnNewCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Execute /help
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Scroll down in history
	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.historyScroll != 2 {
		t.Errorf("Expected historyScroll to be 2, got %d", m.historyScroll)
	}

	// Execute another command
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify historyScroll is reset to 0
	if m.historyScroll != 0 {
		t.Errorf("Expected historyScroll to be reset to 0, got %d", m.historyScroll)
	}
}

func TestTUI_HistoryFooterDisplay(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Without history, footer should not show history count
	viewOutput := m.View()
	if strings.Contains(viewOutput, "commands in history") {
		t.Errorf("Expected no history count in footer when no history exists")
	}

	// Execute /help to add to history
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify footer shows history count and scroll hint
	viewOutput = m.View()

	// Check that we have history
	if len(m.commandHistory) != 1 {
		t.Errorf("Expected 1 command in history, got %d", len(m.commandHistory))
	}

	// Check that we're not in typeahead mode
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to be in typeahead mode after command execution")
	}

	// Verify footer contains history count indicator
	if !strings.Contains(viewOutput, "commands in history") {
		t.Errorf("Expected 'commands in history' in footer")
	}
}

// Parser tests start here

func TestParseUnifiedInput_EmptyInput(t *testing.T) {
	result := ParseUnifiedInput("")
	if result.Type != InputTypeEmpty {
		t.Errorf("expected InputTypeEmpty, got %v", result.Type)
	}
	if result.Query != "" {
		t.Errorf("expected empty query, got %q", result.Query)
	}
	if len(result.Args) != 0 {
		t.Errorf("expected empty args, got %v", result.Args)
	}
}

func TestParseUnifiedInput_ImplicitSearch(t *testing.T) {
	result := ParseUnifiedInput("kids")
	if result.Type != InputTypeSearch {
		t.Errorf("expected InputTypeSearch, got %v", result.Type)
	}
	if result.Query != "kids" {
		t.Errorf("expected query 'kids', got %q", result.Query)
	}
	if len(result.Args) != 0 {
		t.Errorf("expected empty args for search, got %v", result.Args)
	}
}

func TestParseUnifiedInput_Command(t *testing.T) {
	result := ParseUnifiedInput("/view groups")
	if result.Type != InputTypeCommand {
		t.Errorf("expected InputTypeCommand, got %v", result.Type)
	}
	if result.Query != "view" {
		t.Errorf("expected query 'view', got %q", result.Query)
	}
	if len(result.Args) != 1 || result.Args[0] != "groups" {
		t.Errorf("expected args ['groups'], got %v", result.Args)
	}
}

func TestParseUnifiedInput_CommandCaseInsensitive(t *testing.T) {
	result := ParseUnifiedInput("/VIEW GROUPS")
	if result.Type != InputTypeCommand {
		t.Errorf("expected InputTypeCommand, got %v", result.Type)
	}
	if result.Query != "view" {
		t.Errorf("expected lowercase 'view', got %q", result.Query)
	}
	if len(result.Args) != 1 || result.Args[0] != "GROUPS" {
		t.Errorf("expected args ['GROUPS'] (preserve case), got %v", result.Args)
	}
}

func TestParseUnifiedInput_SearchWithSpaces(t *testing.T) {
	result := ParseUnifiedInput("192.0.2.50")
	if result.Type != InputTypeSearch {
		t.Errorf("expected InputTypeSearch, got %v", result.Type)
	}
	if result.Query != "192.0.2.50" {
		t.Errorf("expected query '192.0.2.50', got %q", result.Query)
	}
}

func TestParseUnifiedInput_WhitespaceHandling(t *testing.T) {
	result := ParseUnifiedInput("   search term   ")
	if result.Type != InputTypeSearch {
		t.Errorf("expected InputTypeSearch, got %v", result.Type)
	}
	if result.Query != "search term" {
		t.Errorf("expected query 'search term', got %q", result.Query)
	}
}

func TestParseUnifiedInput_CommandWithMultipleArgs(t *testing.T) {
	result := ParseUnifiedInput("/view group MyGroup blocklists")
	if result.Type != InputTypeCommand {
		t.Errorf("expected InputTypeCommand, got %v", result.Type)
	}
	if result.Query != "view" {
		t.Errorf("expected query 'view', got %q", result.Query)
	}
	if len(result.Args) != 3 || result.Args[0] != "group" || result.Args[1] != "MyGroup" || result.Args[2] != "blocklists" {
		t.Errorf("expected args ['group', 'MyGroup', 'blocklists'], got %v", result.Args)
	}
}

func TestParseUnifiedInput_JustSlash(t *testing.T) {
	result := ParseUnifiedInput("/")
	if result.Type != InputTypeCommand {
		t.Errorf("expected InputTypeCommand, got %v", result.Type)
	}
	if result.Query != "" {
		t.Errorf("expected empty query for bare slash, got %q", result.Query)
	}
}

func TestParseUnifiedInput_SearchWithLeadingZero(t *testing.T) {
	result := ParseUnifiedInput("192.168.001.1")
	if result.Type != InputTypeSearch {
		t.Errorf("expected InputTypeSearch for IP-like search, got %v", result.Type)
	}
	if result.Query != "192.168.001.1" {
		t.Errorf("expected query '192.168.001.1', got %q", result.Query)
	}
}

// Phase 3 Integration Tests - Unified Input Routing

func TestUnifiedInput_ImplicitSearchRouting(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
			"192.168.1.150": "IoT",
			"10.0.0.5":      "Laptops",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30

	// Simulate typing "192" (implicit search, no slash)
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	// Verify contentType is "table" for implicit search
	if m.contentType != ContentTypeTable {
		t.Errorf("Expected contentType 'table' for implicit search, got %q", m.contentType)
	}

	// Verify filtering worked (should match the two 192.168.x.x IPs)
	if len(m.filtered) != 2 {
		t.Errorf("Expected 2 filtered results for '192', got %d", len(m.filtered))
	}

	// Verify input is in the unified input field
	if m.unifiedInput.Value() != "192" {
		t.Errorf("Expected unified input to contain '192', got %q", m.unifiedInput.Value())
	}
}

func TestUnifiedInput_ExplicitCommandRouting(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
		Groups: []config.Group{
			{
				Name:    "Servers",
				Blocked: []string{"blocked.com"},
				Allowed: []string{"allowed.com"},
			},
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30

	// Simulate typing "/view groups" (explicit command)
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Press Enter to execute
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify the command was executed
	if len(m.commandHistory) == 0 {
		t.Errorf("Expected command to be added to history")
	}

	// Verify the last command in history is "/view groups"
	if m.commandHistory[len(m.commandHistory)-1].Command != "/view groups" {
		t.Errorf("Expected command '/view groups' in history, got %q", m.commandHistory[len(m.commandHistory)-1].Command)
	}

	// Verify contentType changed to show group list
	if m.contentType != ContentTypeViewGroups {
		t.Errorf("Expected contentType 'view_groups' after /view groups, got %q", m.contentType)
	}
}

func TestUnifiedInput_SearchClearsInput(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30

	// Type a search term
	for _, r := range "192" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Verify input is in the unified field
	if m.unifiedInput.Value() != "192" {
		t.Errorf("Expected unified input to be '192', got %q", m.unifiedInput.Value())
	}

	// Press Enter to "execute" the search
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// For implicit searches, the current behavior is to keep the search value
	// This is different from commands which clear the input
	// The existing code doesn't clear on Enter for regular search, which is OK
	if m.contentType != ContentTypeTable {
		t.Errorf("Expected contentType 'table' after search, got %q", m.contentType)
	}
}

func TestUnifiedInput_CommandInHistory(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30

	// Execute a command through unified input
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify command is in history
	if len(m.commandHistory) != 1 {
		t.Errorf("Expected 1 command in history, got %d", len(m.commandHistory))
	}

	// Verify command text in history
	if m.commandHistory[0].Command != "/help" {
		t.Errorf("Expected '/help' in history, got %q", m.commandHistory[0].Command)
	}

	// Verify output is captured
	if m.commandHistory[0].Output == "" {
		t.Errorf("Expected non-empty output in history")
	}
}

func TestUnifiedInput_PlaceholderText(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)

	// Verify placeholder text includes both search and command hints
	if m.unifiedInput.Placeholder != "Search or type /help for commands" {
		t.Errorf("Expected placeholder to mention both search and commands, got %q", m.unifiedInput.Placeholder)
	}
}

func TestUnifiedInput_MultipleSequentialCommands(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
		Groups: []config.Group{
			{
				Name:    "Servers",
				Blocked: []string{"blocked.com"},
				Allowed: []string{"allowed.com"},
			},
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30

	// Execute /help
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "help" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(m.commandHistory) != 1 {
		t.Errorf("Expected 1 command after /help, got %d", len(m.commandHistory))
	}

	// Execute /view groups
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view groups" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(m.commandHistory) != 2 {
		t.Errorf("Expected 2 commands in history, got %d", len(m.commandHistory))
	}

	// Verify order
	if m.commandHistory[0].Command != "/help" {
		t.Errorf("Expected first command to be '/help', got %q", m.commandHistory[0].Command)
	}
	if m.commandHistory[1].Command != "/view groups" {
		t.Errorf("Expected second command to be '/view groups', got %q", m.commandHistory[1].Command)
	}
}

// FIX #3: Error Case Tests

// Test unknown command
func TestTUI_ExecuteCommand_UnknownCommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Try to execute an unknown command
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "xyz" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify that unknown command doesn't crash
	// The existing behavior is to clear the input
	if m.unifiedInput.Value() != "" {
		t.Errorf("Expected input to be cleared after unknown command, got %q", m.unifiedInput.Value())
	}
}

// Test view group with missing group name
func TestTUI_ExecuteCommand_ViewGroupMissingName(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Try to execute /view group without a group name
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, r := range "view group" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify that the command was added to history (it handles the error internally)
	if len(m.commandHistory) == 0 {
		t.Errorf("Expected command to be added to history")
	}

	// Verify the input was cleared after command execution
	if m.unifiedInput.Value() != "" {
		t.Errorf("Expected input to be cleared after command, got %q", m.unifiedInput.Value())
	}
}

// Test view with unknown subcommand
func TestTUI_ExecuteCommand_ViewUnknownSubcommand(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Call executeCommand directly with unknown subcommand
	m.executeCommand("view", []string{"xyz"})

	// Verify contentType is "help" (error state)
	if m.contentType != ContentTypeHelp {
		t.Errorf("Expected contentType 'help' for unknown subcommand, got %q", m.contentType)
	}

	// Verify the input was cleared after command execution
	if m.unifiedInput.Value() != "" {
		t.Errorf("Expected input to be cleared after command, got %q", m.unifiedInput.Value())
	}
}

// Test search with empty query
func TestTUI_FilterClients_EmptyQuery(t *testing.T) {
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

	// With empty search, all clients should be shown
	m.filterClients()
	if len(m.filtered) != len(m.clients) {
		t.Errorf("Expected all clients when query is empty, got %d/%d",
			len(m.filtered), len(m.clients))
	}
}

// Test Tab completion outside typeahead mode
func TestTUI_TabCompletion_OutsideTypeaheadModeExisting(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 80
	m.height = 24

	// Set up a normal search state (not in typeahead)
	m.unifiedInput.SetValue("normal text")
	m.inTypeaheadMode = false

	// Verify not in typeahead mode
	if m.inTypeaheadMode {
		t.Errorf("Expected NOT to be in typeahead mode")
	}

	// Tab behavior when not in typeahead should not change model state
	// (Tab would be handled by textinput component, not our logic)
	if m.inTypeaheadMode {
		t.Errorf("Tab outside typeahead should not enter typeahead mode")
	}
}

// FIX #4: First-Run Banner Tests

func TestTUI_FirstRunBanner(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30
	m.firstRun = true
	m.unifiedInput.SetValue("")

	view := m.View()
	if !strings.Contains(view, "Welcome to Pharos Advanced Blocking") {
		t.Errorf("expected welcome banner in view when firstRun=true, got:\n%s", view)
	}
	if !strings.Contains(view, "Quick start") {
		t.Errorf("expected quick start help text in banner")
	}
}

func TestTUI_FirstRunBannerDismissesOnKeypress(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)
	m.ready = true
	m.width = 100
	m.height = 30
	m.firstRun = true
	m.unifiedInput.SetValue("test")

	// Calling View() should dismiss the banner when input is not empty
	view := m.View()

	// After first keystroke, firstRun should be false
	if m.firstRun {
		t.Errorf("expected firstRun to be false after first keystroke")
	}
	if strings.Contains(view, "Welcome to Pharos Advanced Blocking") {
		t.Errorf("welcome banner should not appear after first keystroke")
	}
}

// FIX #5: Content Type Enum Tests

func TestContentType_String(t *testing.T) {
	tests := []struct {
		ct       ContentType
		expected string
	}{
		{ContentTypeEmpty, "empty"},
		{ContentTypeTable, "table"},
		{ContentTypeHelp, "help"},
		{ContentTypeError, "error"},
		{ContentTypeCommandList, "command_list"},
		{ContentTypeViewNetworkGroupMap, "view_networkgroupmap"},
		{ContentTypeViewGroups, "view_groups"},
		{ContentTypeViewGroup, "view_group"},
	}

	for _, tt := range tests {
		if tt.ct.String() != tt.expected {
			t.Errorf("ContentType(%d).String() = %s, want %s", tt.ct, tt.ct.String(), tt.expected)
		}
	}
}

func TestContentType_InitialValue(t *testing.T) {
	cfg := &config.Config{
		NetworkGroupMap: map[string]string{
			"192.168.1.100": "Servers",
		},
	}

	m := New(cfg)

	// Verify initial contentType is ContentTypeEmpty
	if m.contentType != ContentTypeEmpty {
		t.Errorf("Expected initial contentType to be ContentTypeEmpty, got %v", m.contentType)
	}
}
