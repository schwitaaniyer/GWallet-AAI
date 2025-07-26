# Project Raseed - AI-Powered Personal Assistant with Google Wallet Integration

## Overview
Project Raseed is an AI-powered personal assistant integrated with Google Wallet, designed to revolutionize receipt management and financial planning. Using Gemini AI on Vertex AI, Raseed digitizes receipts from photos, videos, or live streams in any language, extracting details like items, prices, taxes, and fees.

## Features
- **Multimodal Receipt Ingestion**: Process photos, videos, or live streams using Gemini AI
- **Google Wallet Integration**: Create and manage passes for receipts, shopping lists, and insights
- **Local Language Queries**: Handle queries in any language with Vertex AI Agent Builder
- **Spending Analysis**: Analyze expenses and suggest savings
- **Stock Management**: Track items and notify about expiry dates
- **Third-Party Integration**: Fetch bills from apps like Zomato/Blinkit
- **Push Notifications**: Update users via Google Wallet API
- **Location-Based Insights**: Use Google Maps Platform for store-specific suggestions
- **Predictive Analysis**: Forecast expenditure patterns

## Architecture
- **Backend**: Cloud Run (Go)
- **AI Processing**: Vertex AI with Gemini AI
- **Database**: Firestore
- **Storage**: Cloud Storage
- **Functions**: Cloud Functions (Go)
- **Messaging**: Pub/Sub
- **Monitoring**: Cloud Logging
- **Maps**: Google Maps Platform

## Project Structure
```
GWallet-AAI/
├── backend/                 # Cloud Run backend (Go)
├── functions/              # Cloud Functions (Go)
├── ai-agent/              # Vertex AI Agent Builder
├── database/              # Firestore schemas and rules
├── storage/               # Cloud Storage configuration
├── pubsub/                # Pub/Sub topics and subscriptions
├── monitoring/            # Cloud Logging and monitoring
├── deployment/            # Deployment configurations
├── docs/                  # Documentation
└── tests/                 # Test files
```

## Quick Start
1. Set up Google Cloud Project
2. Enable required APIs
3. Deploy backend to Cloud Run
4. Deploy Cloud Functions
5. Configure Vertex AI Agent
6. Set up Firestore and Cloud Storage
7. Configure Pub/Sub topics

## API Documentation
See `docs/api.md` for detailed API documentation.

## Deployment
See `deployment/` directory for deployment instructions.
