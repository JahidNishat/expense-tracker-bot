package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/masudur-rahman/expense-tracker-bot/models"
)

type GeneratorAI string

const (
	GeneratorGemini     GeneratorAI = "gemini"
	GeneratorOpenRouter GeneratorAI = "open-router"
)

type ClassificationResult struct {
	Category    string  `json:"category_id"`
	Subcategory string  `json:"subcategory_id"`
	Confidence  float64 `json:"confidence"` // Optional: ask for reasoning/confidence
}

func TxnCategoryGenerator(ctx context.Context, userInput string, ai ...GeneratorAI) (subCatID string, err error) {
	var result *ClassificationResult
	generator := GeneratorOpenRouter
	if len(ai) > 0 {
		generator = ai[0]
	}

	taxonomyJson, err := json.MarshalIndent(models.TxnSubcategories, "", "  ")
	if err != nil {
		return "", err
	}

	switch generator {
	case GeneratorGemini:
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return userInput, nil
		}
		result, err = TxnSubcategoryClassifier(ctx, apiKey, userInput, string(taxonomyJson))
		if err != nil {
			return "", err
		}
	default:
		apiKey := os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			return userInput, nil
		}
		client := NewClient(apiKey)
		result, err = client.TxnSubcategoryClassifier(ctx, userInput, string(taxonomyJson))
		if err != nil {
			return "", err
		}
	}

	fmt.Printf("Matched: %s > %s (Confidence: %v)\n", result.Category, result.Subcategory, result.Confidence)
	if result.Subcategory == "" {
		return "", errors.New("transaction category can't be determined")
	}
	return result.Subcategory, nil
}
