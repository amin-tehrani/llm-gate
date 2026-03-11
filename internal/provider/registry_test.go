package provider

import "testing"

func TestLookupByName(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
	}{
		{"openai", "openai"},
		{"anthropic", "anthropic"},
		{"ollama", "ollama"},
		{"deepseek", "deepseek"},
	}
	for _, tt := range tests {
		p := Lookup(tt.input)
		if p == nil {
			t.Fatalf("Lookup(%q) returned nil", tt.input)
		}
		if p.Name != tt.wantName {
			t.Errorf("Lookup(%q).Name = %q, want %q", tt.input, p.Name, tt.wantName)
		}
	}
}

func TestLookupByAlias(t *testing.T) {
	tests := []struct {
		alias    string
		wantName string
	}{
		{"grok", "xai"},
		{"google", "gemini"},
		{"google-gemini", "gemini"},
		{"codex", "openai-codex"},
		{"aws-bedrock", "bedrock"},
		{"github-copilot", "copilot"},
		{"baidu", "qianfan"},
		{"dashscope", "qwen"},
		{"together-ai", "together"},
		{"kimi", "moonshot"},
	}
	for _, tt := range tests {
		p := Lookup(tt.alias)
		if p == nil {
			t.Fatalf("Lookup(%q) returned nil", tt.alias)
		}
		if p.Name != tt.wantName {
			t.Errorf("Lookup(%q).Name = %q, want %q", tt.alias, p.Name, tt.wantName)
		}
	}
}

func TestLookupCaseInsensitive(t *testing.T) {
	p := Lookup("OpenAI")
	if p == nil {
		t.Fatal("Lookup(\"OpenAI\") returned nil")
	}
	if p.Name != "openai" {
		t.Errorf("got %q, want \"openai\"", p.Name)
	}
}

func TestLookupUnknown(t *testing.T) {
	if p := Lookup("nonexistent"); p != nil {
		t.Errorf("Lookup(\"nonexistent\") = %v, want nil", p)
	}
}

func TestMustLookupError(t *testing.T) {
	_, err := MustLookup("nonexistent")
	if err == nil {
		t.Error("MustLookup(\"nonexistent\") returned nil error")
	}
}

func TestAllReturnsAllProviders(t *testing.T) {
	all := All()
	if len(all) != 30 {
		t.Errorf("All() returned %d providers, want 30", len(all))
	}
}

func TestAllIsSorted(t *testing.T) {
	all := All()
	for i := 1; i < len(all); i++ {
		if all[i].Name < all[i-1].Name {
			t.Errorf("All() not sorted: %q comes after %q", all[i].Name, all[i-1].Name)
		}
	}
}

func TestNamesMatchesAll(t *testing.T) {
	names := Names()
	all := All()
	if len(names) != len(all) {
		t.Errorf("Names() length %d != All() length %d", len(names), len(all))
	}
	for i, n := range names {
		if n != all[i].Name {
			t.Errorf("Names()[%d] = %q, All()[%d].Name = %q", i, n, i, all[i].Name)
		}
	}
}
