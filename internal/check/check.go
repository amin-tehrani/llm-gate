package check

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/amin-tehrani/llm-gate/internal/provider"
)

// Result holds the outcome of a connectivity check.
type Result struct {
	Provider *provider.Provider
	OK       bool
	Latency  time.Duration
	Error    string
}

// Check tests connectivity to a provider's API.
func Check(p *provider.Provider, apiKey string) Result {
	if p.CheckEndpoint == "" {
		return Result{Provider: p, OK: false, Error: "no check endpoint defined"}
	}

	url := p.BaseURL + p.CheckEndpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Result{Provider: p, OK: false, Error: fmt.Sprintf("creating request: %v", err)}
	}

	// Set auth header based on provider type
	switch p.AuthType {
	case provider.AuthAPIKey:
		if p.Name == "anthropic" {
			req.Header.Set("x-api-key", apiKey)
			req.Header.Set("anthropic-version", "2023-06-01")
		} else if p.Name == "cloudflare" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		} else if p.Name == "gemini" {
			// Gemini uses query param
			q := req.URL.Query()
			q.Set("key", apiKey)
			req.URL.RawQuery = q.Encode()
		} else {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
	case provider.AuthOAuth:
		if p.Name == "copilot" {
			req.Header.Set("Authorization", "token "+apiKey)
		} else {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
	case provider.AuthLocal:
		// No auth needed for local providers
	}

	client := &http.Client{Timeout: 10 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return Result{Provider: p, OK: false, Latency: latency, Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return Result{Provider: p, OK: true, Latency: latency}
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	errMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
	if len(bodyBytes) > 0 {
		// Truncate the body if it's too long
		bodyStr := string(bodyBytes)
		// Clean up newlines for CLI output
		bodyStr = strings.ReplaceAll(bodyStr, "\n", " ")
		bodyStr = strings.ReplaceAll(bodyStr, "\r", "")
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}
		errMsg += fmt.Sprintf(": %s", bodyStr)
	}

	return Result{Provider: p, OK: false, Latency: latency, Error: errMsg}
}
