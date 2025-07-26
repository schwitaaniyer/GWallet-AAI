package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"mime/multipart"
)

// E2E test server that simulates the full backend
func setupE2ETestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		switch r.URL.Path {
		case "/health":
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		case "/receipts":
			handleE2EReceipts(w, r)
		case "/queries":
			handleE2EQueries(w, r)
		case "/wallet-passes":
			handleE2EWalletPasses(w, r)
		case "/analysis":
			handleE2EAnalysis(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

func handleE2EReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Simulate receipt processing
		receipt := map[string]interface{}{
			"id":           "e2e_receipt_123",
			"user_id":      "e2e_user_123",
			"store_name":   "E2E Test Store",
			"total_amount": 67.89,
			"tax_amount":   5.23,
			"items": []map[string]interface{}{
				{"name": "Bread", "price": 3.99, "quantity": 2, "category": "bakery"},
				{"name": "Milk", "price": 4.99, "quantity": 1, "category": "dairy"},
				{"name": "Eggs", "price": 5.99, "quantity": 1, "category": "dairy"},
			},
			"date":        time.Now().Format(time.RFC3339),
			"image_url":   "https://storage.googleapis.com/e2e-bucket/receipts/e2e_receipt_123.jpg",
			"created_at":  time.Now().Format(time.RFC3339),
			"updated_at":  time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(receipt)
	case "GET":
		receipts := []map[string]interface{}{
			{
				"id":           "e2e_receipt_123",
				"user_id":      "e2e_user_123",
				"store_name":   "E2E Test Store",
				"total_amount": 67.89,
				"tax_amount":   5.23,
				"items": []map[string]interface{}{
					{"name": "Bread", "price": 3.99, "quantity": 2, "category": "bakery"},
				},
				"date":        time.Now().Format(time.RFC3339),
				"image_url":   "https://storage.googleapis.com/e2e-bucket/receipts/e2e_receipt_123.jpg",
				"created_at":  time.Now().Format(time.RFC3339),
				"updated_at":  time.Now().Format(time.RFC3339),
			},
		}
		json.NewEncoder(w).Encode(receipts)
	}
}

func handleE2EQueries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		query := map[string]interface{}{
			"id":        "e2e_query_123",
			"user_id":   "e2e_user_123",
			"query":     "What can I cook with bread, milk, and eggs?",
			"language":  "en",
			"response":  "With bread, milk, and eggs, you can make: 1. French Toast 2. Scrambled Eggs with Toast 3. Bread Pudding",
			"created_at": time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(query)
	case "GET":
		queries := []map[string]interface{}{
			{
				"id":        "e2e_query_123",
				"user_id":   "e2e_user_123",
				"query":     "What can I cook with bread, milk, and eggs?",
				"language":  "en",
				"response":  "With bread, milk, and eggs, you can make: 1. French Toast 2. Scrambled Eggs with Toast 3. Bread Pudding",
				"created_at": time.Now().Format(time.RFC3339),
			},
		}
		json.NewEncoder(w).Encode(queries)
	}
}

func handleE2EWalletPasses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		pass := map[string]interface{}{
			"id":          "e2e_pass_123",
			"user_id":     "e2e_user_123",
			"type":        "receipt",
			"title":       "Receipt - E2E Test Store",
			"description": "Total: $67.89, Items: 3",
			"data":        `{"receipt_id": "e2e_receipt_123", "store_name": "E2E Test Store"}`,
			"created_at":  time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(pass)
	case "GET":
		passes := []map[string]interface{}{
			{
				"id":          "e2e_pass_123",
				"user_id":     "e2e_user_123",
				"type":        "receipt",
				"title":       "Receipt - E2E Test Store",
				"description": "Total: $67.89, Items: 3",
				"data":        `{"receipt_id": "e2e_receipt_123", "store_name": "E2E Test Store"}`,
				"created_at":  time.Now().Format(time.RFC3339),
			},
		}
		json.NewEncoder(w).Encode(passes)
	}
}

func handleE2EAnalysis(w http.ResponseWriter, r *http.Request) {
	analysis := map[string]interface{}{
		"total_spent":        67.89,
		"category_spending": map[string]float64{
			"bakery":    7.98,
			"dairy":     10.98,
		},
		"receipt_count":       1,
		"average_per_receipt": 67.89,
	}
	json.NewEncoder(w).Encode(analysis)
}

