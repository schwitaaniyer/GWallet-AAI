# Raseed API Documentation

## Overview
The Raseed API provides endpoints for receipt management, AI-powered query processing, and Google Wallet integration. All endpoints return JSON responses and use standard HTTP status codes.

## Base URL
```
https://raseed-backend-PROJECT_ID-uc.a.run.app
```

## Authentication
All API endpoints require authentication. Include the user ID in the request headers or body as specified for each endpoint.

## Endpoints

### Health Check
**GET** `/health`

Check if the service is running.

**Response:**
```json
{
  "status": "healthy"
}
```

---

### Receipt Management

#### Upload Receipt
**POST** `/receipts`

Upload a receipt image for processing.

**Content-Type:** `multipart/form-data`

**Form Data:**
- `user_id` (string, required): User identifier
- `receipt` (file, required): Receipt image file (JPEG, PNG, up to 32MB)

**Response:**
```json
{
  "id": "1703123456789",
  "user_id": "user123",
  "store_name": "",
  "total_amount": 0,
  "tax_amount": 0,
  "items": [],
  "date": "2023-12-21T10:30:45Z",
  "image_url": "https://storage.googleapis.com/bucket/receipts/user123/receipt.jpg",
  "location": {
    "latitude": 0,
    "longitude": 0,
    "address": ""
  },
  "created_at": "2023-12-21T10:30:45Z",
  "updated_at": "2023-12-21T10:30:45Z"
}
```

#### Get User Receipts
**GET** `/receipts?user_id={user_id}`

Retrieve all receipts for a user.

**Query Parameters:**
- `user_id` (string, required): User identifier

**Response:**
```json
[
  {
    "id": "1703123456789",
    "user_id": "user123",
    "store_name": "Walmart",
    "total_amount": 45.99,
    "tax_amount": 3.50,
    "items": [
      {
        "name": "Milk",
        "price": 4.99,
        "quantity": 2,
        "category": "dairy"
      }
    ],
    "date": "2023-12-21T10:30:45Z",
    "image_url": "https://storage.googleapis.com/bucket/receipts/user123/receipt.jpg",
    "location": {
      "latitude": 37.7749,
      "longitude": -122.4194,
      "address": "123 Main St, San Francisco, CA"
    },
    "created_at": "2023-12-21T10:30:45Z",
    "updated_at": "2023-12-21T10:35:12Z"
  }
]
```

---

### Query Processing

#### Submit Query
**POST** `/queries`

Submit a natural language query for AI processing.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "user_id": "user123",
  "query": "What can I cook with my recent purchases?",
  "language": "en"
}
```

**Response:**
```json
{
  "id": "1703123456790",
  "user_id": "user123",
  "query": "What can I cook with my recent purchases?",
  "language": "en",
  "response": "",
  "created_at": "2023-12-21T10:30:45Z"
}
```

#### Get User Queries
**GET** `/queries?user_id={user_id}`

Retrieve all queries for a user.

**Query Parameters:**
- `user_id` (string, required): User identifier

**Response:**
```json
[
  {
    "id": "1703123456790",
    "user_id": "user123",
    "query": "What can I cook with my recent purchases?",
    "language": "en",
    "response": "Based on your recent purchases, you can make: 1. Scrambled eggs with toast 2. Pasta with tomato sauce 3. Grilled cheese sandwich",
    "created_at": "2023-12-21T10:30:45Z"
  }
]
```

---

### Wallet Pass Management

#### Create Wallet Pass
**POST** `/wallet-passes`

Create a new Google Wallet pass.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "user_id": "user123",
  "type": "receipt",
  "title": "Receipt - Walmart",
  "description": "Total: $45.99, Items: 5",
  "data": "{\"receipt_id\": \"1703123456789\", \"store_name\": \"Walmart\"}"
}
```

**Response:**
```json
{
  "id": "receipt_1703123456789",
  "user_id": "user123",
  "type": "receipt",
  "title": "Receipt - Walmart",
  "description": "Total: $45.99, Items: 5",
  "data": "{\"receipt_id\": \"1703123456789\", \"store_name\": \"Walmart\"}",
  "created_at": "2023-12-21T10:30:45Z"
}
```

#### Get User Wallet Passes
**GET** `/wallet-passes?user_id={user_id}`

Retrieve all wallet passes for a user.

**Query Parameters:**
- `user_id` (string, required): User identifier

**Response:**
```json
[
  {
    "id": "receipt_1703123456789",
    "user_id": "user123",
    "type": "receipt",
    "title": "Receipt - Walmart",
    "description": "Total: $45.99, Items: 5",
    "data": "{\"receipt_id\": \"1703123456789\", \"store_name\": \"Walmart\"}",
    "created_at": "2023-12-21T10:30:45Z"
  }
]
```

---

### Spending Analysis

#### Get Spending Analysis
**GET** `/analysis?user_id={user_id}`

Get spending analysis for a user.

**Query Parameters:**
- `user_id` (string, required): User identifier

