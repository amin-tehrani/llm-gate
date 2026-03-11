package shell

import (
	"fmt"
	"os"
)

// ExportCommand returns a shell export statement.
func ExportCommand(envVar, value string) string {
	return fmt.Sprintf("export %s=%q", envVar, value)
}

// UnsetCommand returns a shell unset statement.
func UnsetCommand(envVar string) string {
	return fmt.Sprintf("unset %s", envVar)
}

// ShellInit returns a shell function that wraps the llm-gate binary
// so that activate/deactivate commands can modify the parent shell's env.
func ShellInit() string {
	binPath, err := os.Executable()
	if err != nil || binPath == "" {
		binPath = "command llm-gate"
	} else {
		binPath = fmt.Sprintf("%q", binPath)
	}

	return fmt.Sprintf(`# Add this to your .bashrc or .zshrc:
# eval "$(llm-gate shell-init)"

_llm_gate() {
    local cmd="${1:-}"
    if [ "$cmd" = "activate" ] || [ "$cmd" = "deactivate" ]; then
        eval "$(LLM_GATE_WRAPPER=1 %s "$@")"
        eval $(_llm_gate current)
    else
        %s "$@"
    fi
}
eval $(_llm_gate current)
alias llm-gate='_llm_gate'
`, binPath, binPath)
}
