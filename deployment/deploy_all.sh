#!/bin/bash

# Project Raseed - Complete Deployment Script
# This script deploys the entire Raseed AI system to Google Cloud Platform

set -e

# Configuration
PROJECT_ID=${1:-"your-project-id"}
REGION=${2:-"us-central1"}
ZONE=${3:-"us-central1-a"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Project Raseed - Complete Deployment${NC}"
echo -e "${BLUE}=====================================${NC}"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Zone: $ZONE"
echo ""

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed. Please install it first."
    exit 1
fi

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install it first."
    exit 1
fi

# Set project
print_info "Setting up Google Cloud project..."
gcloud config set project $PROJECT_ID

# Enable required APIs
print_info "Enabling required APIs..."
APIS=(
    "cloudbuild.googleapis.com"
    "run.googleapis.com"
    "cloudfunctions.googleapis.com"
    "firestore.googleapis.com"
    "storage.googleapis.com"
    "pubsub.googleapis.com"
    "vertexai.googleapis.com"
    "aiplatform.googleapis.com"
    "maps-backend.googleapis.com"
    "logging.googleapis.com"
    "monitoring.googleapis.com"
    "iam.googleapis.com"
)

for api in "${APIS[@]}"; do
    gcloud services enable $api --quiet
    print_status "Enabled $api"
done

# Create service accounts
print_info "Creating service accounts..."
SERVICE_ACCOUNTS=(
    "raseed-backend@$PROJECT_ID.iam.gserviceaccount.com"
    "raseed-functions@$PROJECT_ID.iam.gserviceaccount.com"
    "raseed-cli@$PROJECT_ID.iam.gserviceaccount.com"
    "raseed-web@$PROJECT_ID.iam.gserviceaccount.com"
)

for sa in "${SERVICE_ACCOUNTS[@]}"; do
    if ! gcloud iam service-accounts describe $sa &> /dev/null; then
        gcloud iam service-accounts create $(echo $sa | cut -d@ -f1) \
            --display-name="Raseed $(echo $sa | cut -d@ -f1 | cut -d- -f2) Service Account"
        print_status "Created service account: $sa"
    else
        print_warning "Service account already exists: $sa"
    fi
done

# Grant IAM roles
print_info "Granting IAM roles..."
ROLES=(
    "roles/datastore.user"
    "roles/storage.admin"
    "roles/pubsub.publisher"
    "roles/pubsub.subscriber"
    "roles/aiplatform.user"
    "roles/logging.logWriter"
    "roles/monitoring.metricWriter"
    "roles/run.invoker"
)

for sa in "${SERVICE_ACCOUNTS[@]}"; do
    for role in "${ROLES[@]}"; do
        gcloud projects add-iam-policy-binding $PROJECT_ID \
            --member="serviceAccount:$sa" \
            --role="$role" --quiet
    done
    print_status "Granted roles to $sa"
done

# Create Cloud Storage buckets
print_info "Creating Cloud Storage buckets..."
BUCKETS=(
    "raseed-receipts-$PROJECT_ID"
    "raseed-assets-$PROJECT_ID"
)

for bucket in "${BUCKETS[@]}"; do
    if ! gsutil ls -b gs://$bucket &> /dev/null; then
        gsutil mb -l $REGION gs://$bucket
        gsutil iam ch allUsers:objectViewer gs://$bucket
        print_status "Created bucket: $bucket"
    else
        print_warning "Bucket already exists: $bucket"
    fi
done

# Initialize Firestore
print_info "Initializing Firestore..."
if ! gcloud firestore databases describe --database="(default)" &> /dev/null; then
    gcloud firestore databases create --region=$REGION
    print_status "Created Firestore database"
else
    print_warning "Firestore database already exists"
fi

# Deploy Firestore security rules
print_info "Deploying Firestore security rules..."
gcloud firestore rules deploy database/firestore_rules.rules
print_status "Deployed Firestore security rules"

# Create Pub/Sub topics and subscriptions
print_info "Creating Pub/Sub topics and subscriptions..."
TOPICS=(
    "receipt-processing"
    "query-processing"
    "wallet-pass-creation"
    "third-party-integration"
    "spending-analysis"
    "stock-management"
    "notifications"
)

for topic in "${TOPICS[@]}"; do
    gcloud pubsub topics create $topic --quiet || true
    gcloud pubsub subscriptions create ${topic}-sub \
        --topic=$topic \
        --ack-deadline=60 \
        --expiration-period=never --quiet || true
    print_status "Created topic and subscription: $topic"
done

# Build and deploy backend
print_info "Building and deploying backend..."
cd backend
gcloud builds submit --tag gcr.io/$PROJECT_ID/raseed-backend
gcloud run deploy raseed-backend \
    --image gcr.io/$PROJECT_ID/raseed-backend \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-backend@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,CLOUD_STORAGE_BUCKET=raseed-receipts-$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"
cd ..
BACKEND_URL=$(gcloud run services describe raseed-backend --region=$REGION --format="value(status.url)")
print_status "Deployed backend: $BACKEND_URL"

# Deploy Cloud Functions
print_info "Deploying Cloud Functions..."
cd functions

# Receipt Processor
cd receipt_processor
gcloud functions deploy receipt-processor \
    --runtime go121 \
    --region $REGION \
    --trigger-topic receipt-processing \
    --entry-point ProcessReceipt \
    --service-account=raseed-functions@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"
cd ..

# Query Processor
cd query_processor
gcloud functions deploy query-processor \
    --runtime go121 \
    --region $REGION \
    --trigger-topic query-processing \
    --entry-point ProcessQuery \
    --service-account=raseed-functions@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,VERTEX_AI_LOCATION=$REGION"
cd ..

# Third Party Integration
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
print_status "Deployed all Cloud Functions"

# Build and deploy CLI
print_info "Building and deploying CLI..."
cd cli
gcloud builds submit --tag gcr.io/$PROJECT_ID/raseed-cli
gcloud run deploy raseed-cli \
    --image gcr.io/$PROJECT_ID/raseed-cli \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-cli@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,BACKEND_URL=$BACKEND_URL"
cd ..
CLI_URL=$(gcloud run services describe raseed-cli --region=$REGION --format="value(status.url)")
print_status "Deployed CLI: $CLI_URL"

# Build and deploy Flutter web app
print_info "Building and deploying Flutter web app..."
cd flutter_web

# Update API endpoint in constants
sed -i "s|http://localhost:8080|$BACKEND_URL|g" lib/utils/constants.dart

# Build Flutter web app
flutter build web --release

# Create Docker image
gcloud builds submit --tag gcr.io/$PROJECT_ID/raseed-web

# Deploy to Cloud Run
gcloud run deploy raseed-web \
    --image gcr.io/$PROJECT_ID/raseed-web \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --service-account=raseed-web@$PROJECT_ID.iam.gserviceaccount.com \
    --set-env-vars="GOOGLE_CLOUD_PROJECT=$PROJECT_ID,BACKEND_URL=$BACKEND_URL"
cd ..
WEB_URL=$(gcloud run services describe raseed-web --region=$REGION --format="value(status.url)")
print_status "Deployed Flutter web app: $WEB_URL"

# Set up Vertex AI Agent (manual step required)
print_warning "Vertex AI Agent setup requires manual configuration:"
echo "1. Go to Google Cloud Console > Vertex AI > Agent Builder"
echo "2. Create a new agent using the configuration in ai-agent/agent_config.yaml"
echo "3. Configure the agent with your project settings"
echo "4. Deploy the agent and note the endpoint URL"

# Set up Google Wallet API (manual step required)
print_warning "Google Wallet API setup requires manual configuration:"
echo "1. Go to Google Cloud Console > APIs & Services > Credentials"
echo "2. Create a new API key for Google Wallet API"
echo "3. Configure the API key in your application"
echo "4. Set up OAuth 2.0 credentials if needed"

# Create monitoring dashboard
print_info "Setting up monitoring..."
gcloud monitoring dashboards create --config-from-file=monitoring/dashboard.json || true
print_status "Created monitoring dashboard"

# Set up logging
print_info "Setting up logging..."
gcloud logging sinks create raseed-logs \
    storage.googleapis.com/gs://raseed-logs-$PROJECT_ID \
    --log-filter="resource.type=\"cloud_run_revision\" OR resource.type=\"cloud_function\"" || true
print_status "Created logging sink"

# Final summary
echo ""
echo -e "${GREEN}🎉 Deployment Complete!${NC}"
echo -e "${GREEN}=====================${NC}"
echo ""
echo -e "${BLUE}Service URLs:${NC}"
echo "Backend API: $BACKEND_URL"
echo "CLI Interface: $CLI_URL"
echo "Flutter Web App: $WEB_URL"
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo "1. Configure Vertex AI Agent manually in Google Cloud Console"
echo "2. Set up Google Wallet API credentials"
echo "3. Test the CLI interface:"
echo "   curl $CLI_URL/health"
echo "4. Access the web interface: $WEB_URL"
echo "5. Monitor logs: gcloud logging read 'resource.type=\"cloud_run_revision\"'"
echo ""
echo -e "${YELLOW}Important Notes:${NC}"
echo "- The AI agent needs manual configuration in Vertex AI"
echo "- Google Wallet API requires separate setup"
echo "- All services are deployed with public access for testing"
echo "- Consider setting up proper authentication for production"
echo ""
print_status "Project Raseed is now deployed and ready for testing!" 