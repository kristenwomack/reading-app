const { app } = require('@azure/functions');
const { OpenAIClient, AzureKeyCredential } = require('@azure/openai');
const { 
  searchBooksByTitle, 
  searchAuthorsByName, 
  getBookById, 
  getAuthorInfo,
  getBookCoverUrl,
  getAuthorPhotoUrl,
  enrichBookData 
} = require('../services/openLibraryService');

// Initialize Azure OpenAI client for GPT-4o
const openaiClient = new OpenAIClient(
  process.env.AZURE_OPENAI_ENDPOINT,
  new AzureKeyCredential(process.env.AZURE_OPENAI_API_KEY)
);

// In-memory storage for demo (in production, use Azure Storage/CosmosDB)
let booksData = [];

app.http('books', {
  methods: ['GET', 'POST', 'PUT'],
  authLevel: 'anonymous',
  route: 'books/{action?}',
  handler: async (request, context) => {
    const action = request.params.action;
    
    // Enable CORS
    const headers = {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
      'Content-Type': 'application/json'
    };

    if (request.method === 'OPTIONS') {
      return { status: 200, headers };
    }

    try {
      switch (action) {
        case 'sync':
          return await handleSync(request, headers);
        case 'chat':
          return await handleChat(request, headers);
        case 'search':
          return await handleSearch(request, headers);
        case 'enrich':
          return await handleEnrich(request, headers);
        case 'author':
          return await handleAuthorSearch(request, headers);
        case 'cover':
          return await handleCoverUrl(request, headers);
        default:
          return await handleBooks(request, headers);
      }
    } catch (error) {
      context.log('Error:', error);
      return {
        status: 500,
        headers,
        body: JSON.stringify({ error: error.message })
      };
    }
  }
});

// Health check endpoint
app.http('health', {
  methods: ['GET'],
  authLevel: 'anonymous',
  route: 'health',
  handler: async (request, context) => {
    const headers = {
      'Access-Control-Allow-Origin': '*',
      'Content-Type': 'application/json'
    };

    try {
      const health = {
        status: 'healthy',
        timestamp: new Date().toISOString(),
        version: '1.0.0',
        environment: process.env.AZURE_FUNCTIONS_ENVIRONMENT || 'development',
        azureOpenAI: {
          configured: !!(process.env.AZURE_OPENAI_ENDPOINT && process.env.AZURE_OPENAI_API_KEY),
          endpoint: process.env.AZURE_OPENAI_ENDPOINT ? 'configured' : 'missing',
          model: 'GPT-4o'
        }
      };

      return {
        status: 200,
        headers,
        body: JSON.stringify(health)
      };
    } catch (error) {
      return {
        status: 503,
        headers,
        body: JSON.stringify({
          status: 'unhealthy',
          error: error.message,
          timestamp: new Date().toISOString()
        })
      };
    }
  }
});

async function handleBooks(request, headers) {
  switch (request.method) {
    case 'GET':
      return {
        status: 200,
        headers,
        body: JSON.stringify(booksData)
      };
      
    case 'POST':
      const newBooks = await request.json();
      booksData = newBooks;
      return {
        status: 200,
        headers,
        body: JSON.stringify({ success: true, count: booksData.length })
      };
      
    default:
      return {
        status: 405,
        headers,
        body: JSON.stringify({ error: 'Method not allowed' })
      };
  }
}

async function handleSync(request, headers) {
  // In a real implementation, this would sync with a database
  return {
    status: 200,
    headers,
    body: JSON.stringify({ 
      success: true, 
      lastSync: new Date().toISOString(),
      bookCount: booksData.length 
    })
  };
}

