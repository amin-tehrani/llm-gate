# llm-gate

> A central hub CLI for managing multiple LLM provider API keys.

Store, activate, deactivate, and check connectivity for 30+ LLM providers from a single tool.

## Install

```bash
# Build from source
make build

# Install to PATH
make install
```

## Quick Start

```bash
# Set up shell integration (add to .bashrc / .zshrc)
eval "$(llm-gate shell-init)"

# Store an API key
llm-gate auth openai      # interactive (hidden input)
llm-gate set openai sk-... # direct

# Activate — exports the env var
llm-gate activate openai
# → export OPENAI_API_KEY="sk-..."

# Check connectivity
llm-gate check openai
llm-gate check --all

# View all providers and status
llm-gate current

# Deactivate — unsets the env var
llm-gate deactivate openai

# Remove a stored key
llm-gate remove openai
```

## Commands

| Command                   | Description                                        |
| ------------------------- | -------------------------------------------------- |
| `auth <provider>`         | Interactive API key prompt with connectivity check |
| `set <provider> <key>`    | Store an API key directly                          |
| `update <provider> <key>` | Update an API key (alias for set)                  |
| `activate <provider>`     | Export the env var for this provider               |
| `deactivate <provider>`   | Unset the env var                                  |
| `check [provider]`        | Test connectivity (`--all` for all configured)     |
| `current`                 | Show all providers with active/configured status   |
| `list`                    | List all supported providers and aliases           |
| `remove <provider>`       | Remove a stored key                                |
| `shell-init`              | Output shell wrapper function                      |
| `config`                  | Show config file path and info                     |

## Shell Integration

Since child processes can't modify the parent shell's environment, `activate` and `deactivate` print shell statements that need to be `eval`'d. The `shell-init` command provides a wrapper that handles this automatically:

```bash
# Add to your .bashrc or .zshrc:
eval "$(llm-gate shell-init)"
```

After that, `llm-gate activate` and `llm-gate deactivate` work seamlessly.

## Supported Providers

Run `llm-gate list` to see all 30+ supported providers including OpenAI, Anthropic, Google Gemini, Groq, Mistral, DeepSeek, xAI, Together AI, and many more.

## Config

Keys are stored in `~/.config/llm-gate/config.yaml` with `0600` permissions.

```bash
llm-gate config  # show config path and info
```

## License

MIT
