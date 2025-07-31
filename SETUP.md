# Reading Tracker - Azure Functions + MCP Server Integration

This project combines a **Reading Tracker web app** with **Azure Functions backend** and an **MCP Server** for AI-powered book management.

## Project Structure

```
reading-app/
├── api/                          # Azure Functions API
│   ├── src/
│   │   ├── functions/books.js    # Main API endpoints
│   │   └── services/openLibraryService.js  # Open Library integration
│   ├── package.json
│   └── host.json
├── mcp-server/                   # MCP Server for AI integration
│   ├── src/index.ts              # MCP server implementation
│   ├── package.json
│   └── tsconfig.json
├── index.html                    # Frontend web app
├── script.js                     # Frontend JavaScript
├── style.css                     # Styling
└── books.json                    # Sample data
```

## Features

### ✅ Implemented Features

#### Azure Functions API (`/api`)
- **Chat endpoint** (`POST /api/books/chat`) - OpenAI integration for reading recommendations
- **Book sync** (`POST /api/books/sync`) - Sync reading data
- **Book search** (`POST /api/books/search`) - Search Open Library by title
- **Book enrichment** (`POST /api/books/enrich`) - Enhance book data with Open Library metadata
- **Author search** (`POST /api/books/author`) - Find authors by name
- **Cover URLs** (`POST /api/books/cover`) - Get book cover image URLs

#### Open Library Service (`/api/src/services/openLibraryService.js`)
- Book search by title with metadata
- Author search and information retrieval
- Book data enrichment
- Cover image URL generation
- Detailed book information by ISBN/OLID

#### MCP Server (`/mcp-server`)
- **7 AI Tools** for book discovery and management
- TypeScript implementation with proper validation
- Open Library API integration
- Integration hooks for Reading Tracker API

#### Frontend Web App
- Reading list management
- Book progress tracking
- Clean, responsive design
- Local storage for demo data

### 🔄 Integration Points

1. **Frontend ↔ Azure Functions** - REST API calls
2. **Azure Functions ↔ Open Library** - External API integration
3. **MCP Server ↔ Open Library** - Direct API integration
4. **MCP Server ↔ Azure Functions** - Optional enhanced integration

## Setup Instructions

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Azure Functions Core Tools (for local development)
- TypeScript (for MCP server development)

### 1. Azure Functions API Setup

```bash
# Install Azure Functions Core Tools (one-time setup)
npm install -g azure-functions-core-tools@4 --unsafe-perm true

# Navigate to API directory
cd api

# Install dependencies
npm install

# Start the local development server
func host start
```

The API will be available at `http://localhost:7071`

### 2. MCP Server Setup

```bash
# Navigate to MCP server directory
cd mcp-server

# Install dependencies
npm install

# Build TypeScript code
npm run build

# Test the server
npm start
```

### 3. Frontend Setup

Simply open `index.html` in a web browser, or serve it locally:

```bash
# Simple HTTP server (Python)
python3 -m http.server 8000

# Or with Node.js
npx serve .
```

## API Endpoints

### Azure Functions Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/books/chat` | Get AI reading recommendations |
| POST | `/api/books/sync` | Sync reading progress data |
| POST | `/api/books/search` | Search books by title |
| POST | `/api/books/enrich` | Enrich book with Open Library data |
| POST | `/api/books/author` | Search authors by name |
| POST | `/api/books/cover` | Get book cover URL |

### Example API Calls

#### Search for books:
```bash
curl -X POST http://localhost:7071/api/books/search \
  -H "Content-Type: application/json" \
  -d '{"title": "Harry Potter"}'
```

#### Enrich book data:
```bash
curl -X POST http://localhost:7071/api/books/enrich \
  -H "Content-Type: application/json" \
  -d '{"title": "The Hobbit", "author": "J.R.R. Tolkien"}'
```

## MCP Server Tools

The MCP server provides these tools for AI assistants:

1. **search_books_by_title** - Find books by title
2. **search_authors_by_name** - Find authors by name  
3. **get_book_by_id** - Get detailed book info by ID
4. **get_author_info** - Get author details
5. **get_book_cover** - Get cover image URLs
6. **get_author_photo** - Get author photo URLs
7. **enrich_book_data** - Enhance book metadata

## Technology Stack

- **Frontend**: Vanilla JavaScript, HTML5, CSS3
- **Backend**: Azure Functions (Node.js 18)
- **MCP Server**: TypeScript, Model Context Protocol SDK
- **External APIs**: Open Library API, OpenAI API
- **Validation**: Zod schema validation
- **HTTP Client**: Axios

## Data Sources

- **Open Library API** - Book and author metadata
- **OpenAI API** - AI-powered reading recommendations
- **Local Storage** - Demo reading list data

## Environment Variables

### Azure Functions
- `AZURE_OPENAI_ENDPOINT` - Your Azure OpenAI endpoint
- `AZURE_OPENAI_API_KEY` - Your Azure OpenAI API key
- `AZURE_OPENAI_DEPLOYMENT_NAME` - Your deployment name (default: gpt-4)

### MCP Server
- `READING_TRACKER_API_URL` - Reading Tracker API URL (default: http://localhost:7071/api)

## Next Steps

### Immediate Enhancements
1. **Deploy to Azure** - Use Azure Functions for production deployment
2. **Database Integration** - Replace JSON file with Azure Cosmos DB or SQL Database
3. **Authentication** - Add user authentication and authorization
4. **Enhanced UI** - Add more interactive features and better styling

### Advanced Features
1. **AI Recommendations** - Enhanced book recommendation engine
2. **Social Features** - Reading groups and book sharing
3. **Progress Analytics** - Reading statistics and insights
4. **Mobile App** - React Native or Flutter mobile application

## Contributing

This is a learning project that demonstrates:
- Azure Functions development
- Model Context Protocol (MCP) server implementation
- Open Library API integration
- Modern web development patterns

Feel free to explore, modify, and enhance any part of the codebase!

## License

MIT License - feel free to use this code for learning and development.
