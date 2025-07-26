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

// ThirdPartyIntegrationEvent represents the event data from Pub/Sub
type ThirdPartyIntegrationEvent struct {
	UserID       string `json:"user_id"`
	Service      string `json:"service"` // zomato, blinkit, etc.
	Action       string `json:"action"`  // fetch_bills, create_pass, etc.
	ServiceData  string `json:"service_data"`
	RequestedAt  string `json:"requested_at"`
}

// ThirdPartyBill represents a bill from third-party service
type ThirdPartyBill struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Service     string    `json:"service"`
	OrderID     string    `json:"order_id"`
	Restaurant  string    `json:"restaurant"`
	TotalAmount float64   `json:"total_amount"`
	Items       []BillItem `json:"items"`
	OrderDate   time.Time `json:"order_date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// BillItem represents an item in a third-party bill
type BillItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

var firestoreClient *firestore.Client

func init() {
	ctx := context.Background()
	
	// Initialize Firestore client
	var err error
	firestoreClient, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
}

// ProcessThirdPartyIntegration is the Cloud Function entry point
func ProcessThirdPartyIntegration(ctx context.Context, msg pubsub.Message) error {
	var event ThirdPartyIntegrationEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %v", err)
	}

	log.Printf("Processing third-party integration for user %s, service %s, action %s", 
		event.UserID, event.Service, event.Action)

	switch event.Action {
	case "fetch_bills":
		return fetchThirdPartyBills(ctx, event)
	case "create_pass":
		return createThirdPartyWalletPass(ctx, event)
	default:
		return fmt.Errorf("unknown action: %s", event.Action)
	}
}

func fetchThirdPartyBills(ctx context.Context, event ThirdPartyIntegrationEvent) error {
	var bills []ThirdPartyBill

	switch event.Service {
	case "zomato":
		bills = fetchZomatoBills(ctx, event.UserID)
	case "blinkit":
		bills = fetchBlinkitBills(ctx, event.UserID)
	default:
		return fmt.Errorf("unsupported service: %s", event.Service)
	}

	// Save bills to Firestore
	for _, bill := range bills {
		err := saveThirdPartyBill(ctx, bill)
		if err != nil {
			log.Printf("Failed to save bill %s: %v", bill.ID, err)
			continue
		}

		// Create wallet pass for the bill
		err = createThirdPartyBillWalletPass(ctx, bill)
		if err != nil {
			log.Printf("Failed to create wallet pass for bill %s: %v", bill.ID, err)
		}
	}

	log.Printf("Successfully fetched %d bills from %s", len(bills), event.Service)
	return nil
}

func fetchZomatoBills(ctx context.Context, userID string) []ThirdPartyBill {
	// In a real implementation, this would integrate with Zomato's API
	// For now, we'll create mock data
	
	bills := []ThirdPartyBill{
		{
			ID:          fmt.Sprintf("zomato_%d", time.Now().Unix()),
			UserID:      userID,
			Service:     "zomato",
			OrderID:     "ZOM123456",
			Restaurant:  "Pizza Palace",
			TotalAmount: 45.99,
			Items: []BillItem{
				{Name: "Margherita Pizza", Price: 25.99, Quantity: 1, Category: "food"},
				{Name: "Garlic Bread", Price: 8.99, Quantity: 1, Category: "food"},
				{Name: "Coke", Price: 3.99, Quantity: 2, Category: "beverage"},
				{Name: "Delivery Fee", Price: 4.99, Quantity: 1, Category: "service"},
				{Name: "Tax", Price: 2.03, Quantity: 1, Category: "tax"},
			},
			OrderDate: time.Now().Add(-24 * time.Hour),
			Status:    "delivered",
			CreatedAt: time.Now(),
		},
		{
			ID:          fmt.Sprintf("zomato_%d", time.Now().Unix()+1),
			UserID:      userID,
			Service:     "zomato",
			OrderID:     "ZOM123457",
			Restaurant:  "Burger House",
			TotalAmount: 32.50,
			Items: []BillItem{
				{Name: "Chicken Burger", Price: 18.99, Quantity: 1, Category: "food"},
				{Name: "French Fries", Price: 6.99, Quantity: 1, Category: "food"},
				{Name: "Milkshake", Price: 4.99, Quantity: 1, Category: "beverage"},
				{Name: "Delivery Fee", Price: 3.99, Quantity: 1, Category: "service"},
				{Name: "Tax", Price: 1.54, Quantity: 1, Category: "tax"},
			},
			OrderDate: time.Now().Add(-48 * time.Hour),
			Status:    "delivered",
			CreatedAt: time.Now(),
		},
	}

	return bills
}

func fetchBlinkitBills(ctx context.Context, userID string) []ThirdPartyBill {
	// In a real implementation, this would integrate with Blinkit's API
	// For now, we'll create mock data
	
	bills := []ThirdPartyBill{
		{
			ID:          fmt.Sprintf("blinkit_%d", time.Now().Unix()),
			UserID:      userID,
			Service:     "blinkit",
			OrderID:     "BLK789012",
			Restaurant:  "Quick Mart",
			TotalAmount: 67.25,
			Items: []BillItem{
				{Name: "Milk", Price: 4.99, Quantity: 2, Category: "dairy"},
				{Name: "Bread", Price: 3.99, Quantity: 1, Category: "bakery"},
				{Name: "Eggs", Price: 5.99, Quantity: 1, Category: "dairy"},
				{Name: "Bananas", Price: 2.99, Quantity: 1, Category: "fruits"},
				{Name: "Rice", Price: 12.99, Quantity: 1, Category: "grains"},
				{Name: "Tomatoes", Price: 3.99, Quantity: 1, Category: "vegetables"},
				{Name: "Delivery Fee", Price: 2.99, Quantity: 1, Category: "service"},
				{Name: "Tax", Price: 3.32, Quantity: 1, Category: "tax"},
			},
			OrderDate: time.Now().Add(-12 * time.Hour),
			Status:    "delivered",
			CreatedAt: time.Now(),
		},
	}

	return bills
}

func saveThirdPartyBill(ctx context.Context, bill ThirdPartyBill) error {
	_, err := firestoreClient.Collection("third_party_bills").Doc(bill.ID).Set(ctx, bill)
	return err
}

func createThirdPartyBillWalletPass(ctx context.Context, bill ThirdPartyBill) error {
	// Create wallet pass data
	passData := map[string]interface{}{
		"bill_id":      bill.ID,
		"service":      bill.Service,
		"order_id":     bill.OrderID,
		"restaurant":   bill.Restaurant,
		"total_amount": bill.TotalAmount,
		"items_count":  len(bill.Items),
		"order_date":   bill.OrderDate.Format("2006-01-02"),
		"status":       bill.Status,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	// Create wallet pass document
	pass := map[string]interface{}{
		"id":          fmt.Sprintf("bill_%s", bill.ID),
		"user_id":     bill.UserID,
		"type":        "third_party_bill",
		"title":       fmt.Sprintf("%s - %s", bill.Service, bill.Restaurant),
		"description": fmt.Sprintf("Order: %s, Total: $%.2f", bill.OrderID, bill.TotalAmount),
		"data":        string(passDataJSON),
		"created_at":  firestore.ServerTimestamp,
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(fmt.Sprintf("bill_%s", bill.ID)).Set(ctx, pass)
	return err
}

func createThirdPartyWalletPass(ctx context.Context, event ThirdPartyIntegrationEvent) error {
	// Parse service data
	var serviceData map[string]interface{}
	if err := json.Unmarshal([]byte(event.ServiceData), &serviceData); err != nil {
		return fmt.Errorf("failed to parse service data: %v", err)
	}

	// Create wallet pass data
	passData := map[string]interface{}{
		"service":     event.Service,
		"action":      event.Action,
		"service_data": serviceData,
		"requested_at": event.RequestedAt,
	}

	passDataJSON, err := json.Marshal(passData)
	if err != nil {
		return fmt.Errorf("failed to marshal pass data: %v", err)
	}

	// Create wallet pass document
	pass := map[string]interface{}{
		"id":          fmt.Sprintf("integration_%s_%d", event.Service, time.Now().Unix()),
		"user_id":     event.UserID,
		"type":        "third_party_integration",
		"title":       fmt.Sprintf("%s Integration", event.Service),
		"description": fmt.Sprintf("Action: %s", event.Action),
		"data":        string(passDataJSON),
		"created_at":  firestore.ServerTimestamp,
	}

	_, err = firestoreClient.Collection("wallet_passes").Doc(pass["id"].(string)).Set(ctx, pass)
	return err
} 