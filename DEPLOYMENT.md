# Reading Tracker - Azure Functions Backend

## Quick Setup (5 minutes!)

### 1. Create Azure Resources
```bash
# Login to Azure
az login

# Create resource group
az group create --name reading-tracker-rg --location eastus

# Create storage account (required for Functions)
az storage account create \
  --name readingtrackerstorage \
  --resource-group reading-tracker-rg \
  --location eastus \
  --sku Standard_LRS

# Create Function App
az functionapp create \
  --resource-group reading-tracker-rg \
  --consumption-plan-location eastus \
  --runtime node \
  --runtime-version 18 \
  --functions-version 4 \
  --name reading-tracker-api \
  --storage-account readingtrackerstorage
```

### 2. Set OpenAI Configuration
```bash
# Add your OpenAI API key
az functionapp config appsettings set \
  --name reading-tracker-api \
  --resource-group reading-tracker-rg \
  --settings "OPENAI_API_KEY=your_openai_key_here"
```

### 3. Deploy Functions
```bash
cd api
npm install
npm run build
func azure functionapp publish reading-tracker-api
```

### 4. Setup GitHub Pages
1. Go to your GitHub repository settings
2. Navigate to "Pages" in the sidebar
3. Set source to "GitHub Actions"
4. Add these secrets in repository settings > Secrets and variables > Actions:
   - `AZURE_CREDENTIALS`: Service principal JSON
   - `AZURE_RESOURCE_GROUP`: reading-tracker-rg  
   - `AZURE_FUNCTION_APP_NAME`: reading-tracker-api

### 5. Get Your Function URL
```bash
az functionapp function show \
  --name reading-tracker-api \
  --resource-group reading-tracker-rg \
  --function-name books \
  --query "invokeUrlTemplate"
```

Update the `API_BASE_URL` in `script.js` with this URL.

## Cost Estimate
- **Azure Functions**: FREE (under 1M executions/month)
- **Storage**: ~$0.05/month for small data
- **GitHub Pages**: FREE
- **Total**: ~$0.05/month

## Architecture
```
GitHub Pages (Frontend) → Azure Functions (API) → Azure Storage (Data)
```

Your site will be available at: `https://yourusername.github.io/reading-app`
