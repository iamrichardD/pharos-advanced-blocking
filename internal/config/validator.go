package config

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
)

// ValidationError represents a specific schema or configuration validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation failed: %s", e.Message)
}

// ValidationErrors collects multiple validation errors.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString("validation failed:")
	for _, e := range ve {
		sb.WriteString("\n- ")
		sb.WriteString(e.Error())
	}
	return sb.String()
}

// Unwrap returns the list of individual errors to support errors.As/Is.
func (ve ValidationErrors) Unwrap() []error {
	errs := make([]error, len(ve))
	for i := range ve {
		errs[i] = &ve[i]
	}
	return errs
}

// ValidateConfig verifies that the configuration structure is valid.
// It ensures that:
// 1. All keys in NetworkGroupMap are strictly valid IP addresses or CIDR blocks (rejecting hostnames/MAC addresses).
// 2. All target groups in NetworkGroupMap actually exist in the Groups definition list.
func ValidateConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("configuration is nil")
	}

	var errs []ValidationError

	// Build a fast lookup set of defined group names
	definedGroups := make(map[string]bool)
	for _, g := range cfg.Groups {
		if g.Name == "" {
			errs = append(errs, ValidationError{
				Field:   "groups",
				Message: "group name cannot be empty",
			})
		} else {
			definedGroups[g.Name] = true
		}
	}

	// Extract and sort NetworkGroupMap keys to ensure deterministic error ordering
	clientKeys := make([]string, 0, len(cfg.NetworkGroupMap))
	for k := range cfg.NetworkGroupMap {
		clientKeys = append(clientKeys, k)
	}
	slices.Sort(clientKeys)

	// Validate NetworkGroupMap keys and target groups in sorted order
	for _, clientKey := range clientKeys {
		groupName := cfg.NetworkGroupMap[clientKey]
		// 1. Verify target group exists
		if !definedGroups[groupName] {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("networkGroupMap[%s]", clientKey),
				Message: fmt.Sprintf("target group %q does not exist", groupName),
			})
		}

		// 2. Verify key is a valid IP or CIDR (No MACs, no hostnames)
		if strings.Contains(clientKey, "/") {
			_, _, err := net.ParseCIDR(clientKey)
			if err != nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("networkGroupMap[%s]", clientKey),
					Message: fmt.Sprintf("invalid CIDR network block: %v", err),
				})
			}
		} else {
			ip := net.ParseIP(clientKey)
			if ip == nil {
				errs = append(errs, ValidationError{
					Field:   fmt.Sprintf("networkGroupMap[%s]", clientKey),
					Message: "must be a valid IPv4 or IPv6 address (MAC addresses and hostnames are rejected)",
				})
			}
		}
	}

	if len(errs) > 0 {
		return ValidationErrors(errs)
	}

	return nil
}
