package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func callLLM(system_prompt, user_prompt string) string {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT4o,
			Temperature: 0.3,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system_prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: user_prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Println("Error calling OpenAI API:", err)
	}

	return response.Choices[0].Message.Content
}

func getSentiment(text string) string {
	system_prompt := "You will be presented with text from a news article. You need to determine the sentiment of the text. The sentiment should be one of the following: 'positive', 'neutral', 'negative'. If you are unsure, return neutral. Do not return anything else."

	sentiment := callLLM(system_prompt, text)

	return strings.ToLower(sentiment)
}
