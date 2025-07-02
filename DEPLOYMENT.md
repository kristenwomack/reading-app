# Reading Tracker - Modern Azure Deployment

This guide covers deploying the Reading Tracker app using modern Azure tooling (azd) with Infrastructure as Code.

## 🚀 Quick Setup with Azure Developer CLI (Recommended)

### Prerequisites
- [Azure Developer CLI (azd)](https://aka.ms/install-azd)
- [Node.js 18+](https://nodejs.org/)
- [Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli) (optional, for manual tasks)
- Azure subscription

### 1. Clone and Initialize
```bash
git clone <your-repo-url>
cd reading-app
azd init
```

### 2. Deploy Everything at Once
```bash
# This provisions all Azure resources and deploys the application
azd up
```

That's it! The `azd up` command will:
- Create resource group
- Provision Azure Functions
- Set up storage account
- Deploy Azure OpenAI
- Configure Application Insights
- Deploy your functions
- Output the API URL

### 3. Configure GitHub Actions (Auto-Deployment)

Set these secrets in your GitHub repository (Settings > Secrets and variables > Actions):

#### Required Secrets:
- `AZURE_CREDENTIALS`: Service principal JSON (see below)
- `AZURE_OPENAI_ENDPOINT`: Your Azure OpenAI endpoint URL
- `AZURE_OPENAI_API_KEY`: Your Azure OpenAI API key
- `AZURE_OPENAI_DEPLOYMENT_NAME`: Model deployment name (e.g., "gpt-4")

#### Get Azure Credentials:
```bash
# Create service principal for GitHub Actions
az ad sp create-for-rbac --name "reading-tracker-deploy" \
  --role contributor \
  --scopes /subscriptions/{subscription-id} \
  --sdk-auth
```

Copy the JSON output to the `AZURE_CREDENTIALS` secret.

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   GitHub Pages  │───▶│  Azure Functions │───▶│  Azure Storage  │
│   (Frontend)    │    │     (API)        │    │     (Data)      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Azure OpenAI   │
                       │    (Chat AI)    │
                       └─────────────────┘
```

### Resources Created:
- **Azure Functions**: Serverless API backend
- **Storage Account**: Book data and function storage
- **Azure OpenAI**: AI-powered book recommendations
- **Application Insights**: Monitoring and analytics
- **Resource Group**: Container for all resources

## 💰 Cost Estimate

- **Azure Functions**: FREE (first 1M executions/month)
- **Storage**: ~$0.05/month for small data
- **Azure OpenAI**: Pay-per-use (~$0.10 per conversation)
- **Application Insights**: FREE (first 5GB/month)
- **GitHub Pages**: FREE
- **Total**: ~$0.15/month + AI usage

## 🔧 Manual Deployment (Alternative)

If you prefer manual control:

### 1. Provision Infrastructure
```bash
azd provision
```

### 2. Deploy Application
```bash
azd deploy
```

### 3. Set Environment Variables
```bash
azd env set AZURE_OPENAI_ENDPOINT "https://your-instance.openai.azure.com/"
azd env set AZURE_OPENAI_API_KEY "your-api-key"
azd env set AZURE_OPENAI_DEPLOYMENT_NAME "gpt-4"
```

## 🛠️ Development Workflow

### Local Development
```bash
# Start functions locally
cd api
npm install
npm start

# In another terminal, serve frontend
python -m http.server 8000
```

### Testing Deployment
```bash
# Preview changes before deploying
azd provision --preview

# Deploy with confirmation
azd up
```

### Environment Management
```bash
# List environments
azd env list

# Switch environments
azd env select <environment-name>

# View environment variables
azd env get-values
```

## 📋 CI/CD Pipeline

The GitHub Actions workflow automatically:
1. Deploys backend to Azure Functions using `azd up`
2. Updates frontend with API URL
3. Deploys frontend to GitHub Pages

Triggered on every push to `main` branch.

## 🔍 Monitoring

- **Application Insights**: Monitor function performance and errors
- **Function Logs**: `azd logs` or Azure portal
- **Cost Management**: Azure Cost Management + Billing

## 🆘 Troubleshooting

### Common Issues:

**Functions not starting:**
```bash
azd logs
```

**Environment variables missing:**
```bash
azd env get-values
azd env set <KEY> "<VALUE>"
```

**Deployment failures:**
```bash
# Check deployment status
azd show

# Re-run deployment
azd up
```

**Resource naming conflicts:**
```bash
# Use custom environment name
azd env new <unique-name>
```

## 🔗 Useful Links

- [Azure Developer CLI Documentation](https://docs.microsoft.com/azure/developer/azure-developer-cli/)
- [Azure Functions Documentation](https://docs.microsoft.com/azure/azure-functions/)
- [Azure OpenAI Documentation](https://docs.microsoft.com/azure/cognitive-services/openai/)

---

Your app will be available at:
- **Frontend**: `https://yourusername.github.io/reading-app`
- **API**: Get URL from `azd env get-values | grep API_URI`
