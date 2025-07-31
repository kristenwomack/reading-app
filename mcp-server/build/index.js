#!/usr/bin/env node
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { CallToolRequestSchema, ErrorCode, ListToolsRequestSchema, McpError, } from "@modelcontextprotocol/sdk/types.js";
import axios from "axios";
import { z } from "zod";
// Validation schemas
const GetBookByTitleSchema = z.object({
    title: z.string().min(1, { message: "Title cannot be empty" }),
});
const GetAuthorsByNameSchema = z.object({
    name: z.string().min(1, { message: "Author name cannot be empty" }),
});
const GetBookByIdSchema = z.object({
    idType: z
        .string()
        .transform((val) => val.toLowerCase())
        .pipe(z.enum(["isbn", "lccn", "oclc", "olid"], {
        errorMap: () => ({
            message: "idType must be one of: isbn, lccn, oclc, olid",
        }),
    })),
    idValue: z.string().min(1, { message: "idValue cannot be empty" }),
});
const GetAuthorInfoSchema = z.object({
    author_key: z
        .string()
        .min(1, { message: "Author key cannot be empty" })
        .regex(/^OL\d+A$/, {
        message: "Author key must be in the format OL<number>A",
    }),
});
const GetBookCoverSchema = z.object({
    key: z.enum(["ISBN", "OCLC", "LCCN", "OLID", "ID"], {
        errorMap: () => ({
            message: "Key must be one of ISBN, OCLC, LCCN, OLID, ID",
        }),
    }),
    value: z.string().min(1, { message: "Value cannot be empty" }),
    size: z
        .nullable(z.enum(["S", "M", "L"]))
        .optional()
        .transform((val) => val || "L"),
});
const GetAuthorPhotoSchema = z.object({
    olid: z
        .string()
        .min(1, { message: "OLID cannot be empty" })
        .regex(/^OL\d+A$/, {
        message: "OLID must be in the format OL<number>A",
    }),
});
const EnrichBookDataSchema = z.object({
    title: z.string().min(1, { message: "Title cannot be empty" }),
    author: z.string().optional(),
    isbn: z.string().optional(),
    publisher: z.string().optional(),
});
class ReadingTrackerMCPServer {
    server;
    openLibraryApi;
    readingTrackerApi;
    constructor() {
        this.server = new Server({
            name: "reading-tracker-mcp-server",
            version: "1.0.0",
        }, {
            capabilities: {
                resources: {},
                tools: {},
            },
        });
        // Initialize API clients
        this.openLibraryApi = axios.create({
            baseURL: "https://openlibrary.org",
            timeout: 10000,
            headers: {
                'User-Agent': 'ReadingTracker-MCP/1.0 (https://github.com/kristenwomack/reading-app)'
            }
        });
        // Reading Tracker API client (for integration with Azure Functions)
        this.readingTrackerApi = axios.create({
            baseURL: process.env.READING_TRACKER_API_URL || "http://localhost:7071/api",
            timeout: 10000,
        });
        this.setupToolHandlers();
        this.server.onerror = (error) => console.error("[MCP Error]", error);
        process.on("SIGINT", async () => {
            await this.server.close();
            process.exit(0);
        });
    }
    setupToolHandlers() {
        this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
            tools: [
                {
                    name: "search_books_by_title",
                    description: "Search for books by title using Open Library",
                    inputSchema: {
                        type: "object",
                        properties: {
                            title: {
                                type: "string",
                                description: "The title of the book to search for",
                            },
                        },
                        required: ["title"],
                    },
                },
                {
                    name: "search_authors_by_name",
                    description: "Search for authors by name using Open Library",
                    inputSchema: {
                        type: "object",
                        properties: {
                            name: {
                                type: "string",
                                description: "The name of the author to search for",
                            },
                        },
                        required: ["name"],
                    },
                },
                {
                    name: "get_book_by_id",
                    description: "Get detailed book information using an identifier (ISBN, LCCN, OCLC, OLID)",
                    inputSchema: {
                        type: "object",
                        properties: {
                            idType: {
                                type: "string",
                                enum: ["isbn", "lccn", "oclc", "olid"],
                                description: "The type of identifier (ISBN, LCCN, OCLC, OLID)",
                            },
                            idValue: {
                                type: "string",
                                description: "The value of the identifier",
                            },
                        },
                        required: ["idType", "idValue"],
                    },
                },
                {
                    name: "get_author_info",
                    description: "Get detailed author information using Open Library author key",
                    inputSchema: {
                        type: "object",
                        properties: {
                            author_key: {
                                type: "string",
                                description: "The Open Library author key (e.g., OL23919A)",
                            },
                        },
                        required: ["author_key"],
                    },
                },
                {
                    name: "get_book_cover",
                    description: "Get book cover URL using various identifiers",
                    inputSchema: {
                        type: "object",
                        properties: {
                            key: {
                                type: "string",
                                enum: ["ISBN", "OCLC", "LCCN", "OLID", "ID"],
                                description: "The type of identifier",
                            },
                            value: {
                                type: "string",
                                description: "The value of the identifier",
                            },
                            size: {
                                type: "string",
                                enum: ["S", "M", "L"],
                                description: "The size of the cover image",
                            },
                        },
                        required: ["key", "value"],
                    },
                },
                {
                    name: "get_author_photo",
                    description: "Get author photo URL using Open Library Author ID",
                    inputSchema: {
                        type: "object",
                        properties: {
                            olid: {
                                type: "string",
                                description: "The Open Library Author ID (e.g., OL23919A)",
                            },
                        },
                        required: ["olid"],
                    },
                },
                {
                    name: "enrich_book_data",
                    description: "Enrich book data with Open Library information",
                    inputSchema: {
                        type: "object",
                        properties: {
                            title: {
                                type: "string",
                                description: "The title of the book",
                            },
                            author: {
                                type: "string",
                                description: "The author of the book (optional)",
                            },
                            isbn: {
                                type: "string",
                                description: "The ISBN of the book (optional)",
                            },
                            publisher: {
                                type: "string",
                                description: "The publisher of the book (optional)",
                            },
                        },
                        required: ["title"],
                    },
                },
            ],
        }));
        this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
            const { name, arguments: args } = request.params;
            try {
                switch (name) {
                    case "search_books_by_title":
                        return await this.handleSearchBooksByTitle(args);
                    case "search_authors_by_name":
                        return await this.handleSearchAuthorsByName(args);
                    case "get_book_by_id":
                        return await this.handleGetBookById(args);
                    case "get_author_info":
                        return await this.handleGetAuthorInfo(args);
                    case "get_book_cover":
                        return await this.handleGetBookCover(args);
                    case "get_author_photo":
                        return await this.handleGetAuthorPhoto(args);
                    case "enrich_book_data":
                        return await this.handleEnrichBookData(args);
                    default:
                        throw new McpError(ErrorCode.MethodNotFound, `Unknown tool: ${name}`);
                }
            }
            catch (error) {
                if (error instanceof McpError) {
                    throw error;
                }
                throw new McpError(ErrorCode.InternalError, `Tool execution failed: ${error?.message || 'Unknown error'}`);
            }
        });
    }
    async handleSearchBooksByTitle(args) {
        const parseResult = GetBookByTitleSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { title } = parseResult.data;
        try {
            const response = await this.openLibraryApi.get('/search.json', {
                params: { title, limit: 10 }
            });
            if (!response.data.docs || response.data.docs.length === 0) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `No books found matching title: "${title}"`
                        }
                    ]
                };
            }
            const books = response.data.docs.map((book) => ({
                title: book.title,
                authors: book.author_name || [],
                first_publish_year: book.first_publish_year || null,
                open_library_work_key: book.key,
                edition_count: book.edition_count || 0,
                cover_url: book.cover_i ? `https://covers.openlibrary.org/b/id/${book.cover_i}-M.jpg` : undefined,
                isbn: book.isbn ? book.isbn.slice(0, 3) : [],
                publisher: book.publisher ? book.publisher.slice(0, 3) : [],
                language: book.language || [],
                subject: book.subject ? book.subject.slice(0, 5) : []
            }));
            return {
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(books, null, 2)
                    }
                ]
            };
        }
        catch (error) {
            throw new McpError(ErrorCode.InternalError, `Search failed: ${error?.message || 'Unknown error'}`);
        }
    }
    async handleSearchAuthorsByName(args) {
        const parseResult = GetAuthorsByNameSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { name } = parseResult.data;
        try {
            const response = await this.openLibraryApi.get('/search/authors.json', {
                params: { q: name }
            });
            if (!response.data.docs || response.data.docs.length === 0) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `No authors found matching name: "${name}"`
                        }
                    ]
                };
            }
            const authors = response.data.docs.map((author) => ({
                key: author.key,
                name: author.name,
                alternate_names: author.alternate_names || [],
                birth_date: author.birth_date || undefined,
                death_date: author.death_date || undefined,
                top_work: author.top_work || undefined,
                work_count: author.work_count || 0,
                top_subjects: author.top_subjects || []
            }));
            return {
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(authors, null, 2)
                    }
                ]
            };
        }
        catch (error) {
            throw new McpError(ErrorCode.InternalError, `Author search failed: ${error?.message || 'Unknown error'}`);
        }
    }
    async handleGetBookById(args) {
        const parseResult = GetBookByIdSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { idType, idValue } = parseResult.data;
        try {
            const response = await this.openLibraryApi.get(`/api/volumes/brief/${idType}/${idValue}.json`);
            if (!response.data || !response.data.records || Object.keys(response.data.records).length === 0) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `No book found for ${idType}: ${idValue}`
                        }
                    ]
                };
            }
            const recordKey = Object.keys(response.data.records)[0];
            const record = response.data.records[recordKey];
            const data = record.data;
            const bookDetails = {
                title: data.title,
                authors: data.authors ? data.authors.map((a) => a.name) : [],
                publishers: data.publishers ? data.publishers.map((p) => p.name) : [],
                publish_date: data.publish_date || undefined,
                number_of_pages: data.number_of_pages || undefined,
                isbn_13: data.identifiers?.isbn_13 || [],
                isbn_10: data.identifiers?.isbn_10 || [],
                lccn: data.identifiers?.lccn || [],
                oclc: data.identifiers?.oclc || [],
                olid: data.identifiers?.openlibrary || [],
                open_library_edition_key: data.key,
                cover_url: data.cover?.medium || undefined,
                info_url: data.url || undefined,
                preview_url: data.ebooks?.[0]?.preview_url || undefined
            };
            return {
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(bookDetails, null, 2)
                    }
                ]
            };
        }
        catch (error) {
            if (axios.isAxiosError(error) && error.response?.status === 404) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `No book found for ${idType}: ${idValue}`
                        }
                    ]
                };
            }
            throw new McpError(ErrorCode.InternalError, `Failed to get book details: ${error?.message || 'Unknown error'}`);
        }
    }
    async handleGetAuthorInfo(args) {
        const parseResult = GetAuthorInfoSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { author_key } = parseResult.data;
        try {
            const response = await this.openLibraryApi.get(`/authors/${author_key}.json`);
            if (!response.data) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `No author found for key: ${author_key}`
                        }
                    ]
                };
            }
            const author = response.data;
            let bio = null;
            if (author.bio) {
                if (typeof author.bio === 'string') {
                    bio = author.bio;
                }
                else if (author.bio.value) {
                    bio = author.bio.value;
                }
            }
            const authorDetails = {
                name: author.name,
                personal_name: author.personal_name || undefined,
                birth_date: author.birth_date || undefined,
                death_date: author.death_date || undefined,
                bio: bio || undefined,
                alternate_names: author.alternate_names || [],
                photos: author.photos || [],
                key: author.key,
                remote_ids: author.remote_ids || {},
                links: author.links || []
            };
            return {
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(authorDetails, null, 2)
                    }
                ]
            };
        }
        catch (error) {
            if (axios.isAxiosError(error) && error.response?.status === 404) {
                return {
                    content: [
                        {
                            type: "text",
                            text: `Author with key "${author_key}" not found.`
                        }
                    ]
                };
            }
            throw new McpError(ErrorCode.InternalError, `Failed to get author info: ${error?.message || 'Unknown error'}`);
        }
    }
    async handleGetBookCover(args) {
        const parseResult = GetBookCoverSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { key, value, size } = parseResult.data;
        const coverUrl = `https://covers.openlibrary.org/b/${key.toLowerCase()}/${value}-${size}.jpg`;
        return {
            content: [
                {
                    type: "text",
                    text: coverUrl
                }
            ]
        };
    }
    async handleGetAuthorPhoto(args) {
        const parseResult = GetAuthorPhotoSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const { olid } = parseResult.data;
        const photoUrl = `https://covers.openlibrary.org/a/olid/${olid}-L.jpg`;
        return {
            content: [
                {
                    type: "text",
                    text: photoUrl
                }
            ]
        };
    }
    async handleEnrichBookData(args) {
        const parseResult = EnrichBookDataSchema.safeParse(args);
        if (!parseResult.success) {
            const errorMessages = parseResult.error.errors
                .map((e) => `${e.path.join(".")}: ${e.message}`)
                .join(", ");
            throw new McpError(ErrorCode.InvalidParams, `Invalid arguments: ${errorMessages}`);
        }
        const bookData = parseResult.data;
        try {
            // Use the Reading Tracker API's enrich endpoint
            const response = await this.readingTrackerApi.post('/books/enrich', bookData);
            return {
                content: [
                    {
                        type: "text",
                        text: JSON.stringify(response.data.book, null, 2)
                    }
                ]
            };
        }
        catch (error) {
            // Fallback to direct Open Library search if the API is not available
            try {
                const searchResponse = await this.openLibraryApi.get('/search.json', {
                    params: { title: bookData.title, limit: 1 }
                });
                if (searchResponse.data.docs && searchResponse.data.docs.length > 0) {
                    const book = searchResponse.data.docs[0];
                    const enrichedBook = {
                        ...bookData,
                        authors: book.author_name || [],
                        first_publish_year: book.first_publish_year || null,
                        open_library_work_key: book.key,
                        edition_count: book.edition_count || 0,
                        cover_url: book.cover_i ? `https://covers.openlibrary.org/b/id/${book.cover_i}-M.jpg` : null,
                        isbn: book.isbn ? book.isbn[0] : null,
                        publisher: book.publisher ? book.publisher[0] : null,
                        subjects: book.subject ? book.subject.slice(0, 5) : []
                    };
                    return {
                        content: [
                            {
                                type: "text",
                                text: JSON.stringify(enrichedBook, null, 2)
                            }
                        ]
                    };
                }
                else {
                    return {
                        content: [
                            {
                                type: "text",
                                text: JSON.stringify({ ...bookData, enriched: false, reason: "No matches found in Open Library" }, null, 2)
                            }
                        ]
                    };
                }
            }
            catch (fallbackError) {
                throw new McpError(ErrorCode.InternalError, `Book enrichment failed: ${fallbackError?.message || 'Unknown error'}`);
            }
        }
    }
    async run() {
        const transport = new StdioServerTransport();
        await this.server.connect(transport);
        console.error("Reading Tracker MCP Server running on stdio");
    }
}
// Start the server if this file is run directly
if (process.argv[1] === new URL(import.meta.url).pathname) {
    const server = new ReadingTrackerMCPServer();
    server.run().catch(console.error);
}
export { ReadingTrackerMCPServer };
//# sourceMappingURL=index.js.map