package config

import (
	"encoding/json"
	"testing"
)

func TestConfigUnmarshalJSON_EmptyGroups(t *testing.T) {
	testCases := []struct {
		name          string
		json          string
		wantErr       bool
		wantGroupsLen int
	}{
		{
			name:          "empty object format",
			json:          `{"enableBlocking": true, "groups": {}}`,
			wantErr:       false,
			wantGroupsLen: 0,
		},
		{
			name:          "empty array format",
			json:          `{"enableBlocking": true, "groups": []}`,
			wantErr:       false,
			wantGroupsLen: 0,
		},
		{
			name:          "groups null",
			json:          `{"enableBlocking": true, "groups": null}`,
			wantErr:       false,
			wantGroupsLen: 0,
		},
		{
			name:          "malformed groups (string)",
			json:          `{"enableBlocking": true, "groups": "invalid"}`,
			wantErr:       true,
			wantGroupsLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cfg Config
			err := json.Unmarshal([]byte(tc.json), &cfg)
			if (err != nil) != tc.wantErr {
				t.Fatalf("want err=%v, got %v", tc.wantErr, err)
			}
			if len(cfg.Groups) != tc.wantGroupsLen {
				t.Errorf("want %d groups, got %d", tc.wantGroupsLen, len(cfg.Groups))
			}
		})
	}
}
