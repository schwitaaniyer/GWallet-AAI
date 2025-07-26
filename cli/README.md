# Raseed CLI Interface

A command-line interface for testing and interacting with the Project Raseed AI Agent.

## Features

- **Health Check**: Verify backend connectivity
- **Receipt Upload**: Upload receipt images for AI processing
- **Query Processing**: Submit natural language queries to the AI
- **Data Retrieval**: View receipts, queries, and wallet passes
- **Spending Analysis**: Get spending insights and analytics
- **Interactive Mode**: Command-line interface for easy testing

## Installation

### Prerequisites

- Go 1.21 or later
- Access to the deployed Raseed backend

### Local Development

1. Navigate to the CLI directory:
```bash
cd cli
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the CLI:
```bash
go build -o raseed-cli
```

4. Run the CLI:
```bash
./raseed-cli --help
```

## Usage

### Basic Commands

```bash
# Check backend health
./raseed-cli health

# Upload a receipt image
./raseed-cli upload-receipt /path/to/receipt.jpg

# Submit a query to the AI
./raseed-cli query "What can I cook with my recent purchases?"

# Get all receipts
./raseed-cli receipts

# Get all queries
./raseed-cli queries

# Get wallet passes
./raseed-cli passes

# Get spending analysis
./raseed-cli analyze
```

### Interactive Mode

Start interactive mode for easier testing:

```bash
./raseed-cli interactive
```

Available commands in interactive mode:
- `upload <image-path>` - Upload receipt image
- `query <question>` - Ask a question
- `receipts` - List all receipts
- `queries` - List all queries
- `passes` - List wallet passes
- `analyze` - Get spending analysis
- `health` - Check backend health
- `exit` - Exit interactive mode

### Configuration

Set the backend URL:

```bash
./raseed-cli --url https://your-backend-url.com health
```

## API Integration

The CLI communicates with the Raseed backend API endpoints:

- `GET /health` - Health check
- `POST /receipts` - Upload receipt
- `POST /queries` - Submit query
- `GET /receipts` - Get receipts
- `GET /queries` - Get queries
- `GET /wallet-passes` - Get wallet passes
- `GET /analysis` - Get spending analysis

## Testing Examples

### 1. Health Check
```bash
./raseed-cli health
```
Expected output: `Health Check: {"status":"healthy","timestamp":"..."}`

### 2. Upload Receipt
```bash
./raseed-cli upload-receipt sample_receipt.jpg
```
Expected output: Receipt processing confirmation with receipt ID

### 3. AI Query
```bash
./raseed-cli query "What's my total spending this month?"
```
Expected output: Query submission confirmation

### 4. View Data
```bash
./raseed-cli receipts
./raseed-cli passes
./raseed-cli analyze
```

## Error Handling

The CLI provides clear error messages for:
- Network connectivity issues
- Invalid file paths
- API errors
- Authentication issues

## Development

### Adding New Commands

1. Define the command in `main.go`:
```go
var newCmd = &cobra.Command{
    Use:   "new-command",
    Short: "Description of new command",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}
```

2. Add to root command:
```go
rootCmd.AddCommand(newCmd)
```

### Testing

Run tests:
```bash
go test ./...
```

## Deployment

The CLI is deployed as a Cloud Run service for web-based access:

```bash
# Deploy to Cloud Run
gcloud run deploy raseed-cli \
    --image gcr.io/PROJECT_ID/raseed-cli \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated
```

## Troubleshooting

### Common Issues

1. **Connection refused**: Check if backend is running and URL is correct
2. **File not found**: Verify file path exists and is accessible
3. **Authentication error**: Ensure proper credentials are configured
4. **Timeout errors**: Check network connectivity and backend response times

### Debug Mode

Enable verbose logging:
```bash
./raseed-cli --verbose health
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is part of Project Raseed and follows the same license terms. 