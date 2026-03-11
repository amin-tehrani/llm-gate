<div align="center">
  <h1>🚪 LLM Gate</h1>
  <p><strong>The Ultimate Control Plane for your LLM API Keys</strong></p>
  
  <p>
    <a href="https://golang.org/doc/install"><img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
    <a href="https://github.com/spf13/cobra"><img src="https://img.shields.io/badge/CLI-Cobra-blue?style=flat-square" alt="Cobra"></a>
    <a href="#"><img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License"></a>
  </p>
</div>

<br/>

Working with multiple LLM providers means juggling dozens of API keys, environment variables, and authentication flows. **`llm-gate`** brings order to the chaos. It serves as a unified, central hub CLI to efficiently store, manage, verify, and seamlessly inject API keys for over **30 LLM providers**.

## ✨ Features

- **🌐 30+ Supported Providers**: OpenAI, Anthropic, Google Gemini, Groq, Mistral, Code Llama, xAI (Grok), and many more.
- **🛡️ Secure Storage**: Keys are permanently stored in `~/.config/llm-gate/config.yaml` using strict `0600` file permissions.
- **🚀 Seamless Shell Integration**: Activate or deactivate any provider's environment variables directly within your current terminal session via `eval`.
- **✅ Real-Time Connectivity**: Ping any network to verify if your API credentials are active and healthy, retrieving latencies and specific JSON error payloads if failures occur.
- **🧭 Interactive Onboarding**: Provides interactive, guided flows that will automatically open the exact API-key generation URLs in your browser for standard providers.
- **🔄 Smart Aliasing**: Resolves popular names (e.g. `grok` -> `xai`, `github-copilot` -> `copilot`).

---

## 📦 Installation

To get up and running, clone the repository and build from the source:

```bash
# Clone the repository
git clone https://github.com/amintehrani/llm-gate.git
cd llm-gate

# Build and install to your PATH
make install
```

### 🔌 Enable Shell Integration (Required)
Since child processes cannot modify the parent shell's environment variables, `llm-gate` uses an `eval` wrapper function. You must add the following snippet to your shell configuration (`~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`):

```bash
# Add this line to the end of your shell rc file:
eval "$(llm-gate shell-init)"
```
Restart your terminal, and `llm-gate activate` will now work seamlessly! Furthermore, opening new terminal tabs or splits will automatically pull your persisted configuration from `llm-gate current` and keep your API keys active across sessions.

---

## 🛠️ Quick Start

#### 1. Add your first provider
Start by authenticating an LLM provider. The interactive CLI will prompt you and even open your browser directly to the provider's API Key generation page:
```bash
llm-gate auth openai
```

*(Alternatively, you can skip the interactive prompt by running: `llm-gate set openai sk-...`)*

#### 2. Check your connection
Make sure your key works:
```bash
llm-gate check openai
# Output: ✓ OpenAI connected (1338ms)
```

#### 3. Activate the key
Export the API key natively into your shell's current session:
```bash
llm-gate activate openai
# Your shell now has $OPENAI_API_KEY exported perfectly!
```

---

## 🕹️ Command Reference

| Command                   | Description                                                                                            |
| ------------------------- | ------------------------------------------------------------------------------------------------------ |
| `help`                    | Standard help & usage command list.                                                                    |
| `list`                    | Show all 30+ fully supported LLM providers and their mapped aliases.                                   |
| `status`                  | Print a clean table of all providers and their localized configuration statuses.                       |
| `current`                 | Echoes the active environment variable `export` and `unset` statements tailored for your shell.        |
| `auth <provider>`         | Launch the interactive authentication wizard, guiding you with browser URLs.                           |
| `set <provider> <key>`    | Save an API key directly without interactive prompts.                                                  |
| `update <provider> <key>` | Update/overwrite an API key (alias for `set`).                                                         |
| `activate <provider>`     | Export the environment variable mapping specifically for this vendor into your active CLI environment. |
| `deactivate <provider>`   | Unset the exported environment variable.                                                               |
| `check [provider]`        | Test provider health. Optionally pass `--all` or `-a` to ping all customized providers sequentially.   |
| `remove <provider>`       | Delete a stored key locally.                                                                           |
| `config`                  | Echo out the location strings for where configurations exist locally.                                  |
| `completion <shell>`      | Generate the autocompletion script (bash, zsh, fish, powershell).                                      |

---

## 🤖 Supported Providers

Run `llm-gate list` to view the comprehensive list. Highlights include:

- **OpenAI** (`openai`, `openai-codex`)
- **Anthropic** (`anthropic`)
- **Google Gemini** (`gemini`, aliases: `google`)
- **Meta/Local** (`ollama`)
- **xAI/Grok** (`xai`, aliases: `grok`)
- **Cloud Infrastructure Models** (`cloudflare`, `vercel`, `doubao`, `qwen`, `qianfan`, `bedrock`)
- **High-Performance inference hubs** (`groq`, `together`, `fireworks`, `openrouter`)
- **Asian Leading Models** (`moonshot`, `kimi-code`, `glm`, `zai`, `minimax`, `deepseek`)
- **Others** (`mistral`, `synthetic`, `opencode`, `novita`, `perplexity`, `cohere`, `copilot`)

---

## 🤝 Contributing

Contributions, issues, and feature requests are always welcome!
Feel free to check the [issues page](https://github.com/amintehrani/llm-gate/issues). If you want to add a new provider, just append it to the `registry.go` array with the designated configuration parameters.

## 📄 License
This architecture is licensed under **[MIT](LICENSE)**. 
