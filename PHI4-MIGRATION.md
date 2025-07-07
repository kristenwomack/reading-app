# Phi-4 Integration Migration Guide

This document explains the migration from Azure OpenAI to Azure AI (Phi-4) for the Reading Tracker app.

## Changes Made

### 1. Package Dependencies
- **Changed**: `@azure/openai` → `@azure/ai-inference`
- **File**: `api/package.json`

### 2. Azure Functions Code
- **Changed**: OpenAI client → Azure AI Inference client
- **Model**: Now using `Phi-4` instead of GPT models
- **File**: `api/src/functions/books.js`

### 3. Environment Variables
- **Removed**: 
  - `AZURE_OPENAI_ENDPOINT`
  - `AZURE_OPENAI_API_KEY`
  - `AZURE_OPENAI_DEPLOYMENT_NAME`
- **Added**:
  - `AZURE_AI_ENDPOINT`
  - `AZURE_AI_API_KEY`
- **Files**: `.env.example`, `azure.yaml`, `infra/main.bicep`

### 4. Infrastructure (Bicep)
- **Changed**: Azure OpenAI service → Azure AI service (kind: 'AIServices')
- **Removed**: Model deployment resource (Phi-4 is available by default)
- **Files**: `infra/main.bicep`, `infra/main.parameters.json`

## Deployment Instructions

### Local Development

1. **Install dependencies**:
   ```bash
   cd api
   npm install
   ```

2. **Configure environment**:
   - Copy `.env.example` to `.env.local`
   - Set `AZURE_AI_ENDPOINT` and `AZURE_AI_API_KEY` with your Azure AI service values

3. **Run locally**:
   ```bash
   npm run dev
   ```

### Azure Deployment

1. **Deploy infrastructure**:
   ```bash
   azd up
   ```

2. **Set environment variables** (if not using azd):
   ```bash
   azd env set AZURE_AI_ENDPOINT "https://your-ai-resource.cognitiveservices.azure.com/"
   azd env set AZURE_AI_API_KEY "your-api-key"
   ```

## Azure AI Service Configuration

The infrastructure will create an Azure AI service that supports:
- **Phi-4** (our primary model)
- **Other available models**: Meta Llama, Mistral, Cohere, etc.

## Benefits of Phi-4

1. **Performance**: Optimized for Azure infrastructure
2. **Cost**: More cost-effective than GPT-4 family models
3. **Capabilities**: Excellent for conversational AI and structured responses
4. **Integration**: Native Azure AI service integration

## Testing the Chat Feature

The chat endpoint (`/api/books/chat`) now uses Phi-4 for:
- Book recommendations
- Adding books to library (JSON response format)
- General reading assistance

Example request:
```json
{
  "message": "Recommend a sci-fi book",
  "books": []
}
```

## Troubleshooting

1. **"Azure AI not configured" error**: Check environment variables
2. **Deployment failures**: Ensure Azure AI service is available in your region
3. **Model errors**: Verify Phi-4 is available in your subscription/region

## Model Alternatives

If Phi-4 isn't suitable, other recommended models in Azure AI:
- `Mistral-Large-2411`
- `Meta-Llama-3.1-8B-Instruct`
- `Cohere-command-r-plus`

To change models, update the `model` parameter in `api/src/functions/books.js`.
