#!/bin/bash
set -e

echo "🚀 Setting up Reading Tracker on Azure..."

# Check if Azure CLI is installed
if ! command -v az &> /dev/null; then
    echo "❌ Azure CLI not found. Please install it first:"
    echo "   https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
    exit 1
fi

# Check if logged in
if ! az account show &> /dev/null; then
    echo "🔑 Please login to Azure first:"
    az login
fi

# Configuration
RESOURCE_GROUP="reading-tracker-rg"
LOCATION="eastus"
STORAGE_ACCOUNT="readingtrackerstorage$(date +%s)"
FUNCTION_APP="reading-tracker-api-$(date +%s)"

echo "📝 Creating resources..."
echo "   Resource Group: $RESOURCE_GROUP"
echo "   Function App: $FUNCTION_APP"
echo "   Storage: $STORAGE_ACCOUNT"

# Create resource group
echo "🏗️ Creating resource group..."
az group create --name $RESOURCE_GROUP --location $LOCATION

# Create storage account
echo "💾 Creating storage account..."
az storage account create \
  --name $STORAGE_ACCOUNT \
  --resource-group $RESOURCE_GROUP \
  --location $LOCATION \
  --sku Standard_LRS

# Create Function App
echo "⚡ Creating Function App..."
az functionapp create \
  --resource-group $RESOURCE_GROUP \
  --consumption-plan-location $LOCATION \
  --runtime node \
  --runtime-version 18 \
  --functions-version 4 \
  --name $FUNCTION_APP \
  --storage-account $STORAGE_ACCOUNT

# Get Function App URL
FUNCTION_URL=$(az functionapp show --name $FUNCTION_APP --resource-group $RESOURCE_GROUP --query "defaultHostName" -o tsv)

echo "✅ Azure resources created successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Set up OpenAI API key:"
echo "   az functionapp config appsettings set --name $FUNCTION_APP --resource-group $RESOURCE_GROUP --settings \"OPENAI_API_KEY=your_key_here\""
echo ""
echo "2. Deploy functions:"
echo "   cd api && npm install && func azure functionapp publish $FUNCTION_APP"
echo ""
echo "3. Update script.js with your Function URL:"
echo "   const API_BASE_URL = 'https://$FUNCTION_URL/api';"
echo ""
echo "4. Push to GitHub to trigger deployment!"
echo ""
echo "🌐 Your app will be available at: https://yourusername.github.io/reading-app"
