#!/usr/bin/env node

// Simple test server for the Azure Functions API
const http = require('http');
const url = require('url');
const querystring = require('querystring');

// Import our books function
const booksFunction = require('./api/src/functions/books.js');

const PORT = 7071;

const server = http.createServer(async (req, res) => {
    // Enable CORS
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
    
    if (req.method === 'OPTIONS') {
        res.writeHead(200);
        res.end();
        return;
    }

    const parsedUrl = url.parse(req.url, true);
    const path = parsedUrl.pathname;

    console.log(`${req.method} ${path}`);

    // Handle API routes
    if (path.startsWith('/api/books')) {
        let body = '';
        
        req.on('data', chunk => {
            body += chunk.toString();
        });
        
        req.on('end', async () => {
            try {
                // Create a mock Azure Functions context
                const context = {
                    log: console.log,
                    res: {}
                };

                // Create a mock request object
                const request = {
                    method: req.method,
                    url: req.url,
                    headers: req.headers,
                    body: body ? JSON.parse(body) : undefined,
                    params: {},
                    query: parsedUrl.query || {}
                };

                // Call the books function
                await booksFunction.handler(request, context);

                // Send the response
                res.writeHead(context.res.status || 200, {
                    'Content-Type': 'application/json'
                });
                res.end(JSON.stringify(context.res.body || context.res));

            } catch (error) {
                console.error('Error:', error);
                res.writeHead(500, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ error: error.message }));
            }
        });
    } else {
        // Serve static files for frontend testing
        const fs = require('fs');
        const path = require('path');
        
        let filePath = path.join(__dirname, path === '/' ? 'index.html' : parsedUrl.pathname);
        
        try {
            const data = fs.readFileSync(filePath);
            const ext = path.extname(filePath);
            const contentType = {
                '.html': 'text/html',
                '.js': 'application/javascript',
                '.css': 'text/css',
                '.json': 'application/json'
            }[ext] || 'text/plain';
            
            res.writeHead(200, { 'Content-Type': contentType });
            res.end(data);
        } catch (err) {
            res.writeHead(404);
            res.end('File not found');
        }
    }
});

server.listen(PORT, () => {
    console.log(`🚀 Reading Tracker API Server running on http://localhost:${PORT}`);
    console.log(`📚 API endpoints available at http://localhost:${PORT}/api/books/*`);
    console.log(`🌐 Frontend available at http://localhost:${PORT}`);
    console.log('\n📋 Available API endpoints:');
    console.log('  POST /api/books/search - Search books by title');
    console.log('  POST /api/books/enrich - Enrich book data');
    console.log('  POST /api/books/author - Search authors by name');
    console.log('  POST /api/books/cover - Get book cover URLs');
    console.log('  POST /api/books/chat - AI reading recommendations');
    console.log('  POST /api/books/sync - Sync reading data');
    console.log('\n🧪 Test with: curl -X POST http://localhost:7071/api/books/search -H "Content-Type: application/json" -d \'{"title": "Harry Potter"}\'');
});
