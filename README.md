# Docker AI

<p align="center">
  <img src="https://raw.githubusercontent.com/Aj7Ay/docker-ai/main/.github/logo.svg" alt="Docker AI Logo" width="150"/>
</p>

<p align="center">
  <strong>An AI-powered CLI for Docker.</strong>
  <br />
  Chat with an AI to generate Docker commands, or use the local learning mode to master Docker concepts offline.
</p>

---

`docker-ai` is a command-line tool that makes it easier to work with Docker by translating natural language into executable commands. It also features a built-in learning mode to help beginners understand core Docker concepts without an internet connection.

## Features

-   **AI Mode**: Translates natural language queries into Docker commands using providers like Groq or OpenAI.
-   **Learning Mode**: An offline, interactive shell to learn Docker fundamentals.
-   **Interactive Shell**: A chat-like interface for a fluid user experience.
-   **Command History**: Use up/down arrow keys to recall previous commands in both modes.
-   **Context-Aware**: The AI can recognize running containers to correctly target commands.

## Installation

`docker-ai` is built with Go and requires Go 1.18 or higher to be installed on your system.

1.  **Clone the Repository**
    ```sh
    git clone https://github.com/Aj7Ay/docker-ai.git
    cd docker-ai
    ```

2.  **Build the Binary**
    This command compiles the source code into a single executable file named `docker-ai`.
    ```sh
    go build -o docker-ai ./cmd/docker-ai/main.go
    ```

3.  **Make it Executable (Optional but Recommended)**
    To run `docker-ai` from anywhere on your system, move it to a directory that is in your system's `PATH`.
    ```sh
    # First, make sure the file is executable
    chmod +x docker-ai

    # Then, move it to a common location for binaries
    sudo mv docker-ai /usr/local/bin/
    ```

## Usage

### AI Mode (Default)

This is the default mode. It connects to an LLM provider to generate Docker commands.

1.  **Set Your API Key**
    You need to provide an API key for the LLM provider you want to use.
    ```sh
    # For Groq
    export GROQ_API_KEY="your-groq-api-key"

    # For OpenAI
    export OPENAI_API_KEY="your-openai-api-key"
    ```

2.  **Run the Tool**
    ```sh
    # Run with default provider (Groq)
    docker-ai

    # Or specify a provider and model
    docker-ai --llm-provider=openai --model=gpt-4
    ```

3.  **Example Session**
    ```
    docker-ai> show running containers
    ➜ executing: docker ps
    CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES

    docker-ai> show me the last 10 logs for my_container
    ➜ executing: docker logs --tail 10 my_container
    ...
    ```

### Learning Mode

This mode works completely offline and is designed to teach Docker concepts.

1.  **Run in Learning Mode**
    Set the `DOCKER_AI_MODE` environment variable to `learn`.
    ```sh
    DOCKER_AI_MODE=learn docker-ai
    ```

2.  **Example Session**
    ```
    Docker AI Learning Mode. Ask me questions about Docker!
    learn-ai> what is a dockerfile?
    A Dockerfile is like a recipe for baking a cake. The "cake" is your application's container image...

    learn-ai> how to check all containers
    To see your containers, you use the 'docker ps' command...
    ```

### Special Commands

These commands work in both modes:

-   `version`: Prints the current version of the tool.
-   `clear`: Clears the terminal screen.
-   `exit` or `quit`: Exits the interactive shell.

## Contributing

Contributions are welcome! If you have suggestions or find a bug, please open an issue.

## License

This project is licensed under the Apache-2.0 License. 