package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Receipt represents a receipt document in Firestore
type Receipt struct {
	ID          string    `json:"id" firestore:"id"`
	UserID      string    `json:"user_id" firestore:"user_id"`
	StoreName   string    `json:"store_name" firestore:"store_name"`
	TotalAmount float64   `json:"total_amount" firestore:"total_amount"`
	TaxAmount   float64   `json:"tax_amount" firestore:"tax_amount"`
	Items       []Item    `json:"items" firestore:"items"`
	Date        time.Time `json:"date" firestore:"date"`
	ImageURL    string    `json:"image_url" firestore:"image_url"`
	Location    Location  `json:"location" firestore:"location"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
}

// Item represents an item in a receipt
type Item struct {
	Name     string  `json:"name" firestore:"name"`
	Price    float64 `json:"price" firestore:"price"`
	Quantity int     `json:"quantity" firestore:"quantity"`
	Category string  `json:"category" firestore:"category"`
}

// Location represents store location
type Location struct {
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
	Address   string  `json:"address" firestore:"address"`
}

// Query represents a user query
type Query struct {
	ID        string    `json:"id" firestore:"id"`
	UserID    string    `json:"user_id" firestore:"user_id"`
	Query     string    `json:"query" firestore:"query"`
	Language  string    `json:"language" firestore:"language"`
	Response  string    `json:"response" firestore:"response"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}

// WalletPass represents a Google Wallet pass
type WalletPass struct {
	ID          string    `json:"id" firestore:"id"`
	UserID      string    `json:"user_id" firestore:"user_id"`
	Type        string    `json:"type" firestore:"type"` // receipt, shopping_list, insight
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Data        string    `json:"data" firestore:"data"` // JSON string
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
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
	Status       string    `json:"status" firestore:"status"` // fresh, expiring_soon, expired
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" firestore:"updated_at"`
}

// Global clients
var (
	firestoreClient *firestore.Client
	pubsubClient    *pubsub.Client
	storageClient   *storage.Client
)

