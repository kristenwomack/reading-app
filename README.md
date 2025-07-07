# 📚 Reading Tracker App

A reading dashboard hosted on GitHub Pages with an Azure backend integration and AI-powered book management. Track your reading progress, manage your library, and interact with an intelligent reading assistant.

![Reading Tracker Dashboard](https://img.shields.io/badge/Status-Active-success)
![GitHub Pages](https://img.shields.io/badge/Deploy-GitHub%20Pages-blue)
![Azure Functions](https://img.shields.io/badge/Backend-Azure%20Functions-0078d4)
![Azure OpenAI](https://img.shields.io/badge/AI-Azure%20OpenAI-0078d4)

## ✨ Features

- **📊 Reading Dashboard**: GitHub-styled interface with progress tracking
- **🎯 Goal Management**: Set and track annual reading goals with visual progress bars
- **📖 Book Management**: Add, track, and manage your reading library (admin mode)
- **👁️ Read-Only Mode**: Clean, public dashboard perfect for sharing
- **🔓 Admin Mode**: Full functionality with URL parameter (`?admin=true`)
- **🤖 AI Assistant**: Natural language book management powered by Azure OpenAI
- **☁️ Cloud Sync**: Azure Functions backend for data persistence
- **📱 Responsive Design**: Works seamlessly on desktop and mobile browsers
- **🚀 Easy Deployment**: One-click deployment to GitHub Pages + Azure

## 🏗️ Architecture

### Frontend (GitHub Pages)
```
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Pages Frontend                    │
├─────────────────────────────────────────────────────────────┤
│  • index.html - Main dashboard interface                    │
│  • style.css - GitHub-inspired dark theme                   │
│  • script.js - JavaScript functionality & Azure integration │
│  • books.json - Local data fallback                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTPS API Calls
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Azure Functions Backend                  │
├─────────────────────────────────────────────────────────────┤
│  • books.js - Main API endpoints                            │
│  • /books - CRUD operations for book data                   │
│  • /books/sync - Data synchronization                       │
│  • /books/chat - OpenAI integration                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ OpenAI API
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Azure OpenAI Service                   │
├─────────────────────────────────────────────────────────────┤
│  • Natural language book management                         │
│  • Reading recommendations                                  │
│  • Intelligent data extraction                              │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow
```
User Input → Frontend Dashboard → Azure Functions → Azure OpenAI API
    ↓              ↓                    ↓
Local Storage ← Book Data ← Azure Storage ← AI Processing
```

### Technology Stack

**Frontend:**
- HTML5 + CSS3 (GitHub-styled dark theme)
- Vanilla JavaScript (ES6+)
- GitHub Pages hosting
- Responsive design with CSS Grid/Flexbox

**Backend:**
- Azure Functions (Node.js 18)
- Azure OpenAI Service
- Azure Storage (future enhancement)
- CORS-enabled REST API

**DevOps:**
- GitHub Actions for CI/CD
- Azure CLI for infrastructure management
- Automated deployment pipeline

## 🚀 Quick Start

### Prerequisites
- [Azure Developer CLI (azd)](https://aka.ms/install-azd) installed
- [Node.js 18+](https://nodejs.org/) installed
- [Git](https://git-scm.com/) installed
- Azure subscription
- Azure OpenAI Service deployed

### 1. Clone & Setup
```bash
git clone https://github.com/yourusername/reading-app.git
cd reading-app
```

### 2. Deploy Everything with azd (2 minutes)
```bash
# Initialize azd project (if not already done)
azd init

# Deploy all Azure resources and application
azd up
```

This single command will:
- Create Azure resource group
- Provision Azure Functions
- Set up storage account
- Deploy Azure OpenAI integration
- Configure Application Insights
- Deploy your functions
- Output the API URL

### 3. Configure Environment Variables
```bash
# Set Azure OpenAI configuration
azd env set AZURE_OPENAI_ENDPOINT "https://your-openai-resource.openai.azure.com/"
azd env set AZURE_OPENAI_API_KEY "your_azure_openai_key"
azd env set AZURE_OPENAI_DEPLOYMENT_NAME "gpt-4"

# Redeploy with new configuration
azd up
```

### 5. Deploy to GitHub Pages
1. Push to your GitHub repository
2. Configure GitHub Actions secrets:
   - `AZURE_CREDENTIALS`: Azure service principal JSON
   - `AZURE_OPENAI_ENDPOINT`: Your Azure OpenAI endpoint
   - `AZURE_OPENAI_API_KEY`: Your Azure OpenAI key  
   - `AZURE_OPENAI_DEPLOYMENT_NAME`: Model deployment name
3. The workflow will automatically deploy using azd

🎉 **Your app will be live at:** `https://yourusername.github.io/reading-app`

## 📖 Documentation

- [📋 Deployment Guide](DEPLOYMENT.md) - Complete deployment instructions
- [🔧 Development Guide](DEVELOPMENT.md) - Local development and API docs
- [🔓 Admin Mode Guide](ADMIN-MODE.md) - Admin vs read-only configuration
- [⚡ azd Migration Summary](AZD-MIGRATION.md) - Modern deployment approach

## 📁 Project Structure

```
reading-app/
├── 📄 index.html              # Main dashboard HTML
├── 🎨 style.css               # GitHub-styled CSS theme
├── ⚡ script.js               # Frontend JavaScript logic
├── 📚 books.json              # Sample/fallback book data
├── 📋 prompt.txt              # AI assistant prompt template
├── 🔧 setup-azure.sh          # Azure deployment script
├── 📖 DEPLOYMENT.md           # Detailed deployment guide
├── 📊 reading-data-manager-prompt.txt  # Data management workflow
├── 🚀 .github/workflows/      # GitHub Actions CI/CD
│   └── deploy.yml             # Automated deployment pipeline
└── ☁️ api/                    # Azure Functions backend
    ├── 📦 package.json        # Node.js dependencies
    ├── ⚙️ host.json           # Azure Functions configuration
    └── 🔧 src/functions/
        └── books.js           # Main API endpoints
```

## 🛠️ API Endpoints

### Books Management
```
GET    /api/books              # Get all books
POST   /api/books              # Save books data
GET    /api/books/sync         # Sync status
POST   /api/books/chat         # AI chat assistant
```

### Example Usage
```javascript
// Get books
const response = await fetch('https://your-app.azurewebsites.net/api/books');
const books = await response.json();

// Chat with AI
const response = await fetch('https://your-app.azurewebsites.net/api/books/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    message: "I finished reading Orbital by Samantha Harvey",
    books: currentBooks
  })
});
```

## 🎯 Dashboard Features

### 📊 Progress Tracking
- Annual reading goal with visual progress bars
- Monthly statistics and trends
- Pages per day calculations
- Reading activity timeline

### 📚 Book Management
- Add books manually or via AI chat
- Track reading status (want-to-read, currently-reading, read)
- Book details with metadata
- Reading progress tracking

### 🤖 AI Assistant
- Natural language book management
- Smart book recommendations
- Automatic data extraction
- Conversational interface

### 🎨 Interface Design
- GitHub-inspired dark theme
- Responsive grid layout
- Smooth animations and transitions
- Mobile-optimized experience

## 💰 Cost Estimate

**Monthly Operating Costs:**
- **Azure Functions**: FREE (under 1M executions)
- **Azure Storage**: ~$0.05/month (small data)
- **GitHub Pages**: FREE
- **Azure OpenAI**: ~$2-5/month (moderate usage)

**Total: ~$2-5/month** for a fully-featured reading tracker!

## 🔧 Configuration

### Environment Variables (Azure Functions)
```bash
AZURE_OPENAI_ENDPOINT=https://your-openai-resource.openai.azure.com/
AZURE_OPENAI_API_KEY=your_azure_openai_key
AZURE_OPENAI_DEPLOYMENT_NAME=gpt-35-turbo
```

### Local Development
```bash
# Start local server
python3 -m http.server 8000

# Access at http://localhost:8000
```

## 🚀 Deployment Options

### Option 1: GitHub Actions (Recommended)
- Automatic deployment on push to main
- Integrated CI/CD pipeline
- Zero-configuration setup

### Option 2: Manual Deployment
```bash
# Deploy functions
cd api && func azure functionapp publish your-app-name

# Deploy frontend
# Push to GitHub and enable Pages
```

## 🔍 Troubleshooting

### Common Issues

**1. CORS Errors**
- Ensure Azure Functions has CORS enabled
- Check API endpoint URLs in `script.js`

**2. Azure OpenAI API Errors**
- Verify Azure OpenAI resource is deployed and accessible
- Check deployment name matches your configuration
- Ensure API key and endpoint are correctly set

**3. GitHub Pages Not Updating**
- Check Actions tab for deployment status
- Ensure `index.html` is in repository root

**4. Books Not Loading**
- Check browser console for errors
- Verify `books.json` format is valid
- Test Azure Functions endpoint directly

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- GitHub for free hosting via GitHub Pages
- Microsoft Azure for serverless backend infrastructure
- Azure OpenAI for intelligent book management capabilities
- The open-source community for inspiration and tools

---

**Built for book lovers who want to track their reading journey with AI.**

### 🔗 Quick Links
- [Live Demo](https://yourusername.github.io/reading-app)
- [Deployment Guide](DEPLOYMENT.md)
- [Azure Functions Setup](azure-functions-setup.md)
- [Data Management Workflow](reading-data-manager-prompt.txt)