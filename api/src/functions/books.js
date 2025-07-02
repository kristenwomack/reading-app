const { app } = require('@azure/functions');
const OpenAI = require('openai');

// Initialize OpenAI client
const openai = new OpenAI({
  apiKey: process.env.OPENAI_API_KEY
});

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
  
  if (!process.env.OPENAI_API_KEY) {
    return {
      status: 400,
      headers,
      body: JSON.stringify({ error: 'OpenAI API key not configured' })
    };
  }

  try {
    const completion = await openai.chat.completions.create({
      model: "gpt-3.5-turbo",
      messages: [
        {
          role: "system",
          content: `You are a helpful reading assistant. You help manage a personal reading library and provide book recommendations. 

Current library contains ${books?.length || 0} books. When adding books, respond with a JSON object containing the book details in this format:
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

For other queries, provide helpful conversational responses about books and reading.`
        },
        {
          role: "user",
          content: message
        }
      ],
      max_tokens: 500,
      temperature: 0.7
    });

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
      body: JSON.stringify({ error: 'Failed to process chat request' })
    };
  }
}
