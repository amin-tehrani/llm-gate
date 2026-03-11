package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSaveRoundTrip(t *testing.T) {
	// Use a temp dir for config
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Version != 1 {
		t.Errorf("default version = %d, want 1", cfg.Version)
	}

	cfg.SetKey("openai", "sk-test-123")
	cfg.SetActive("openai", true)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Reload and verify
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load() after save error: %v", err)
	}
	if cfg2.GetKey("openai") != "sk-test-123" {
		t.Errorf("GetKey = %q, want %q", cfg2.GetKey("openai"), "sk-test-123")
	}
	if !cfg2.IsActive("openai") {
		t.Error("IsActive(openai) = false, want true")
	}
	if !cfg2.IsConfigured("openai") {
		t.Error("IsConfigured(openai) = false, want true")
	}
}

func TestRemoveKey(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, _ := Load()
	cfg.SetKey("anthropic", "sk-ant-test")
	cfg.RemoveKey("anthropic")

	if cfg.IsConfigured("anthropic") {
		t.Error("IsConfigured(anthropic) = true after RemoveKey")
	}
	if cfg.GetKey("anthropic") != "" {
		t.Error("GetKey(anthropic) not empty after RemoveKey")
	}
}

func TestFilePermissions(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg, _ := Load()
	cfg.SetKey("test", "key")
	cfg.Save()

	info, err := os.Stat(filepath.Join(tmp, "llm-gate", "config.yaml"))
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}
}

func TestConfigDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	if got := ConfigDir(); got != "/custom/config/llm-gate" {
		t.Errorf("ConfigDir() = %q, want %q", got, "/custom/config/llm-gate")
	}
}
