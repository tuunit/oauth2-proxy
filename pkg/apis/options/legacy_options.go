package options

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/spf13/pflag"
)

type LegacyOptions struct {
	// Legacy options related to upstream servers
	LegacyUpstreams LegacyUpstreams `cfg:",squash"`

	// Legacy options for injecting request/response headers
	LegacyHeaders LegacyHeaders `cfg:",squash"`

	// Legacy options for the server address and TLS
	LegacyServer LegacyServer `cfg:",squash"`

	// Legacy options for single provider
	LegacyProvider LegacyProvider `cfg:",squash"`

	LegacyCookie LegacyCookie `cfg:",squash"`

	LegacySessionOptions LegacySessionOptions `cfg:",squash"`

	LegacyLogging LegacyLogging `cfg:",squash"`

	LegacyPageTemplates LegacyPageTemplates `cfg:",squash"`

	LegacyProxyOptions LegacyProxyOptions `cfg:",squash"`

	LegacyProbeOptions LegacyProbeOptions `cfg:",squash"`
}

func NewLegacyOptions() *LegacyOptions {
	return &LegacyOptions{
		LegacyProxyOptions: LegacyProxyOptions{
			ProxyPrefix:        "/oauth2",
			RealClientIPHeader: "X-Real-IP",
			ForceHTTPS:         false,
			SkipAuthPreflight:  false,
		},

		LegacyUpstreams: LegacyUpstreams{
			PassHostHeader:  true,
			ProxyWebSockets: true,
			FlushInterval:   DefaultUpstreamFlushInterval,
			Timeout:         DefaultUpstreamTimeout,
		},

		LegacyHeaders: LegacyHeaders{
			PassBasicAuth:        true,
			PassUserHeaders:      true,
			SkipAuthStripHeaders: true,
		},

		LegacyServer: LegacyServer{
			HTTPAddress:  "127.0.0.1:4180",
			HTTPSAddress: ":443",
		},

		LegacyProvider: LegacyProvider{
			ProviderType:          "google",
			AzureTenant:           "common",
			ApprovalPrompt:        "force",
			UserIDClaim:           "email",
			OIDCEmailClaim:        "email",
			OIDCGroupsClaim:       "groups",
			OIDCAudienceClaims:    []string{"aud"},
			OIDCExtraAudiences:    []string{},
			InsecureOIDCSkipNonce: true,
		},

		LegacyCookie: LegacyCookie{
			Name:           "_oauth2_proxy",
			Secret:         "",
			Domains:        nil,
			Path:           "/",
			Expire:         time.Duration(168) * time.Hour,
			Refresh:        time.Duration(0),
			Secure:         true,
			HTTPOnly:       true,
			SameSite:       "",
			CSRFPerRequest: false,
			CSRFExpire:     time.Duration(15) * time.Minute,
		},

		LegacySessionOptions: LegacySessionOptions{
			Type: "cookie",
			Cookie: LegacyCookieStoreOptions{
				Minimal: false,
			},
		},

		LegacyPageTemplates: LegacyPageTemplates{
			Path:             "",
			CustomLogo:       "",
			Banner:           "",
			Footer:           "",
			DisplayLoginForm: true,
			Debug:            false,
		},

		LegacyProbeOptions: LegacyProbeOptions{
			PingPath:        "/ping",
			PingUserAgent:   "",
			ReadyPath:       "/ready",
			GCPHealthChecks: false,
		},

		LegacyLogging: LegacyLogging{
			ExcludePaths:    nil,
			LocalTime:       true,
			SilencePing:     false,
			RequestIDHeader: "X-Request-Id",
			AuthEnabled:     true,
			AuthFormat:      logger.DefaultAuthLoggingFormat,
			RequestEnabled:  true,
			RequestFormat:   logger.DefaultRequestLoggingFormat,
			StandardEnabled: true,
			StandardFormat:  logger.DefaultStandardLoggingFormat,
			ErrToInfo:       false,
			File: LegacyLogFileOptions{
				Filename:   "",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 0,
				Compress:   false,
			},
		},
	}
}

