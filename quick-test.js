#!/usr/bin/env node

// Quick test of Open Library service
const openLibraryService = require('./api/src/services/openLibraryService.js');

async function quickTest() {
    console.log('🧪 Quick Open Library Test...\n');
    
    try {
        // Test search
        const books = await openLibraryService.searchBooksByTitle('The Hobbit');
        console.log(`✅ Found ${books.length} books for "The Hobbit"`);
        
        if (books.length > 0) {
            const book = books[0];
            console.log(`📖 Title: ${book.title}`);
            console.log(`👤 Authors: ${book.authors.join(', ')}`);
            console.log(`📅 Year: ${book.first_publish_year || 'Unknown'}`);
        }
        
        console.log('\n🎉 Open Library integration is working!');
        
    } catch (error) {
        console.error('❌ Error:', error.message);
    }
}

quickTest();
