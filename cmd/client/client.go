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

func ChatCompletion(prompt string, context bool) {

	if GROQ_API_KEY == "" {
		log.Fatal("Warning: GROQ_API_KEY wasn't found in the environemnt!")
	}

	fullPrompt := prompt
	if context {
		ctx := getContext()
		fullPrompt = prompt + ctx
	}

	body := GroqRequest{
		Model:  "llama-3.3-70b-versatile",
		Stream: true,
		Messages: []GroqMessage{
			{
				Role:    "user",
				Content: fullPrompt,
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

	response := ""

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
			response += formatMarkdown(streamRes.Choices[0].Delta.Content)
		}
	}
	fmt.Println()

	var message HistoryMessage
	message.Time = time.Now().Format(time.RFC3339)
	message.Conversation = append(message.Conversation, GroqMessage{
		Role:    "user: ",
		Content: splitContextAndPrompt(prompt),
	},
		GroqMessage{
			Role:    "Llama: ",
			Content: response,
		},
	)

	saveToHistory(message)
}

func formatMarkdown(text string) string {
	re := regexp.MustCompile(`\*\*|__|\*|_`)
	return re.ReplaceAllString(text, "")
}

func saveToHistory(message HistoryMessage) {
	filename := "ollama_history.json"

	history, _ := loadHistory(filename)

	history.Conversations = append(history.Conversations, message)
	jsonData, err := json.Marshal(history)
	if err != nil {
		fmt.Println("Error marshaling json data:", err)
	}

	err = os.WriteFile(filename, jsonData, 0664)
	if err != nil {
		fmt.Println("Error saving history to file:", err)
		return
	}

}

func loadHistory(filename string) (History, error) {
	var history History
	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &history)
		if err != nil {
			return History{}, err
		}
	} else {
		history = History{
			Conversations: []HistoryMessage{},
		}
	}
	return history, nil
}

func getContext() string {
	history, _ := loadHistory("ollama_history.json")
	var context string
	context += "\n\nContext: \n"

	now := time.Now()
	contextWindow := now.Add(-10 * time.Minute)

	for _, convos := range history.Conversations {
		convoTime, _ := time.Parse(time.RFC3339, convos.Time)
		if convoTime.Before(contextWindow) {
			continue
		}

		for _, ctx := range convos.Conversation {
			context += fmt.Sprintf("%s%s. ", ctx.Role, ctx.Content)
		}
	}

	return context
}

func DisplayHistory(filter int16) {
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("\t\tYour History with ollama-cli for %d days", filter)
	fmt.Println("\n" + strings.Repeat("-", 70) + "\n\n")

	history, err := loadHistory("ollama_history.json")
	if err != nil {
		fmt.Println("Nothing to show right now.")
	}

	thresholdDate := time.Now().AddDate(0, 0, -int(filter))

	for _, convos := range history.Conversations {
		convoTime, _ := time.Parse(time.RFC3339, convos.Time)

		if convoTime.After(thresholdDate) {
			fmt.Printf("date: %v", convoTime)
			for _, convo := range convos.Conversation {
				fmt.Printf("\n%s%s", convo.Role, convo.Content)
			}
			fmt.Println("\n\n" + strings.Repeat("-", 70) + "\n")
		}
	}
}

func splitContextAndPrompt(fullPrompt string) string {
	if strings.Contains(fullPrompt, "Context:") {
		parts := strings.Split(fullPrompt, "Context:")
		if len(parts) >= 2 {
			contextAndPrompt := parts[len(parts)-1]
			if lastDot := strings.LastIndex(contextAndPrompt, ". "); lastDot != -1 {
				return strings.TrimSpace(contextAndPrompt[lastDot+2:])
			}
			return strings.TrimSpace(contextAndPrompt)
		}
	}
	return fullPrompt
}
