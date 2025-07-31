const axios = require('axios');

// Create axios instance for Open Library API
const openLibraryApi = axios.create({
  baseURL: 'https://openlibrary.org',
  timeout: 10000,
  headers: {
    'User-Agent': 'ReadingTracker/1.0 (https://github.com/kristenwomack/reading-app)'
  }
});

/**
 * Search for books by title
 * @param {string} title - Book title to search for
 * @returns {Promise<Array>} Array of book results
 */
async function searchBooksByTitle(title) {
  try {
    const response = await openLibraryApi.get('/search.json', {
      params: {
        title: title,
        limit: 10
      }
    });

    if (!response.data.docs || response.data.docs.length === 0) {
      return [];
    }

    return response.data.docs.map(book => ({
      title: book.title,
      authors: book.author_name || [],
      first_publish_year: book.first_publish_year || null,
      open_library_work_key: book.key,
      edition_count: book.edition_count || 0,
      cover_url: book.cover_i ? `https://covers.openlibrary.org/b/id/${book.cover_i}-M.jpg` : null,
      isbn: book.isbn ? book.isbn.slice(0, 3) : [], // Limit to first 3 ISBNs
      publisher: book.publisher ? book.publisher.slice(0, 3) : [], // Limit to first 3 publishers
      language: book.language || [],
      subject: book.subject ? book.subject.slice(0, 5) : [] // Limit to first 5 subjects
    }));
  } catch (error) {
    console.error('Error searching books by title:', error.message);
    throw new Error(`Failed to search books: ${error.message}`);
  }
}

/**
 * Search for authors by name
 * @param {string} name - Author name to search for
 * @returns {Promise<Array>} Array of author results
 */
async function searchAuthorsByName(name) {
  try {
    const response = await openLibraryApi.get('/search/authors.json', {
      params: {
        q: name
      }
    });

    if (!response.data.docs || response.data.docs.length === 0) {
      return [];
    }

    return response.data.docs.map(author => ({
      key: author.key,
      name: author.name,
      alternate_names: author.alternate_names || [],
      birth_date: author.birth_date || null,
      death_date: author.death_date || null,
      top_work: author.top_work || null,
      work_count: author.work_count || 0,
      top_subjects: author.top_subjects || []
    }));
  } catch (error) {
    console.error('Error searching authors by name:', error.message);
    throw new Error(`Failed to search authors: ${error.message}`);
  }
}

/**
 * Get detailed book information by identifier
 * @param {string} idType - Type of identifier (isbn, lccn, oclc, olid)
 * @param {string} idValue - Value of the identifier
 * @returns {Promise<Object|null>} Book details or null if not found
 */
async function getBookById(idType, idValue) {
  try {
    const response = await openLibraryApi.get(`/api/volumes/brief/${idType}/${idValue}.json`);

    if (!response.data || !response.data.records || Object.keys(response.data.records).length === 0) {
      return null;
    }

    // Get the first record
    const recordKey = Object.keys(response.data.records)[0];
    const record = response.data.records[recordKey];
    const data = record.data;

    return {
      title: data.title,
      authors: data.authors ? data.authors.map(a => a.name) : [],
      publishers: data.publishers ? data.publishers.map(p => p.name) : [],
      publish_date: data.publish_date || null,
      number_of_pages: data.number_of_pages || null,
      isbn_13: data.identifiers?.isbn_13 || [],
      isbn_10: data.identifiers?.isbn_10 || [],
      lccn: data.identifiers?.lccn || [],
      oclc: data.identifiers?.oclc || [],
      olid: data.identifiers?.openlibrary || [],
      open_library_edition_key: data.key,
      cover_url: data.cover?.medium || null,
      info_url: data.url || null,
      preview_url: data.ebooks?.[0]?.preview_url || null
    };
  } catch (error) {
    if (error.response && error.response.status === 404) {
      return null;
    }
    console.error('Error getting book by ID:', error.message);
    throw new Error(`Failed to get book details: ${error.message}`);
  }
}

/**
 * Get author information by Open Library key
 * @param {string} authorKey - Open Library author key (e.g., OL23919A)
 * @returns {Promise<Object|null>} Author details or null if not found
 */
async function getAuthorInfo(authorKey) {
  try {
    const response = await openLibraryApi.get(`/authors/${authorKey}.json`);

    if (!response.data) {
      return null;
    }

    const author = response.data;
    
    // Handle bio field which might be an object
    let bio = null;
    if (author.bio) {
      if (typeof author.bio === 'string') {
        bio = author.bio;
      } else if (author.bio.value) {
        bio = author.bio.value;
      }
    }

    return {
      name: author.name,
      personal_name: author.personal_name || null,
      birth_date: author.birth_date || null,
      death_date: author.death_date || null,
      bio: bio,
      alternate_names: author.alternate_names || [],
      photos: author.photos || [],
      key: author.key,
      remote_ids: author.remote_ids || {},
      links: author.links || []
    };
  } catch (error) {
    if (error.response && error.response.status === 404) {
      return null;
    }
    console.error('Error getting author info:', error.message);
    throw new Error(`Failed to get author details: ${error.message}`);
  }
}

/**
 * Get book cover URL
 * @param {string} key - Type of identifier (ISBN, OCLC, LCCN, OLID, ID)
 * @param {string} value - Value of the identifier
 * @param {string} size - Size of the cover (S, M, L)
 * @returns {string} Cover URL
 */
function getBookCoverUrl(key, value, size = 'L') {
  return `https://covers.openlibrary.org/b/${key.toLowerCase()}/${value}-${size}.jpg`;
}

/**
 * Get author photo URL
 * @param {string} olid - Open Library author ID (e.g., OL23919A)
 * @returns {string} Author photo URL
 */
function getAuthorPhotoUrl(olid) {
  return `https://covers.openlibrary.org/a/olid/${olid}-L.jpg`;
}

/**
 * Enrich book data with Open Library information
 * @param {Object} book - Book object with at least title and author
 * @returns {Promise<Object>} Enriched book object
 */
async function enrichBookData(book) {
  try {
    let enrichedData = { ...book };

    // Try to search by title first
    if (book.title) {
      const searchResults = await searchBooksByTitle(book.title);
      
      if (searchResults.length > 0) {
        const match = searchResults[0]; // Take the first result
        
        enrichedData = {
          ...enrichedData,
          open_library_work_key: match.open_library_work_key,
          cover_url: match.cover_url,
          first_publish_year: match.first_publish_year,
          edition_count: match.edition_count,
          subjects: match.subject
        };

        // Add ISBNs if not already present
        if (!book.isbn && match.isbn.length > 0) {
          enrichedData.isbn = match.isbn[0];
        }

        // Add publisher if not already present
        if (!book.publisher && match.publisher.length > 0) {
          enrichedData.publisher = match.publisher[0];
        }
      }
    }

    // Try to get more details by ISBN if available
    if (book.isbn && book.isbn.length >= 10) {
      const bookDetails = await getBookById('isbn', book.isbn);
      if (bookDetails) {
        enrichedData = {
          ...enrichedData,
          ...bookDetails,
          // Preserve original data
          title: book.title || bookDetails.title,
          authors: book.author ? [book.author] : bookDetails.authors
        };
      }
    }

    return enrichedData;
  } catch (error) {
    console.error('Error enriching book data:', error.message);
    // Return original book data if enrichment fails
    return book;
  }
}

module.exports = {
  searchBooksByTitle,
  searchAuthorsByName,
  getBookById,
  getAuthorInfo,
  getBookCoverUrl,
  getAuthorPhotoUrl,
  enrichBookData
};