async function handleChat(request, headers) {
  const { message, books } = await request.json();
  
  if (!process.env.AZURE_OPENAI_ENDPOINT || !process.env.AZURE_OPENAI_API_KEY) {
    return {
      status: 400,
      headers,
      body: JSON.stringify({ error: 'Azure OpenAI not configured' })
    };
  }

  try {
    const completion = await openaiClient.getChatCompletions(
      process.env.AZURE_OPENAI_DEPLOYMENT_NAME || "gpt-4o",
      [
        {
          role: "system",
          content: `You are a helpful reading assistant with access to Open Library data. You help manage a personal reading library and provide book recommendations. 

Current library contains ${books?.length || 0} books. When adding books, you can search Open Library for accurate book details. 

For book-related queries, you can:
1. Search for books by title using Open Library
2. Find author information
3. Get book covers and metadata
4. Provide reading recommendations

When adding books, respond with a JSON object containing the book details in this format:
{
  "action": "add_book",
  "book": {
    "title": "Book Title",
    "author": "Author Name", 
    "isbn": "1234567890",
    "pages": 300,
    "publisher": "Publisher Name",
    "year": 2024,
    "status": "to-read",
    "dateAdded": "2024-01-15",
    "dateStarted": null,
    "dateCompleted": null,
    "rating": null,
    "notes": ""
  }
}

For searches, respond with:
{
  "action": "search_books",
  "query": "search term"
}

For other queries, provide helpful conversational responses about books and reading.`
        },
        {
          role: "user",
          content: message
        }
      ],
      {
        maxTokens: 500,
        temperature: 0.7
      }
    );

    const response = completion.choices[0].message.content;
    
    return {
      status: 200,
      headers,
      body: JSON.stringify({ response })
    };
  } catch (error) {
    return {
      status: 500,
      headers,
      body: JSON.stringify({ error: `Failed to process chat request: ${error.message}` })
    };
  }
}

// New Open Library integration handlers
async function handleSearch(request, headers) {
  try {
    const url = new URL(request.url);
    const title = url.searchParams.get('title');
    const author = url.searchParams.get('author');

    if (!title && !author) {
      return {
        status: 400,
        headers,
        body: JSON.stringify({ error: 'Either title or author parameter is required' })
      };
    }

    let results = [];

    if (title) {
      results = await searchBooksByTitle(title);
    } else if (author) {
      const authorResults = await searchAuthorsByName(author);
      results = authorResults;
    }

    return {
      status: 200,
      headers,
      body: JSON.stringify({ results, total: results.length })
    };
  } catch (error) {
    return {
      status: 500,
      headers,
      body: JSON.stringify({ error: `Search failed: ${error.message}` })
    };
  }
}

async function handleEnrich(request, headers) {
  try {
    const bookData = await request.json();
    
    if (!bookData.title) {
      return {
        status: 400,
        headers,
        body: JSON.stringify({ error: 'Book title is required for enrichment' })
      };
    }

    const enrichedBook = await enrichBookData(bookData);

    return {
      status: 200,
      headers,
      body: JSON.stringify({ book: enrichedBook })
    };
  } catch (error) {
    return {
      status: 500,
      headers,
      body: JSON.stringify({ error: `Enrichment failed: ${error.message}` })
    };
  }
}

async function handleAuthorSearch(request, headers) {
  try {
    const url = new URL(request.url);
    const name = url.searchParams.get('name');
    const key = url.searchParams.get('key');

    if (!name && !key) {
      return {
        status: 400,
        headers,
        body: JSON.stringify({ error: 'Either name or key parameter is required' })
      };
    }

    let result;

    if (key) {
      result = await getAuthorInfo(key);
      if (!result) {
        return {
          status: 404,
          headers,
          body: JSON.stringify({ error: 'Author not found' })
        };
      }
    } else {
      result = await searchAuthorsByName(name);
    }

    return {
      status: 200,
      headers,
      body: JSON.stringify({ result })
    };
  } catch (error) {
    return {
      status: 500,
      headers,
      body: JSON.stringify({ error: `Author search failed: ${error.message}` })
    };
  }
}

async function handleCoverUrl(request, headers) {
  try {
    const url = new URL(request.url);
    const key = url.searchParams.get('key');
    const value = url.searchParams.get('value');
    const size = url.searchParams.get('size') || 'L';
    const type = url.searchParams.get('type'); // 'book' or 'author'

    if (!key || !value) {
      return {
        status: 400,
        headers,
        body: JSON.stringify({ error: 'Both key and value parameters are required' })
      };
    }

    let coverUrl;
    
    if (type === 'author') {
      coverUrl = getAuthorPhotoUrl(value);
    } else {
      coverUrl = getBookCoverUrl(key, value, size);
    }

    return {
      status: 200,
      headers,
      body: JSON.stringify({ coverUrl })
    };
  } catch (error) {
    return {
      status: 500,
      headers,
      body: JSON.stringify({ error: `Cover URL generation failed: ${error.message}` })
    };
  }
}
