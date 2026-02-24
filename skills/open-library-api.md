# Open Library API Skills

A comprehensive guide for agents and developers working with the Open Library API in this reading-app repository.

## Overview

[Open Library](https://openlibrary.org) is an open, editable library catalog that provides free access to book metadata, covers, and full-text content. This skills document enables agents to enrich book data in `books.json` with additional metadata from Open Library.

## Core Concepts

### Identifiers

Open Library uses several identifier types:

| Identifier | Format | Example | Description |
|------------|--------|---------|-------------|
| **ISBN-10** | 10 digits | `0446603775` | Standard book identifier |
| **ISBN-13** | 13 digits | `9780446603775` | Extended ISBN format |
| **OLID (Work)** | `OL{number}W` | `OL27448W` | Open Library Work ID (abstract work) |
| **OLID (Edition)** | `OL{number}M` | `OL7353617M` | Open Library Edition ID (specific edition) |
| **OLID (Author)** | `OL{number}A` | `OL23919A` | Open Library Author ID |

### Works vs Editions

- **Work**: The abstract concept of a book (e.g., "Dawn" by Octavia Butler)
- **Edition**: A specific published version (e.g., 1997 Warner Books paperback)

## API Endpoints

### 1. Books API

#### Lookup by ISBN

```
GET https://openlibrary.org/isbn/{ISBN}.json
```

**Example** (from books.json - Dawn by Octavia Butler):
```
GET https://openlibrary.org/isbn/0446603775.json
```

**Response includes:**
- `title`: Book title
- `authors`: Array of author references `[{key: "/authors/OL23919A"}]`
- `publishers`: Array of publisher names
- `number_of_pages`: Page count
- `publish_date`: Publication date string
- `covers`: Array of cover IDs for Covers API
- `works`: Reference to parent work

#### Lookup by OLID

```
GET https://openlibrary.org/books/{OLID}.json
```

**Example:**
```
GET https://openlibrary.org/books/OL7353617M.json
```

#### Works Endpoint

Get the parent work (for series, editions list, etc.):

```
GET https://openlibrary.org/works/{OLID}.json
```

**Response includes:**
- `title`: Work title
- `description`: Book description/summary
- `subjects`: Subject headings
- `first_publish_date`: Original publication year

#### Get All Editions of a Work

```
GET https://openlibrary.org/works/{OLID}/editions.json
```

Returns paginated list of all editions.

### 2. Search API

#### General Search

```
GET https://openlibrary.org/search.json?q={query}
```

**Query Parameters:**
| Parameter | Description | Example |
|-----------|-------------|---------|
| `q` | General search query | `octavia butler dawn` |
| `title` | Title filter | `Dawn` |
| `author` | Author name filter | `Octavia Butler` |
| `isbn` | ISBN lookup | `0446603775` |
| `page` | Page number (1-indexed) | `1` |
| `limit` | Results per page (max 100) | `10` |
| `sort` | Sort order | `new`, `old`, `editions` |

**Example - Search by Author from books.json:**
```
GET https://openlibrary.org/search.json?author=Octavia+Butler&limit=10
```

**Response Fields (useful for enriching books.json):**
- `docs[].title`: Book title
- `docs[].author_name`: Array of author names
- `docs[].first_publish_year`: Original publication year
- `docs[].number_of_pages_median`: Median page count across editions
- `docs[].isbn`: Array of ISBNs for this work
- `docs[].cover_i`: Cover ID for Covers API
- `docs[].key`: Work key (e.g., `/works/OL27448W`)
- `docs[].subject`: Subject tags

### 3. Authors API

#### Get Author Details

```
GET https://openlibrary.org/authors/{OLID}.json
```

**Example - Octavia Butler:**
```
GET https://openlibrary.org/authors/OL23919A.json
```

**Response includes:**
- `name`: Author's full name
- `bio`: Biography (string or object with `value`)
- `birth_date`: Birth date
- `death_date`: Death date (if applicable)
- `photos`: Array of photo IDs

#### Search Authors

```
GET https://openlibrary.org/search/authors.json?q={name}
```

**Example:**
```
GET https://openlibrary.org/search/authors.json?q=Marie+Benedict
```

#### Get Author's Works

```
GET https://openlibrary.org/authors/{OLID}/works.json
```

### 4. Covers API

Retrieve book cover images. No API key required.

#### URL Format

```
https://covers.openlibrary.org/b/{key}/{value}-{size}.jpg
```

**Key types:**
- `isbn`: ISBN-10 or ISBN-13
- `olid`: Edition OLID
- `id`: Cover ID from API responses

**Sizes:**
- `S`: Small (45×68 px)
- `M`: Medium (90×136 px)
- `L`: Large (180×271 px)

#### Examples

**By ISBN (from books.json - Dawn):**
```
https://covers.openlibrary.org/b/isbn/0446603775-L.jpg
```

**By Cover ID:**
```
https://covers.openlibrary.org/b/id/8314143-M.jpg
```

**Handle Missing Covers:**
Add `?default=false` to get 404 instead of blank image:
```
https://covers.openlibrary.org/b/isbn/0446603775-L.jpg?default=false
```

#### Author Photos

```
https://covers.openlibrary.org/a/olid/{AUTHOR_OLID}-{size}.jpg
```

**Example:**
```
https://covers.openlibrary.org/a/olid/OL23919A-M.jpg
```

### 5. Subjects API

Get books by subject:

```
GET https://openlibrary.org/subjects/{subject}.json
```

**Example:**
```
GET https://openlibrary.org/subjects/science_fiction.json
```

## Usage Examples for Reading App

### Enrich Book Data with Missing ISBNs

Many entries in `books.json` have empty ISBN fields. Use search to find them:

```javascript
// Example: Find ISBN for "Co-Intelligence" by Ethan Mollick
const response = await fetch(
  'https://openlibrary.org/search.json?title=Co-Intelligence&author=Ethan+Mollick&limit=1'
);
const data = await response.json();
const isbn = data.docs[0]?.isbn?.[0];
```

### Get Book Cover URLs

```javascript
// Generate cover URL for a book with ISBN
function getCoverUrl(isbn, size = 'M') {
  return `https://covers.openlibrary.org/b/isbn/${isbn}-${size}.jpg`;
}

