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
			Model:       openai.GPT4oMini,
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
	system_prompt := "Du vil bli presentert med et utdrag fra en nyhetsartikkel. Du skal bestemme sentimentet i artikkelen. Sentimentet skal være en av følgende: 'positiv', 'nøytral', 'negativ'. Svaret ditt skal være et av disse ordene, og ikke noe mer."

	sentiment := callLLM(system_prompt, text)

	return strings.ToLower(sentiment)
}
