package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"https://github.com/amin-tehrani/llm-gate/internal/browser"
	"https://github.com/amin-tehrani/llm-gate/internal/check"
	"https://github.com/amin-tehrani/llm-gate/internal/config"
	"https://github.com/amin-tehrani/llm-gate/internal/provider"
	"https://github.com/amin-tehrani/llm-gate/internal/shell"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ── Root ──────────────────────────────────────────────────────────────────────

var rootCmd = &cobra.Command{
	Use:     "llm-gate",
	Short:   "Central hub for managing LLM provider API keys",
	Long:    "llm-gate is a CLI tool to store, activate, deactivate, and check connectivity for multiple LLM provider API keys.",
	Version: version,
}

// ── auth ──────────────────────────────────────────────────────────────────────

var authCmd = &cobra.Command{
	Use:               "auth <provider>",
	Short:             "Authenticate with an LLM provider (interactive key prompt)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}

		if p.AuthType == provider.AuthLocal {
			fmt.Printf("  %s is a local provider — no API key needed.\n", p.DisplayName)
			fmt.Printf("  Use 'llm-gate check %s' to verify it's running.\n", p.Name)
			return nil
		}

		bold := color.New(color.Bold)
		bold.Printf("  Authenticating with %s\n", p.DisplayName)

		if p.APIKeyURL != "" {
			dim := color.New(color.FgHiBlack)
			fmt.Printf("  To get an API key, visit: ")
			color.New(color.FgCyan, color.Underline).Printf("%s\n", p.APIKeyURL)
			dim.Printf("  Press Enter to open this URL in your browser, or type 'n' to skip: ")

			var b [1]byte
			os.Stdin.Read(b[:])
			if b[0] == '\n' || b[0] == '\r' || b[0] == 'y' || b[0] == 'Y' {
				if err := browser.Open(p.APIKeyURL); err != nil {
					dim.Printf("  Could not open browser: %v\n", err)
				}
			} else {
				// Consume the rest of the line if they typed 'n' + Enter
				if b[0] != '\n' && b[0] != '\r' {
					var discard [1024]byte
					os.Stdin.Read(discard[:])
				}
			}
			fmt.Println()
		}

		fmt.Printf("  Enter your API key: ")

		reader := bufio.NewReader(os.Stdin)
		keyStr, err := reader.ReadString('\n')
		fmt.Println()
		if err != nil {
			return fmt.Errorf("reading key: %w", err)
		}
		key := strings.TrimSpace(keyStr)
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.SetKey(p.Name, key)
		if err := cfg.Save(); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		green.Printf("  ✓ API key stored for %s\n", p.DisplayName)

		// Quick connectivity check
		result := check.Check(p, key)
		if result.OK {
			green.Printf("  ✓ Connectivity verified (%dms)\n", result.Latency.Milliseconds())
		} else {
			yellow := color.New(color.FgYellow)
			yellow.Printf("  ⚠ Could not verify connectivity: %s\n", result.Error)
		}

		fmt.Printf("\n  To activate: llm-gate activate %s\n", p.Name)
		return nil
	},
}

// ── set / update ──────────────────────────────────────────────────────────────

var setCmd = &cobra.Command{
	Use:               "set <provider> <api-key>",
	Short:             "Store an API key for a provider",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.SetKey(p.Name, args[1])
		if err := cfg.Save(); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		green.Printf("  ✓ API key stored for %s\n", p.DisplayName)
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:               "update <provider> <api-key>",
	Short:             "Update the API key for a provider (alias for set)",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: providerCompletion,
	RunE:              setCmd.RunE,
}

// ── activate ──────────────────────────────────────────────────────────────────

