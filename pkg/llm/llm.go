package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"google.golang.org/genai"
)

// LLMResponse represents a minimal OpenAI-compatible response
type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// queryGemini handles the request using the Gemini API
func queryGemini(ctx context.Context, model, prompt, systemPrompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", errors.New("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create a new Gemini client: %w", err)
	}

	// Combine system prompt and user prompt
	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\nAssistant:", systemPrompt, prompt)

	result, err := client.Models.GenerateContent(ctx, model, genai.Text(fullPrompt), nil)
	if err != nil {
		return "", fmt.Errorf("gemini content generation failed: %w", err)
	}

	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil || len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("gemini returned no content")
	}

	return result.Text(), nil
}

// QueryLLM sends a prompt to the configured LLM and returns the response.
func QueryLLM(prompt, provider, model string) (string, error) {
	systemPrompt := `You are an expert-level CLI tool that translates natural language into a single, executable Docker command.

**Primary Directive:** NEVER respond conversationally. Your only purpose is to provide a single, valid Docker command.

**Rules:**
1.  **No Explanations:** Do not provide any explanation, context, or markdown. Output only the raw command.
2.  **Use Provided Context:** You will be given a list of containers. You MUST use the names or IDs from this list. Do not use placeholders.
3.  **Be Specific:** If a user's request is ambiguous (e.g., "delete the container" when multiple exist), you MUST ask for clarification. Do not guess which container to use.
4.  **Detached Mode:** When starting containers, ALWAYS use detached mode ('-d') unless explicitly told otherwise.
5.  **Logging:** For 'docker logs', do not use the '-f' (follow) flag unless requested. Default to showing the last 20 lines (e.g., 'docker logs --tail 20 <container>').
6.  **Scout and Model Runner:**
    *   The user has 'docker scout' and 'docker model' commands.
    *   For 'docker scout', the primary subcommands are 'cves', 'recommendations', and 'quickview', which are used with an image name (e.g., 'docker scout cves nginx').
    *   If the user asks to "install" or "update" 'docker scout', you MUST respond with only this exact text: To update Docker Scout, please run this command in your terminal: curl -sSfL https://raw.githubusercontent.com/docker/scout-cli/main/install.sh | sh -s --
7.  **No Guesses:** If you cannot determine a valid Docker command from the user's request, ask a clarifying question. Do not make up a command.

**Examples:**
- User: "show me all running containers" -> "docker ps"
- User: "list all images" -> "docker images"
- User: "delete the 'web-server' container" -> "docker rm web-server"
- User: "show me the logs for 'api-gateway'" -> "docker logs --tail 20 api-gateway"
- User: "what's the docker scout command to find vulnerabilities in the latest ubuntu image" -> "docker scout cves ubuntu:latest"
`
	var apiKey, endpoint string

	switch provider {
	case "groq":
		apiKey = os.Getenv("GROQ_API_KEY")
		if apiKey == "" {
			return "", errors.New("GROQ_API_KEY not set")
		}
		endpoint = "https://api.groq.com/openai/v1/chat/completions"
		if model == "" {
			model = "gemma-3n-e4b-it"
		}
	case "gemini":
		if model == "" {
			model = "gemini-1.5-flash"
		}
		// Gemini uses its own Go SDK, so we'll call its function and return
		ctx := context.Background()
		return queryGemini(ctx, model, prompt, systemPrompt)
	case "openai":
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return "", errors.New("OPENAI_API_KEY not set")
		}
		endpoint = "https://api.openai.com/v1/chat/completions"
		if model == "" {
			model = "gpt-4o"
		}
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", provider)
	}

	// This part is for Groq and OpenAI (OpenAI-compatible APIs)
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.1,
		"max_tokens":  1024,
		"top_p":       1,
		"stream":      false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var llmResponse LLMResponse
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		return "", err
	}

	if len(llmResponse.Choices) == 0 {
		return "", errors.New("no response from LLM")
	}

	// Clean up the response to remove markdown and extra quotes
	response := llmResponse.Choices[0].Message.Content
	response = strings.TrimSpace(response)
	response = regexp.MustCompile("`{3}(bash|sh)?").ReplaceAllString(response, "")
	response = strings.Trim(response, "`\n ")

	return response, nil
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