**Response:**
```json
{
  "total_spent": 245.67,
  "category_spending": {
    "groceries": 120.50,
    "restaurants": 85.25,
    "transportation": 40.00
  },
  "receipt_count": 15,
  "average_per_receipt": 16.38
}
```

---

### Stock Management

#### Create Stock Item
**POST** `/stock-items`

Create a new stock item for inventory management.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "user_id": "user123",
  "name": "Milk",
  "category": "dairy",
  "quantity": 2,
  "unit": "liters",
  "purchase_date": "2023-12-21T10:30:45Z",
  "expiry_date": "2023-12-28T10:30:45Z"
}
```

**Response:**
```json
{
  "id": "1703123456791",
  "user_id": "user123",
  "name": "Milk",
  "category": "dairy",
  "quantity": 2,
  "unit": "liters",
  "purchase_date": "2023-12-21T10:30:45Z",
  "expiry_date": "2023-12-28T10:30:45Z",
  "status": "fresh",
  "created_at": "2023-12-21T10:30:45Z",
  "updated_at": "2023-12-21T10:30:45Z"
}
```

#### Get User Stock Items
**GET** `/stock-items?user_id={user_id}&status={status}`

Retrieve all stock items for a user.

**Query Parameters:**
- `user_id` (string, required): User identifier
- `status` (string, optional): Filter by status (fresh, expiring_soon, expired)

**Response:**
```json
[
  {
    "id": "1703123456791",
    "user_id": "user123",
    "name": "Milk",
    "category": "dairy",
    "quantity": 2,
    "unit": "liters",
    "purchase_date": "2023-12-21T10:30:45Z",
    "expiry_date": "2023-12-28T10:30:45Z",
    "status": "fresh",
    "created_at": "2023-12-21T10:30:45Z",
    "updated_at": "2023-12-21T10:30:45Z"
  }
]
```

#### Update Stock Item
**PUT** `/stock-items?id={item_id}`

Update an existing stock item.

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "name": "Milk",
  "category": "category",
  "quantity": 1,
  "unit": "liters",
  "expiry_date": "2023-12-25T10:30:45Z"
}
```

**Response:**
```json
{
  "id": "1703123456791",
  "user_id": "user123",
  "name": "Milk",
  "category": "dairy",
  "quantity": 1,
  "unit": "liters",
  "purchase_date": "2023-12-21T10:30:45Z",
  "expiry_date": "2023-12-25T10:30:45Z",
  "status": "expiring_soon",
  "created_at": "2023-12-21T10:30:45Z",
  "updated_at": "2023-12-21T11:30:45Z"
}
```

#### Delete Stock Item
**DELETE** `/stock-items?id={item_id}`

Delete a stock item.

**Response:**
```
204 No Content
```

---
## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "user_id is required"
}
```

### 401 Unauthorized
```json
{
  "error": "Authentication required"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Rate Limits

- **Receipt uploads:** 10 requests per minute per user
- **Queries:** 20 requests per minute per user
- **Analysis requests:** 5 requests per minute per user
- **Wallet pass creation:** 10 requests per minute per user

---

## Webhooks

The system publishes events to Pub/Sub topics that can be consumed by webhooks:

### Receipt Processing Events
**Topic:** `receipt-processing`

**Message Format:**
```json
{
  "receipt_id": "1703123456789",
  "user_id": "user123",
  "image_url": "https://storage.googleapis.com/bucket/receipts/user123/receipt.jpg"
}
```

### Query Processing Events
**Topic:** `query-processing`

**Message Format:**
```json
{
  "query_id": "1703123456790",
  "user_id": "user123",
  "query": "What can I cook with my recent purchases?",
  "language": "en"
}
```

### Wallet Pass Creation Events
**Topic:** `wallet-pass-creation`

**Message Format:**
```json
{
  "pass_id": "receipt_1703123456789",
  "user_id": "user123",
  "type": "receipt",
  "title": "Receipt - Walmart"
}
```

### Stock Management Events
**Topic:** `stock-management`

**Message Format:**
```json
{
  "item_id": "1703123456791",
  "user_id": "user123",
  "action": "created",
  "status": "fresh"
}
```

---

## SDKs and Libraries

### Go
```go
import "github.com/your-org/raseed-go-sdk"

client := raseed.NewClient("https://raseed-backend-PROJECT_ID-uc.a.run.app")

// Upload receipt
receipt, err := client.UploadReceipt("user123", file)
```

### Python
```python
from raseed import RaseedClient

client = RaseedClient("https://raseed-backend-PROJECT_ID-uc.a.run.app")

# Upload receipt
receipt = client.upload_receipt("user123", file)
```

### JavaScript
```javascript
import { RaseedClient } from '@raseed/sdk';

const client = new RaseedClient('https://raseed-backend-PROJECT_ID-uc.a.run.app');

// Upload receipt
const receipt = await client.uploadReceipt('user123', file);
```

---

## Support

For API support and questions:
- Email: api-support@raseed.com
- Documentation: https://docs.raseed.com
- Status page: https://status.raseed.com 