func NewLegacyFlagSet() *pflag.FlagSet {
	flagSet := NewFlagSet()

	flagSet.AddFlagSet(legacyUpstreamsFlagSet())
	flagSet.AddFlagSet(legacyHeadersFlagSet())
	flagSet.AddFlagSet(legacyServerFlagset())
	flagSet.AddFlagSet(legacyProxyOptionsFlagSet())
	flagSet.AddFlagSet(legacyCookieFlagSet())
	flagSet.AddFlagSet(legacySessionFlagSet())
	flagSet.AddFlagSet(legacyLoggingFlagSet())
	flagSet.AddFlagSet(legacyPageTemplatesFlagSet())
	flagSet.AddFlagSet(legacyProbeOptionsFlagSet())
	flagSet.AddFlagSet(legacyProviderFlagSet())
	flagSet.AddFlagSet(legacyGoogleFlagSet())

	return flagSet
}

func (l *LegacyOptions) ToOptions() (*Options, error) {
	opts := NewOptions()

	opts.ProxyOptions = l.LegacyProxyOptions.convert()

	opts.ProxyOptions.LegacyPreferEmailToUser = l.LegacyHeaders.PreferEmailToUser

	opts.ProbeOptions = l.LegacyProbeOptions.convert()

	upstreams, err := l.LegacyUpstreams.convert()
	if err != nil {
		return nil, fmt.Errorf("error converting upstreams: %v", err)
	}
	opts.UpstreamServers = upstreams

	opts.InjectRequestHeaders, opts.InjectResponseHeaders = l.LegacyHeaders.convert()

	opts.Server, opts.MetricsServer = l.LegacyServer.convert()

	opts.Cookie = l.LegacyCookie.convert()

	opts.Session = l.LegacySessionOptions.convert()

	opts.Logging = l.LegacyLogging.convert()

	opts.PageTemplates = l.LegacyPageTemplates.convert()

	providers, err := l.LegacyProvider.convert()
	if err != nil {
		return nil, fmt.Errorf("error converting provider: %v", err)
	}
	opts.Providers = providers

	return opts, nil
}

type LegacyUpstreams struct {
	FlushInterval                 time.Duration `flag:"flush-interval" cfg:"flush_interval"`
	PassHostHeader                bool          `flag:"pass-host-header" cfg:"pass_host_header"`
	ProxyWebSockets               bool          `flag:"proxy-websockets" cfg:"proxy_websockets"`
	SSLUpstreamInsecureSkipVerify bool          `flag:"ssl-upstream-insecure-skip-verify" cfg:"ssl_upstream_insecure_skip_verify"`
	Upstreams                     []string      `flag:"upstream" cfg:"upstreams"`
	Timeout                       time.Duration `flag:"upstream-timeout" cfg:"upstream_timeout"`
}

func legacyUpstreamsFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("upstreams", pflag.ExitOnError)

	flagSet.Duration("flush-interval", DefaultUpstreamFlushInterval, "period between response flushing when streaming responses")
	flagSet.Bool("pass-host-header", true, "pass the request Host Header to upstream")
	flagSet.Bool("proxy-websockets", true, "enables WebSocket proxying")
	flagSet.Bool("ssl-upstream-insecure-skip-verify", false, "skip validation of certificates presented when using HTTPS upstreams")
	flagSet.StringSlice("upstream", []string{}, "the http url(s) of the upstream endpoint, file:// paths for static files or static://<status_code> for static response. Routing is based on the path")
	flagSet.Duration("upstream-timeout", DefaultUpstreamTimeout, "maximum amount of time the server will wait for a response from the upstream")

	return flagSet
}

