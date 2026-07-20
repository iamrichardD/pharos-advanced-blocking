package config

// UrlEntry represents a URL block/allow list entry, which can either be a simple string or a detailed object.
type UrlEntry struct {
	URL               string   `json:"url"`
	BlockAsNxDomain   *bool    `json:"blockAsNxDomain,omitempty"`
	BlockingAddresses []string `json:"blockingAddresses,omitempty"`
}

// Group defines specific filtering settings, lists, and rules applied to a matching client.
type Group struct {
	Name                   string     `json:"name"`
	EnableBlocking         bool       `json:"enableBlocking"`
	AllowTxtBlockingReport bool       `json:"allowTxtBlockingReport"`
	BlockAsNxDomain        bool       `json:"blockAsNxDomain"`
	BlockingAddresses      []string   `json:"blockingAddresses,omitempty"`
	Allowed                []string   `json:"allowed,omitempty"`
	Blocked                []string   `json:"blocked,omitempty"`
	AllowListUrls          []string   `json:"allowListUrls,omitempty"`
	BlockListUrls          []UrlEntry `json:"blockListUrls,omitempty"`
	AllowedRegex           []string   `json:"allowedRegex,omitempty"`
	BlockedRegex           []string   `json:"blockedRegex,omitempty"`
	RegexAllowListUrls     []string   `json:"regexAllowListUrls,omitempty"`
	RegexBlockListUrls     []UrlEntry `json:"regexBlockListUrls,omitempty"`
	AdblockListUrls        []UrlEntry `json:"adblockListUrls,omitempty"`
}

// Config represents the complete schema for the Technitium Advanced Blocking App configuration (dnsApp.config).
type Config struct {
	EnableBlocking                      bool              `json:"enableBlocking"`
	BlockingAnswerTtl                   uint32            `json:"blockingAnswerTtl"`
	BlockListUrlUpdateIntervalHours     int               `json:"blockListUrlUpdateIntervalHours"`
	BlockListUrlUpdateIntervalMinutes   int               `json:"blockListUrlUpdateIntervalMinutes"`
	LocalEndPointGroupMap               map[string]string `json:"localEndPointGroupMap,omitempty"`
	NetworkGroupMap                     map[string]string `json:"networkGroupMap,omitempty"`
	Groups                              []Group           `json:"groups"`
}
