package config

import (
	"encoding/json"
)

// UrlEntry represents a URL block/allow list entry, which can either be a simple string or a detailed object.
type UrlEntry struct {
	URL               string   `json:"url"`
	BlockAsNxDomain   *bool    `json:"blockAsNxDomain,omitempty"`
	BlockingAddresses []string `json:"blockingAddresses,omitempty"`
}

// Group defines specific filtering settings, lists, and rules applied to a matching client.
type Group struct {
	Name                   string   `json:"name"`
	EnableBlocking         bool     `json:"enableBlocking"`
	AllowTxtBlockingReport bool     `json:"allowTxtBlockingReport"`
	BlockAsNxDomain        bool     `json:"blockAsNxDomain"`
	BlockingAddresses      []string `json:"blockingAddresses,omitempty"`
	Allowed                []string `json:"allowedDomains,omitempty"`
	Blocked                []string `json:"blockedDomains,omitempty"`
	AllowListUrls          []string `json:"allowListUrls,omitempty"`
	BlockListUrls          []string `json:"blockLists,omitempty"`
	AllowedRegex           []string `json:"allowedRegex,omitempty"`
	BlockedRegex           []string `json:"blockedRegex,omitempty"`
	RegexAllowListUrls     []string `json:"regexAllowListUrls,omitempty"`
	RegexBlockListUrls     []string `json:"regexBlockListUrls,omitempty"`
	AdblockListUrls        []string `json:"adblockListUrls,omitempty"`
}

// Config represents the complete schema for the Technitium Advanced Blocking App configuration (dnsApp.config).
type Config struct {
	EnableBlocking                    bool              `json:"enableBlocking"`
	BlockingAnswerTtl                 uint32            `json:"blockingAnswerTtl"`
	BlockListUrlUpdateIntervalHours   int               `json:"blockListUrlUpdateIntervalHours"`
	BlockListUrlUpdateIntervalMinutes int               `json:"blockListUrlUpdateIntervalMinutes"`
	LocalEndPointGroupMap             map[string]string `json:"localEndPointGroupMap,omitempty"`
	NetworkGroupMap                   map[string]string `json:"networkGroupMap,omitempty"`
	Groups                            []Group           `json:"groups"`
}

// UnmarshalJSON handles both Technitium's object format (groups as map) and
// Go's standard array format. This custom unmarshaler converts Technitium's
// object format {"groups": {"GroupName": {...}}} to our array format with
// explicit Name fields, enabling seamless config loading without conversion.
func (c *Config) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal as an object-format config (groups as map)
	type rawConfigObject struct {
		EnableBlocking                    bool              `json:"enableBlocking"`
		BlockingAnswerTtl                 uint32            `json:"blockingAnswerTtl"`
		BlockListUrlUpdateIntervalHours   int               `json:"blockListUrlUpdateIntervalHours"`
		BlockListUrlUpdateIntervalMinutes int               `json:"blockListUrlUpdateIntervalMinutes"`
		LocalEndPointGroupMap             map[string]string `json:"localEndPointGroupMap,omitempty"`
		NetworkGroupMap                   map[string]string `json:"networkGroupMap,omitempty"`
		Groups                            map[string]Group  `json:"groups"` // Object format: group name as key
	}

	// Try parsing as object format first
	var rawObj rawConfigObject
	if err := json.Unmarshal(data, &rawObj); err == nil {
		// Successfully parsed as object format (even if groups is empty)
		c.Groups = make([]Group, 0, len(rawObj.Groups))
		for name, group := range rawObj.Groups {
			group.Name = name // Add name from object key
			c.Groups = append(c.Groups, group)
		}

		// Copy other fields
		c.EnableBlocking = rawObj.EnableBlocking
		c.BlockingAnswerTtl = rawObj.BlockingAnswerTtl
		c.BlockListUrlUpdateIntervalHours = rawObj.BlockListUrlUpdateIntervalHours
		c.BlockListUrlUpdateIntervalMinutes = rawObj.BlockListUrlUpdateIntervalMinutes
		c.LocalEndPointGroupMap = rawObj.LocalEndPointGroupMap
		c.NetworkGroupMap = rawObj.NetworkGroupMap

		return nil
	}

	// If object format didn't work, try standard format (groups as array)
	type rawConfigArray struct {
		EnableBlocking                    bool              `json:"enableBlocking"`
		BlockingAnswerTtl                 uint32            `json:"blockingAnswerTtl"`
		BlockListUrlUpdateIntervalHours   int               `json:"blockListUrlUpdateIntervalHours"`
		BlockListUrlUpdateIntervalMinutes int               `json:"blockListUrlUpdateIntervalMinutes"`
		LocalEndPointGroupMap             map[string]string `json:"localEndPointGroupMap,omitempty"`
		NetworkGroupMap                   map[string]string `json:"networkGroupMap,omitempty"`
		Groups                            []Group           `json:"groups"` // Array format
	}

	var rawArr rawConfigArray
	if err := json.Unmarshal(data, &rawArr); err != nil {
		return err
	}

	// Copy all fields from array format
	c.EnableBlocking = rawArr.EnableBlocking
	c.BlockingAnswerTtl = rawArr.BlockingAnswerTtl
	c.BlockListUrlUpdateIntervalHours = rawArr.BlockListUrlUpdateIntervalHours
	c.BlockListUrlUpdateIntervalMinutes = rawArr.BlockListUrlUpdateIntervalMinutes
	c.LocalEndPointGroupMap = rawArr.LocalEndPointGroupMap
	c.NetworkGroupMap = rawArr.NetworkGroupMap
	c.Groups = rawArr.Groups

	return nil
}
