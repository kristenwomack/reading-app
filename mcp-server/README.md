# Reading Tracker MCP Server

A Model Context Protocol (MCP) server that provides AI assistants with powerful book search and management capabilities using the Open Library API.

## Features

The MCP server provides 7 powerful tools for book discovery and management:

### Book Search & Discovery
- **`search_books_by_title`** - Search for books by title
- **`search_authors_by_name`** - Find authors by name
- **`get_book_by_id`** - Get detailed book info using ISBN, LCCN, OCLC, or OLID
- **`get_author_info`** - Get detailed author information

### Visual Content
- **`get_book_cover`** - Get book cover image URLs
- **`get_author_photo`** - Get author photo URLs

### Data Enhancement
- **`enrich_book_data`** - Enhance basic book information with Open Library metadata

## Installation

1. Install dependencies:
```bash
npm install
```

2. Build the TypeScript code:
```bash
npm run build
```

3. Test the server:
```bash
npm start
```

## Integration with Reading Tracker

This MCP server is designed to work with the Reading Tracker Azure Functions API located in the `../api` directory. It can:

- Enrich book data using the `/books/enrich` endpoint
- Integrate with your existing reading list functionality
- Provide rich metadata for books in your collection

## Environment Variables

- `READING_TRACKER_API_URL` - URL of your Reading Tracker API (defaults to `http://localhost:7071/api`)

## Usage with AI Assistants

Once connected to an MCP-compatible AI assistant, you can:

- Ask to search for books: "Find books by Stephen King"
- Get detailed information: "Get details for ISBN 9780439708180"
- Enrich your reading data: "Enrich information for 'The Hobbit'"
- Find author information: "Tell me about author OL23919A"

## API Integration

The server integrates with:
- **Open Library API** - For book and author data
- **Reading Tracker API** - For enhanced book management (optional)

## Tools Reference

### search_books_by_title
```json
{
  "title": "Harry Potter"
}
```

### search_authors_by_name
```json
{
  "name": "J.K. Rowling"
}
```

### get_book_by_id
```json
{
  "idType": "isbn",
  "idValue": "9780439708180"
}
```

### get_author_info
```json
{
  "author_key": "OL23919A"
}
```

### get_book_cover
```json
{
  "key": "ISBN",
  "value": "9780439708180",
  "size": "L"
}
```

### get_author_photo
```json
{
  "olid": "OL23919A"
}
```

### enrich_book_data
```json
{
  "title": "The Hobbit",
  "author": "J.R.R. Tolkien",
  "isbn": "9780547928227"
}
```

## License

MIT License - see the main project for details.