func (l *LegacyUpstreams) convert() (UpstreamConfig, error) {
	upstreams := UpstreamConfig{}

	for _, upstreamString := range l.Upstreams {
		u, err := url.Parse(upstreamString)
		if err != nil {
			return UpstreamConfig{}, fmt.Errorf("could not parse upstream %q: %v", upstreamString, err)
		}

		if u.Path == "" {
			u.Path = "/"
		}

		flushInterval := Duration(l.FlushInterval)
		timeout := Duration(l.Timeout)
		upstream := Upstream{
			ID:                    u.Path,
			Path:                  u.Path,
			URI:                   upstreamString,
			InsecureSkipTLSVerify: l.SSLUpstreamInsecureSkipVerify,
			PassHostHeader:        &l.PassHostHeader,
			ProxyWebSockets:       &l.ProxyWebSockets,
			FlushInterval:         &flushInterval,
			Timeout:               &timeout,
		}

		switch u.Scheme {
		case "file":
			if u.Fragment != "" {
				upstream.ID = u.Fragment
				upstream.Path = u.Fragment
				// Trim the fragment from the end of the URI
				upstream.URI = strings.SplitN(upstreamString, "#", 2)[0]
			}
		case "static":
			responseCode, err := strconv.Atoi(u.Host)
			if err != nil {
				logger.Errorf("unable to convert %q to int, use default \"200\"", u.Host)
				responseCode = 200
			}
			upstream.Static = true
			upstream.StaticCode = &responseCode

			// This is not allowed to be empty and must be unique
			upstream.ID = upstreamString

			// We only support the root path in the legacy config
			upstream.Path = "/"

			// Force defaults compatible with static responses
			upstream.URI = ""
			upstream.InsecureSkipTLSVerify = false
			upstream.PassHostHeader = nil
			upstream.ProxyWebSockets = nil
			upstream.FlushInterval = nil
			upstream.Timeout = nil
		case "unix":
			upstream.Path = "/"
		}

		upstreams.Upstreams = append(upstreams.Upstreams, upstream)
	}

	return upstreams, nil
}

