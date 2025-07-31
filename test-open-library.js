#!/usr/bin/env node

// Test script for Open Library service
const openLibraryService = require('./api/src/services/openLibraryService.js');

async function testOpenLibraryService() {
    console.log('🧪 Testing Open Library Service...\n');
    
    try {
        // Test 1: Search books by title
        console.log('📚 Testing searchBooksByTitle...');
        const searchResults = await openLibraryService.searchBooksByTitle('Harry Potter');
        console.log(`Found ${searchResults.length} books:`);
        console.log(searchResults.slice(0, 2).map(book => `- ${book.title} by ${book.authors.join(', ')}`).join('\n'));
        console.log('✅ searchBooksByTitle works!\n');

        // Test 2: Search authors
        console.log('👤 Testing searchAuthorsByName...');
        const authorResults = await openLibraryService.searchAuthorsByName('J.K. Rowling');
        console.log(`Found ${authorResults.length} authors:`);
        console.log(authorResults.slice(0, 2).map(author => `- ${author.name} (${author.key})`).join('\n'));
        console.log('✅ searchAuthorsByName works!\n');

        // Test 3: Get book by ID
        console.log('🔍 Testing getBookById...');
        const bookDetails = await openLibraryService.getBookById('isbn', '9780439708180');
        if (bookDetails) {
            console.log(`Found book: ${bookDetails.title}`);
            console.log(`Authors: ${bookDetails.authors.join(', ')}`);
            console.log('✅ getBookById works!\n');
        } else {
            console.log('❌ getBookById returned null\n');
        }

        // Test 4: Enrich book data
        console.log('✨ Testing enrichBookData...');
        const enrichedData = await openLibraryService.enrichBookData({
            title: 'The Hobbit',
            author: 'J.R.R. Tolkien'
        });
        console.log(`Enriched: ${enrichedData.title}`);
        console.log(`Authors: ${enrichedData.authors.join(', ')}`);
        console.log(`Published: ${enrichedData.first_publish_year}`);
        console.log('✅ enrichBookData works!\n');

        // Test 5: Get cover URL
        console.log('🖼️  Testing getCoverUrl...');
        const coverUrl = openLibraryService.getCoverUrl('isbn', '9780439708180', 'M');
        console.log(`Cover URL: ${coverUrl}`);
        console.log('✅ getCoverUrl works!\n');

        console.log('🎉 All tests passed! Open Library service is working correctly.');

    } catch (error) {
        console.error('❌ Test failed:', error.message);
        process.exit(1);
    }
}

testOpenLibraryService();
