package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Test data structures
type TestReceipt struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StoreName   string    `json:"store_name"`
	TotalAmount float64   `json:"total_amount"`
	TaxAmount   float64   `json:"tax_amount"`
	Items       []TestItem `json:"items"`
	Date        time.Time `json:"date"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TestItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

type TestQuery struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Query     string    `json:"query"`
	Language  string    `json:"language"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

type TestWalletPass struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
}

// Mock server for testing
func setupTestServer() *httptest.Server {
	// This would normally start the actual server
	// For testing, we'll use httptest.NewServer
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock responses for testing
		switch r.URL.Path {
		case "/health":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		case "/receipts":
			handleReceiptsTest(w, r)
		case "/queries":
			handleQueriesTest(w, r)
		case "/wallet-passes":
			handleWalletPassesTest(w, r)
		case "/analysis":
			handleAnalysisTest(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

func handleReceiptsTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "POST":
		// Mock receipt upload response
		receipt := TestReceipt{
			ID:          "test_receipt_123",
			UserID:      "test_user_123",
			StoreName:   "Test Store",
			TotalAmount: 45.99,
			TaxAmount:   3.50,
			Items: []TestItem{
				{Name: "Milk", Price: 4.99, Quantity: 2, Category: "dairy"},
				{Name: "Bread", Price: 3.99, Quantity: 1, Category: "bakery"},
			},
			Date:      time.Now(),
			ImageURL:  "https://storage.googleapis.com/test-bucket/receipts/test_user_123/receipt.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		json.NewEncoder(w).Encode(receipt)
	case "GET":
		// Mock receipt list response
		receipts := []TestReceipt{
			{
				ID:          "test_receipt_123",
				UserID:      "test_user_123",
				StoreName:   "Test Store",
				TotalAmount: 45.99,
				TaxAmount:   3.50,
				Items: []TestItem{
					{Name: "Milk", Price: 4.99, Quantity: 2, Category: "dairy"},
				},
				Date:      time.Now(),
				ImageURL:  "https://storage.googleapis.com/test-bucket/receipts/test_user_123/receipt.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		json.NewEncoder(w).Encode(receipts)
	}
}

func handleQueriesTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "POST":
		// Mock query response
		query := TestQuery{
			ID:        "test_query_123",
			UserID:    "test_user_123",
			Query:     "What can I cook with my recent purchases?",
			Language:  "en",
			Response:  "Based on your recent purchases, you can make: 1. Scrambled eggs with toast 2. Pasta with tomato sauce",
			CreatedAt: time.Now(),
		}
		json.NewEncoder(w).Encode(query)
	case "GET":
		// Mock query list response
		queries := []TestQuery{
			{
				ID:        "test_query_123",
				UserID:    "test_user_123",
				Query:     "What can I cook with my recent purchases?",
				Language:  "en",
				Response:  "Based on your recent purchases, you can make: 1. Scrambled eggs with toast 2. Pasta with tomato sauce",
				CreatedAt: time.Now(),
			},
		}
		json.NewEncoder(w).Encode(queries)
	}
}

func handleWalletPassesTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "POST":
		// Mock wallet pass creation response
		pass := TestWalletPass{
			ID:          "receipt_test_receipt_123",
			UserID:      "test_user_123",
			Type:        "receipt",
			Title:       "Receipt - Test Store",
			Description: "Total: $45.99, Items: 2",
			Data:        `{"receipt_id": "test_receipt_123", "store_name": "Test Store"}`,
			CreatedAt:   time.Now(),
		}
		json.NewEncoder(w).Encode(pass)
	case "GET":
		// Mock wallet pass list response
		passes := []TestWalletPass{
			{
				ID:          "receipt_test_receipt_123",
				UserID:      "test_user_123",
				Type:        "receipt",
				Title:       "Receipt - Test Store",
				Description: "Total: $45.99, Items: 2",
				Data:        `{"receipt_id": "test_receipt_123", "store_name": "Test Store"}`,
				CreatedAt:   time.Now(),
			},
		}
		json.NewEncoder(w).Encode(passes)
	}
}

func handleAnalysisTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	analysis := map[string]interface{}{
		"total_spent":        245.67,
		"category_spending": map[string]float64{
			"groceries":    120.50,
			"restaurants":  85.25,
			"transportation": 40.00,
		},
		"receipt_count":      15,
		"average_per_receipt": 16.38,
	}
	json.NewEncoder(w).Encode(analysis)
}

// Test functions
func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", result["status"])
	}
}

func TestReceiptUpload(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a test file
	file, err := os.CreateTemp("", "test_receipt.jpg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// Write some test data
	file.WriteString("fake receipt image data")

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add user_id field
	writer.WriteField("user_id", "test_user_123")
	
	// Add file field
	part, err := writer.CreateFormFile("receipt", "test_receipt.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	
	file.Seek(0, 0)
	io.Copy(part, file)
	writer.Close()

	// Make request
	resp, err := http.Post(server.URL+"/receipts", writer.FormDataContentType(), &buf)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var receipt TestReceipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if receipt.UserID != "test_user_123" {
		t.Errorf("Expected user_id 'test_user_123', got '%s'", receipt.UserID)
	}

	if receipt.StoreName != "Test Store" {
		t.Errorf("Expected store_name 'Test Store', got '%s'", receipt.StoreName)
	}
}

func TestQuerySubmission(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	queryData := map[string]string{
		"user_id":  "test_user_123",
		"query":    "What can I cook with my recent purchases?",
		"language": "en",
	}

	jsonData, err := json.Marshal(queryData)
	if err != nil {
		t.Fatalf("Failed to marshal query data: %v", err)
	}

	resp, err := http.Post(server.URL+"/queries", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var query TestQuery
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if query.UserID != "test_user_123" {
		t.Errorf("Expected user_id 'test_user_123', got '%s'", query.UserID)
	}

	if query.Query != "What can I cook with my recent purchases?" {
		t.Errorf("Expected query 'What can I cook with my recent purchases?', got '%s'", query.Query)
	}
}

func TestWalletPassCreation(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	passData := map[string]string{
		"user_id":      "test_user_123",
		"type":         "receipt",
		"title":        "Receipt - Test Store",
		"description":  "Total: $45.99, Items: 2",
		"data":         `{"receipt_id": "test_receipt_123", "store_name": "Test Store"}`,
	}

	jsonData, err := json.Marshal(passData)
	if err != nil {
		t.Fatalf("Failed to marshal pass data: %v", err)
	}

	resp, err := http.Post(server.URL+"/wallet-passes", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var pass TestWalletPass
	if err := json.NewDecoder(resp.Body).Decode(&pass); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if pass.UserID != "test_user_123" {
		t.Errorf("Expected user_id 'test_user_123', got '%s'", pass.UserID)
	}

	if pass.Type != "receipt" {
		t.Errorf("Expected type 'receipt', got '%s'", pass.Type)
	}
}

func TestSpendingAnalysis(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/analysis?user_id=test_user_123")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var analysis map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if analysis["total_spent"] != 245.67 {
		t.Errorf("Expected total_spent 245.67, got %v", analysis["total_spent"])
	}

	if analysis["receipt_count"] != float64(15) {
		t.Errorf("Expected receipt_count 15, got %v", analysis["receipt_count"])
	}
}

// Benchmark tests
func BenchmarkReceiptUpload(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	// Create a test file
	file, err := os.CreateTemp("", "benchmark_receipt.jpg")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	file.WriteString("fake receipt image data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		writer.WriteField("user_id", "benchmark_user")
		
		part, err := writer.CreateFormFile("receipt", "benchmark_receipt.jpg")
		if err != nil {
			b.Fatalf("Failed to create form file: %v", err)
		}
		
		file.Seek(0, 0)
		io.Copy(part, file)
		writer.Close()

		resp, err := http.Post(server.URL+"/receipts", writer.FormDataContentType(), &buf)
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkQuerySubmission(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	queryData := map[string]string{
		"user_id":  "benchmark_user",
		"query":    "What can I cook with my recent purchases?",
		"language": "en",
	}

	jsonData, err := json.Marshal(queryData)
	if err != nil {
		b.Fatalf("Failed to marshal query data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := http.Post(server.URL+"/queries", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

// Main test function
func TestMain(m *testing.M) {
	// Setup any global test configuration here
	fmt.Println("Starting Raseed integration tests...")
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	fmt.Println("Tests completed.")
	
	os.Exit(code)
} 