func main() {
	ctx := context.Background()

	// Initialize Firestore
	var err error
	firestoreClient, err = firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// Initialize Pub/Sub
	pubsubClient, err = pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer pubsubClient.Close()

	// Initialize Cloud Storage
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Storage client: %v", err)
	}
	defer storageClient.Close()

	// Set up HTTP routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/receipts", receiptsHandler)
	http.HandleFunc("/queries", queriesHandler)
	http.HandleFunc("/wallet-passes", walletPassesHandler)
	http.HandleFunc("/analysis", analysisHandler)
	http.HandleFunc("/stock-items", stockItemsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func receiptsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		uploadReceipt(w, r)
	case "GET":
		getReceipts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func uploadReceipt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get user ID from form
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("receipt")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload to Cloud Storage
	bucketName := os.Getenv("CLOUD_STORAGE_BUCKET")
	bucket := storageClient.Bucket(bucketName)
	
	objectName := fmt.Sprintf("receipts/%s/%s", userID, header.Filename)
	obj := bucket.Object(objectName)
	writer := obj.NewWriter(ctx)
	
	if _, err := io.Copy(writer, file); err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	writer.Close()

	// Get public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	// Create receipt document
	receipt := Receipt{
		ID:        generateID(),
		UserID:    userID,
		ImageURL:  imageURL,
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to Firestore
	_, err = firestoreClient.Collection("receipts").Doc(receipt.ID).Set(ctx, receipt)
	if err != nil {
		http.Error(w, "Failed to save receipt", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub for AI processing
	topic := pubsubClient.Topic("receipt-processing")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"receipt_id": "%s", "user_id": "%s", "image_url": "%s"}`, receipt.ID, userID, imageURL)),
	}
	topic.Publish(ctx, msg)

	json.NewEncoder(w).Encode(receipt)
}

func getReceipts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	iter := firestoreClient.Collection("receipts").Where("user_id", "==", userID).Documents(ctx)
	var receipts []Receipt
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to fetch receipts", http.StatusInternalServerError)
			return
		}

		var receipt Receipt
		if err := doc.DataTo(&receipt); err != nil {
			continue
		}
		receipts = append(receipts, receipt)
	}

	json.NewEncoder(w).Encode(receipts)
}

func queriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		processQuery(w, r)
	case "GET":
		getQueries(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func processQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID   string `json:"user_id"`
		Query    string `json:"query"`
		Language string `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Query == "" {
		http.Error(w, "user_id and query are required", http.StatusBadRequest)
		return
	}

	// Create query document
	query := Query{
		ID:        generateID(),
		UserID:    req.UserID,
		Query:     req.Query,
		Language:  req.Language,
		CreatedAt: time.Now(),
	}

	// Save to Firestore
	_, err := firestoreClient.Collection("queries").Doc(query.ID).Set(ctx, query)
	if err != nil {
		http.Error(w, "Failed to save query", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub for AI processing
	topic := pubsubClient.Topic("query-processing")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"query_id": "%s", "user_id": "%s", "query": "%s", "language": "%s"}`, query.ID, req.UserID, req.Query, req.Language)),
	}
	topic.Publish(ctx, msg)

	json.NewEncoder(w).Encode(query)
}

func getQueries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	iter := firestoreClient.Collection("queries").Where("user_id", "==", userID).Documents(ctx)
	var queries []Query
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to fetch queries", http.StatusInternalServerError)
			return
		}

		var query Query
		if err := doc.DataTo(&query); err != nil {
			continue
		}
		queries = append(queries, query)
	}

	json.NewEncoder(w).Encode(queries)
}

func walletPassesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		createWalletPass(w, r)
	case "GET":
		getWalletPasses(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createWalletPass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID      string `json:"user_id"`
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Data        string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Type == "" || req.Title == "" {
		http.Error(w, "user_id, type, and title are required", http.StatusBadRequest)
		return
	}

	// Create wallet pass
	pass := WalletPass{
		ID:          generateID(),
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Data:        req.Data,
		CreatedAt:   time.Now(),
	}

	// Save to Firestore
	_, err := firestoreClient.Collection("wallet_passes").Doc(pass.ID).Set(ctx, pass)
	if err != nil {
		http.Error(w, "Failed to save wallet pass", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub for Google Wallet API integration
	topic := pubsubClient.Topic("wallet-pass-creation")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"pass_id": "%s", "user_id": "%s", "type": "%s", "title": "%s"}`, pass.ID, req.UserID, req.Type, req.Title)),
	}
	topic.Publish(ctx, msg)

	json.NewEncoder(w).Encode(pass)
}

func getWalletPasses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	iter := firestoreClient.Collection("wallet_passes").Where("user_id", "==", userID).Documents(ctx)
	var passes []WalletPass
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to fetch wallet passes", http.StatusInternalServerError)
			return
		}

		var pass WalletPass
		if err := doc.DataTo(&pass); err != nil {
			continue
		}
		passes = append(passes, pass)
	}

	json.NewEncoder(w).Encode(passes)
}

func analysisHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		getSpendingAnalysis(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getSpendingAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Get user's receipts
	iter := firestoreClient.Collection("receipts").Where("user_id", "==", userID).Documents(ctx)
	var receipts []Receipt
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to fetch receipts", http.StatusInternalServerError)
			return
		}

		var receipt Receipt
		if err := doc.DataTo(&receipt); err != nil {
			continue
		}
		receipts = append(receipts, receipt)
	}

	// Calculate basic analytics
	totalSpent := 0.0
	categorySpending := make(map[string]float64)
	
	for _, receipt := range receipts {
		totalSpent += receipt.TotalAmount
		for _, item := range receipt.Items {
			categorySpending[item.Category] += item.Price * float64(item.Quantity)
		}
	}

	analysis := map[string]interface{}{
		"total_spent":        totalSpent,
		"category_spending":  categorySpending,
		"receipt_count":      len(receipts),
		"average_per_receipt": totalSpent / float64(len(receipts)),
	}

	json.NewEncoder(w).Encode(analysis)
}

func stockItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		createStockItem(w, r)
	case "GET":
		getStockItems(w, r)
	case "PUT":
		updateStockItem(w, r)
	case "DELETE":
		deleteStockItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID       string    `json:"user_id"`
		Name         string    `json:"name"`
		Category     string    `json:"category"`
		Quantity     int       `json:"quantity"`
		Unit         string    `json:"unit"`
		PurchaseDate time.Time `json:"purchase_date"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Name == "" {
		http.Error(w, "user_id and name are required", http.StatusBadRequest)
		return
	}

	// Determine status based on expiry date
	status := "fresh"
	now := time.Now()
	if req.ExpiryDate.Before(now) {
		status = "expired"
	} else if req.ExpiryDate.Sub(now) < 7*24*time.Hour { // 7 days
		status = "expiring_soon"
	}

	// Create stock item
	item := StockItem{
		ID:           generateID(),
		UserID:       req.UserID,
		Name:         req.Name,
		Category:     req.Category,
		Quantity:     req.Quantity,
		Unit:         req.Unit,
		PurchaseDate: req.PurchaseDate,
		ExpiryDate:   req.ExpiryDate,
		Status:       status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to Firestore
	_, err := firestoreClient.Collection("stock_items").Doc(item.ID).Set(ctx, item)
	if err != nil {
		http.Error(w, "Failed to save stock item", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub for stock management processing
	topic := pubsubClient.Topic("stock-management")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"item_id": "%s", "user_id": "%s", "action": "created", "status": "%s"}`, item.ID, req.UserID, status)),
	}
	topic.Publish(ctx, msg)

	json.NewEncoder(w).Encode(item)
}

func getStockItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Get status filter if provided
	status := r.URL.Query().Get("status")

	var iter *firestore.DocumentIterator
	if status != "" {
		iter = firestoreClient.Collection("stock_items").Where("user_id", "==", userID).Where("status", "==", status).Documents(ctx)
	} else {
		iter = firestoreClient.Collection("stock_items").Where("user_id", "==", userID).Documents(ctx)
	}

	var items []StockItem
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to fetch stock items", http.StatusInternalServerError)
			return
		}

		var item StockItem
		if err := doc.DataTo(&item); err != nil {
			continue
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func updateStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name       string    `json:"name"`
		Category   string    `json:"category"`
		Quantity   int       `json:"quantity"`
		Unit       string    `json:"unit"`
		ExpiryDate time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing item
	doc, err := firestoreClient.Collection("stock_items").Doc(itemID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Stock item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch stock item", http.StatusInternalServerError)
		}
		return
	}

	var item StockItem
	if err := doc.DataTo(&item); err != nil {
		http.Error(w, "Failed to parse stock item", http.StatusInternalServerError)
		return
	}

	// Update fields
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Category != "" {
		item.Category = req.Category
	}
	if req.Quantity > 0 {
		item.Quantity = req.Quantity
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if !req.ExpiryDate.IsZero() {
		item.ExpiryDate = req.ExpiryDate
	}

	// Update status based on new expiry date
	now := time.Now()
	if item.ExpiryDate.Before(now) {
		item.Status = "expired"
	} else if item.ExpiryDate.Sub(now) < 7*24*time.Hour {
		item.Status = "expiring_soon"
	} else {
		item.Status = "fresh"
	}

	item.UpdatedAt = time.Now()

	// Save updated item
	_, err = firestoreClient.Collection("stock_items").Doc(itemID).Set(ctx, item)
	if err != nil {
		http.Error(w, "Failed to update stock item", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub
	topic := pubsubClient.Topic("stock-management")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"item_id": "%s", "user_id": "%s", "action": "updated", "status": "%s"}`, itemID, item.UserID, item.Status)),
	}
	topic.Publish(ctx, msg)

	json.NewEncoder(w).Encode(item)
}

func deleteStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Get item to get user_id for event
	doc, err := firestoreClient.Collection("stock_items").Doc(itemID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Stock item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch stock item", http.StatusInternalServerError)
		}
		return
	}

	var item StockItem
	if err := doc.DataTo(&item); err != nil {
		http.Error(w, "Failed to parse stock item", http.StatusInternalServerError)
		return
	}

	// Delete item
	_, err = firestoreClient.Collection("stock_items").Doc(itemID).Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete stock item", http.StatusInternalServerError)
		return
	}

	// Publish event to Pub/Sub
	topic := pubsubClient.Topic("stock-management")
	msg := &pubsub.Message{
		Data: []byte(fmt.Sprintf(`{"item_id": "%s", "user_id": "%s", "action": "deleted"}`, itemID, item.UserID)),
	}
	topic.Publish(ctx, msg)

	w.WriteHeader(http.StatusNoContent)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
} 