var activateCmd = &cobra.Command{
	Use:               "activate <provider>",
	Short:             "Export the API key as an environment variable",
	Long:              "Prints an export statement. Use with: eval \"$(llm-gate activate <provider>)\"",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("LLM_GATE_WRAPPER") != "1" {
			return fmt.Errorf("shell integration not active.\n\n" +
				"  To use activate/deactivate, add this to your .zshrc or .bashrc:\n" +
				"  eval \"$(llm-gate shell-init)\"\n\n" +
				"  Then restart your terminal.")
		}

		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := cfg.GetKey(p.Name)
		if key == "" {
			return fmt.Errorf("no API key stored for %s — use 'llm-gate auth %s' first", p.DisplayName, p.Name)
		}

		// Deactivate any currently active provider that shares the same env var
		for name, pc := range cfg.Providers {
			if pc.Active {
				other := provider.Lookup(name)
				if other != nil && other.EnvVar == p.EnvVar && other.Name != p.Name {
					cfg.SetActive(name, false)
				}
			}
		}

		cfg.SetActive(p.Name, true)
		if err := cfg.Save(); err != nil {
			return err
		}

		// Print the export statement for shell eval
		fmt.Println(shell.ExportCommand(p.EnvVar, key))

		// Print status to stderr so it doesn't interfere with eval
		fmt.Fprintf(os.Stderr, "  ✓ %s exported\n", p.EnvVar)
		fmt.Fprintf(os.Stderr, "  %s is now active.\n", p.DisplayName)
		return nil
	},
}

// ── deactivate ────────────────────────────────────────────────────────────────

var deactivateCmd = &cobra.Command{
	Use:               "deactivate <provider>",
	Short:             "Unset the environment variable for a provider",
	Long:              "Prints an unset statement. Use with: eval \"$(llm-gate deactivate <provider>)\"",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("LLM_GATE_WRAPPER") != "1" {
			return fmt.Errorf("shell integration not active.\n\n" +
				"  To use activate/deactivate, add this to your .zshrc or .bashrc:\n" +
				"  eval \"$(llm-gate shell-init)\"\n\n" +
				"  Then restart your terminal.")
		}

		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.SetActive(p.Name, false)
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Println(shell.UnsetCommand(p.EnvVar))

		fmt.Fprintf(os.Stderr, "  ✓ %s unset\n", p.EnvVar)
		fmt.Fprintf(os.Stderr, "  %s deactivated.\n", p.DisplayName)
		return nil
	},
}

// ── check ─────────────────────────────────────────────────────────────────────

var checkAllFlag bool

var checkCmd = &cobra.Command{
	Use:               "check [provider]",
	Short:             "Test connectivity to a provider's API",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		red := color.New(color.FgRed)
		dim := color.New(color.FgHiBlack)

		if checkAllFlag || len(args) == 0 {
			fmt.Println()
			bold := color.New(color.Bold)
			bold.Println("  Checking connectivity...")
			fmt.Println()

			var ok, total int
			for _, p := range provider.All() {
				key := cfg.GetKey(p.Name)
				if key == "" && p.AuthType != provider.AuthLocal {
					dim.Printf("  ○ %-18s not configured\n", p.Name)
					continue
				}
				total++
				result := check.Check(&p, key)
				if result.OK {
					ok++
					green.Printf("  ✓ %-18s connected    (%dms)\n", p.Name, result.Latency.Milliseconds())
				} else {
					red.Printf("  ✗ %-18s %-12s (%s)\n", p.Name, "failed", result.Error)
				}
			}
			fmt.Printf("\n  %d/%d providers connected\n", ok, total)
			return nil
		}

		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}
		key := cfg.GetKey(p.Name)
		if key == "" && p.AuthType != provider.AuthLocal {
			return fmt.Errorf("no API key stored for %s — use 'llm-gate auth %s' first", p.DisplayName, p.Name)
		}

		result := check.Check(p, key)
		if result.OK {
			green.Printf("  ✓ %s connected (%dms)\n", p.DisplayName, result.Latency.Milliseconds())
		} else {
			red.Printf("  ✗ %s failed: %s\n", p.DisplayName, result.Error)
		}
		return nil
	},
}

