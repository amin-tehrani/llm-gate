package provider

import (
	"fmt"
	"sort"
	"strings"
)

// registry holds all known providers keyed by canonical name.
var registry []Provider

func init() {
	registry = []Provider{
		{
			Name: "openrouter", DisplayName: "OpenRouter",
			EnvVar: "OPENROUTER_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://openrouter.ai/api", CheckEndpoint: "/v1/models",
		},
		{
			Name: "anthropic", DisplayName: "Anthropic",
			EnvVar: "ANTHROPIC_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.anthropic.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "openai", DisplayName: "OpenAI",
			EnvVar: "OPENAI_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.openai.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "openai-codex", DisplayName: "OpenAI Codex (OAuth)",
			Aliases: []string{"openai_codex", "codex"},
			EnvVar:  "OPENAI_API_KEY", AuthType: AuthOAuth,
			BaseURL: "https://api.openai.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "ollama", DisplayName: "Ollama",
			EnvVar: "OLLAMA_HOST", AuthType: AuthLocal,
			BaseURL: "http://localhost:11434", CheckEndpoint: "/api/tags",
			Tags: []string{"local"},
		},
		{
			Name: "gemini", DisplayName: "Google Gemini",
			Aliases: []string{"google", "google-gemini"},
			EnvVar:  "GEMINI_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://generativelanguage.googleapis.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "venice", DisplayName: "Venice",
			EnvVar: "VENICE_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.venice.ai", CheckEndpoint: "/api/v1/models",
		},
		{
			Name: "vercel", DisplayName: "Vercel AI Gateway",
			Aliases: []string{"vercel-ai"},
			EnvVar:  "VERCEL_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.vercel.ai", CheckEndpoint: "/v1/models",
		},
		{
			Name: "cloudflare", DisplayName: "Cloudflare AI",
			Aliases: []string{"cloudflare-ai"},
			EnvVar:  "CLOUDFLARE_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.cloudflare.com", CheckEndpoint: "/client/v4/user/tokens/verify",
		},
		{
			Name: "moonshot", DisplayName: "Moonshot",
			Aliases: []string{"kimi"},
			EnvVar:  "MOONSHOT_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.moonshot.cn", CheckEndpoint: "/v1/models",
		},
		{
			Name: "kimi-code", DisplayName: "Kimi Code",
			Aliases: []string{"kimi_coding", "kimi_for_coding"},
			EnvVar:  "KIMI_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.moonshot.cn", CheckEndpoint: "/v1/models",
		},
		{
			Name: "synthetic", DisplayName: "Synthetic",
			EnvVar: "SYNTHETIC_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.synthetic.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "opencode", DisplayName: "OpenCode Zen",
			Aliases: []string{"opencode-zen"},
			EnvVar:  "OPENCODE_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.opencode.ai", CheckEndpoint: "/v1/models",
		},
		{
			Name: "zai", DisplayName: "Z.AI",
			Aliases: []string{"z.ai"},
			EnvVar:  "ZAI_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.z.ai", CheckEndpoint: "/v1/models",
		},
		{
			Name: "glm", DisplayName: "GLM (Zhipu)",
			Aliases: []string{"zhipu"},
			EnvVar:  "GLM_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://open.bigmodel.cn", CheckEndpoint: "/api/paas/v4/models",
		},
		{
			Name: "minimax", DisplayName: "MiniMax",
			Aliases: []string{"minimax-intl", "minimax-io", "minimax-global", "minimax-cn", "minimaxi", "minimax-oauth", "minimax-oauth-cn", "minimax-portal", "minimax-portal-cn"},
			EnvVar:  "MINIMAX_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.minimax.chat", CheckEndpoint: "/v1/models",
		},
		{
			Name: "bedrock", DisplayName: "Amazon Bedrock",
			Aliases: []string{"aws-bedrock"},
			EnvVar:  "AWS_ACCESS_KEY_ID", AuthType: AuthAPIKey,
			BaseURL: "https://bedrock.us-east-1.amazonaws.com", CheckEndpoint: "/",
		},
		{
			Name: "qianfan", DisplayName: "Qianfan (Baidu)",
			Aliases: []string{"baidu"},
			EnvVar:  "QIANFAN_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://aip.baidubce.com", CheckEndpoint: "/oauth/2.0/token",
		},
		{
			Name: "doubao", DisplayName: "Doubao (Volcengine)",
			Aliases: []string{"volcengine", "ark", "doubao-cn"},
			EnvVar:  "DOUBAO_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://ark.cn-beijing.volces.com", CheckEndpoint: "/api/v3/models",
		},
		{
			Name: "qwen", DisplayName: "Qwen (DashScope)",
			Aliases: []string{"dashscope", "qwen-intl", "dashscope-intl", "qwen-us", "dashscope-us", "qwen-code", "qwen-oauth", "qwen_oauth"},
			EnvVar:  "DASHSCOPE_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://dashscope.aliyuncs.com", CheckEndpoint: "/api/v1/models",
		},
		{
			Name: "groq", DisplayName: "Groq",
			EnvVar: "GROQ_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.groq.com", CheckEndpoint: "/openai/v1/models",
		},
		{
			Name: "mistral", DisplayName: "Mistral",
			EnvVar: "MISTRAL_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.mistral.ai", CheckEndpoint: "/v1/models",
		},
		{
			Name: "xai", DisplayName: "xAI (Grok)",
			Aliases: []string{"grok"},
			EnvVar:  "XAI_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.x.ai", CheckEndpoint: "/v1/models",
		},
		{
			Name: "deepseek", DisplayName: "DeepSeek",
			EnvVar: "DEEPSEEK_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.deepseek.com", CheckEndpoint: "/v1/models",
		},
		{
			Name: "together", DisplayName: "Together AI",
			Aliases: []string{"together-ai"},
			EnvVar:  "TOGETHER_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.together.xyz", CheckEndpoint: "/v1/models",
		},
		{
			Name: "fireworks", DisplayName: "Fireworks AI",
			Aliases: []string{"fireworks-ai"},
			EnvVar:  "FIREWORKS_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.fireworks.ai", CheckEndpoint: "/inference/v1/models",
		},
		{
			Name: "novita", DisplayName: "Novita AI",
			EnvVar: "NOVITA_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.novita.ai", CheckEndpoint: "/v3/openai/models",
		},
		{
			Name: "perplexity", DisplayName: "Perplexity",
			EnvVar: "PERPLEXITY_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.perplexity.ai", CheckEndpoint: "/models",
		},
		{
			Name: "cohere", DisplayName: "Cohere",
			EnvVar: "COHERE_API_KEY", AuthType: AuthAPIKey,
			BaseURL: "https://api.cohere.com", CheckEndpoint: "/v2/models",
		},
		{
			Name: "copilot", DisplayName: "GitHub Copilot",
			Aliases: []string{"github-copilot"},
			EnvVar:  "GITHUB_TOKEN", AuthType: AuthOAuth,
			BaseURL: "https://api.github.com", CheckEndpoint: "/user",
		},
	}
}

