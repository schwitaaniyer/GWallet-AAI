package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

// ReceiptProcessingEvent represents the event data from Pub/Sub
type ReceiptProcessingEvent struct {
	ReceiptID string `json:"receipt_id"`
	UserID    string `json:"user_id"`
	ImageURL  string `json:"image_url"`
}

// ExtractedReceiptData represents the data extracted from receipt
type ExtractedReceiptData struct {
	StoreName   string  `json:"store_name"`
	TotalAmount float64 `json:"total_amount"`
	TaxAmount   float64 `json:"tax_amount"`
	Items       []Item  `json:"items"`
	Date        string  `json:"date"`
}

// Item represents an item in a receipt
type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
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

// ProcessReceipt is the Cloud Function entry point
func ProcessReceipt(ctx context.Context, msg pubsub.Message) error {
	var event ReceiptProcessingEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	log.Printf("Processing receipt %s for user %s", event.ReceiptID, event.UserID)

	// Extract data from receipt image using Gemini AI
	extractedData, err := extractReceiptData(ctx, event.ImageURL)
	if err != nil {
		log.Printf("Failed to extract receipt data: %v", err)
		return err
	}

	// Update receipt document in Firestore
	err = updateReceiptDocument(ctx, event.ReceiptID, extractedData)
	if err != nil {
		log.Printf("Failed to update receipt document: %v", err)
		return err
	}

	// Create wallet pass for the receipt
	err = createReceiptWalletPass(ctx, event.UserID, event.ReceiptID, extractedData)
	if err != nil {
		log.Printf("Failed to create wallet pass: %v", err)
		return err
	}

	log.Printf("Successfully processed receipt %s", event.ReceiptID)
	return nil
}

func extractReceiptData(ctx context.Context, imageURL string) (*ExtractedReceiptData, error) {
	model := vertexClient.GenerativeModel("gemini-pro-vision")
	
	prompt := `Analyze this receipt image and extract the following information in JSON format:
	{
		"store_name": "Store name",
		"total_amount": 0.00,
		"tax_amount": 0.00,
		"items": [
			{
				"name": "Item name",
				"price": 0.00,
				"quantity": 1,
				"category": "Category (e.g., groceries, electronics, etc.)"
			}
		],
		"date": "YYYY-MM-DD"
	}
	
	Please ensure all monetary values are numbers, quantities are integers, and categorize items appropriately.`

	// Create image part
	img := genai.ImageData{
		MimeType: "image/jpeg",
		Data:     []byte(imageURL), // In production, download the image first
	}

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), img)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// Parse the response
	var extractedData ExtractedReceiptData
	responseText := resp.Candidates[0].Content.Parts[0].Text()
	
	// Clean the response (remove markdown if present)
	cleanResponse := responseText
	if len(cleanResponse) > 0 && cleanResponse[0] == '`' {
		cleanResponse = cleanResponse[3 : len(cleanResponse)-3] // Remove ```json and ```
	}

	if err := json.Unmarshal([]byte(cleanResponse), &extractedData); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	return &extractedData, nil
}

func updateReceiptDocument(ctx context.Context, receiptID string, data *ExtractedReceiptData) error {
	// Convert ExtractedReceiptData to Firestore document
	docData := map[string]interface{}{
		"store_name":   data.StoreName,
		"total_amount": data.TotalAmount,
		"tax_amount":   data.TaxAmount,
		"items":        data.Items,
		"date":         data.Date,
		"updated_at":   firestore.ServerTimestamp,
	}

	_, err := firestoreClient.Collection("receipts").Doc(receiptID).Update(ctx, []firestore.Update{
		{Path: "store_name", Value: data.StoreName},
		{Path: "total_amount", Value: data.TotalAmount},
		{Path: "tax_amount", Value: data.TaxAmount},
		{Path: "items", Value: data.Items},
		{Path: "date", Value: data.Date},
		{Path: "updated_at", Value: firestore.ServerTimestamp},
	})

	return err
}

func createReceiptWalletPass(ctx context.Context, userID, receiptID string, data *ExtractedReceiptData) error {
	// Create wallet pass data
	passData := map[string]interface{}{
		"receipt_id":   receiptID,
		"store_name":   data.StoreName,
		"total_amount": data.TotalAmount,
		"items_count":  len(data.Items),
		"date":         data.Date,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	// Create wallet pass document
	pass := map[string]interface{}{
		"id":          fmt.Sprintf("receipt_%s", receiptID),
		"user_id":     userID,
		"type":        "receipt",
		"title":       fmt.Sprintf("Receipt - %s", data.StoreName),
		"description": fmt.Sprintf("Total: $%.2f, Items: %d", data.TotalAmount, len(data.Items)),
		"data":        string(passDataJSON),
		"created_at":  firestore.ServerTimestamp,
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("receipt_%s", receiptID)).Set(ctx, pass)
	return err
} 