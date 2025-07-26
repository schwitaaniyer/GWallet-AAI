# Project Raseed - Complete Deployment Guide

This guide provides step-by-step instructions for deploying the entire Project Raseed AI system to Google Cloud Platform, including the backend, CLI interface, Flutter web app, and AI agent.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Project Setup](#project-setup)
3. [Backend Deployment](#backend-deployment)
4. [CLI Interface Deployment](#cli-interface-deployment)
5. [Flutter Web App Deployment](#flutter-web-app-deployment)
6. [AI Agent Configuration](#ai-agent-configuration)
7. [Testing and Validation](#testing-and-validation)
8. [Monitoring and Maintenance](#monitoring-and-maintenance)
9. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools

- **Google Cloud CLI**: [Install gcloud CLI](https://cloud.google.com/sdk/docs/install)
- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Flutter SDK**: [Install Flutter](https://flutter.dev/docs/get-started/install)
- **Go**: [Install Go 1.21+](https://golang.org/doc/install)

### Google Cloud Requirements

- **Google Cloud Project**: Create or use existing project
- **Billing Enabled**: Ensure billing is enabled for the project
- **IAM Permissions**: Owner or Editor role on the project

### Verify Installation

```bash
# Check gcloud
gcloud version

# Check Docker
docker --version

# Check Flutter
flutter --version

# Check Go
go version
```

## Project Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd GWallet-AAI
```

### 2. Set Project Configuration

```bash
# Set your project ID
export PROJECT_ID="your-project-id"
export REGION="us-central1"
export ZONE="us-central1-a"

# Configure gcloud
gcloud config set project $PROJECT_ID
gcloud config set compute/region $REGION
gcloud config set compute/zone $ZONE
```

### 3. Enable Required APIs

```bash
# Run the deployment script to enable all APIs
chmod +x deployment/deploy_all.sh
./deployment/deploy_all.sh $PROJECT_ID $REGION $ZONE
```

## Backend Deployment

### 1. Build Backend Image

```bash
cd backend

# Build Docker image
docker build -t gcr.io/$PROJECT_ID/raseed-backend .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/raseed-backend
```

### 2. Deploy to Cloud Run

```bash
# Deploy backend service
gcloud run deploy raseed-backend \
    --image gcr.io/$PROJECT_ID/raseed-backend \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-backend@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,CLOUD_STORAGE_BUCKET=raseed-receipts-$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"

# Get backend URL
BACKEND_URL=$(gcloud run services describe raseed-backend --region=$REGION --format="value(status.url)")
echo "Backend URL: $BACKEND_URL"
```

### 3. Deploy Cloud Functions

```bash
cd functions

# Deploy receipt processor
cd receipt_processor
gcloud functions deploy receipt-processor \
    --runtime go121 \
    --region $REGION \
    --trigger-topic receipt-processing \
    --entry-point ProcessReceipt \
    --service-account=raseed-functions@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"
cd ..

# Deploy query processor
cd query_processor
gcloud functions deploy query-processor \
    --runtime go121 \
    --region $REGION \
    --trigger-topic query-processing \
    --entry-point ProcessQuery \
    --service-account=raseed-functions@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"
cd ..

# Deploy third-party integration
cd third_party_integration
gcloud functions deploy third-party-integration \
    --runtime go121 \
    --region $REGION \
    --trigger-topic third-party-integration \
    --entry-point ProcessThirdPartyBill \
    --service-account=raseed-functions@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID"
cd ..

cd ..
```

## CLI Interface Deployment

### 1. Build CLI Image

```bash
cd cli

# Build Docker image
docker build -t gcr.io/$PROJECT_ID/raseed-cli .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/raseed-cli
```

### 2. Deploy CLI to Cloud Run

```bash
# Deploy CLI service
gcloud run deploy raseed-cli \
    --image gcr.io/$PROJECT_ID/raseed-cli \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-cli@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,BACKEND_URL=$BACKEND_URL"

# Get CLI URL
CLI_URL=$(gcloud run services describe raseed-cli --region=$REGION --format="value(status.url)")
echo "CLI URL: $CLI_URL"
cd ..
```

## Flutter Web App Deployment

### 1. Configure Flutter App

```bash
cd flutter_web

# Update backend URL in constants
sed -i "s|http://localhost:8080|$BACKEND_URL|g" lib/utils/constants.dart

# Install dependencies
flutter pub get
```

### 2. Build Flutter Web App

```bash
# Build for web
flutter build web --release

# Verify build output
ls -la build/web/
```

### 3. Build and Deploy Docker Image

```bash
# Build Docker image
docker build -t gcr.io/$PROJECT_ID/raseed-web .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/raseed-web
```

### 4. Deploy to Cloud Run

```bash
# Deploy web app
gcloud run deploy raseed-web \
    --image gcr.io/$PROJECT_ID/raseed-web \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-web@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,BACKEND_URL=$BACKEND_URL"

# Get web app URL
WEB_URL=$(gcloud run services describe raseed-web --region=$REGION --format="value(status.url)")
echo "Web App URL: $WEB_URL"
cd ..
```

## AI Agent Configuration

### 1. Vertex AI Agent Setup

1. **Go to Google Cloud Console**
   - Navigate to Vertex AI > Agent Builder
   - Click "Create Agent"

2. **Configure Agent**
   - **Display Name**: Raseed AI Agent
   - **Description**: AI-powered personal assistant for receipt management
   - **Model**: Gemini Pro
   - **Location**: Same as your project region

3. **Upload Configuration**
   - Use the configuration from `ai-agent/agent_config.yaml`
   - Update project-specific values

4. **Deploy Agent**
   - Click "Deploy"
   - Note the endpoint URL for integration

### 2. Google Wallet API Setup

1. **Enable Google Wallet API**
   ```bash
   gcloud services enable walletobjects.googleapis.com
   ```

2. **Create API Credentials**
   - Go to APIs & Services > Credentials
   - Create API Key
   - Restrict to Google Wallet API

3. **Configure OAuth 2.0**
   - Create OAuth 2.0 Client ID
   - Add authorized origins and redirect URIs

## Testing and Validation

### 1. Health Check

```bash
# Test backend
curl $BACKEND_URL/health

# Test CLI
curl $CLI_URL/health

# Test web app
curl $WEB_URL/health
```

### 2. CLI Testing

```bash
# Test CLI commands
./cli/raseed-cli --url $BACKEND_URL health
./cli/raseed-cli --url $BACKEND_URL receipts
./cli/raseed-cli --url $BACKEND_URL interactive
```

### 3. Web App Testing

1. **Open Web App**: Navigate to `$WEB_URL`
2. **Login**: Use demo credentials (demo@raseed.com / demo123)
3. **Test Features**:
   - Dashboard overview
   - Receipt upload (mock)
   - AI query submission
   - Analytics view

### 4. API Testing

```bash
# Test receipt upload
curl -X POST $BACKEND_URL/receipts \
  -F "receipt=@sample_receipt.jpg" \
  -F "user_id=test-user-123"

# Test query submission
curl -X POST $BACKEND_URL/queries \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test-user-123","query":"What can I cook?","language":"en"}'

# Test spending analysis
curl "$BACKEND_URL/analysis?user_id=test-user-123"
```

## Monitoring and Maintenance

### 1. Set Up Monitoring

```bash
# Create monitoring dashboard
gcloud monitoring dashboards create --config-from-file=monitoring/dashboard.json

# Set up alerting policies
gcloud alpha monitoring policies create --policy-from-file=monitoring/alerting-policy.yaml
```

### 2. Logging

```bash
# View logs
gcloud logging read 'resource.type="cloud_run_revision"' --limit=50

# Set up log export
gcloud logging sinks create raseed-logs \
    storage.googleapis.com/gs://raseed-logs-$PROJECT_ID \
    --log-filter="resource.type=\"cloud_run_revision\" OR resource.type=\"cloud_function\""
```

### 3. Performance Monitoring

```bash
# Monitor service metrics
gcloud monitoring metrics list --filter="metric.type:run.googleapis.com"

# Check resource usage
gcloud compute instances list --filter="name~raseed"
```

## Troubleshooting

### Common Issues

#### 1. Build Failures

```bash
# Check Docker build logs
docker build --no-cache -t test-image .

# Check Flutter build
flutter doctor
flutter clean && flutter pub get
```

#### 2. Deployment Failures

```bash
# Check service status
gcloud run services list --region=$REGION

# View deployment logs
gcloud run services logs read raseed-backend --region=$REGION
```

#### 3. API Errors

```bash
# Check API enablement
gcloud services list --enabled --filter="name:raseed"

# Test API endpoints
curl -v $BACKEND_URL/health
```

#### 4. Authentication Issues

```bash
# Check service account permissions
gcloud projects get-iam-policy $PROJECT_ID \
    --flatten="bindings[].members" \
    --filter="bindings.members:raseed"

# Verify service account exists
gcloud iam service-accounts list --filter="email:raseed"
```

### Debug Commands

```bash
# Check all services
gcloud run services list --region=$REGION

# Check Cloud Functions
gcloud functions list --region=$REGION

# Check Pub/Sub topics
gcloud pubsub topics list

# Check Firestore
gcloud firestore databases list

# Check Cloud Storage
gsutil ls gs://raseed-*
```

### Performance Optimization

#### 1. Scaling

```bash
# Update scaling configuration
gcloud run services update raseed-backend \
    --region=$REGION \
    --min-instances=1 \
    --max-instances=10
```

#### 2. Resource Limits

```bash
# Update resource limits
gcloud run services update raseed-backend \
    --region=$REGION \
    --cpu=2 \
    --memory=2Gi
```

## Security Considerations

### 1. Authentication

- Implement proper authentication for production
- Use Firebase Authentication or Google Identity
- Set up proper CORS policies

### 2. Network Security

- Use VPC for internal communication
- Implement proper firewall rules
- Use HTTPS for all external communication

### 3. Data Protection

- Encrypt data at rest and in transit
- Implement proper access controls
- Regular security audits

## Cost Optimization

### 1. Resource Management

- Use appropriate instance sizes
- Implement auto-scaling
- Monitor resource usage

### 2. Storage Optimization

- Use appropriate storage classes
- Implement lifecycle policies
- Regular cleanup of unused resources

### 3. API Usage

- Monitor API quotas
- Implement caching strategies
- Optimize API calls

## Backup and Recovery

### 1. Data Backup

```bash
# Export Firestore data
gcloud firestore export gs://raseed-backup-$PROJECT_ID

# Backup Cloud Storage
gsutil -m rsync -r gs://raseed-receipts-$PROJECT_ID gs://raseed-backup-$PROJECT_ID/receipts
```

### 2. Disaster Recovery

- Document recovery procedures
- Test backup restoration
- Maintain multiple backup locations

## Support and Maintenance

### 1. Regular Maintenance

- Update dependencies regularly
- Monitor security patches
- Review and optimize performance

### 2. Monitoring

- Set up alerts for critical issues
- Monitor costs and usage
- Track user feedback and issues

### 3. Documentation

- Keep deployment documentation updated
- Document configuration changes
- Maintain runbooks for common issues

## Conclusion

This deployment guide provides a comprehensive approach to deploying Project Raseed on Google Cloud Platform. The system is designed to be scalable, secure, and maintainable. Regular monitoring and maintenance will ensure optimal performance and reliability.

For additional support:
- Check the project documentation
- Review Google Cloud documentation
- Create issues in the project repository 