# Command-Line Reference

`docker-ai` can be run in two modes: interactive (default) and single-command.

## Interactive Mode

To start the interactive shell, simply run the `docker-ai` command.

```bash
docker-ai
```

This will open a `docker-ai>` prompt where you can type your requests in natural language.

```
docker-ai> list all running containers
...
docker-ai> exit
```

You can use the `--llm-provider` and `--model` flags when starting the interactive shell to configure the session.

### Special Commands

-   `exit` or `quit`: Exits the interactive shell.
-   `reset confirm`: If you previously selected "don't ask again" for cleanup command warnings, this command will reset that preference, and you will be prompted for confirmation again.

## Single-Command Mode

You can execute a single command and exit immediately by using the `-c` flag.

```bash
docker-ai -c "delete all unused docker images"
```

## Flags

| Flag             | Argument      | Description                                     | Default            |
| ---------------- | ------------- | ----------------------------------------------- | ------------------ |
| `-c`             | `"command"`   | Execute a single command and exit.              | `""`               |
| `--llm-provider` | `provider`    | Specify the LLM provider to use.                | `groq`             |
|                  | *Allowed:*    | `groq`, `gemini`, `openai`                      |                    |
| `--model`        | `model_name`  | Specify the exact model name to use.            | `gemma-3n-e4b-it`  | 