// Usage with books.json entry
const dawnCover = getCoverUrl('0446603775', 'L');
```

### Get Missing Page Counts

```javascript
// Fetch page count for book
async function getPageCount(isbn) {
  const response = await fetch(`https://openlibrary.org/isbn/${isbn}.json`);
  const data = await response.json();
  return data.number_of_pages;
}
```

### Get Book Descriptions

```javascript
// Descriptions are on the Work, not Edition
async function getBookDescription(isbn) {
  // First get edition to find work
  const edition = await fetch(`https://openlibrary.org/isbn/${isbn}.json`).then(r => r.json());
  const workKey = edition.works?.[0]?.key;
  
  if (workKey) {
    const work = await fetch(`https://openlibrary.org${workKey}.json`).then(r => r.json());
    return typeof work.description === 'string' 
      ? work.description 
      : work.description?.value;
  }
  return null;
}
```

## Go Implementation Patterns

For the backend in this repository, use Go's standard library:

```go
package openlibrary

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "time"
)

const baseURL = "https://openlibrary.org"

// Client for Open Library API
type Client struct {
    httpClient *http.Client
    userAgent  string
}

// NewClient creates a new Open Library client
func NewClient(userAgent string) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 10 * time.Second},
        userAgent:  userAgent,
    }
}

// Edition represents an Open Library book edition
type Edition struct {
    Title         string   `json:"title"`
    Authors       []Author `json:"authors"`
    Publishers    []string `json:"publishers"`
    NumberOfPages int      `json:"number_of_pages"`
    PublishDate   string   `json:"publish_date"`
    Covers        []int    `json:"covers"`
    Works         []Work   `json:"works"`
}

type Author struct {
    Key string `json:"key"`
}

type Work struct {
    Key string `json:"key"`
}

// GetByISBN retrieves book data by ISBN
func (c *Client) GetByISBN(isbn string) (*Edition, error) {
    url := fmt.Sprintf("%s/isbn/%s.json", baseURL, isbn)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", c.userAgent)
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == http.StatusNotFound {
        return nil, fmt.Errorf("book not found: %s", isbn)
    }
    
    var edition Edition
    if err := json.NewDecoder(resp.Body).Decode(&edition); err != nil {
        return nil, err
    }
    
    return &edition, nil
}

// CoverURL generates a cover image URL
func CoverURL(isbn, size string) string {
    return fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-%s.jpg", isbn, size)
}
```

## JavaScript Implementation Patterns

For the frontend in this repository:

```javascript
// src/openlibrary.js

const BASE_URL = 'https://openlibrary.org';
const COVERS_URL = 'https://covers.openlibrary.org';

