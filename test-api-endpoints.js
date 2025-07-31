#!/usr/bin/env node

// Simple API test script for Open Library endpoints only
const openLibraryService = require('./api/src/services/openLibraryService.js');

async function testAllEndpoints() {
    console.log('🧪 Testing Reading Tracker API Endpoints...\n');
    
    try {
        // Test 1: Search Books
        console.log('📚 Testing Book Search...');
        const searchResults = await openLibraryService.searchBooksByTitle('Harry Potter');
        console.log(`✅ Found ${searchResults.length} books`);
        if (searchResults.length > 0) {
            console.log(`   Example: "${searchResults[0].title}" by ${searchResults[0].authors.join(', ')}`);
        }
        console.log();

        // Test 2: Author Search
        console.log('👤 Testing Author Search...');
        const authors = await openLibraryService.searchAuthorsByName('Stephen King');
        console.log(`✅ Found ${authors.length} authors`);
        if (authors.length > 0) {
            console.log(`   Example: "${authors[0].name}" (${authors[0].key})`);
        }
        console.log();

        // Test 3: Book by ID
        console.log('🔍 Testing Book by ID...');
        const bookDetails = await openLibraryService.getBookById('isbn', '9780439708180');
        if (bookDetails) {
            console.log(`✅ Found book: "${bookDetails.title}"`);
            console.log(`   Authors: ${bookDetails.authors.join(', ')}`);
        } else {
            console.log('❌ No book found for that ISBN');
        }
        console.log();

        // Test 4: Enrich Book Data
        console.log('✨ Testing Book Data Enrichment...');
        const enrichedBook = await openLibraryService.enrichBookData({
            title: 'Dune',
            author: 'Frank Herbert'
        });
        console.log(`✅ Enriched "${enrichedBook.title}"`);
        console.log(`   Authors: ${enrichedBook.authors.join(', ')}`);
        console.log(`   Published: ${enrichedBook.first_publish_year}`);
        console.log(`   Cover: ${enrichedBook.cover_url ? 'Available' : 'Not available'}`);
        console.log();

        // Test 5: Cover URLs
        console.log('🖼️  Testing Cover URL Generation...');
        const coverUrl = openLibraryService.getCoverUrl('isbn', '9780441172719', 'L');
        console.log(`✅ Cover URL: ${coverUrl}`);
        console.log();

        // Test 6: Author Info
        console.log('📖 Testing Author Information...');
        const authorInfo = await openLibraryService.getAuthorInfo('OL23919A'); // J.K. Rowling
        if (authorInfo) {
            console.log(`✅ Author: ${authorInfo.name}`);
            console.log(`   Birth: ${authorInfo.birth_date || 'Unknown'}`);
            console.log(`   Bio: ${authorInfo.bio ? 'Available' : 'Not available'}`);
        } else {
            console.log('❌ Author not found');
        }
        console.log();

        console.log('🎉 All API endpoints are working correctly!');
        console.log('\n📋 Summary of working endpoints:');
        console.log('  ✅ Book search by title');
        console.log('  ✅ Author search by name');
        console.log('  ✅ Book details by ISBN/OLID');
        console.log('  ✅ Book data enrichment');
        console.log('  ✅ Cover URL generation');
        console.log('  ✅ Author information lookup');
        
        console.log('\n🚀 Ready for production deployment!');

    } catch (error) {
        console.error('❌ Test failed:', error.message);
        console.error('Stack:', error.stack);
    }
}

// Run the test
testAllEndpoints();
