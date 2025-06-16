package main

import (
	"bufio"
	"bytes"
	"docker-ai/pkg/config"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"docker-ai/pkg/learning"
	"docker-ai/pkg/llm"

	"github.com/peterh/liner"
)

func main() {
	if os.Getenv("DOCKER_AI_MODE") == "learn" {
		runLearningMode()
		return
	}

	runAIMode()
}

func isCleanupCommand(command string) bool {
	// Regex to detect potentially destructive docker commands
	cleanupRegex := regexp.MustCompile(`docker\s+(system|container|image|volume|network)\s+prune|docker\s+rm|docker\s+rmi|docker\s+volume\s+rm|docker\s+network\s+rm`)
	return cleanupRegex.MatchString(command)
}

func runLearningMode() {
	fmt.Println("Docker AI Learning Mode. Ask me questions about Docker!")
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	for {
		input, err := line.Prompt("learn-ai> ")
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		if input == "exit" || input == "quit" {
			break
		}

		if input == "" {
			continue
		}

		response, err := learning.HandleQuery(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(response)
	}
}

func runAIMode() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: could not load config file: %v\n", err)
	}

	llmProvider := flag.String("llm-provider", "groq", "LLM provider to use (groq, gemini, openai)")
	model := flag.String("model", "gemma-3n-e4b-it", "Model to use")
	command := flag.String("c", "", "Execute a single command and exit")
	flag.Parse()

	if *command != "" {
		runSingleCommand(&appConfig, *command, *llmProvider, *model)
		return
	}

	runInteractiveMode(&appConfig, *llmProvider, *model)
}

