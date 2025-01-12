package client

import "os"

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

type HistoryMessage struct {
	Time         string        `json:"time"`
	Conversation []GroqMessage `json:"conversation"`
}

type History struct {
	Conversations []HistoryMessage `json:"conversations"`
}

var (
	GROQ_API_KEY = os.Getenv("GROQ_API_KEY")
)
