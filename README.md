## OLLAMA-CLI - Chat Completion with context

![preview](https://github.com/user-attachments/assets/9c9c9f61-6fd8-4986-bc06-6860ce17f051)

---

Get the `GROQ_API_KEY` from [GRQO CLOUD](https://console.groq.com/keys).

- Set the key to environment variable

```bash
export GROQ_API_KEY=<your_key>
```

- Verify your key

```bash
echo $GROQ_API_KEY
```

- check your PATH

```bash
echo $PATH
```

- move binary to your PATH

```bash
mv ./ollama ~/.local/bin
# ./local/bin is my PATH
```

- use the CLI

```bash
ollama ask "your prompt...""
```

- Context

```bash
ollama ask -c "prompt which needs context of previous chats..."
```

- History

```bash
ollama history        # all history
ollama history -d 2   # past 2 days
```

---

```
>> ollama --help
CLI Wrapper for llama.

Usage:
  ollama [command]

Available Commands:
  ask         prompt the LLM
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  history     Show the history of prompts

Flags:
  -h, --help     help for ollama
  -t, --toggle   Help message for toggle

Use "ollama [command] --help" for more information about a command.
```

```
>> ollama ask --help
prompt the LLM

Usage:
  ollama ask <message> [flags]

Flags:
  -c, --ctx    Use Context Mode. [Context window: 10 mins]
  -h, --help   help for ask
```
