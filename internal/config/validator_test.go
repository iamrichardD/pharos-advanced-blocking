package config

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		wantErr     bool
		errField    string
		errContains string
	}{
		{
			name: "valid configuration with IP and CIDR keys",
			cfg: &Config{
				NetworkGroupMap: map[string]string{
					"192.168.1.50":   "default-group",
					"10.0.0.0/24":    "secure-group",
					"2001:db8::/32":  "secure-group",
					"fe80::1ff:fe23": "default-group",
				},
				Groups: []Group{
					{Name: "default-group", EnableBlocking: true},
					{Name: "secure-group", EnableBlocking: true},
				},
			},
			wantErr: false,
		},
		{
			name: "missing group name",
			cfg: &Config{
				Groups: []Group{
					{Name: ""},
				},
			},
			wantErr:     true,
			errField:    "groups",
			errContains: "group name cannot be empty",
		},
		{
			name: "missing target group in networkGroupMap",
			cfg: &Config{
				NetworkGroupMap: map[string]string{
					"192.168.1.1": "non-existent",
				},
				Groups: []Group{
					{Name: "default-group"},
				},
			},
			wantErr:     true,
			errField:    "networkGroupMap[192.168.1.1]",
			errContains: `target group "non-existent" does not exist`,
		},
		{
			name: "invalid IP key (MAC address rejected)",
			cfg: &Config{
				NetworkGroupMap: map[string]string{
					"00:1A:2B:3C:4D:5E": "default-group",
				},
				Groups: []Group{
					{Name: "default-group"},
				},
			},
			wantErr:     true,
			errField:    "networkGroupMap[00:1A:2B:3C:4D:5E]",
			errContains: "must be a valid IPv4 or IPv6 address",
		},
		{
			name: "invalid IP key (hostname rejected)",
			cfg: &Config{
				NetworkGroupMap: map[string]string{
					"my-client.local": "default-group",
				},
				Groups: []Group{
					{Name: "default-group"},
				},
			},
			wantErr:     true,
			errField:    "networkGroupMap[my-client.local]",
			errContains: "must be a valid IPv4 or IPv6 address",
		},
		{
			name: "invalid CIDR key",
			cfg: &Config{
				NetworkGroupMap: map[string]string{
					"192.168.1.0/33": "default-group",
				},
				Groups: []Group{
					{Name: "default-group"},
				},
			},
			wantErr:     true,
			errField:    "networkGroupMap[192.168.1.0/33]",
			errContains: "invalid CIDR network block",
		},
		{
			name:    "nil configuration",
			cfg:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				var valErr *ValidationError
				if errors.As(err, &valErr) {
					if tt.errField != "" && valErr.Field != tt.errField {
						t.Errorf("expected error field %q, got %q", tt.errField, valErr.Field)
					}
					if tt.errContains != "" && !strings.Contains(valErr.Message, tt.errContains) {
						t.Errorf("expected error to contain %q, got %q", tt.errContains, valErr.Message)
					}
				}
			}
		})
	}
}

func TestValidateConfigAggregation(t *testing.T) {
	cfg := &Config{
		NetworkGroupMap: map[string]string{
			"invalid-ip":     "non-existent-1",
			"192.168.1.0/33": "non-existent-2",
		},
		Groups: []Group{
			{Name: ""},
		},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected validation errors, got nil")
	}

	var valErrs ValidationErrors
	if !errors.As(err, &valErrs) {
		t.Fatalf("expected error to be of type ValidationErrors, got: %T (%v)", err, err)
	}

	if len(valErrs) != 5 {
		t.Errorf("expected exactly 5 validation errors, got %d: %v", len(valErrs), valErrs)
	}

	expectedFields := map[string]bool{
		"groups":                          true,
		"networkGroupMap[invalid-ip]":     true,
		"networkGroupMap[192.168.1.0/33]": true,
	}

	for _, e := range valErrs {
		if !expectedFields[e.Field] {
			t.Errorf("unexpected error field found: %q (message: %q)", e.Field, e.Message)
		}
	}
}
