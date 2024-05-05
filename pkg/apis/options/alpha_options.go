package options

// AlphaOptions contains alpha structured configuration options.
// Usage of these options allows users to access alpha features that are not
// available as part of the primary configuration structure for OAuth2 Proxy.
//
// :::warning
// The options within this structure are considered alpha.
// They may change between releases without notice.
// :::
type AlphaOptions struct {
	// ProxyOptions is used to configure the proxy behaviour.
	// This includes things like the prefix for protected paths, authentication
	// and routing options.
	ProxyOptions ProxyOptions `json:"proxyOptions,omitempty"`

	// ProbeOptions is used to configure the probe endpoint for health and readiness checks.
	ProbeOptions ProbeOptions `yaml:"probeOptions,omitempty"`

	// UpstreamConfig is used to configure upstream servers.
	// Once a user is authenticated, requests to the server will be proxied to
	// these upstream servers based on the path mappings defined in this list.
	UpstreamConfig UpstreamConfig `json:"upstreamConfig,omitempty"`

	// InjectRequestHeaders is used to configure headers that should be added
	// to requests to upstream servers.
	// Headers may source values from either the authenticated user's session
	// or from a static secret value.
	InjectRequestHeaders []Header `json:"injectRequestHeaders,omitempty"`

	// InjectResponseHeaders is used to configure headers that should be added
	// to responses from the proxy.
	// This is typically used when using the proxy as an external authentication
	// provider in conjunction with another proxy such as NGINX and its
	// auth_request module.
	// Headers may source values from either the authenticated user's session
	// or from a static secret value.
	InjectResponseHeaders []Header `json:"injectResponseHeaders,omitempty"`

	// Server is used to configure the HTTP(S) server for the proxy application.
	// You may choose to run both HTTP and HTTPS servers simultaneously.
	// This can be done by setting the BindAddress and the SecureBindAddress simultaneously.
	// To use the secure server you must configure a TLS certificate and key.
	Server Server `json:"server,omitempty"`

	// MetricsServer is used to configure the HTTP(S) server for metrics.
	// You may choose to run both HTTP and HTTPS servers simultaneously.
	// This can be done by setting the BindAddress and the SecureBindAddress simultaneously.
	// To use the secure server you must configure a TLS certificate and key.
	MetricsServer Server `json:"metricsServer,omitempty"`

	// Providers is used to configure multiple providers.
	// As of yet multiple providers aren't supported only the first entry is actually used.
	Providers Providers `yaml:"providers,omitempty"`

	// Cookie is used to configure the cookie used to store the session state.
	// This includes options such as the cookie name, its expiry and its domain.
	Cookie Cookie `json:"cookie,omitempty"`

	// Session is used to configure the session storage.
	// To either use a cookie or a redis store.
	Session SessionOptions `yaml:"session,omitempty"`

	// PageTemplates is used to configure custom page templates.
	// This includes the sign in and error pages.
	PageTemplates PageTemplates `yaml:"pageTemplates,omitempty"`
}

// Initialize alpha options with default values and settings of the core options
func NewAlphaOptions(opts *Options) *AlphaOptions {
	aOpts := &AlphaOptions{}
	aOpts.ExtractFrom(opts)
	return aOpts
}

// ExtractFrom populates the fields in the AlphaOptions with the values from
// the Options
func (a *AlphaOptions) ExtractFrom(opts *Options) {
	a.ProxyOptions = opts.ProxyOptions
	a.ProbeOptions = opts.ProbeOptions
	a.UpstreamConfig = opts.UpstreamServers
	a.InjectRequestHeaders = opts.InjectRequestHeaders
	a.InjectResponseHeaders = opts.InjectResponseHeaders
	a.Server = opts.Server
	a.MetricsServer = opts.MetricsServer
	a.Providers = opts.Providers
	a.Cookie = opts.Cookie
	a.Session = opts.Session
	a.PageTemplates = opts.PageTemplates
}

// MergeInto replaces alpha options in the Options struct with the values
// from the AlphaOptions
func (a *AlphaOptions) MergeInto(opts *Options) {
	opts.ProxyOptions = a.ProxyOptions
	opts.ProbeOptions = a.ProbeOptions
	opts.UpstreamServers = a.UpstreamConfig
	opts.InjectRequestHeaders = a.InjectRequestHeaders
	opts.InjectResponseHeaders = a.InjectResponseHeaders
	opts.Server = a.Server
	opts.MetricsServer = a.MetricsServer
	opts.Providers = a.Providers
	opts.Cookie = a.Cookie
	opts.Session = a.Session
	opts.PageTemplates = a.PageTemplates
}
