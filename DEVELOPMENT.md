# Development Guide

## 🏃‍♂️ Quick Start for Developers

### Prerequisites
- [Node.js 18+](https://nodejs.org/)
- [Azure Functions Core Tools](https://docs.microsoft.com/azure/azure-functions/functions-run-local)
- [Azure Developer CLI](https://aka.ms/install-azd)

### Local Development Setup

1. **Clone and install dependencies:**
   ```bash
   git clone <your-repo>
   cd reading-app
   cd api && npm install && cd ..
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env.local
   # Edit .env.local with your Azure OpenAI credentials
   ```

3. **Start the backend:**
   ```bash
   cd api
   npm run dev
   ```

4. **Start the frontend:**
   ```bash
   # In another terminal, from the root directory
   python -m http.server 8000
   # Or use any local server: npx serve, live-server, etc.
   ```

5. **Open the app:**
   Navigate to `http://localhost:8000`

### Project Structure

```
reading-app/
├── 📁 Frontend (Static Files)
│   ├── index.html          # Main dashboard
│   ├── style.css           # GitHub-inspired styling
│   ├── script.js           # Frontend JavaScript
│   └── books.json          # Sample data
├── 📁 Backend (Azure Functions)
│   └── api/
│       ├── package.json    # Dependencies
│       ├── host.json       # Function app config
│       └── src/functions/
│           └── books.js    # API endpoints
├── 📁 Infrastructure (azd)
│   ├── azure.yaml          # azd project config
│   └── infra/
│       ├── main.bicep      # Infrastructure as Code
│       └── main.parameters.json
└── 📁 DevOps
    └── .github/workflows/
        └── deploy.yml      # CI/CD pipeline
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/books` | Get all books |
| POST | `/api/books` | Save books data |
| GET | `/api/books/sync` | Check sync status |
| POST | `/api/books/chat` | AI assistant chat |

### Development Workflow

1. **Make changes** to frontend or backend code
2. **Test locally** using the development setup
3. **Commit changes** to your feature branch
4. **Create pull request** - CI/CD will test deployment
5. **Merge to main** - Auto-deploys to production

### Debugging

#### Frontend Issues
- Check browser console for JavaScript errors
- Verify API_BASE URL in script.js
- Test API endpoints directly: `curl http://localhost:7071/api/books`

#### Backend Issues
- Check Azure Functions logs: `func logs` or `azd logs`
- Verify environment variables: `azd env get-values`
- Test function locally: `cd api && npm run dev`

#### Deployment Issues
- Preview changes: `azd provision --preview`
- Check deployment status: `azd show`
- View deployment logs in GitHub Actions

### Testing

#### Manual Testing Checklist
- [ ] Dashboard loads correctly
- [ ] Books display with proper styling
- [ ] Add/edit/delete books functionality works
- [ ] AI chat responds (requires Azure OpenAI setup)
- [ ] Sync functionality works
- [ ] Responsive design on mobile

#### API Testing
```bash
# Test local API
curl http://localhost:7071/api/books
curl -X POST http://localhost:7071/api/books \
  -H "Content-Type: application/json" \
  -d '[]'

# Test production API
curl https://your-app.azurewebsites.net/api/books
```

### Environment Variables

#### Required for Local Development
```env
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
AZURE_OPENAI_API_KEY=your_api_key
AZURE_OPENAI_DEPLOYMENT_NAME=gpt-4
```

#### Automatically Set by azd
- `FUNCTIONS_EXTENSION_VERSION`
- `WEBSITE_NODE_DEFAULT_VERSION` 
- `FUNCTIONS_WORKER_RUNTIME`
- `AzureWebJobsStorage`
- `APPINSIGHTS_INSTRUMENTATIONKEY`

### Troubleshooting Common Issues

#### "Function not found" Error
- Ensure `api/src/functions/books.js` exists
- Check `host.json` configuration
- Restart Azure Functions runtime

#### CORS Errors
- Bicep template includes CORS configuration
- For local development, CORS is handled by Azure Functions runtime

#### Environment Variables Not Loading
```bash
# Check current values
azd env get-values

# Set missing values
azd env set KEY "value"

# Redeploy
azd up
```

#### Build Failures
- Ensure Node.js 18+ is installed
- Clear npm cache: `npm cache clean --force`
- Delete node_modules: `rm -rf api/node_modules && cd api && npm install`

### Performance Tips

#### Frontend Optimization
- Books data is cached in localStorage
- Graceful offline fallback to local JSON
- Minimal JavaScript bundle (vanilla JS)

#### Backend Optimization
- Azure Functions automatically scale
- Cold start mitigation with keep-alive
- Efficient OpenAI API usage

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and test locally
4. Commit with clear messages: `git commit -m "Add: new feature description"`
5. Push and create a pull request

### Additional Resources

- [Azure Functions Documentation](https://docs.microsoft.com/azure/azure-functions/)
- [Azure Developer CLI Guide](https://docs.microsoft.com/azure/developer/azure-developer-cli/)
- [Azure OpenAI Service](https://docs.microsoft.com/azure/cognitive-services/openai/)
