package shell

import "fmt"

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
	return `# Add this to your .bashrc or .zshrc:
# eval "$(llm-gate shell-init)"

llm-gate() {
    local cmd="${1:-}"
    if [[ "$cmd" == "activate" || "$cmd" == "deactivate" ]]; then
        eval "$(command llm-gate "$@")"
    else
        command llm-gate "$@"
    fi
}`
}