// aliasIndex is built lazily for fast alias lookups.
var aliasIndex map[string]*Provider

func buildAliasIndex() {
	aliasIndex = make(map[string]*Provider, len(registry)*3)
	for i := range registry {
		p := &registry[i]
		aliasIndex[p.Name] = p
		for _, a := range p.Aliases {
			aliasIndex[strings.ToLower(a)] = p
		}
	}
}

// Lookup finds a provider by canonical name or alias. Returns nil if not found.
func Lookup(nameOrAlias string) *Provider {
	if aliasIndex == nil {
		buildAliasIndex()
	}
	return aliasIndex[strings.ToLower(nameOrAlias)]
}

// All returns all providers sorted by canonical name.
func All() []Provider {
	sorted := make([]Provider, len(registry))
	copy(sorted, registry)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

// Names returns all canonical provider names sorted.
func Names() []string {
	providers := All()
	names := make([]string, len(providers))
	for i, p := range providers {
		names[i] = p.Name
	}
	return names
}

// MustLookup is like Lookup but returns an error if the provider is not found.
func MustLookup(nameOrAlias string) (*Provider, error) {
	p := Lookup(nameOrAlias)
	if p == nil {
		return nil, fmt.Errorf("unknown provider: %s", nameOrAlias)
	}
	return p, nil
}
