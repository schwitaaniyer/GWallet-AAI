package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// Mock server for CLI testing
func setupCLITestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		switch r.URL.Path {
		case "/health":
			json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		case "/receipts":
			handleCLIReceiptsTest(w, r)
		case "/queries":
			handleCLIQueriesTest(w, r)
		case "/wallet-passes":
			handleCLIWalletPassesTest(w, r)
		case "/analysis":
			handleCLIAnalysisTest(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

func handleCLIReceiptsTest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		receipt := map[string]interface{}{
			"id":           "test_receipt_123",
			"user_id":      "test-user-123",
			"store_name":   "Test Store",
			"total_amount": 45.99,
			"items": []map[string]interface{}{
				{"name": "Milk", "price": 4.99, "quantity": 2},
			},
		}
		json.NewEncoder(w).Encode(receipt)
	} else {
		receipts := []map[string]interface{}{
			{
				"id":           "test_receipt_123",
				"store_name":   "Test Store",
				"total_amount": 45.99,
			},
		}
		json.NewEncoder(w).Encode(receipts)
	}
}

func handleCLIQueriesTest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		query := map[string]interface{}{
			"id":       "test_query_123",
			"user_id":  "test-user-123",
			"query":    "What can I cook?",
			"response": "You can make pasta with your ingredients.",
		}
		json.NewEncoder(w).Encode(query)
	} else {
		queries := []map[string]interface{}{
			{
				"id":       "test_query_123",
				"query":    "What can I cook?",
				"response": "You can make pasta with your ingredients.",
			},
		}
		json.NewEncoder(w).Encode(queries)
	}
}

func handleCLIWalletPassesTest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		pass := map[string]interface{}{
			"id":          "test_pass_123",
			"user_id":     "test-user-123",
			"type":        "receipt",
			"title":       "Test Receipt",
			"description": "Total: $45.99",
		}
		json.NewEncoder(w).Encode(pass)
	} else {
		passes := []map[string]interface{}{
			{
				"id":          "test_pass_123",
				"type":        "receipt",
				"title":       "Test Receipt",
				"description": "Total: $45.99",
			},
		}
		json.NewEncoder(w).Encode(passes)
	}
}

func handleCLIAnalysisTest(w http.ResponseWriter, r *http.Request) {
	analysis := map[string]interface{}{
		"total_spent":        245.67,
		"category_spending": map[string]float64{
			"groceries":     120.50,
			"restaurants":   85.25,
			"transportation": 40.00,
		},
		"receipt_count":       15,
		"average_per_receipt": 16.38,
	}
	json.NewEncoder(w).Encode(analysis)
}

// Test helper functions
func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	
	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	
	err := cmd.Execute()
	return buf.String(), err
}

// Test functions
func TestHealthCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, healthCmd)
	
	if err != nil {
		t.Errorf("Health command failed: %v", err)
	}
	
	if !strings.Contains(output, "healthy") {
		t.Errorf("Expected 'healthy' in output, got: %s", output)
	}
}

func TestUploadReceiptCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_receipt.jpg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	// Write test data
	tempFile.WriteString("fake receipt image data")
	
	output, err := executeCommand(t, uploadReceiptCmd, tempFile.Name())
	
	if err != nil {
		t.Errorf("Upload receipt command failed: %v", err)
	}
	
	if !strings.Contains(output, "Test Store") {
		t.Errorf("Expected 'Test Store' in output, got: %s", output)
	}
}

func TestUploadReceiptCommandFileNotFound(t *testing.T) {
	output, err := executeCommand(t, uploadReceiptCmd, "nonexistent_file.jpg")
	
	if err == nil {
		t.Error("Expected error for non-existent file, but got none")
	}
	
	if !strings.Contains(output, "does not exist") {
		t.Errorf("Expected 'does not exist' in output, got: %s", output)
	}
}

func TestQueryCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, queryCmd, "What can I cook?")
	
	if err != nil {
		t.Errorf("Query command failed: %v", err)
	}
	
	if !strings.Contains(output, "pasta") {
		t.Errorf("Expected 'pasta' in output, got: %s", output)
	}
}

func TestListReceiptsCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, listReceiptsCmd)
	
	if err != nil {
		t.Errorf("List receipts command failed: %v", err)
	}
	
	if !strings.Contains(output, "Test Store") {
		t.Errorf("Expected 'Test Store' in output, got: %s", output)
	}
}

func TestListQueriesCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, listQueriesCmd)
	
	if err != nil {
		t.Errorf("List queries command failed: %v", err)
	}
	
	if !strings.Contains(output, "What can I cook?") {
		t.Errorf("Expected 'What can I cook?' in output, got: %s", output)
	}
}

func TestCreateWalletPassCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, createWalletPassCmd, "receipt", "test_receipt_123")
	
	if err != nil {
		t.Errorf("Create wallet pass command failed: %v", err)
	}
	
	if !strings.Contains(output, "Test Receipt") {
		t.Errorf("Expected 'Test Receipt' in output, got: %s", output)
	}
}

func TestAnalysisCommand(t *testing.T) {
	server := setupCLITestServer()
	defer server.Close()
	
	// Override baseURL for testing
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	output, err := executeCommand(t, analysisCmd)
	
	if err != nil {
		t.Errorf("Analysis command failed: %v", err)
	}
	
	if !strings.Contains(output, "245.67") {
		t.Errorf("Expected '245.67' in output, got: %s", output)
	}
}

func TestRootCommand(t *testing.T) {
	output, err := executeCommand(t, rootCmd, "--help")
	
	if err != nil {
		t.Errorf("Root command help failed: %v", err)
	}
	
	if !strings.Contains(output, "CLI interface for testing") {
		t.Errorf("Expected help text in output, got: %s", output)
	}
}

// Benchmark tests
func BenchmarkHealthCommand(b *testing.B) {
	server := setupCLITestServer()
	defer server.Close()
	
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeCommand(b, healthCmd)
	}
}

func BenchmarkUploadReceiptCommand(b *testing.B) {
	server := setupCLITestServer()
	defer server.Close()
	
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = originalBaseURL }()
	
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "benchmark_receipt.jpg")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	tempFile.WriteString("fake receipt image data")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeCommand(b, uploadReceiptCmd, tempFile.Name())
	}
}

// Test main function
func TestMain(m *testing.M) {
	fmt.Println("Starting CLI tests...")
	code := m.Run()
	fmt.Println("CLI tests completed.")
	os.Exit(code)
} 