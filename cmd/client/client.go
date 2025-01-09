package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/joho/godotenv"
)

type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseStream struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

var (
	GROQ_API_KEY = os.Getenv("GROQ_API_KEY")
)

func ChatCompletion(prompt string) {

	if GROQ_API_KEY == "" {
		log.Fatal("Warning: GROQ_API_KEY wasn't found in the environemnt!")
	}

	body := GroqRequest{
		Model:  "llama-3.3-70b-versatile",
		Stream: true,
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Print("Error Parsing Request body: ", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		log.Print("Error creating request: ", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", GROQ_API_KEY))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Print("Error making request: ", err)
	}
	defer res.Body.Close()

	fmt.Print("OLLAMA: ")
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			// Server Sent Event format (a standard protocol used by many streaming APIs) requires each data chunk to start with "data: ".
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			// [DONE] is a SSE standard message indicating the stream is complete.
			break
		}

		var streamRes ResponseStream
		if err := json.Unmarshal([]byte(data), &streamRes); err != nil {
			continue
		}

		if len(streamRes.Choices) > 0 {
			fmt.Print(formatMarkdown(streamRes.Choices[0].Delta.Content))
		}
	}
	fmt.Println()
}

func formatMarkdown(text string) string {
	re := regexp.MustCompile(`\*\*|__|\*|_`)
	return re.ReplaceAllString(text, "")
}
