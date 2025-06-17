# LLM Providers

`docker-ai` supports multiple Large Language Model (LLM) providers. You can choose which provider you want to use by setting an environment variable and using the `--llm-provider` flag.

## OpenAI

By default, `docker-ai` uses OpenAI models.

### Configuration

1.  **Get an API Key**: If you don't have one, get an API key from the [OpenAI Platform](https://platform.openai.com/api-keys).
2.  **Set Environment Variable**: Export your API key as an environment variable.

    ```bash
    export OPENAI_API_KEY="sk-..."
    ```

    To make this setting permanent, add the line to your shell's configuration file (e.g., `~/.zshrc`, `~/.bash_profile`, or `~/.bashrc`).

### Usage

Run `docker-ai` with your prompt:

```bash
docker-ai "delete all dangling images"
```

You can also specify a particular model using the `--model` flag:

```bash
docker-ai --model=gpt-4 "list all running containers"
```

## Gemini

You can also use Google's Gemini models with `docker-ai`.

### Configuration

1.  **Get an API Key**: Generate a free API key from [Google AI Studio](https://aistudio.google.com/app/apikey).
2.  **Set Environment Variable**: Export your API key as an environment variable.

    ```bash
    export GEMINI_API_KEY="your-api-key"
    ```

    Add this line to your shell's configuration file to make it permanent.

### Usage

To use Gemini, you must specify it with the `--llm-provider` flag:

```bash
docker-ai --llm-provider=gemini "show me all images related to ubuntu"
```

You can also specify a Gemini model:

```bash
docker-ai --llm-provider=gemini --model=gemini-1.5-flash "stop the container named web-server"
``` 