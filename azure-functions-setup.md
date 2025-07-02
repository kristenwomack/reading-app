# Reading Dashboard - Azure Functions Backend

This folder contains the Azure Functions that power the reading dashboard backend.

## Features

- **Book Data Storage**: Store and retrieve reading data using Azure Storage
- **AI Chat Assistant**: OpenAI integration for natural language book management
- **Auto-sync**: Automatic data synchronization with the frontend

## Setup

### 1. Create Azure Function App

```bash
# Install Azure Functions Core Tools
npm install -g azure-functions-core-tools@4

# Create function app
func init ReadingDashboardAPI --javascript
cd ReadingDashboardAPI

# Add CORS settings for GitHub Pages
func settings add CORS_ORIGINS "https://yourusername.github.io"
```

### 2. Create Functions

```bash
# Create book management function
func new --name books --template "HTTP trigger"

# Create AI chat function  
func new --name chat --template "HTTP trigger"
```

### 3. Environment Variables

Add these to your Function App settings:

```
OPENAI_API_KEY=your_openai_api_key
OPENAI_ENDPOINT=https://your-openai-resource.openai.azure.com/
STORAGE_CONNECTION_STRING=your_storage_connection_string
COSMOS_DB_CONNECTION_STRING=your_cosmos_connection_string
```

### 4. Deploy

```bash
# Deploy to Azure
func azure functionapp publish your-function-app-name
```

## Function Endpoints

### Books API (`/api/books`)

**GET** - Retrieve all books
**POST** - Save/update books data

### Chat API (`/api/chat`)

**POST** - Send message to AI assistant
```json
{
  "message": "I finished reading Orbital",
  "books": [...],
  "context": "reading_assistant"
}
```

Response:
```json
{
  "response": "Great! I've marked Orbital as completed...",
  "bookUpdates": [
    {
      "action": "update",
      "title": "Orbital",
      "changes": {
        "Date Read": "2025/07/01",
        "Shelf": "read"
      }
    }
  ]
}
```

## Local Development

```bash
# Start local development server
func start

# Your functions will be available at:
# http://localhost:7071/api/books
# http://localhost:7071/api/chat
```

Update `script.js` to use localhost URLs for development:

```javascript
const CONFIG = {
    AZURE_FUNCTION_URL: 'http://localhost:7071/api', // Local development
    // AZURE_FUNCTION_URL: 'https://your-function-app.azurewebsites.net/api', // Production
};
```
