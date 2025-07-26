package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

// QueryProcessingEvent represents the event data from Pub/Sub
type QueryProcessingEvent struct {
	QueryID  string `json:"query_id"`
	UserID   string `json:"user_id"`
	Query    string `json:"query"`
	Language string `json:"language"`
}

// QueryResponse represents the AI-generated response
type QueryResponse struct {
	Response    string                 `json:"response"`
	Intent      string                 `json:"intent"`
	Confidence  float64                `json:"confidence"`
	Suggestions []string               `json:"suggestions"`
	Data        map[string]interface{} `json:"data"`
}

var (
	firestoreClient *firestore.Client
	vertexClient    *genai.Client
)

func init() {
	ctx := context.Background()
	
	// Initialize Firestore client
	var err error
	firestoreClient, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	// Initialize Vertex AI client
	vertexClient, err = genai.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"), option.WithLocation("us-central1"))
	if err != nil {
		log.Fatalf("Failed to create Vertex AI client: %v", err)
	}
}

// ProcessQuery is the Cloud Function entry point
func ProcessQuery(ctx context.Context, msg pubsub.Message) error {
	var event QueryProcessingEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	log.Printf("Processing query %s for user %s: %s", event.QueryID, event.UserID, event.Query)

	// Get user's receipt data for context
	userReceipts, err := getUserReceipts(ctx, event.UserID)
	if err != nil {
		log.Printf("Failed to get user receipts: %v", err)
		return err
	}

	// Process query with AI
	response, err := processQueryWithAI(ctx, event.Query, event.Language, userReceipts)
	if err != nil {
		log.Printf("Failed to process query with AI: %v", err)
		return err
	}

	// Update query document with response
	err = updateQueryDocument(ctx, event.QueryID, response)
	if err != nil {
		log.Printf("Failed to update query document: %v", err)
		return err
	}

	// Create wallet pass if needed
	if shouldCreateWalletPass(response.Intent) {
		err = createQueryWalletPass(ctx, event.UserID, event.QueryID, response)
		if err != nil {
			log.Printf("Failed to create wallet pass: %v", err)
			return err
		}
	}

	log.Printf("Successfully processed query %s", event.QueryID)
	return nil
}

func getUserReceipts(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	iter := firestoreClient.Collection("receipts").Where("user_id", "==", userID).Documents(ctx)
	var receipts []map[string]interface{}
	
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		receipts = append(receipts, doc.Data())
	}

	return receipts, nil
}

func processQueryWithAI(ctx context.Context, query, language string, receipts []map[string]interface{}) (*QueryResponse, error) {
	model := vertexClient.GenerativeModel("gemini-pro")
	
	// Create context from user's receipts
	receiptContext := createReceiptContext(receipts)
	
	prompt := fmt.Sprintf(`You are Raseed, an AI-powered personal assistant for financial management and receipt analysis. 

User's Receipt History:
%s

User Query (Language: %s): %s

Please analyze this query and provide a helpful response. Consider the user's spending patterns, recent purchases, and financial context.

Respond in JSON format:
{
	"response": "Your helpful response to the user",
	"intent": "cooking_suggestion|spending_analysis|shopping_list|financial_insight|general_help",
	"confidence": 0.95,
	"suggestions": ["suggestion1", "suggestion2"],
	"data": {
		"relevant_items": ["item1", "item2"],
		"total_spent": 0.00,
		"category_breakdown": {"category": "amount"}
	}
}

Focus on being helpful, actionable, and personalized based on the user's receipt history.`, receiptContext, language, query)

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// Parse the response
	var queryResponse QueryResponse
	responseText := resp.Candidates[0].Content.Parts[0].Text()
	
	// Clean the response (remove markdown if present)
	cleanResponse := responseText
	if strings.Contains(cleanResponse, "```json") {
		start := strings.Index(cleanResponse, "```json") + 7
		end := strings.LastIndex(cleanResponse, "```")
		if end > start {
			cleanResponse = cleanResponse[start:end]
		}
	}

	if err := json.Unmarshal([]byte(cleanResponse), &queryResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	return &queryResponse, nil
}

func createReceiptContext(receipts []map[string]interface{}) string {
	if len(receipts) == 0 {
		return "No receipt history available."
	}

	context := "Recent Receipts:\n"
	for i, receipt := range receipts {
		if i >= 10 { // Limit to last 10 receipts
			break
		}
		
		storeName, _ := receipt["store_name"].(string)
		totalAmount, _ := receipt["total_amount"].(float64)
		date, _ := receipt["date"].(string)
		
		context += fmt.Sprintf("- %s: $%.2f on %s\n", storeName, totalAmount, date)
		
		// Add items if available
		if items, ok := receipt["items"].([]interface{}); ok {
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					name, _ := itemMap["name"].(string)
					category, _ := itemMap["category"].(string)
					context += fmt.Sprintf("  * %s (%s)\n", name, category)
				}
			}
		}
	}
	
	return context
}

func updateQueryDocument(ctx context.Context, queryID string, response *QueryResponse) error {
	_, err := firestoreClient.Collection("queries").Doc(queryID).Update(ctx, []firestore.Update{
		{Path: "response", Value: response.Response},
		{Path: "updated_at", Value: firestore.ServerTimestamp},
	})

	return err
}

func shouldCreateWalletPass(intent string) bool {
	walletPassIntents := []string{"cooking_suggestion", "shopping_list", "financial_insight"}
	for _, validIntent := range walletPassIntents {
		if intent == validIntent {
			return true
		}
	}
	return false
}

func createQueryWalletPass(ctx context.Context, userID, queryID string, response *QueryResponse) error {
	// Create wallet pass data
	passData := map[string]interface{}{
		"query_id":    queryID,
		"intent":      response.Intent,
		"suggestions": response.Suggestions,
		"data":        response.Data,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	// Determine pass type and title based on intent
	passType := "insight"
	title := "Financial Insight"
	
	switch response.Intent {
	case "cooking_suggestion":
		passType = "cooking"
		title = "Cooking Suggestions"
	case "shopping_list":
		passType = "shopping"
		title = "Shopping List"
	case "financial_insight":
		passType = "insight"
		title = "Financial Insight"
	}

	// Create wallet pass document
	pass := map[string]interface{}{
		"id":          fmt.Sprintf("query_%s", queryID),
		"user_id":     userID,
		"type":        passType,
		"title":       title,
		"description": response.Response[:100] + "...", // Truncate if too long
		"data":        string(passDataJSON),
		"created_at":  firestore.ServerTimestamp,
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("query_%s", queryID)).Set(ctx, pass)
	return err
} 