package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	baseURL = "http://localhost:8080" // Default local URL, can be overridden
	client  = &http.Client{Timeout: 30 * time.Second}
)

// CLI Commands
var rootCmd = &cobra.Command{
	Use:   "raseed-cli",
	Short: "CLI interface for testing Project Raseed AI Agent",
	Long:  `A command-line interface for testing the AI-powered personal assistant integrated with Google Wallet.`,
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check backend health",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Health Check: %s\n", string(body))
	},
}

var uploadReceiptCmd = &cobra.Command{
	Use:   "upload-receipt [image-path]",
	Short: "Upload a receipt image for processing",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		
		// Check if file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist\n", imagePath)
			return
		}

		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		
		// Add file
		file, err := os.Open(imagePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()
		
		part, err := writer.CreateFormFile("receipt", filepath.Base(imagePath))
		if err != nil {
			fmt.Printf("Error creating form file: %v\n", err)
			return
		}
		
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return
		}
		
		// Add user ID (mock for testing)
		writer.WriteField("user_id", "test-user-123")
		
		writer.Close()
		
		// Send request
		req, err := http.NewRequest("POST", baseURL+"/receipts", &buf)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Upload Response: %s\n", string(body))
	},
}

var submitQueryCmd = &cobra.Command{
	Use:   "query [question]",
	Short: "Submit a natural language query",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]
		
		queryData := map[string]interface{}{
			"user_id": "test-user-123",
			"query":   question,
			"language": "en",
		}
		
		jsonData, _ := json.Marshal(queryData)
		
		req, err := http.NewRequest("POST", baseURL+"/queries", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Query Response: %s\n", string(body))
	},
}

var getReceiptsCmd = &cobra.Command{
	Use:   "receipts",
	Short: "Get all receipts for the user",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.Get(baseURL + "/receipts?user_id=test-user-123")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Receipts: %s\n", string(body))
	},
}

var getQueriesCmd = &cobra.Command{
	Use:   "queries",
	Short: "Get all queries for the user",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(baseURL + "/queries?user_id=test-user-123")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Queries: %s\n", string(body))
	},
}

var getWalletPassesCmd = &cobra.Command{
	Use:   "passes",
	Short: "Get all wallet passes for the user",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(baseURL + "/wallet-passes?user_id=test-user-123")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Wallet Passes: %s\n", string(body))
	},
}

var analyzeSpendingCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Get spending analysis for the user",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(baseURL + "/analysis?user_id=test-user-123")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Spending Analysis: %s\n", string(body))
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Raseed CLI Interactive Mode!")
		fmt.Println("Type 'help' for available commands, 'exit' to quit")
		
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("raseed> ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			
			if input == "exit" || input == "quit" {
				break
			}
			
			if input == "help" {
				fmt.Println("Available commands:")
				fmt.Println("  upload <image-path> - Upload receipt image")
				fmt.Println("  query <question>   - Ask a question")
				fmt.Println("  receipts           - List all receipts")
				fmt.Println("  queries            - List all queries")
				fmt.Println("  passes             - List wallet passes")
				fmt.Println("  analyze            - Get spending analysis")
				fmt.Println("  health             - Check backend health")
				fmt.Println("  exit               - Exit interactive mode")
				continue
			}
			
			parts := strings.SplitN(input, " ", 2)
			command := parts[0]
			
			switch command {
			case "upload":
				if len(parts) < 2 {
					fmt.Println("Usage: upload <image-path>")
					continue
				}
				uploadReceiptCmd.Run(cmd, []string{parts[1]})
			case "query":
				if len(parts) < 2 {
					fmt.Println("Usage: query <question>")
					continue
				}
				submitQueryCmd.Run(cmd, []string{parts[1]})
			case "receipts":
				getReceiptsCmd.Run(cmd, []string{})
			case "queries":
				getQueriesCmd.Run(cmd, []string{})
			case "passes":
				getWalletPassesCmd.Run(cmd, []string{})
			case "analyze":
				analyzeSpendingCmd.Run(cmd, []string{})
			case "health":
				healthCmd.Run(cmd, []string{})
			default:
				fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "http://localhost:8080", "Backend API URL")
	
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(uploadReceiptCmd)
	rootCmd.AddCommand(submitQueryCmd)
	rootCmd.AddCommand(getReceiptsCmd)
	rootCmd.AddCommand(getQueriesCmd)
	rootCmd.AddCommand(getWalletPassesCmd)
	rootCmd.AddCommand(analyzeSpendingCmd)
	rootCmd.AddCommand(interactiveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 