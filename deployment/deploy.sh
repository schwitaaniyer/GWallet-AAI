#!/bin/bash

# Raseed Deployment Script
# This script deploys the entire Raseed application to Google Cloud Platform

set -e

# Configuration
PROJECT_ID=${1:-"your-project-id"}
REGION="us-central1"
SERVICE_ACCOUNT="raseed-backend@${PROJECT_ID}.iam.gserviceaccount.com"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting Raseed deployment to project: ${PROJECT_ID}${NC}"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo -e "${YELLOW}📋 Checking prerequisites...${NC}"

if ! command_exists gcloud; then
    echo -e "${RED}❌ Google Cloud SDK is not installed. Please install it first.${NC}"
    exit 1
fi

if ! command_exists docker; then
    echo -e "${RED}❌ Docker is not installed. Please install it first.${NC}"
    exit 1
fi

# Set project
echo -e "${YELLOW}🔧 Setting up Google Cloud project...${NC}"
gcloud config set project $PROJECT_ID

# Enable required APIs
echo -e "${YELLOW}🔌 Enabling required APIs...${NC}"
apis=(
    "cloudbuild.googleapis.com"
    "run.googleapis.com"
    "cloudfunctions.googleapis.com"
    "firestore.googleapis.com"
    "storage.googleapis.com"
    "pubsub.googleapis.com"
    "aiplatform.googleapis.com"
    "maps.googleapis.com"
    "logging.googleapis.com"
    "monitoring.googleapis.com"
)

for api in "${apis[@]}"; do
    echo "Enabling $api..."
    gcloud services enable $api
done

# Create service account
echo -e "${YELLOW}👤 Creating service account...${NC}"
gcloud iam service-accounts create raseed-backend \
    --display-name="Raseed Backend Service Account" \
    --description="Service account for Raseed backend services" \
    || echo "Service account already exists"

# Grant necessary roles
echo -e "${YELLOW}🔐 Granting IAM roles...${NC}"
roles=(
    "roles/aiplatform.user"
    "roles/datastore.user"
    "roles/pubsub.publisher"
    "roles/pubsub.subscriber"
    "roles/storage.objectViewer"
    "roles/storage.objectCreator"
    "roles/logging.logWriter"
    "roles/monitoring.metricWriter"
)

for role in "${roles[@]}"; do
    echo "Granting $role..."
    gcloud projects add-iam-policy-binding $PROJECT_ID \
        --member="serviceAccount:$SERVICE_ACCOUNT" \
        --role="$role"
done

# Create Cloud Storage bucket
echo -e "${YELLOW}📦 Creating Cloud Storage bucket...${NC}"
BUCKET_NAME="raseed-receipts-${PROJECT_ID}"
gsutil mb -l $REGION gs://$BUCKET_NAME || echo "Bucket already exists"

# Set bucket permissions
gsutil iam ch serviceAccount:$SERVICE_ACCOUNT:objectViewer gs://$BUCKET_NAME
gsutil iam ch serviceAccount:$SERVICE_ACCOUNT:objectCreator gs://$BUCKET_NAME

# Create Firestore database
echo -e "${YELLOW}🗄️ Setting up Firestore...${NC}"
gcloud firestore databases create --region=$REGION || echo "Firestore database already exists"

# Deploy Firestore rules
echo -e "${YELLOW}📜 Deploying Firestore security rules...${NC}"
gcloud firestore rules deploy database/firestore_rules.rules

# Create Pub/Sub topics and subscriptions
echo -e "${YELLOW}📡 Creating Pub/Sub topics and subscriptions...${NC}"
topics=(
    "receipt-processing"
    "query-processing"
    "wallet-pass-creation"
    "third-party-integration"
    "spending-analysis"
    "stock-management"
    "notification-events"
)

for topic in "${topics[@]}"; do
    echo "Creating topic: $topic"
    gcloud pubsub topics create $topic || echo "Topic already exists"
done

# Create subscriptions
subscriptions=(
    "receipt-processor-sub:receipt-processing"
    "query-processor-sub:query-processing"
    "wallet-pass-creator-sub:wallet-pass-creation"
    "third-party-integration-sub:third-party-integration"
    "spending-analyzer-sub:spending-analysis"
    "stock-manager-sub:stock-management"
    "notification-processor-sub:notification-events"
)

for subscription in "${subscriptions[@]}"; do
    IFS=':' read -r sub_name topic_name <<< "$subscription"
    echo "Creating subscription: $sub_name for topic: $topic_name"
    gcloud pubsub subscriptions create $sub_name \
        --topic=$topic_name \
        --ack-deadline=60 \
        --message-retention-duration=7d \
        || echo "Subscription already exists"