func runInteractiveMode(appConfig *config.Config, llmProvider, model string) {
	fmt.Println("Docker AI interactive shell. Type 'exit' or 'quit' to leave.")

	historyFile := filepath.Join(os.Getenv("HOME"), ".docker-ai-history")
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	for {
		input, err := line.Prompt("docker-ai> ")
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		if input == "" {
			continue
		}

		line.AppendHistory(input)

		if input == "exit" || input == "quit" {
			break
		}

		// Before running the command, close the liner to restore the terminal
		line.Close()

		runSingleCommand(appConfig, input, llmProvider, model)

		// After a command that takes over stdin, the terminal can be left in a
		// "raw" state. We use `stty` to force it back to a sane mode before
		// re-initializing the liner.
		sttyCmd := exec.Command("stty", "-raw", "echo")
		sttyCmd.Stdin = os.Stdin
		_ = sttyCmd.Run()

		// Re-initialize the liner to take back control of the terminal
		line = liner.NewLiner()
		line.SetCtrlCAborts(true)
		if f, err := os.Open(historyFile); err == nil {
			line.ReadHistory(f)
			f.Close()
		}
	}

	if f, err := os.Create(historyFile); err != nil {
		fmt.Printf("Failed to save history: %v\n", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	if err := config.SaveConfig(*appConfig); err != nil {
		fmt.Println("Failed to save configuration:", err)
	} else {
		fmt.Println("Cleanup confirmation has been reset. You will be prompted before cleanup commands are run.")
	}
}

func runSingleCommand(appConfig *config.Config, input, llmProvider, model string) {
	if input == "reset confirm" {
		appConfig.SkipCleanupWarning = false
		if err := config.SaveConfig(*appConfig); err != nil {
			fmt.Println("Failed to save configuration:", err)
		} else {
			fmt.Println("Cleanup confirmation has been reset. You will be prompted before cleanup commands are run.")
		}
		return
	}

	// Replace 'that container' with lastContainerName if present
	userInput := input
	if strings.Contains(input, "that container") && appConfig.LastContainerName != "" {
		userInput = strings.ReplaceAll(input, "that container", appConfig.LastContainerName)
	}

	// Get running containers to provide context to the LLM
	var fullPrompt string
	psCmd := exec.Command("sh", "-c", "docker ps -a --format '{{.Names}} ({{.Status}})'")
	psOutput, err := psCmd.CombinedOutput()
	if err == nil && len(psOutput) > 0 {
		containerList := strings.TrimSpace(string(psOutput))
		containerList = strings.ReplaceAll(containerList, "\n", ", ")
		fullPrompt = fmt.Sprintf("The user wants to perform a Docker command. Here is a list of all containers (running and stopped): %s.\n\nUser's request: %s", containerList, userInput)
	} else {
		fullPrompt = userInput
	}

	if !strings.Contains(userInput, "model") {
		fullPrompt += "\n\nNote: The 'docker model' command is not available on this system."
	}

	response, err := llm.QueryLLM(fullPrompt, llmProvider, model)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !strings.HasPrefix(response, "docker ") {
		fmt.Println(response)
		return
	}

	// Cleanup command confirmation
	if isCleanupCommand(response) && !appConfig.SkipCleanupWarning {
		fmt.Printf("WARNING: The generated command is a cleanup command:\n%s\n\n", response)

		// For single-command mode, we need a way to confirm.
		// We'll use a simple prompt here, but this could be improved.
		// A liner isn't running, so we use fmt.
		fmt.Print("Are you sure you want to execute? [y]es, [n]o, [d]on't ask again: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')

		answer = strings.ToLower(strings.TrimSpace(answer))

		switch answer {
		case "y", "yes":
			// continue to execution
		case "d", "dont", "don't ask again":
			appConfig.SkipCleanupWarning = true
			if err := config.SaveConfig(*appConfig); err != nil {
				fmt.Println("Failed to save configuration:", err)
			}
			// continue to execution
		default:
			fmt.Println("Execution cancelled.")
			return
		}
	}

	fmt.Printf("➜ executing: %s\n", response)

	// Execute the Docker command
	var stderrBuf bytes.Buffer
	cmd := exec.Command("sh", "-c", response)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The command exited with a non-zero status.
			// Stderr is already printed. We can analyze it from the buffer.
			stderrString := stderrBuf.String()
			if strings.Contains(response, "docker model") {
				if strings.Contains(stderrString, "is not a docker command") {
					fmt.Println("\nError: The 'docker model' command is not available.")
					fmt.Println("Please ensure you have installed the Docker Model Runner plugin.")
					fmt.Println("You can try installing it with: 'sudo apt-get update && sudo apt-get install docker-model-plugin'")
				}
			}
			// Check for docker scout update message
			if strings.Contains(stderrString, "New version") && strings.Contains(stderrString, "available") {
				fmt.Println("\nHint: A new version of Docker Scout is available.")
				fmt.Println("To update it, you can run the following command in your terminal:")
				fmt.Println("curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s --")
			}
			fmt.Printf("\nCommand finished with error: %s\n", exitErr)
		} else {
			fmt.Printf("Error executing command: %v\n", err)
		}
		return
	}

	// If the command was a 'docker run', get and store the container name
	if strings.HasPrefix(response, "docker run") {
		// With streaming, we can't just get the container ID from the output of this command.
		// We will need a new way to get the container name.
		// A simple approach is to list recent containers and assume the latest one is it.
		getLatestContainerCmd := `docker ps -l -q`
		getLatestCmd := exec.Command("sh", "-c", getLatestContainerCmd)
		latestID, latestErr := getLatestCmd.Output()
		if latestErr == nil && len(latestID) > 0 {
			containerID := strings.TrimSpace(string(latestID))
			inspectCmd := exec.Command("sh", "-c", fmt.Sprintf("docker inspect --format '{{.Name}}' %s | sed 's/^\\///'", containerID))
			inspectOut, inspectErr := inspectCmd.CombinedOutput()
			if inspectErr == nil {
				appConfig.LastContainerName = strings.TrimSpace(string(inspectOut))
				fmt.Printf("\nContainer '%s' started.\n", appConfig.LastContainerName)
			}
		}

		statusCmd := exec.Command("sh", "-c", "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'")
		statusOut, statusErr := statusCmd.CombinedOutput()
		if statusErr == nil {
			fmt.Println("\nCurrently running containers:")
			fmt.Print(string(statusOut))
		}
	}
} 