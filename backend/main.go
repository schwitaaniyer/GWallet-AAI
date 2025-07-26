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

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
} 