package tokenizer

import (
	"os"
)

// Approximate number of characters per token for GPT models
const charactersPerToken = 4

// CountFileTokens estimates the number of tokens in a file
func CountFileTokens(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Simple heuristic: ~4 characters per token
	chars := len(content)
	tokens := chars / charactersPerToken

	// Add 1 token if there's a remainder
	if chars%charactersPerToken > 0 {
		tokens++
	}

	return tokens, nil
}