type LegacyProvider struct {
	ClientID         string `flag:"client-id" cfg:"client_id"`
	ClientSecret     string `flag:"client-secret" cfg:"client_secret"`
	ClientSecretFile string `flag:"client-secret-file" cfg:"client_secret_file"`

	KeycloakGroups                         []string `flag:"keycloak-group" cfg:"keycloak_groups"`
	AzureTenant                            string   `flag:"azure-tenant" cfg:"azure_tenant"`
	AzureGraphGroupField                   string   `flag:"azure-graph-group-field" cfg:"azure_graph_group_field"`
	BitbucketTeam                          string   `flag:"bitbucket-team" cfg:"bitbucket_team"`
	BitbucketRepository                    string   `flag:"bitbucket-repository" cfg:"bitbucket_repository"`
	GitHubOrg                              string   `flag:"github-org" cfg:"github_org"`
	GitHubTeam                             string   `flag:"github-team" cfg:"github_team"`
	GitHubRepo                             string   `flag:"github-repo" cfg:"github_repo"`
	GitHubToken                            string   `flag:"github-token" cfg:"github_token"`
	GitHubUsers                            []string `flag:"github-user" cfg:"github_users"`
	GitLabGroup                            []string `flag:"gitlab-group" cfg:"gitlab_groups"`
	GitLabProjects                         []string `flag:"gitlab-project" cfg:"gitlab_projects"`
	GoogleGroupsLegacy                     []string `flag:"google-group" cfg:"google_group"`
	GoogleGroups                           []string `flag:"google-group" cfg:"google_groups"`
	GoogleAdminEmail                       string   `flag:"google-admin-email" cfg:"google_admin_email"`
	GoogleServiceAccountJSON               string   `flag:"google-service-account-json" cfg:"google_service_account_json"`
	GoogleUseApplicationDefaultCredentials bool     `flag:"google-use-application-default-credentials" cfg:"google_use_application_default_credentials"`
	GoogleTargetPrincipal                  string   `flag:"google-target-principal" cfg:"google_target_principal"`

	// These options allow for other providers besides Google, with
	// potential overrides.
	ProviderType                       string   `flag:"provider" cfg:"provider"`
	ProviderName                       string   `flag:"provider-display-name" cfg:"provider_display_name"`
	ProviderCAFiles                    []string `flag:"provider-ca-file" cfg:"provider_ca_files"`
	UseSystemTrustStore                bool     `flag:"use-system-trust-store" cfg:"use_system_trust_store"`
	OIDCIssuerURL                      string   `flag:"oidc-issuer-url" cfg:"oidc_issuer_url"`
	InsecureOIDCAllowUnverifiedEmail   bool     `flag:"insecure-oidc-allow-unverified-email" cfg:"insecure_oidc_allow_unverified_email"`
	InsecureOIDCSkipIssuerVerification bool     `flag:"insecure-oidc-skip-issuer-verification" cfg:"insecure_oidc_skip_issuer_verification"`
	InsecureOIDCSkipNonce              bool     `flag:"insecure-oidc-skip-nonce" cfg:"insecure_oidc_skip_nonce"`
	SkipOIDCDiscovery                  bool     `flag:"skip-oidc-discovery" cfg:"skip_oidc_discovery"`
	OIDCJwksURL                        string   `flag:"oidc-jwks-url" cfg:"oidc_jwks_url"`
	OIDCEmailClaim                     string   `flag:"oidc-email-claim" cfg:"oidc_email_claim"`
	OIDCGroupsClaim                    string   `flag:"oidc-groups-claim" cfg:"oidc_groups_claim"`
	OIDCAudienceClaims                 []string `flag:"oidc-audience-claim" cfg:"oidc_audience_claims"`
	OIDCExtraAudiences                 []string `flag:"oidc-extra-audience" cfg:"oidc_extra_audiences"`
	LoginURL                           string   `flag:"login-url" cfg:"login_url"`
	RedeemURL                          string   `flag:"redeem-url" cfg:"redeem_url"`
	ProfileURL                         string   `flag:"profile-url" cfg:"profile_url"`
	SkipClaimsFromProfileURL           bool     `flag:"skip-claims-from-profile-url" cfg:"skip_claims_from_profile_url"`
	ProtectedResource                  string   `flag:"resource" cfg:"resource"`
	ValidateURL                        string   `flag:"validate-url" cfg:"validate_url"`
	Scope                              string   `flag:"scope" cfg:"scope"`
	Prompt                             string   `flag:"prompt" cfg:"prompt"`
	ApprovalPrompt                     string   `flag:"approval-prompt" cfg:"approval_prompt"` // Deprecated by OIDC 1.0
	UserIDClaim                        string   `flag:"user-id-claim" cfg:"user_id_claim"`
	AllowedGroups                      []string `flag:"allowed-group" cfg:"allowed_groups"`
	AllowedRoles                       []string `flag:"allowed-role" cfg:"allowed_roles"`
	BackendLogoutURL                   string   `flag:"backend-logout-url" cfg:"backend_logout_url"`

	AcrValues  string `flag:"acr-values" cfg:"acr_values"`
	JWTKey     string `flag:"jwt-key" cfg:"jwt_key"`
	JWTKeyFile string `flag:"jwt-key-file" cfg:"jwt_key_file"`
	PubJWKURL  string `flag:"pubjwk-url" cfg:"pubjwk_url"`
	// PKCE Code Challenge method to use (either S256 or plain)
	CodeChallengeMethod string `flag:"code-challenge-method" cfg:"code_challenge_method"`
	// Provided for legacy reasons, to be dropped in newer version see #1667
	ForceCodeChallengeMethod string `flag:"force-code-challenge-method" cfg:"force_code_challenge_method"`
}

func legacyProviderFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("provider", pflag.ExitOnError)

	flagSet.StringSlice("keycloak-group", []string{}, "restrict logins to members of these groups (may be given multiple times)")
	flagSet.String("azure-tenant", "common", "go to a tenant-specific or common (tenant-independent) endpoint.")
	flagSet.String("azure-graph-group-field", "", "configures the group field to be used when building the groups list(`id` or `displayName`. Default is `id`) from Microsoft Graph(available only for v2.0 oidc url). Based on this value, the `allowed-group` config values should be adjusted accordingly. If using `id` as group field, `allowed-group` should contains groups IDs, if using `displayName` as group field, `allowed-group` should contains groups name")
	flagSet.String("bitbucket-team", "", "restrict logins to members of this team")
	flagSet.String("bitbucket-repository", "", "restrict logins to user with access to this repository")
	flagSet.String("github-org", "", "restrict logins to members of this organisation")
	flagSet.String("github-team", "", "restrict logins to members of this team")
	flagSet.String("github-repo", "", "restrict logins to collaborators of this repository")
	flagSet.String("github-token", "", "the token to use when verifying repository collaborators (must have push access to the repository)")
	flagSet.StringSlice("github-user", []string{}, "allow users with these usernames to login even if they do not belong to the specified org and team or collaborators (may be given multiple times)")
	flagSet.StringSlice("gitlab-group", []string{}, "restrict logins to members of this group (may be given multiple times)")
	flagSet.StringSlice("gitlab-project", []string{}, "restrict logins to members of this project (may be given multiple times) (eg `group/project=accesslevel`). Access level should be a value matching Gitlab access levels (see https://docs.gitlab.com/ee/api/members.html#valid-access-levels), defaulted to 20 if absent")
	flagSet.String("client-id", "", "the OAuth Client ID: ie: \"123456.apps.googleusercontent.com\"")
	flagSet.String("client-secret", "", "the OAuth Client Secret")
	flagSet.String("client-secret-file", "", "the file with OAuth Client Secret")

	flagSet.String("provider", "google", "OAuth provider")
	flagSet.String("provider-display-name", "", "Provider display name")
	flagSet.StringSlice("provider-ca-file", []string{}, "One or more paths to CA certificates that should be used when connecting to the provider.  If not specified, the default Go trust sources are used instead.")
	flagSet.Bool("use-system-trust-store", false, "Determines if 'provider-ca-file' files and the system trust store are used. If set to true, your custom CA files and the system trust store are used otherwise only your custom CA files.")
	flagSet.String("oidc-issuer-url", "", "OpenID Connect issuer URL (ie: https://accounts.google.com)")
	flagSet.Bool("insecure-oidc-allow-unverified-email", false, "Don't fail if an email address in an id_token is not verified")
	flagSet.Bool("insecure-oidc-skip-issuer-verification", false, "Do not verify if issuer matches OIDC discovery URL")
	flagSet.Bool("insecure-oidc-skip-nonce", true, "skip verifying the OIDC ID Token's nonce claim")
	flagSet.Bool("skip-oidc-discovery", false, "Skip OIDC discovery and use manually supplied Endpoints")
	flagSet.String("oidc-jwks-url", "", "OpenID Connect JWKS URL (ie: https://www.googleapis.com/oauth2/v3/certs)")
	flagSet.String("oidc-groups-claim", OIDCGroupsClaim, "which OIDC claim contains the user groups")
	flagSet.String("oidc-email-claim", OIDCEmailClaim, "which OIDC claim contains the user's email")
	flagSet.StringSlice("oidc-audience-claim", OIDCAudienceClaims, "which OIDC claims are used as audience to verify against client id")
	flagSet.StringSlice("oidc-extra-audience", []string{}, "additional audiences allowed to pass audience verification")
	flagSet.String("login-url", "", "Authentication endpoint")
	flagSet.String("redeem-url", "", "Token redemption endpoint")
	flagSet.String("profile-url", "", "Profile access endpoint")
	flagSet.Bool("skip-claims-from-profile-url", false, "Skip loading missing claims from profile URL")
	flagSet.String("resource", "", "The resource that is protected (Azure AD only)")
	flagSet.String("validate-url", "", "Access token validation endpoint")
	flagSet.String("scope", "", "OAuth scope specification")
	flagSet.String("prompt", "", "OIDC prompt")
	flagSet.String("approval-prompt", "force", "OAuth approval_prompt")
	flagSet.String("code-challenge-method", "", "use PKCE code challenges with the specified method. Either 'plain' or 'S256'")
	flagSet.String("force-code-challenge-method", "", "Deprecated - use --code-challenge-method")

	flagSet.String("acr-values", "", "acr values string:  optional")
	flagSet.String("jwt-key", "", "private key in PEM format used to sign JWT, so that you can say something like -jwt-key=\"${OAUTH2_PROXY_JWT_KEY}\": required by login.gov")
	flagSet.String("jwt-key-file", "", "path to the private key file in PEM format used to sign the JWT so that you can say something like -jwt-key-file=/etc/ssl/private/jwt_signing_key.pem: required by login.gov")
	flagSet.String("pubjwk-url", "", "JWK pubkey access endpoint: required by login.gov")

	flagSet.String("user-id-claim", OIDCEmailClaim, "(DEPRECATED for `oidc-email-claim`) which claim contains the user ID")
	flagSet.StringSlice("allowed-group", []string{}, "restrict logins to members of this group (may be given multiple times)")
	flagSet.StringSlice("allowed-role", []string{}, "(keycloak-oidc) restrict logins to members of these roles (may be given multiple times)")
	flagSet.String("backend-logout-url", "", "url to perform a backend logout, {id_token} can be used as placeholder for the id_token")

	return flagSet
}

func legacyGoogleFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("google", pflag.ExitOnError)

	flagSet.StringSlice("google-group", []string{}, "restrict logins to members of this google group (may be given multiple times).")
	flagSet.String("google-admin-email", "", "the google admin to impersonate for api calls")
	flagSet.String("google-service-account-json", "", "the path to the service account json credentials")
	flagSet.String("google-use-application-default-credentials", "", "use application default credentials instead of service account json (i.e. GKE Workload Identity)")
	flagSet.String("google-target-principal", "", "the target principal to impersonate when using ADC")

	return flagSet
}

func (l *LegacyProvider) convert() (Providers, error) {
	providers := Providers{}

	provider := Provider{
		ClientID:                 l.ClientID,
		ClientSecret:             l.ClientSecret,
		ClientSecretFile:         l.ClientSecretFile,
		Type:                     ProviderType(l.ProviderType),
		CAFiles:                  l.ProviderCAFiles,
		UseSystemTrustStore:      l.UseSystemTrustStore,
		LoginURL:                 l.LoginURL,
		RedeemURL:                l.RedeemURL,
		ProfileURL:               l.ProfileURL,
		SkipClaimsFromProfileURL: l.SkipClaimsFromProfileURL,
		ProtectedResource:        l.ProtectedResource,
		ValidateURL:              l.ValidateURL,
		Scope:                    l.Scope,
		AllowedGroups:            l.AllowedGroups,
		CodeChallengeMethod:      l.CodeChallengeMethod,
		BackendLogoutURL:         l.BackendLogoutURL,
	}

	// This part is out of the switch section for all providers that support OIDC
	provider.OIDCConfig = OIDCOptions{
		IssuerURL:                      l.OIDCIssuerURL,
		InsecureAllowUnverifiedEmail:   l.InsecureOIDCAllowUnverifiedEmail,
		InsecureSkipIssuerVerification: l.InsecureOIDCSkipIssuerVerification,
		InsecureSkipNonce:              l.InsecureOIDCSkipNonce,
		SkipDiscovery:                  l.SkipOIDCDiscovery,
		JwksURL:                        l.OIDCJwksURL,
		UserIDClaim:                    l.UserIDClaim,
		EmailClaim:                     l.OIDCEmailClaim,
		GroupsClaim:                    l.OIDCGroupsClaim,
		AudienceClaims:                 l.OIDCAudienceClaims,
		ExtraAudiences:                 l.OIDCExtraAudiences,
	}

	// Support for legacy configuration option
	if l.ForceCodeChallengeMethod != "" && l.CodeChallengeMethod == "" {
		provider.CodeChallengeMethod = l.ForceCodeChallengeMethod
	}

	// This part is out of the switch section because azure has a default tenant
	// that needs to be added from legacy options
	provider.AzureConfig = AzureOptions{
		Tenant:          l.AzureTenant,
		GraphGroupField: l.AzureGraphGroupField,
	}

	switch provider.Type {
	case "github":
		provider.GitHubConfig = GitHubOptions{
			Org:   l.GitHubOrg,
			Team:  l.GitHubTeam,
			Repo:  l.GitHubRepo,
			Token: l.GitHubToken,
			Users: l.GitHubUsers,
		}
	case "keycloak-oidc":
		provider.KeycloakConfig = KeycloakOptions{
			Groups: l.KeycloakGroups,
			Roles:  l.AllowedRoles,
		}
	case "keycloak":
		provider.KeycloakConfig = KeycloakOptions{
			Groups: l.KeycloakGroups,
		}
	case "gitlab":
		provider.GitLabConfig = GitLabOptions{
			Group:    l.GitLabGroup,
			Projects: l.GitLabProjects,
		}
	case "login.gov":
		provider.LoginGovConfig = LoginGovOptions{
			JWTKey:     l.JWTKey,
			JWTKeyFile: l.JWTKeyFile,
			PubJWKURL:  l.PubJWKURL,
		}
	case "bitbucket":
		provider.BitbucketConfig = BitbucketOptions{
			Team:       l.BitbucketTeam,
			Repository: l.BitbucketRepository,
		}
	case "google":
		if len(l.GoogleGroupsLegacy) != 0 && !reflect.DeepEqual(l.GoogleGroupsLegacy, l.GoogleGroups) {
			// Log the deprecation notice
			logger.Error(
				"WARNING: The 'OAUTH2_PROXY_GOOGLE_GROUP' environment variable is deprecated and will likely be removed in the next major release. Use 'OAUTH2_PROXY_GOOGLE_GROUPS' instead.",
			)
			l.GoogleGroups = l.GoogleGroupsLegacy
		}
		provider.GoogleConfig = GoogleOptions{
			Groups:                           l.GoogleGroups,
			AdminEmail:                       l.GoogleAdminEmail,
			ServiceAccountJSON:               l.GoogleServiceAccountJSON,
			UseApplicationDefaultCredentials: l.GoogleUseApplicationDefaultCredentials,
			TargetPrincipal:                  l.GoogleTargetPrincipal,
		}
	}

	if l.ProviderName != "" {
		provider.ID = l.ProviderName
		provider.Name = l.ProviderName
	} else {
		provider.ID = l.ProviderType + "=" + l.ClientID
	}

	// handle AcrValues, Prompt and ApprovalPrompt
	var urlParams []LoginURLParameter
	if l.AcrValues != "" {
		urlParams = append(urlParams, LoginURLParameter{Name: "acr_values", Default: []string{l.AcrValues}})
	}
	switch {
	case l.Prompt != "":
		urlParams = append(urlParams, LoginURLParameter{Name: "prompt", Default: []string{l.Prompt}})
	case l.ApprovalPrompt != "":
		urlParams = append(urlParams, LoginURLParameter{Name: "approval_prompt", Default: []string{l.ApprovalPrompt}})
	default:
		// match legacy behaviour by default - if neither prompt nor approval_prompt
		// specified, use approval_prompt=force
		urlParams = append(urlParams, LoginURLParameter{Name: "approval_prompt", Default: []string{"force"}})
	}

	provider.LoginURLParameters = urlParams

	providers = append(providers, provider)

	return providers, nil
}

type LegacyProbeOptions struct {
	PingPath        string `flag:"ping-path" cfg:"ping_path"`
	PingUserAgent   string `flag:"ping-user-agent" cfg:"ping_user_agent"`
	ReadyPath       string `flag:"ready-path" cfg:"ready_path"`
	GCPHealthChecks bool   `flag:"gcp-healthchecks" cfg:"gcp_healthchecks"`
}

func legacyProbeOptionsFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("probe", pflag.ExitOnError)

	flagSet.String("ping-path", "/ping", "the ping endpoint that can be used for basic health checks")
	flagSet.String("ping-user-agent", "", "special User-Agent that will be used for basic health checks")
	flagSet.String("ready-path", "/ready", "the ready endpoint that can be used for deep health checks")
	flagSet.Bool("gcp-healthchecks", false, "Enable GCP/GKE healthcheck endpoints")

	return flagSet
}

func (l *LegacyProbeOptions) convert() ProbeOptions {
	return ProbeOptions{
		PingPath:              l.PingPath,
		PingUserAgent:         l.PingUserAgent,
		ReadyPath:             l.ReadyPath,
		LegacyGCPHealthChecks: l.GCPHealthChecks,
	}
}
