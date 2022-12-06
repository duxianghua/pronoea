package v1

import (
	"time"

	"github.com/alecthomas/units"
	"github.com/prometheus/common/config"
)

type Module struct {
	Prober  string        `json:"prober"`
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	HTTP    HTTPProbe     `json:"http,omitempty" yaml:"http,omitempty"`
	TCP     TCPProbe      `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	ICMP    ICMPProbe     `json:"icmp,omitempty" yaml:"icmp,omitempty"`
	DNS     DNSProbe      `json:"dns,omitempty" yaml:"dns,omitempty"`
	GRPC    GRPCProbe     `json:"grpc,omitempty" yaml:"grpc,omitempty"`
}

type HTTPProbe struct {
	// Defaults to 2xx.
	ValidStatusCodes             []int             `json:"valid_status_codes,omitempty" yaml:"valid_status_codes,omitempty"`
	ValidHTTPVersions            []string          `json:"valid_http_versions,omitempty" yaml:"valid_http_versions,omitempty"`
	IPProtocol                   string            `json:"preferred_ip_protocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	IPProtocolFallback           bool              `json:"ip_protocol_fallback,omitempty" yaml:"ip_protocol_fallback,omitempty"`
	SkipResolvePhaseWithProxy    bool              `json:"skip_resolve_phase_with_proxy,omitempty" yaml:"skip_resolve_phase_with_proxy,omitempty"`
	NoFollowRedirects            *bool             `json:"no_follow_redirects,omitempty" yaml:"no_follow_redirects,omitempty"`
	FailIfSSL                    bool              `json:"fail_if_ssl,omitempty" yaml:"fail_if_ssl,omitempty"`
	FailIfNotSSL                 bool              `json:"fail_if_not_ssl,omitempty" yaml:"fail_if_not_ssl,omitempty"`
	Method                       string            `json:"method,omitempty" yaml:"method,omitempty"`
	Headers                      map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	FailIfBodyMatchesRegexp      []string          `json:"fail_if_body_matches_regexp,omitempty" yaml:"fail_if_body_matches_regexp,omitempty"`
	FailIfBodyNotMatchesRegexp   []string          `json:"fail_if_body_not_matches_regexp,omitempty" yaml:"fail_if_body_not_matches_regexp,omitempty"`
	FailIfHeaderMatchesRegexp    []HeaderMatch     `json:"fail_if_header_matches,omitempty" yaml:"fail_if_header_matches,omitempty"`
	FailIfHeaderNotMatchesRegexp []HeaderMatch     `json:"fail_if_header_not_matches,omitempty" yaml:"fail_if_header_not_matches,omitempty"`
	Body                         string            `json:"body,omitempty" yaml:"body,omitempty"`
	HTTPClientConfig             `json:",inline" yaml:",inline"`
	Compression                  string           `json:"compression,omitempty" yaml:"compression,omitempty"`
	BodySizeLimit                units.Base2Bytes `json:"body_size_limit,omitempty" yaml:"body_size_limit,omitempty"`
}

type HTTPClientConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth config.BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
	// The HTTP authorization credentials for the targets.
	Authorization config.Authorization `yaml:"authorization,omitempty" json:"authorization,omitempty"`
	// The OAuth2 client credentials used to fetch a token for the targets.
	OAuth2 OAuth2 `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	// The bearer token for the targets. Deprecated in favour of
	// Authorization.Credentials.
	BearerToken config.Secret `yaml:"bearer_token,omitempty" json:"bearer_token,omitempty"`
	// The bearer token file for the targets. Deprecated in favour of
	// Authorization.CredentialsFile.
	BearerTokenFile string `yaml:"bearer_token_file,omitempty" json:"bearer_token_file,omitempty"`
	// HTTP proxy server to use to connect to the targets.
	ProxyURL string `yaml:"proxy_url,omitempty" json:"proxy_url,omitempty"`
	// TLSConfig to use to connect to the targets.
	TLSConfig TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
	// FollowRedirects specifies whether the client should follow HTTP 3xx redirects.
	// The omitempty flag is not set, because it would be hidden from the
	// marshalled configuration when set to false.
	FollowRedirects bool `yaml:"follow_redirects,omitempty" json:"follow_redirects,omitempty"`
}

// OAuth2 is the oauth2 client configuration.
type OAuth2 struct {
	ClientID         string            `yaml:"client_id" json:"client_id"`
	ClientSecret     string            `yaml:"client_secret" json:"client_secret"`
	ClientSecretFile string            `yaml:"client_secret_file" json:"client_secret_file"`
	Scopes           []string          `yaml:"scopes,omitempty" json:"scopes,omitempty"`
	TokenURL         string            `yaml:"token_url" json:"token_url"`
	EndpointParams   map[string]string `yaml:"endpoint_params,omitempty" json:"endpoint_params,omitempty"`

	// TLSConfig is used to connect to the token URL.
	TLSConfig TLSConfig `yaml:"tls_config,omitempty" json:"tls_config"`
}

// TLSConfig configures the options for TLS connections.
type TLSConfig struct {
	// The CA cert to use for the targets.
	CAFile string `yaml:"ca_file,omitempty" json:"ca_file,omitempty"`
	// The client cert file for the targets.
	CertFile string `yaml:"cert_file,omitempty" json:"cert_file,omitempty"`
	// The client key file for the targets.
	KeyFile string `yaml:"key_file,omitempty" json:"key_file,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty" json:"insecure_skip_verify,omitempty"`
	// Minimum TLS version.
	MinVersion string `yaml:"min_version,omitempty" json:"min_version,omitempty"`
}

type HeaderMatch struct {
	Header       string `json:"header,omitempty" yaml:"header,omitempty"`
	Regexp       string `json:"regexp,omitempty" yaml:"regexp,omitempty"`
	AllowMissing bool   `json:"allow_missing,omitempty" yaml:"allow_missing,omitempty"`
}

type TCPProbe struct {
	IPProtocol         string          `json:"preferred_ip_protocol,omitempty" yaml:"preferred_ip_protocol,omitempty"`
	IPProtocolFallback bool            `json:"ip_protocol_fallback,omitempty" yaml:"ip_protocol_fallback,omitempty"`
	SourceIPAddress    string          `json:"source_ip_address,omitempty" yaml:"source_ip_address,omitempty"`
	QueryResponse      []QueryResponse `json:"query_response,omitempty" yaml:"query_response,omitempty"`
	TLS                bool            `json:"tls,omitempty" yaml:"tls,omitempty"`
	TLSConfig          TLSConfig       `json:"tls_config,omitempty" yaml:"tls_config,omitempty"`
}

type ICMPProbe struct {
	IPProtocol         string `json:"preferred_ip_protocol,omitempty"` // Defaults to "ip6".
	IPProtocolFallback bool   `json:"ip_protocol_fallback,omitempty"`
	SourceIPAddress    string `json:"source_ip_address,omitempty"`
	PayloadSize        int    `json:"payload_size,omitempty"`
	DontFragment       bool   `json:"dont_fragment,omitempty"`
	TTL                int    `json:"ttl,omitempty"`
}

type QueryResponse struct {
	Expect   string `json:"expect,omitempty"`
	Send     string `json:"send,omitempty"`
	StartTLS bool   `json:"starttls,omitempty"`
}

type DNSProbe struct {
	IPProtocol         string         `json:"preferred_ip_protocol,omitempty"`
	IPProtocolFallback bool           `json:"ip_protocol_fallback,omitempty"`
	DNSOverTLS         bool           `json:"dns_over_tls,omitempty"`
	TLSConfig          TLSConfig      `json:"tls_config,omitempty"`
	SourceIPAddress    string         `json:"source_ip_address,omitempty"`
	TransportProtocol  string         `json:"transport_protocol,omitempty"`
	QueryClass         string         `json:"query_class,omitempty"` // Defaults to IN.
	QueryName          string         `json:"query_name,omitempty"`
	QueryType          string         `json:"query_type,omitempty"`        // Defaults to ANY.
	Recursion          bool           `json:"recursion_desired,omitempty"` // Defaults to true.
	ValidRcodes        []string       `json:"valid_rcodes,omitempty"`      // Defaults to NOERROR.
	ValidateAnswer     DNSRRValidator `json:"validate_answer_rrs,omitempty"`
	ValidateAuthority  DNSRRValidator `json:"validate_authority_rrs,omitempty"`
	ValidateAdditional DNSRRValidator `json:"validate_additional_rrs,omitempty"`
}

type DNSRRValidator struct {
	FailIfMatchesRegexp     []string `json:"fail_if_matches_regexp,omitempty"`
	FailIfAllMatchRegexp    []string `json:"fail_if_all_match_regexp,omitempty"`
	FailIfNotMatchesRegexp  []string `json:"fail_if_not_matches_regexp,omitempty"`
	FailIfNoneMatchesRegexp []string `json:"fail_if_none_matches_regexp,omitempty"`
}

type GRPCProbe struct {
	Service             string    `json:"service,omitempty"`
	TLS                 bool      `json:"tls,omitempty"`
	TLSConfig           TLSConfig `json:"tls_config,omitempty"`
	IPProtocolFallback  bool      `json:"ip_protocol_fallback,omitempty"`
	PreferredIPProtocol string    `json:"preferred_ip_protocol,omitempty"`
}
