package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
)

// StockManagementEvent represents the event data from Pub/Sub
type StockManagementEvent struct {
	ItemID string `json:"item_id"`
	UserID string `json:"user_id"`
	Action string `json:"action"` // created, updated, deleted
	Status string `json:"status"` // fresh, expiring_soon, expired
}

// StockItem represents a stock item in inventory
type StockItem struct {
	ID           string    `json:"id" firestore:"id"`
	UserID       string    `json:"user_id" firestore:"user_id"`
	Name         string    `json:"name" firestore:"name"`
	Category     string    `json:"category" firestore:"category"`
	Quantity     int       `json:"quantity" firestore:"quantity"`
	Unit         string    `json:"unit" firestore:"unit"`
	PurchaseDate time.Time `json:"purchase_date" firestore:"purchase_date"`
	ExpiryDate   time.Time `json:"expiry_date" firestore:"expiry_date"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" firestore:"updated_at"`
}

var (
	firestoreClient *firestore.Client
	pubsubClient    *pubsub.Client
)

func init() {
	ctx := context.Background()
	
	// Initialize Firestore client
	var err error
	firestoreClient, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	// Initialize Pub/Sub client
	pubsubClient, err = pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
}

// ProcessStockManagement is the Cloud Function entry point
func ProcessStockManagement(ctx context.Context, msg pubsub.Message) error {
	var event StockManagementEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	log.Printf("Processing stock management event for item %s, user %s, action: %s", event.ItemID, event.UserID, event.Action)

	switch event.Action {
	case "created":
		return handleItemCreated(ctx, event)
	case "updated":
		return handleItemUpdated(ctx, event)
	case "deleted":
		return handleItemDeleted(ctx, event)
	default:
		log.Printf("Unknown action: %s", event.Action)
		return nil
	}
}

func handleItemCreated(ctx context.Context, event StockManagementEvent) error {
	// Get the created item
	doc, err := firestoreClient.Collection("stock_items").Doc(event.ItemID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get item %s: %v", event.ItemID, err)
		return err
	}

	var item StockItem
	if err := doc.DataTo(&item); err != nil {
		log.Printf("Failed to parse item %s: %v", event.ItemID, err)
		return err
	}

	// Send notification if item is expiring soon or expired
	if item.Status == "expiring_soon" || item.Status == "expired" {
		err = sendExpiryNotification(ctx, item)
		if err != nil {
			log.Printf("Failed to send expiry notification: %v", err)
			return err
		}
	}

	// Create wallet pass for the item if it's perishable
	if isPerishable(item.Category) {
		err = createStockItemWalletPass(ctx, item)
		if err != nil {
			log.Printf("Failed to create wallet pass: %v", err)
			return err
		}
	}

	log.Printf("Successfully processed item creation for %s", event.ItemID)
	return nil
}

func handleItemUpdated(ctx context.Context, event StockManagementEvent) error {
	// Get the updated item
	doc, err := firestoreClient.Collection("stock_items").Doc(event.ItemID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get item %s: %v", event.ItemID, err)
		return err
	}

	var item StockItem
	if err := doc.DataTo(&item); err != nil {
		log.Printf("Failed to parse item %s: %v", event.ItemID, err)
		return err
	}

	// Send notification if status changed to expiring soon or expired
	if event.Status == "expiring_soon" || event.Status == "expired" {
		err = sendExpiryNotification(ctx, item)
		if err != nil {
			log.Printf("Failed to send expiry notification: %v", err)
			return err
		}
	}

	// Update wallet pass if needed
	if isPerishable(item.Category) {
		err = updateStockItemWalletPass(ctx, item)
		if err != nil {
			log.Printf("Failed to update wallet pass: %v", err)
			return err
		}
	}

	log.Printf("Successfully processed item update for %s", event.ItemID)
	return nil
}

func handleItemDeleted(ctx context.Context, event StockManagementEvent) error {
	// Delete associated wallet pass if exists
	err := deleteStockItemWalletPass(ctx, event.ItemID)
	if err != nil {
		log.Printf("Failed to delete wallet pass: %v", err)
		return err
	}

	log.Printf("Successfully processed item deletion for %s", event.ItemID)
	return nil
}

func sendExpiryNotification(ctx context.Context, item StockItem) error {
	// Create notification event
	notificationData := map[string]interface{}{
		"user_id":  item.UserID,
		"type":     "stock_expiry",
		"title":    "Item Expiry Alert",
		"message":  fmt.Sprintf("%s is %s", item.Name, item.Status),
		"data": map[string]interface{}{
			"item_id":     item.ID,
			"item_name":   item.Name,
			"status":      item.Status,
			"expiry_date": item.ExpiryDate,
		},
	}

	// Publish notification event
	topic := pubsubClient.Topic("notification-events")
	msgData, err := json.Marshal(notificationData)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %v", err)
	}

	msg := &pubsub.Message{Data: msgData}
	topic.Publish(ctx, msg)

	return nil
}

func createStockItemWalletPass(ctx context.Context, item StockItem) error {
	// Create wallet pass data
	passData := map[string]interface{}{
		"item_id":      item.ID,
		"name":         item.Name,
		"category":     item.Category,
		"quantity":     item.Quantity,
		"unit":         item.Unit,
		"expiry_date":  item.ExpiryDate,
		"status":       item.Status,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	// Create wallet pass document
	pass := map[string]interface{}{
		"id":          fmt.Sprintf("stock_%s", item.ID),
		"user_id":     item.UserID,
		"type":        "stock_item",
		"title":       fmt.Sprintf("Stock - %s", item.Name),
		"description": fmt.Sprintf("Quantity: %d %s, Expires: %s", item.Quantity, item.Unit, item.ExpiryDate.Format("2006-01-02")),
		"data":        string(passDataJSON),
		"created_at":  firestore.ServerTimestamp,
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("stock_%s", item.ID)).Set(ctx, pass)
	return err
}

func updateStockItemWalletPass(ctx context.Context, item StockItem) error {
	// Update existing wallet pass
	passData := map[string]interface{}{
		"item_id":      item.ID,
		"name":         item.Name,
		"category":     item.Category,
		"quantity":     item.Quantity,
		"unit":         item.Unit,
		"expiry_date":  item.ExpiryDate,
		"status":       item.Status,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("stock_%s", item.ID)).Update(ctx, []firestore.Update{
		{Path: "title", Value: fmt.Sprintf("Stock - %s", item.Name)},
		{Path: "description", Value: fmt.Sprintf("Quantity: %d %s, Expires: %s", item.Quantity, item.Unit, item.ExpiryDate.Format("2006-01-02"))},
		{Path: "data", Value: string(passDataJSON)},
	})

	return err
}

func deleteStockItemWalletPass(ctx context.Context, itemID string) error {
	// Delete wallet pass if exists
	_, err := firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("stock_%s", itemID)).Delete(ctx)
	if err != nil {
		// Ignore not found errors
		return nil
	}
	return nil
}

func isPerishable(category string) bool {
	perishableCategories := []string{
		"dairy", "produce", "meat", "seafood", "bakery", "frozen", "beverages",
		"dairy_products", "fruits", "vegetables", "meat_products", "fish",
	}
	
	for _, cat := range perishableCategories {
		if category == cat {
			return true
		}
	}
	return false
} 