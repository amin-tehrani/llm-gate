package shell

import (
	"strings"
	"testing"
)

func TestExportCommand(t *testing.T) {
	got := ExportCommand("OPENAI_API_KEY", "sk-test")
	want := `export OPENAI_API_KEY="sk-test"`
	if got != want {
		t.Errorf("ExportCommand() = %q, want %q", got, want)
	}
}

func TestUnsetCommand(t *testing.T) {
	got := UnsetCommand("OPENAI_API_KEY")
	want := "unset OPENAI_API_KEY"
	if got != want {
		t.Errorf("UnsetCommand() = %q, want %q", got, want)
	}
}

func TestShellInitContainsFunction(t *testing.T) {
	s := ShellInit()
	if !strings.Contains(s, "llm-gate()") {
		t.Error("ShellInit() does not contain function definition")
	}
	if !strings.Contains(s, "activate") {
		t.Error("ShellInit() does not mention activate")
	}
	if !strings.Contains(s, "eval") {
		t.Error("ShellInit() does not mention eval")
	}
}