/**
 * Fetch book data by ISBN
 * @param {string} isbn - ISBN-10 or ISBN-13
 * @returns {Promise<Object>} Book edition data
 */
export async function getBookByISBN(isbn) {
  const response = await fetch(`${BASE_URL}/isbn/${isbn}.json`);
  if (!response.ok) {
    throw new Error(`Book not found: ${isbn}`);
  }
  return response.json();
}

/**
 * Search for books
 * @param {Object} params - Search parameters
 * @param {string} params.title - Book title
 * @param {string} params.author - Author name
 * @param {number} params.limit - Max results (default 10)
 * @returns {Promise<Object>} Search results
 */
export async function searchBooks({ title, author, limit = 10 }) {
  const params = new URLSearchParams();
  if (title) params.append('title', title);
  if (author) params.append('author', author);
  params.append('limit', limit.toString());
  
  const response = await fetch(`${BASE_URL}/search.json?${params}`);
  return response.json();
}

/**
 * Get cover image URL
 * @param {string} isbn - Book ISBN
 * @param {'S'|'M'|'L'} size - Image size
 * @returns {string} Cover URL
 */
export function getCoverUrl(isbn, size = 'M') {
  return `${COVERS_URL}/b/isbn/${isbn}-${size}.jpg`;
}

/**
 * Check if a cover exists
 * @param {string} isbn - Book ISBN
 * @returns {Promise<boolean>} True if cover exists
 */
export async function hasCover(isbn) {
  const url = `${COVERS_URL}/b/isbn/${isbn}-S.jpg?default=false`;
  const response = await fetch(url, { method: 'HEAD' });
  return response.ok;
}
```

## Best Practices

### Rate Limiting

Open Library does not enforce strict rate limits but requests responsible use:

- **Add delays**: Wait 100-200ms between requests
- **Cache responses**: Store results locally to avoid repeated calls
- **Batch wisely**: Use search for multiple lookups instead of individual ISBN calls

### User-Agent Header

**Always include a User-Agent header** identifying your application:

```
User-Agent: ReadingApp/1.0 (https://github.com/kristenwomack/reading-app)
```

This helps Open Library track usage and prevents blocking.

### Error Handling

Common HTTP status codes:

| Code | Meaning | Action |
|------|---------|--------|
| 200 | Success | Process response |
| 404 | Not found | Book/author doesn't exist in Open Library |
| 429 | Rate limited | Increase delay between requests |
| 500+ | Server error | Retry with exponential backoff |

### Handling Missing Data

Open Library has gaps in coverage. Gracefully handle:

- Missing page counts (`number_of_pages` may be 0 or absent)
- Missing covers (not all books have cover images)
- Missing ISBNs (older/rare books may lack ISBNs)
- Inconsistent date formats (`publish_date` varies: "1987", "April 1987", etc.)

## Data Enrichment Strategy for books.json

Given the current structure in `books.json`, here's a strategy to enrich book data:

1. **Books with ISBNs**: Use direct ISBN lookup for covers and missing page counts
2. **Books without ISBNs**: Use search API with title + author to find matches
3. **Cover images**: Generate URLs using ISBN or fallback to cover ID from search
4. **Descriptions**: Fetch from Works API (linked from Edition response)
5. **Author details**: Fetch using author OLID from Edition response

### Priority Fields to Enrich

| Current Field | Open Library Source |
|--------------|---------------------|
| `Number of Pages` | `edition.number_of_pages` |
| `ISBN` / `ISBN13` | `search.docs[].isbn` |
| Cover Image (new) | `covers.openlibrary.org/b/isbn/{isbn}-M.jpg` |
| Description (new) | `work.description` |
| Subjects (new) | `work.subjects` or `search.docs[].subject` |

## API Documentation Links

- **Official Developer Center**: https://openlibrary.org/developers/api
- **Books API**: https://openlibrary.org/dev/docs/api/books
- **Search API**: https://openlibrary.org/dev/docs/api/search
- **Covers API**: https://openlibrary.org/dev/docs/api/covers
- **Swagger Documentation**: https://internetarchive.github.io/openlibrary-api/

## Testing API Calls

Use curl to test endpoints locally:

```bash
# Get book by ISBN
curl -s "https://openlibrary.org/isbn/0446603775.json" | jq

# Search for book
curl -s "https://openlibrary.org/search.json?title=Dawn&author=Octavia+Butler&limit=1" | jq

# Test cover exists
curl -I "https://covers.openlibrary.org/b/isbn/0446603775-M.jpg?default=false"
```

---

*Last updated: 2026-02-05*
