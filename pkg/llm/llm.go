package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// LLMResponse represents a minimal OpenAI-compatible response
type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// QueryLLM sends a prompt to the configured LLM and returns the response.
func QueryLLM(prompt, provider, model string) (string, error) {
	var apiKey, endpoint string
	switch provider {
	case "groq":
		apiKey = os.Getenv("GROQ_API_KEY")
		if apiKey == "" {
			return "", errors.New("GROQ_API_KEY not set")
		}
		endpoint = "https://api.groq.com/openai/v1/chat/completions"
		if model == "" {
			model = "gemma2-9b-it"
		}
	case "openai":
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return "", errors.New("OPENAI_API_KEY not set")
		}
		endpoint = os.Getenv("OPENAI_ENDPOINT")
		if endpoint == "" {
			endpoint = "https://api.openai.com/v1/chat/completions"
		}
		if model == "" {
			model = "gpt-3.5-turbo"
		}
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", provider)
	}

	systemPrompt := `You are a Docker expert. Your primary goal is to generate a single, executable Docker command based on the user's request.

**Rules:**
1.  **Always be specific.** If the user's request is ambiguous and could apply to more than one container, you MUST ask for clarification. Do not guess the container.
2.  **Output format.** Respond with ONLY the Docker command, or a clarifying question. Nothing else. No explanations, no markdown.
3.  **Detached mode.** When starting containers, ALWAYS use detached mode ('-d').
4.  **No follow logs.** When showing logs, never use the follow ('-f') flag unless the user explicitly asks to 'follow' or 'stream' the logs.

**Examples:**

*   **Simple Command**
    *   User: 'show running containers'
    *   Assistant: 'docker ps'

*   **Using Context**
    *   Context: 'The currently running containers are: my-nginx, db-postgres'
    *   User: 'show logs for the web server'
    *   Assistant: 'docker logs my-nginx'

*   **Viewing Specific Logs**
    *   Context: 'The currently running containers are: my-api'
    *   User: 'show me the last 20 logs for the api'
    *   Assistant: 'docker logs --tail 20 my-api'

*   **Vulnerability Scan**
    *   User: 'scan the nginx image for vulnerabilities with scout'
    *   Assistant: 'docker scout cves nginx'

*   **List Models**
    *   User: 'list my local models'
    *   Assistant: 'docker model ls'

*   **Run a Model**
    *   User: 'run the smollm2 model'
    *   Assistant: 'docker model run ai/smollm2'

*   **Run an MCP Server**
    *   User: 'run the file-reader mcp server'
    *   Assistant: 'docker run -d --rm -v /:/context docker/mcp-file-reader:latest'

*   **Ambiguous Request (MUST ask for clarification)**
    *   Context: 'The currently running containers are: web-1, web-2, db-1'
    *   User: 'check logs'
    *   Assistant: 'There are multiple containers running: web-1, web-2, db-1. Which container''s logs do you want to see?'

*   **Stopping all containers**
    *   User: 'stop all containers'
    *   Assistant: 'docker stop $(docker ps -q)'`

	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.1,
		"max_tokens": 50,
	}
	b, _ := json.Marshal(body)
	
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		out, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API error: %s", string(out))
	}
	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}
	if len(llmResp.Choices) == 0 {
		return "", errors.New("no LLM response")
	}
	content := strings.TrimSpace(llmResp.Choices[0].Message.Content)
	
	// Extract the Docker command if it's not already a clean command
	if !strings.HasPrefix(content, "docker ") {
		if cmd, found := ExtractDockerCommand(content); found {
			content = cmd
		}
	}
	
	return content, nil
}

// ExtractDockerCommand tries to find a docker command in the LLM response
func ExtractDockerCommand(response string) (string, bool) {
	// First try to find command with backticks
	re := regexp.MustCompile("```(?:docker)?\\s*(docker[^`]*)```")
	match := re.FindStringSubmatch(response)
	if len(match) > 1 {
		cmd := strings.TrimSpace(match[1])
		// Only return if it's a valid docker command
		if strings.HasPrefix(cmd, "docker ") {
			return cmd, true
		}
	}
	
	// If no backticks, try normal command
	re = regexp.MustCompile(`(?m)^(docker[^\n]*)`)
	match = re.FindStringSubmatch(response)
	if len(match) > 1 {
		cmd := strings.TrimSpace(match[1])
		// Only return if it's a valid docker command
		if strings.HasPrefix(cmd, "docker ") {
			return cmd, true
		}
	}
	return "", false
} 