package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	modelPluginAvailable := true
	llmProvider := flag.String("llm-provider", "groq", "LLM provider to use (groq, openai)")
	model := flag.String("model", "gemma2-9b-it", "Model to use")
	flag.Parse()

	fmt.Println("Docker AI interactive shell. Type 'exit' or 'quit' to leave.")

	historyFile := filepath.Join(os.Getenv("HOME"), ".docker-ai-history")
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	lastContainerName := ""
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

		// Replace 'that container' with lastContainerName if present
		userInput := input
		if strings.Contains(input, "that container") && lastContainerName != "" {
			userInput = strings.ReplaceAll(input, "that container", lastContainerName)
		}

		// Get running containers to provide context to the LLM
		var fullPrompt string
		psCmd := exec.Command("sh", "-c", "docker ps --format '{{.Names}}'")
		psOutput, err := psCmd.CombinedOutput()
		if err == nil && len(psOutput) > 0 {
			containerList := strings.TrimSpace(string(psOutput))
			containerList = strings.ReplaceAll(containerList, "\n", ", ")
			fullPrompt = fmt.Sprintf("The user wants to perform a Docker command. Here are the currently running containers: %s.\n\nUser's request: %s", containerList, userInput)
		} else {
			fullPrompt = userInput
		}

		if !modelPluginAvailable && strings.Contains(userInput, "model") {
			fullPrompt += "\n\nNote: The 'docker model' command is not available on this system."
		}

		response, err := llm.QueryLLM(fullPrompt, *llmProvider, *model)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if !strings.HasPrefix(response, "docker ") {
			fmt.Println(response)
			continue
		}

		fmt.Printf("➜ executing: %s\n", response)

		// Execute the Docker command
		cmd := exec.Command("sh", "-c", response)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if strings.Contains(string(output), "is not a docker command") && strings.Contains(response, "docker model") {
					modelPluginAvailable = false
					fmt.Println("Error: The 'docker model' command is not available.")
					fmt.Println("Please ensure you have installed the Docker Model Runner plugin.")
					fmt.Println("You can try installing it with: 'sudo apt-get update && sudo apt-get install docker-model-plugin'")
					continue
				}
				fmt.Printf("Error executing command: %s\n", exitErr)
				fmt.Println(string(output))
			} else {
				fmt.Printf("Error executing command: %v\n", err)
			}
			continue
		}

		fmt.Print(string(output))

		// If the command was a 'docker run', get and store the container name
		if strings.HasPrefix(response, "docker run") {
			containerID := strings.TrimSpace(string(output))
			if containerID != "" {
				inspectCmd := exec.Command("sh", "-c", fmt.Sprintf("docker inspect --format '{{.Name}}' %s | sed 's/^\\///'", containerID))
				inspectOut, inspectErr := inspectCmd.CombinedOutput()
				if inspectErr == nil {
					lastContainerName = strings.TrimSpace(string(inspectOut))
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

	if f, err := os.Create(historyFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing history file:", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
} 