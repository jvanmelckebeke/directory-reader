package tokenizer

import (
	"os"

	"github.com/pkoukk/tiktoken-go"
)

// CountFileTokens counts the number of tokens in a file using OpenAI's tokenizer
func CountFileTokens(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	tokenizer, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		return 0, err
	}

	tokens := tokenizer.Encode(string(content), nil, nil)
	return len(tokens), nil
}
