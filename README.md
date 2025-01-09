## OLLAMA-CLI - Chat Completion

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

```
>> ollama --help
CLI Wrapper for llama.

Usage:
  ollama [command]

Available Commands:
  ask         prompt the LLM
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help     help for ollama

Use "ollama [command] --help" for more information about a command.
```

```
>> ollama ask --help
prompt the LLM

Usage:
  ollama ask <message> [flags]

Flags:
  -h, --help   help for ask
```
