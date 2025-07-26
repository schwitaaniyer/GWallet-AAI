# Project Raseed Setup Guide

## Overview
This guide will help you set up and deploy Project Raseed, an AI-powered personal assistant integrated with Google Wallet for receipt management and financial planning.

## Prerequisites

### Required Software
- Google Cloud SDK (gcloud CLI)
- Docker
- Go 1.21 or later
- Git

### Required Google Cloud Services
- Google Cloud Project with billing enabled
- Cloud Run
- Cloud Functions
- Firestore
- Cloud Storage
- Pub/Sub
- Vertex AI
- Cloud Logging
- Cloud Monitoring

## Step 1: Project Setup

### 1.1 Create Google Cloud Project
```bash
# Create a new project (or use existing)
gcloud projects create raseed-project-123 --name="Project Raseed"

# Set the project as default
gcloud config set project raseed-project-123

# Enable billing (required for AI services)
# Go to: https://console.cloud.google.com/billing
```

### 1.2 Enable Required APIs
```bash
# Enable all required APIs
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable cloudfunctions.googleapis.com
gcloud services enable firestore.googleapis.com
gcloud services enable storage.googleapis.com
gcloud services enable pubsub.googleapis.com
gcloud services enable aiplatform.googleapis.com
gcloud services enable maps.googleapis.com
gcloud services enable logging.googleapis.com
gcloud services enable monitoring.googleapis.com
```

## Step 2: Service Account Setup

### 2.1 Create Service Account
```bash
# Create service account
gcloud iam service-accounts create raseed-backend \
    --display-name="Raseed Backend Service Account" \
    --description="Service account for Raseed backend services"
```

### 2.2 Grant Required Roles
```bash
# Grant necessary IAM roles
gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/aiplatform.user"

gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/datastore.user"

gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/pubsub.publisher"

gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/pubsub.subscriber"

gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/storage.objectViewer"

gcloud projects add-iam-policy-binding raseed-project-123 \
    --member="serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com" \
    --role="roles/storage.objectCreator"
```

## Step 3: Infrastructure Setup

### 3.1 Create Cloud Storage Bucket
```bash
# Create bucket for receipt storage
gsutil mb -l us-central1 gs://raseed-receipts-raseed-project-123

# Set bucket permissions
gsutil iam ch serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com:objectViewer gs://raseed-receipts-raseed-project-123
gsutil iam ch serviceAccount:raseed-backend@raseed-project-123.iam.gserviceaccount.com:objectCreator gs://raseed-receipts-raseed-project-123
```

### 3.2 Create Firestore Database
```bash
# Create Firestore database
gcloud firestore databases create --region=us-central1

# Deploy security rules
gcloud firestore rules deploy database/firestore_rules.rules
```

### 3.3 Create Pub/Sub Topics and Subscriptions
```bash
# Create topics
gcloud pubsub topics create receipt-processing
gcloud pubsub topics create query-processing
gcloud pubsub topics create wallet-pass-creation
gcloud pubsub topics create third-party-integration
gcloud pubsub topics create spending-analysis
gcloud pubsub topics create stock-management
gcloud pubsub topics create notification-events

# Create subscriptions
gcloud pubsub subscriptions create receipt-processor-sub \
    --topic=receipt-processing \
    --ack-deadline=60 \
    --message-retention-duration=7d

gcloud pubsub subscriptions create query-processor-sub \
    --topic=query-processing \
    --ack-deadline=60 \
    --message-retention-duration=7d

gcloud pubsub subscriptions create wallet-pass-creator-sub \
    --topic=wallet-pass-creation \
    --ack-deadline=60 \
    --message-retention-duration=7d

gcloud pubsub subscriptions create third-party-integration-sub \
    --topic=third-party-integration \
    --ack-deadline=60 \
    --message-retention-duration=7d
```

## Step 4: Deploy Backend Service

### 4.1 Build and Deploy Backend
```bash
# Navigate to backend directory
cd backend

# Build container image
gcloud builds submit --tag gcr.io/raseed-project-123/raseed-backend:latest .

# Deploy to Cloud Run
gcloud run deploy raseed-backend \
    --image gcr.io/raseed-project-123/raseed-backend:latest \
    --platform managed \
    --region us-central1 \
    --service-account raseed-backend@raseed-project-123.iam.gserviceaccount.com \
    --allow-unauthenticated \
    --memory 2Gi \
    --cpu 2 \
    --max-instances 100 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=raseed-project-123,CLOUD_STORAGE_BUCKET=raseed-receipts-raseed-project-123,VERTEX_AI_LOCATION=us-central1"

cd ..
```

## Step 5: Deploy Cloud Functions

### 5.1 Deploy Receipt Processor
```bash
cd functions/receipt_processor

gcloud functions deploy receipt-processor \
    --runtime go121 \
    --region us-central1 \
    --entry-point ProcessReceipt \
    --trigger-topic receipt-processing \
    --memory 2GB \
    --timeout 540s \
    --max-instances 100 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=raseed-project-123,VERTEX_AI_LOCATION=us-central1" \
    --service-account raseed-backend@raseed-project-123.iam.gserviceaccount.com

cd ../..
```

