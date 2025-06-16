# Docker AI

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/Aj7Ay/docker-ai/main/.github/logo-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/Aj7Ay/docker-ai/main/.github/se7en-ai.jpeg">
    <img alt="Docker AI Logo" src="https://raw.githubusercontent.com/Aj7Ay/docker-ai/main/.github/se7en-ai.jpeg" width="250">
  </picture>
</p>

<p align="center">
  <strong>An AI-powered CLI for Docker.</strong>
  <br />
  Chat with an AI to generate Docker commands, or use the local learning mode to master Docker concepts offline.
</p>

---

## 🚀 Demo

Watch a quick demonstration of `docker-ai` in action.

<p align="center">
  <img src="https://raw.githubusercontent.com/Aj7Ay/docker-ai/main/.github/demo.gif" alt="Docker AI Demo">
</p>

`docker-ai` is a command-line tool that makes it easier to work with Docker by translating natural language into executable commands. It also features a built-in learning mode to help beginners understand core Docker concepts without an internet connection.

## Features

-   **Interactive Shell**: An intuitive shell for running Docker commands.
-   **AI-Powered Commands**: Generate Docker commands from natural language.
-   **Learning Mode**: Learn Docker concepts without leaving your terminal.
-   **Context-Aware**: The AI knows about your running containers.
-   **Command History**: Easily access your previously used commands.

## Installation

### Homebrew (macOS)

First, tap the repository:

```sh
brew tap Aj7Ay/homebrew-docker-ai
```

Then, install `docker-ai`:

```sh
brew install docker-ai
```

### From Source

1.  Clone the repository:
    ```sh
    git clone https://github.com/Aj7Ay/docker-ai.git
    cd docker-ai
    ```
2.  Build the binary:
    ```sh
    go build -o docker-ai ./cmd/docker-ai
    ```
3.  Move the binary to a directory in your `$PATH`:
    ```sh
    sudo mv docker-ai /usr/local/bin/
    ```

## Usage

1.  **Set your API Key**:

    `docker-ai` supports Groq, Gemini, and OpenAI. Set the appropriate environment variable for your chosen provider.

    For Groq:
    ```sh
    export GROQ_API_KEY="your-groq-api-key"
    ```

    For Gemini:
    ```sh
    export GEMINI_API_KEY="your-gemini-api-key"
    ```

2.  **Run `docker-ai`**:

    ```sh
    docker-ai
    ```

    You can specify the LLM provider and model with flags:

    ```sh
    docker-ai --llm-provider=gemini --model=gemini-1.5-pro
    ```

    By default, `docker-ai` uses Groq with the `gemma2-9b-it` model.

### Learning Mode

To use the offline learning mode, set the `DOCKER_AI_MODE` environment variable:

```sh
export DOCKER_AI_MODE=learn
docker-ai
```

### Special Commands

-   `exit` or `quit`: Exit the interactive shell.
-   `reset confirm`: Reset the confirmation prompt for cleanup commands.

## Configuration

`docker-ai` will store a configuration file at `~/.docker-ai-config.json` to remember your preferences, such as skipping cleanup warnings.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details. 