// ── status ───────────────────────────────────────────────────────────────────

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show all providers with their status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		bold := color.New(color.Bold)
		green := color.New(color.FgGreen)
		cyan := color.New(color.FgCyan)
		dim := color.New(color.FgHiBlack)

		fmt.Println()
		bold.Printf("  %-20s %-40s %s\n", "Provider", "Description", "Status")
		dim.Println("  " + strings.Repeat("─", 72))

		for _, p := range provider.All() {
			name := p.Name
			desc := p.DisplayName
			if len(p.Tags) > 0 {
				desc += " [" + strings.Join(p.Tags, ", ") + "]"
			}
			if len(p.Aliases) > 0 {
				// Show first alias hint
				if len(p.Aliases) == 1 {
					desc += fmt.Sprintf("  (alias: %s)", p.Aliases[0])
				} else {
					desc += fmt.Sprintf("  (aliases: %s)", strings.Join(p.Aliases[:min(len(p.Aliases), 2)], ", "))
					if len(p.Aliases) > 2 {
						desc += fmt.Sprintf(" +%d more", len(p.Aliases)-2)
					}
				}
			}

			if cfg.IsActive(name) {
				green.Printf("  %-20s %-40s ● active\n", name, desc)
			} else if cfg.IsConfigured(name) {
				cyan.Printf("  %-20s %-40s ○ configured\n", name, desc)
			} else {
				dim.Printf("  %-20s %-40s - not configured\n", name, desc)
			}
		}
		fmt.Println()
		return nil
	},
}

// ── list ──────────────────────────────────────────────────────────────────────

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported LLM providers",
	Run: func(cmd *cobra.Command, args []string) {
		bold := color.New(color.Bold)
		dim := color.New(color.FgHiBlack)

		fmt.Println()
		bold.Printf("  Supported providers (%d):\n\n", len(provider.All()))

		for _, p := range provider.All() {
			line := fmt.Sprintf("  %-20s %s", p.Name, p.DisplayName)
			if len(p.Tags) > 0 {
				line += " [" + strings.Join(p.Tags, ", ") + "]"
			}
			fmt.Println(line)
			if len(p.Aliases) > 0 {
				dim.Printf("    aliases: %s\n", strings.Join(p.Aliases, ", "))
			}
		}
		fmt.Println()
	},
}

// ── remove ────────────────────────────────────────────────────────────────────

var removeCmd = &cobra.Command{
	Use:               "remove <provider>",
	Short:             "Remove a stored API key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: providerCompletion,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := provider.MustLookup(args[0])
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.IsConfigured(p.Name) {
			return fmt.Errorf("no API key stored for %s", p.DisplayName)
		}
		cfg.RemoveKey(p.Name)
		if err := cfg.Save(); err != nil {
			return err
		}

		green := color.New(color.FgGreen)
		green.Printf("  ✓ API key removed for %s\n", p.DisplayName)
		return nil
	},
}

// ── shell-init ────────────────────────────────────────────────────────────────

var shellInitCmd = &cobra.Command{
	Use:   "shell-init",
	Short: "Output shell integration function (eval in your .bashrc/.zshrc)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(shell.ShellInit())
	},
}

// ── current (info) ─────────────────────────────────────────────────────────────

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current env vars",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		for name, pc := range cfg.Providers {
			p, err := provider.MustLookup(name)
			if err != nil {
				return err
			}

			if pc.Active {
				fmt.Printf("export %s=%q\n", p.EnvVar, pc.APIKey)
			} else {
				fmt.Printf("unset %s\n", p.EnvVar)
			}
		}
		return nil
	},
}

// ── config (info) ─────────────────────────────────────────────────────────────

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show config file location and info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("  Config file: %s\n", config.ConfigPath())

		if _, err := os.Stat(config.ConfigPath()); err == nil {
			cfg, err := config.Load()
			if err != nil {
				fmt.Printf("  Error loading: %v\n", err)
				return
			}
			configured := 0
			active := 0
			for _, pc := range cfg.Providers {
				if pc.APIKey != "" {
					configured++
				}
				if pc.Active {
					active++
				}
			}
			fmt.Printf("  Providers configured: %d\n", configured)
			fmt.Printf("  Providers active:     %d\n", active)
		} else {
			fmt.Println("  Config file does not exist yet.")
		}
	},
}

// ── helpers ───────────────────────────────────────────────────────────────────

func providerCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return provider.Names(), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	checkCmd.Flags().BoolVarP(&checkAllFlag, "all", "a", false, "Check all configured providers")

	rootCmd.AddCommand(
		authCmd,
		setCmd,
		updateCmd,
		activateCmd,
		deactivateCmd,
		checkCmd,
		currentCmd,
		statusCmd,
		listCmd,
		removeCmd,
		shellInitCmd,
		configCmd,
	)
}
