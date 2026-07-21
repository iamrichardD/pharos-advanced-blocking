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