done

# Build and deploy backend container
echo -e "${YELLOW}🐳 Building and deploying backend container...${NC}"
cd backend

# Build the container
echo "Building container image..."
gcloud builds submit --tag gcr.io/$PROJECT_ID/raseed-backend:latest .

# Deploy to Cloud Run
echo "Deploying to Cloud Run..."
gcloud run deploy raseed-backend \
    --image gcr.io/$PROJECT_ID/raseed-backend:latest \
    --platform managed \
    --region $REGION \
    --service-account $SERVICE_ACCOUNT \
    --allow-unauthenticated \
    --memory 2Gi \
    --cpu 2 \
    --max-instances 100 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=$PROJECT_ID,CLOUD_STORAGE_BUCKET=$BUCKET_NAME,VERTEX_AI_LOCATION=$REGION"

cd ..

# Deploy Cloud Functions
echo -e "${YELLOW}⚡ Deploying Cloud Functions...${NC}"

# Receipt Processor
echo "Deploying receipt processor..."
cd functions/receipt_processor
gcloud functions deploy receipt-processor \
    --runtime go121 \
    --region $REGION \
    --entry-point ProcessReceipt \
    --trigger-topic receipt-processing \
    --memory 2GB \
    --timeout 540s \
    --max-instances 100 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION" \
    --service-account $SERVICE_ACCOUNT
cd ../..

# Query Processor
echo "Deploying query processor..."
cd functions/query_processor
gcloud functions deploy query-processor \
    --runtime go121 \
    --region $REGION \
    --entry-point ProcessQuery \
    --trigger-topic query-processing \
    --memory 1GB \
    --timeout 540s \
    --max-instances 50 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION" \
    --service-account $SERVICE_ACCOUNT
cd ../..

# Third Party Integration
echo "Deploying third-party integration..."
cd functions/third_party_integration
gcloud functions deploy third-party-integration \
    --runtime go121 \
    --region $REGION \
    --entry-point ProcessThirdPartyIntegration \
    --trigger-topic third-party-integration \
    --memory 512MB \
    --timeout 300s \
    --max-instances 20 \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=$PROJECT_ID" \
    --service-account $SERVICE_ACCOUNT
cd ../..

# Set up Vertex AI Agent
echo -e "${YELLOW}🤖 Setting up Vertex AI Agent...${NC}"
# Note: Vertex AI Agent Builder setup requires manual configuration in the console
echo "Please configure the Vertex AI Agent manually in the Google Cloud Console:"
echo "1. Go to Vertex AI > Agent Builder"
echo "2. Create a new agent using the configuration in ai-agent/agent_config.yaml"
echo "3. Deploy the agent"

# Set up monitoring and logging
echo -e "${YELLOW}📊 Setting up monitoring and logging...${NC}"
# Create log-based metrics
gcloud logging metrics create raseed-receipt-uploads \
    --description="Number of receipt uploads" \
    --log-filter='resource.type="cloud_run_revision" AND resource.labels.service_name="raseed-backend" AND textPayload:"receipt"'

gcloud logging metrics create raseed-queries \
    --description="Number of user queries" \
    --log-filter='resource.type="cloud_run_revision" AND resource.labels.service_name="raseed-backend" AND textPayload:"query"'

# Create alerting policies
echo "Creating alerting policies..."
# This would require more complex setup with monitoring policies

echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo ""
echo -e "${GREEN}🎉 Raseed is now deployed and ready to use!${NC}"
echo ""
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. Configure Vertex AI Agent in the Google Cloud Console"
echo "2. Set up Google Wallet API credentials"
echo "3. Configure third-party integrations (Zomato, Blinkit)"
echo "4. Test the API endpoints"
echo "5. Set up monitoring and alerting"
echo ""
echo -e "${GREEN}🌐 Your API endpoint:${NC}"
echo "https://raseed-backend-${PROJECT_ID}-${REGION}.a.run.app"
echo ""
echo -e "${GREEN}📚 API Documentation:${NC}"
echo "See docs/api.md for detailed API documentation"
echo ""
echo -e "${GREEN}🔧 Management:${NC}"
echo "- View logs: gcloud logging read 'resource.type=cloud_run_revision'"
echo "- Monitor metrics: gcloud monitoring metrics list"
echo "- Update deployment: ./deploy.sh $PROJECT_ID" 