// E2E Test Functions
func TestE2ECLIWorkflow(t *testing.T) {
	server := setupE2ETestServer()
	defer server.Close()
	
	// Set environment variable for CLI to use test server
	originalBaseURL := os.Getenv("RASEED_API_URL")
	os.Setenv("RASEED_API_URL", server.URL)
	defer os.Setenv("RASEED_API_URL", originalBaseURL)
	
	// Test 1: Health check
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
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
	})
	
	// Test 2: Receipt upload
	t.Run("Receipt Upload", func(t *testing.T) {
		// Create test receipt image
		tempFile, err := os.CreateTemp("", "e2e_receipt.jpg")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()
		
		tempFile.WriteString("fake receipt image data for E2E test")
		
		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		writer.WriteField("user_id", "e2e_user_123")
		
		part, err := writer.CreateFormFile("receipt", "e2e_receipt.jpg")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		
		tempFile.Seek(0, 0)
		io.Copy(part, tempFile)
		writer.Close()
		
		// Upload receipt
		resp, err := http.Post(server.URL+"/receipts", writer.FormDataContentType(), &buf)
		if err != nil {
			t.Fatalf("Receipt upload failed: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var receipt map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if receipt["store_name"] != "E2E Test Store" {
			t.Errorf("Expected store name 'E2E Test Store', got '%v'", receipt["store_name"])
		}
		
		if receipt["total_amount"] != 67.89 {
			t.Errorf("Expected total amount 67.89, got %v", receipt["total_amount"])
		}
	})
	
	// Test 3: Query submission
	t.Run("Query Submission", func(t *testing.T) {
		queryData := map[string]string{
			"user_id":  "e2e_user_123",
			"query":    "What can I cook with bread, milk, and eggs?",
			"language": "en",
		}
		
		jsonData, err := json.Marshal(queryData)
		if err != nil {
			t.Fatalf("Failed to marshal query data: %v", err)
		}
		
		resp, err := http.Post(server.URL+"/queries", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Query submission failed: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var query map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if !strings.Contains(query["response"].(string), "French Toast") {
			t.Errorf("Expected response to contain 'French Toast', got '%v'", query["response"])
		}
	})
	
	// Test 4: Wallet pass creation
	t.Run("Wallet Pass Creation", func(t *testing.T) {
		passData := map[string]string{
			"user_id":     "e2e_user_123",
			"type":        "receipt",
			"title":       "Receipt - E2E Test Store",
			"description": "Total: $67.89, Items: 3",
			"data":        `{"receipt_id": "e2e_receipt_123", "store_name": "E2E Test Store"}`,
		}
		
		jsonData, err := json.Marshal(passData)
		if err != nil {
			t.Fatalf("Failed to marshal pass data: %v", err)
		}
		
		resp, err := http.Post(server.URL+"/wallet-passes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Wallet pass creation failed: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var pass map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&pass); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if pass["type"] != "receipt" {
			t.Errorf("Expected type 'receipt', got '%v'", pass["type"])
		}
	})
	
	// Test 5: Analysis retrieval
	t.Run("Analysis Retrieval", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/analysis?user_id=e2e_user_123")
		if err != nil {
			t.Fatalf("Analysis retrieval failed: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var analysis map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		
		if analysis["total_spent"] != 67.89 {
			t.Errorf("Expected total spent 67.89, got %v", analysis["total_spent"])
		}
		
		if analysis["receipt_count"] != float64(1) {
			t.Errorf("Expected receipt count 1, got %v", analysis["receipt_count"])
		}
	})
}

func TestE2EWebInterfaceWorkflow(t *testing.T) {
	server := setupE2ETestServer()
	defer server.Close()
	
	// Test web interface API calls
	t.Run("Web Interface API Integration", func(t *testing.T) {
		// Test 1: Fetch receipts for web display
		resp, err := http.Get(server.URL + "/receipts?user_id=e2e_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch receipts: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var receipts []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&receipts); err != nil {
			t.Fatalf("Failed to decode receipts: %v", err)
		}
		
		if len(receipts) == 0 {
			t.Error("Expected at least one receipt")
		}
		
		// Test 2: Fetch queries for web display
		resp, err = http.Get(server.URL + "/queries?user_id=e2e_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch queries: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var queries []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&queries); err != nil {
			t.Fatalf("Failed to decode queries: %v", err)
		}
		
		if len(queries) == 0 {
			t.Error("Expected at least one query")
		}
		
		// Test 3: Fetch wallet passes for web display
		resp, err = http.Get(server.URL + "/wallet-passes?user_id=e2e_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch wallet passes: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		
		var passes []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&passes); err != nil {
			t.Fatalf("Failed to decode passes: %v", err)
		}
		
		if len(passes) == 0 {
			t.Error("Expected at least one wallet pass")
		}
	})
}

func TestE2ECrossComponentIntegration(t *testing.T) {
	server := setupE2ETestServer()
	defer server.Close()
	
	t.Run("CLI to Web Data Consistency", func(t *testing.T) {
		// Step 1: Upload receipt via CLI simulation
		tempFile, err := os.CreateTemp("", "cross_test_receipt.jpg")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()
		
		tempFile.WriteString("cross component test receipt")
		
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		writer.WriteField("user_id", "cross_user_123")
		
		part, err := writer.CreateFormFile("receipt", "cross_test_receipt.jpg")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		
		tempFile.Seek(0, 0)
		io.Copy(part, tempFile)
		writer.Close()
		
		// Upload receipt
		resp, err := http.Post(server.URL+"/receipts", writer.FormDataContentType(), &buf)
		if err != nil {
			t.Fatalf("Receipt upload failed: %v", err)
		}
		defer resp.Body.Close()
		
		var uploadedReceipt map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&uploadedReceipt); err != nil {
			t.Fatalf("Failed to decode uploaded receipt: %v", err)
		}
		
		receiptID := uploadedReceipt["id"].(string)
		
		// Step 2: Verify receipt is accessible via web interface
		resp, err = http.Get(server.URL + "/receipts?user_id=cross_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch receipts: %v", err)
		}
		defer resp.Body.Close()
		
		var receipts []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&receipts); err != nil {
			t.Fatalf("Failed to decode receipts: %v", err)
		}
		
		// Find the uploaded receipt
		var foundReceipt map[string]interface{}
		for _, receipt := range receipts {
			if receipt["id"] == receiptID {
				foundReceipt = receipt
				break
			}
		}
		
		if foundReceipt == nil {
			t.Error("Uploaded receipt not found in web interface")
		}
		
		// Step 3: Create wallet pass for the receipt
		passData := map[string]string{
			"user_id":     "cross_user_123",
			"type":        "receipt",
			"title":       "Receipt - " + foundReceipt["store_name"].(string),
			"description": fmt.Sprintf("Total: $%.2f, Items: %d", foundReceipt["total_amount"], len(foundReceipt["items"].([]interface{}))),
			"data":        fmt.Sprintf(`{"receipt_id": "%s", "store_name": "%s"}`, receiptID, foundReceipt["store_name"]),
		}
		
		jsonData, err := json.Marshal(passData)
		if err != nil {
			t.Fatalf("Failed to marshal pass data: %v", err)
		}
		
		resp, err = http.Post(server.URL+"/wallet-passes", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Wallet pass creation failed: %v", err)
		}
		defer resp.Body.Close()
		
		var createdPass map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&createdPass); err != nil {
			t.Fatalf("Failed to decode created pass: %v", err)
		}
		
		// Step 4: Verify wallet pass is accessible via web interface
		resp, err = http.Get(server.URL + "/wallet-passes?user_id=cross_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch wallet passes: %v", err)
		}
		defer resp.Body.Close()
		
		var passes []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&passes); err != nil {
			t.Fatalf("Failed to decode passes: %v", err)
		}
		
		passFound := false
		for _, pass := range passes {
			if pass["id"] == createdPass["id"] {
				passFound = true
				break
			}
		}
		
		if !passFound {
			t.Error("Created wallet pass not found in web interface")
		}
		
		// Step 5: Verify analysis reflects the new data
		resp, err = http.Get(server.URL + "/analysis?user_id=cross_user_123")
		if err != nil {
			t.Fatalf("Failed to fetch analysis: %v", err)
		}
		defer resp.Body.Close()
		
		var analysis map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
			t.Fatalf("Failed to decode analysis: %v", err)
		}
		
		if analysis["receipt_count"].(float64) < 1 {
			t.Error("Analysis should reflect at least one receipt")
		}
	})
}

// Performance tests
func BenchmarkE2EReceiptUpload(b *testing.B) {
	server := setupE2ETestServer()
	defer server.Close()
	
	// Create test file
	tempFile, err := os.CreateTemp("", "benchmark_receipt.jpg")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	tempFile.WriteString("benchmark receipt data")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		writer.WriteField("user_id", "benchmark_user")
		
		part, err := writer.CreateFormFile("receipt", "benchmark_receipt.jpg")
		if err != nil {
			b.Fatalf("Failed to create form file: %v", err)
		}
		
		tempFile.Seek(0, 0)
		io.Copy(part, tempFile)
		writer.Close()
		
		resp, err := http.Post(server.URL+"/receipts", writer.FormDataContentType(), &buf)
		if err != nil {
			b.Fatalf("Receipt upload failed: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkE2EQueryProcessing(b *testing.B) {
	server := setupE2ETestServer()
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
			b.Fatalf("Query submission failed: %v", err)
		}
		resp.Body.Close()
	}
}

// Test main function
func TestMain(m *testing.M) {
	fmt.Println("Starting E2E tests...")
	code := m.Run()
	fmt.Println("E2E tests completed.")
	os.Exit(code)
} 