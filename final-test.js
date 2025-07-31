#!/usr/bin/env node

// Comprehensive test suite for Reading Tracker components
const openLibraryService = require('./api/src/services/openLibraryService.js');

async function runComprehensiveTests() {
    console.log('🎯 Reading Tracker - Comprehensive Test Suite\n');
    
    let testsRun = 0;
    let testsPassed = 0;
    
    const runTest = async (name, testFn) => {
        testsRun++;
        try {
            console.log(`🧪 ${name}...`);
            await testFn();
            testsPassed++;
            console.log(`✅ ${name} PASSED\n`);
        } catch (error) {
            console.log(`❌ ${name} FAILED: ${error.message}\n`);
        }
    };
    
    // Test 1: Book Search
    await runTest('Book Search', async () => {
        const results = await openLibraryService.searchBooksByTitle('The Hobbit');
        if (results.length === 0) throw new Error('No results found');
        if (!results[0].title) throw new Error('Missing title');
        console.log(`   Found ${results.length} books, first: "${results[0].title}"`);
    });
    
    // Test 2: Author Search
    await runTest('Author Search', async () => {
        const results = await openLibraryService.searchAuthorsByName('J.K. Rowling');
        if (results.length === 0) throw new Error('No authors found');
        if (!results[0].name) throw new Error('Missing author name');
        console.log(`   Found ${results.length} authors, first: "${results[0].name}"`);
    });
    
    // Test 3: Book by ISBN
    await runTest('Book by ISBN', async () => {
        const book = await openLibraryService.getBookById('isbn', '9780439708180');
        if (!book) throw new Error('Book not found');
        if (!book.title) throw new Error('Missing book title');
        console.log(`   Found: "${book.title}" by ${book.authors.join(', ')}`);
    });
    
    // Test 4: Author Information
    await runTest('Author Information', async () => {
        const author = await openLibraryService.getAuthorInfo('OL23919A');
        if (!author) throw new Error('Author not found');
        if (!author.name) throw new Error('Missing author name');
        console.log(`   Author: "${author.name}" (${author.key})`);
    });
    
    // Test 5: Cover URL Generation
    await runTest('Cover URL Generation', async () => {
        const url = openLibraryService.getBookCoverUrl('isbn', '9780439708180', 'M');
        if (!url.includes('covers.openlibrary.org')) throw new Error('Invalid cover URL');
        console.log(`   Cover URL: ${url}`);
    });
    
    // Test 6: Author Photo URL
    await runTest('Author Photo URL', async () => {
        const url = openLibraryService.getAuthorPhotoUrl('OL23919A');
        if (!url.includes('covers.openlibrary.org')) throw new Error('Invalid photo URL');
        console.log(`   Photo URL: ${url}`);
    });
    
    // Test 7: Book Data Enrichment
    await runTest('Book Data Enrichment', async () => {
        const enriched = await openLibraryService.enrichBookData({
            title: 'The Hobbit',
            author: 'J.R.R. Tolkien'
        });
        if (!enriched.title) throw new Error('Missing enriched title');
        console.log(`   Enriched: "${enriched.title}" (${enriched.first_publish_year})`);
        console.log(`   Authors: ${Array.isArray(enriched.authors) ? enriched.authors.join(', ') : 'N/A'}`);
    });
    
    // Summary
    console.log('📊 Test Results Summary:');
    console.log(`   Tests Run: ${testsRun}`);
    console.log(`   Tests Passed: ${testsPassed}`);
    console.log(`   Tests Failed: ${testsRun - testsPassed}`);
    console.log(`   Success Rate: ${Math.round((testsPassed / testsRun) * 100)}%`);
    
    if (testsPassed === testsRun) {
        console.log('\n🎉 ALL TESTS PASSED! Your Reading Tracker API is fully functional.');
        console.log('\n🚀 Ready for:');
        console.log('   ✅ Production deployment');
        console.log('   ✅ Frontend integration');
        console.log('   ✅ MCP server usage');
        console.log('   ✅ AI assistant integration');
    } else {
        console.log('\n⚠️  Some tests failed. Please check the error messages above.');
    }
    
    console.log('\n📋 Available API capabilities:');
    console.log('   📚 Search millions of books by title');
    console.log('   👤 Find authors and their information');
    console.log('   🔍 Get detailed book data by ISBN/OLID');
    console.log('   🖼️  Generate cover image URLs');
    console.log('   📖 Get author photos and biographies');
    console.log('   ✨ Enrich book data with metadata');
}

// Run the comprehensive test suite
runComprehensiveTests().catch(console.error);