### 5.2 Deploy Query Processor
```bash
cd functions/query_processor

gcloud functions deploy query-processor \
    --runtime go121 \
    --region us-central1 \
    --entry-point ProcessQuery \
    --trigger-topic query-processing \
    --memory 1GB \
    --timeout 540s \
    --max-instances 50 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=raseed-project-123,VERTEX_AI_LOCATION=us-central1" \
    --service-account raseed-backend@raseed-project-123.iam.gserviceaccount.com

cd ../..
```

### 5.3 Deploy Third-Party Integration
```bash
cd functions/third_party_integration

gcloud functions deploy third-party-integration \
    --runtime go121 \
    --region us-central1 \
    --entry-point ProcessThirdPartyIntegration \
    --trigger-topic third-party-integration \
    --memory 512MB \
    --timeout 300s \
    --max-instances 20 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=raseed-project-123" \
    --service-account raseed-backend@raseed-project-123.iam.gserviceaccount.com

cd ../..
```

## Step 6: Configure Vertex AI Agent

### 6.1 Set up Vertex AI Agent Builder
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Navigate to Vertex AI > Agent Builder
3. Click "Create Agent"
4. Use the configuration from `ai-agent/agent_config.yaml`
5. Deploy the agent

### 6.2 Configure Gemini AI
1. In Vertex AI, go to Model Garden
2. Enable Gemini Pro and Gemini Pro Vision models
3. Set up API access for the service account

## Step 7: Google Wallet Integration

### 7.1 Set up Google Wallet API
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Navigate to APIs & Services > Library
3. Search for "Google Wallet API" and enable it
4. Create credentials (Service Account Key)
5. Download the JSON key file

### 7.2 Configure Wallet Pass Types
1. Go to [Google Wallet Console](https://wallet.google.com/manager)
2. Create pass classes for:
   - Receipt passes
   - Shopping list passes
   - Financial insight passes
3. Note the pass class IDs for integration

## Step 8: Monitoring and Logging

### 8.1 Set up Cloud Logging
```bash
# Create log-based metrics
gcloud logging metrics create raseed-receipt-uploads \
    --description="Number of receipt uploads" \
    --log-filter='resource.type="cloud_run_revision" AND resource.labels.service_name="raseed-backend" AND textPayload:"receipt"'

gcloud logging metrics create raseed-queries \
    --description="Number of user queries" \
    --log-filter='resource.type="cloud_run_revision" AND resource.labels.service_name="raseed-backend" AND textPayload:"query"'
```

### 8.2 Set up Cloud Monitoring
1. Go to Cloud Monitoring in the console
2. Create dashboards for:
   - API request metrics
   - Function execution metrics
   - Error rates
   - Response times

## Step 9: Testing

### 9.1 Test API Endpoints
```bash
# Test health endpoint
curl https://raseed-backend-raseed-project-123-uc.a.run.app/health

# Test receipt upload (replace with actual file)
curl -X POST https://raseed-backend-raseed-project-123-uc.a.run.app/receipts \
  -F "user_id=test_user" \
  -F "receipt=@test_receipt.jpg"

# Test query submission
curl -X POST https://raseed-backend-raseed-project-123-uc.a.run.app/queries \
  -H "Content-Type: application/json" \
  -d '{"user_id": "test_user", "query": "What can I cook?", "language": "en"}'
```

### 9.2 Run Integration Tests
```bash
# Run tests
cd tests
go test -v
cd ..
```

## Step 10: Production Considerations

### 10.1 Security
- Enable Cloud Armor for DDoS protection
- Set up proper authentication (Firebase Auth recommended)
- Configure CORS policies
- Enable audit logging

### 10.2 Performance
- Set up CDN for static assets
- Configure auto-scaling policies
- Monitor and optimize database queries
- Set up caching strategies

### 10.3 Cost Optimization
- Set up billing alerts
- Monitor resource usage
- Optimize function memory allocation
- Use committed use discounts where applicable

## Troubleshooting

### Common Issues

1. **Permission Denied Errors**
   - Verify service account has correct roles
   - Check IAM policies
   - Ensure API is enabled

2. **Function Deployment Failures**
   - Check Go module dependencies
   - Verify function entry points
   - Check environment variables

3. **AI Processing Errors**
   - Verify Vertex AI setup
   - Check model availability
   - Monitor quota limits

4. **Database Connection Issues**
   - Verify Firestore rules
   - Check network connectivity
   - Monitor database quotas

### Useful Commands

```bash
# View logs
gcloud logging read 'resource.type="cloud_run_revision"'

# Monitor function logs
gcloud functions logs read receipt-processor --limit=50

# Check service status
gcloud run services describe raseed-backend --region=us-central1

# Update deployment
./deployment/deploy.sh raseed-project-123
```

## Support

For additional support:
- Check the [API Documentation](docs/api.md)
- Review [Google Cloud Documentation](https://cloud.google.com/docs)
- Monitor [Google Cloud Status](https://status.cloud.google.com)

## Next Steps

1. **Frontend Development**: Build the user interface using Flutter or web technologies
2. **Mobile App**: Develop native mobile applications
3. **Advanced Features**: Implement advanced AI features like predictive analytics
4. **Third-party Integrations**: Add more service integrations
5. **Analytics**: Implement detailed analytics and reporting 