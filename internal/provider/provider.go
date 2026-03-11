package provider

// AuthType represents how a provider authenticates.
type AuthType string

const (
	AuthAPIKey AuthType = "api_key"
	AuthOAuth  AuthType = "oauth"
	AuthLocal  AuthType = "local"
)

// Provider describes an LLM provider.
type Provider struct {
	Name          string   // Canonical short name (e.g. "openai")
	DisplayName   string   // Human-readable name (e.g. "OpenAI")
	Aliases       []string // Alternative names users can type
	EnvVar        string   // Default environment variable (e.g. "OPENAI_API_KEY")
	AuthType      AuthType
	BaseURL       string // API base URL
	CheckEndpoint string // Relative or absolute URL for a lightweight health check
	APIKeyURL     string // URL to the page where the user can create an API key
	Tags          []string
}
