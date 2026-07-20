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
	if !strings.Contains(viewOutput, "Pharos Advanced Blocking TUI") {
		t.Errorf("Expected view to contain title")
	}
	if !strings.Contains(viewOutput, "ctrl+c / esc: exit | start typing to search") {
		t.Errorf("Expected view to contain footer help status line")
	}
}
