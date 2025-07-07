// Reading Tracker App - Frontend JavaScript
class ReadingTracker {
    constructor() {
        this.books = [];
        this.goal = 90; // 2025 annual reading goal
        this.API_BASE = 'http://localhost:7071/api'; // Default to local development
        this.isOnline = navigator.onLine;
        
        // Configuration: Set to false for read-only public dashboard
        this.isAdminMode = this.checkAdminMode();
        
        this.init();
    }

    checkAdminMode() {
        // Method 1: Check URL parameter (?admin=true)
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.get('admin') === 'true') {
            return true;
        }
        
        // Method 2: Check localStorage for admin flag
        if (localStorage.getItem('reading-tracker-admin') === 'true') {
            return true;
        }
        
        // Method 3: Check if running on localhost (development)
        if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
            return true;
        }
        
        // Default to read-only for public viewers
        return false;
    }

    async init() {
        console.log('ReadingTracker initializing...');
        await this.loadBooksData();
        console.log('Books loaded, setting up event listeners...');
        this.setupEventListeners();
        console.log('Event listeners set up, rendering dashboard...');
        this.renderDashboard();
        console.log('Dashboard rendered, detecting API endpoint...');
        this.detectAPIEndpoint();
        console.log('ReadingTracker initialization complete.');
    }

    async detectAPIEndpoint() {
        // Try to detect if we're running in production
        const hostname = window.location.hostname;
        if (hostname.includes('github.io')) {
            // Try to get API URL from environment or use placeholder
            this.API_BASE = 'https://your-function-app.azurewebsites.net/api';
            console.log('Production mode detected, using Azure Functions API');
        } else {
            console.log('Development mode detected, using local API');
        }
    }

    async loadBooksData() {
        try {
            console.log('loadBooksData: Starting to load books...');
            
            // First check if books data is available from the global script
            if (window.booksData && Array.isArray(window.booksData)) {
                this.books = window.booksData;
                console.log('Loaded books from global data:', this.books.length);
                
                // Test filtering right here
                const completed2025 = this.books.filter(book => 
                    book.Shelf === 'read' && 
                    book['Date Read'] && 
                    book['Date Read'].startsWith('2025')
                );
                console.log('Books completed in 2025 (from global data):', completed2025.length);
                return;
            }
            
            // Try to load from Azure Functions first
            if (this.isOnline) {
                console.log('loadBooksData: Trying API first...');
                const response = await fetch(`${this.API_BASE}/books`);
                if (response.ok) {
                    this.books = await response.json();
                    console.log('Loaded books from API:', this.books.length);
                    return;
                }
            }
            
            // Fallback to local JSON file
            console.log('loadBooksData: Loading from local JSON file...');
            const response = await fetch('./books.json');
            if (response.ok) {
                this.books = await response.json();
                console.log('Loaded books from local JSON:', this.books.length);
                console.log('Sample book:', this.books[0]);
                
                // Test filtering right here
                const completed2025 = this.books.filter(book => 
                    book.Shelf === 'read' && 
                    book['Date Read'] && 
                    book['Date Read'].startsWith('2025')
                );
                console.log('Books completed in 2025 (in loadBooksData):', completed2025.length);
                return; // Successfully loaded
            } else {
                console.error('Failed to load books.json');
                throw new Error('Failed to load books.json');
            }
        } catch (error) {
            console.warn('Failed to load books data:', error);
            console.log('Using default books for testing...');
            this.books = this.getDefaultBooks();
        }
    }

    getDefaultBooks() {
        // Return some test data to verify the rendering works
        return [
            {
                "Title": "Test Book 1",
                "Author": "Test Author",
                "Shelf": "read",
                "Date Read": "2025/01/15",
                "Number of Pages": 300
            },
            {
                "Title": "Test Book 2",
                "Author": "Test Author",
                "Shelf": "read",
                "Date Read": "2025/02/20",
                "Number of Pages": 250
            },
            {
                "Title": "Test Book 3",
                "Author": "Test Author",
                "Shelf": "currently-reading",
                "Date Read": "",
                "Number of Pages": 400
            }
        ];
    }

    setupEventListeners() {
        // Hide admin-only features if in read-only mode
        if (!this.isAdminMode) {
            this.hideAdminFeatures();
        }
        
        // Goal setting (only in admin mode)
        const goalInput = document.getElementById('goal-input');
        if (goalInput && this.isAdminMode) {
            goalInput.addEventListener('change', (e) => {
                this.goal = parseInt(e.target.value) || 90;
                this.renderDashboard();
            });
        }

        // Add book button (admin only)
        const addBookBtn = document.getElementById('add-book-btn');
        if (addBookBtn && this.isAdminMode) {
            addBookBtn.addEventListener('click', () => this.showAddBookModal());
        }

        // AI Chat button (admin only)
        const aiChatBtn = document.getElementById('ai-chat-btn');
        if (aiChatBtn && this.isAdminMode) {
            aiChatBtn.addEventListener('click', () => this.showAIChatModal());
        }

        // Sync button (admin only)
        const syncBtn = document.getElementById('sync-btn');
        if (syncBtn && this.isAdminMode) {
            syncBtn.addEventListener('click', () => this.syncData());
        }
    }

    hideAdminFeatures() {
        // Hide admin action buttons
        const adminButtons = document.querySelectorAll('.action-btn');
        adminButtons.forEach(button => {
            const text = button.textContent.toLowerCase();
            if (text.includes('add book') || text.includes('ai assistant')) {
                button.style.display = 'none';
            }
        });
        
        // Hide modals
        const modals = document.querySelectorAll('.modal, .chat-panel');
        modals.forEach(modal => modal.style.display = 'none');
        
        // Add read-only indicator
        this.addReadOnlyIndicator();
    }

    addReadOnlyIndicator() {
        const header = document.querySelector('.header');
        if (header) {
            const indicator = document.createElement('div');
            indicator.className = 'read-only-indicator';
            indicator.innerHTML = `
                <span class="read-only-badge">👁️ Read-Only View</span>
                <small>This is a public dashboard. <a href="?admin=true">Admin access</a></small>
            `;
            header.appendChild(indicator);
        }
    }

    renderDashboard() {
        console.log('renderDashboard: Starting render...');
        this.renderStats();
        this.renderProgress();
        this.renderBooksList();
        console.log('renderDashboard: Render complete.');
    }

    renderStats() {
        console.log('Total books loaded:', this.books.length);
        
        // Filter books completed in 2025 (using your data format)
        const completed2025 = this.books.filter(book => 
            book.Shelf === 'read' && 
            book['Date Read'] && 
            book['Date Read'].startsWith('2025')
        );
        
        console.log('Books completed in 2025:', completed2025.length);
        console.log('Sample completed book:', completed2025[0]);
        
        const currentlyReading = this.books.filter(book => 
            book.Shelf === 'currently-reading' || 
            (book.Shelf === 'to-read' && book['Date Read'] === '')
        ).length;
        
        const totalPages2025 = completed2025.reduce((sum, book) => sum + (book['Number of Pages'] || 0), 0);

        // Update elements that exist
        const booksCompletedEl = document.getElementById('books-completed');
        const booksReadingEl = document.getElementById('books-reading'); 
        const totalPagesEl = document.getElementById('total-pages');
        const pagesPerDayEl = document.getElementById('pages-per-day');
        const avgPagesPerDayEl = document.getElementById('avg-pages-per-day');
        
        if (booksCompletedEl) booksCompletedEl.textContent = completed2025.length;
        if (booksReadingEl) booksReadingEl.textContent = currentlyReading;
        if (totalPagesEl) totalPagesEl.textContent = totalPages2025.toLocaleString();
        
        // Calculate pages per day this year
        const startOfYear = new Date(2025, 0, 1);
        const daysSinceStart = Math.floor((new Date() - startOfYear) / (1000 * 60 * 60 * 24));
        const pagesPerDay = daysSinceStart > 0 ? Math.round(totalPages2025 / daysSinceStart) : 0;
        if (pagesPerDayEl) pagesPerDayEl.textContent = pagesPerDay;
        if (avgPagesPerDayEl) avgPagesPerDayEl.textContent = pagesPerDay;
        
        // Calculate books this month (July 2025)
        const booksThisMonth = this.books.filter(book => 
            book.Shelf === 'read' && 
            book['Date Read'] && 
            book['Date Read'].startsWith('2025/07')
        ).length;
        
        const booksThisMonthEl = document.getElementById('books-this-month');
        if (booksThisMonthEl) booksThisMonthEl.textContent = booksThisMonth;
        
        // Update currently reading count
        const currentlyReadingCountEl = document.getElementById('currentlyReadingCount');
        if (currentlyReadingCountEl) currentlyReadingCountEl.textContent = `(${currentlyReading})`;
    }

    renderProgress() {
        // Count books completed in 2025
        const completed = this.books.filter(book => 
            book.Shelf === 'read' && 
            book['Date Read'] && 
            book['Date Read'].startsWith('2025')
        ).length;
        const progressPercentage = Math.min((completed / this.goal) * 100, 100);
        
        console.log('renderProgress: completed =', completed, 'goal =', this.goal, 'percentage =', progressPercentage);
        
        const progressBar = document.getElementById('progress-bar');
        const progressText = document.getElementById('progress-text');
        
        console.log('renderProgress: progressBar =', progressBar, 'progressText =', progressText);
        
        if (progressBar) {
            progressBar.style.width = `${progressPercentage}%`;
            console.log('renderProgress: set progressBar width to', `${progressPercentage}%`);
        }
        
        if (progressText) {
            progressText.textContent = `${completed} of ${this.goal} books (${Math.round(progressPercentage)}%)`;
            console.log('renderProgress: set progressText to', `${completed} of ${this.goal} books (${Math.round(progressPercentage)}%)`);
        }
    }

    renderBooksList() {
        const booksList = document.getElementById('books-list');
        if (!booksList) return;

        // Show only recent books (2025 reads, currently reading, or recent additions)
        const relevantBooks = this.books.filter(book => 
            (book.Shelf === 'read' && book['Date Read'] && book['Date Read'].startsWith('2025')) ||
            book.Shelf === 'currently-reading' ||
            (book['Date Added'] && book['Date Added'].startsWith('2025'))
        ).slice(0, 20); // Limit to 20 most relevant books

        booksList.innerHTML = relevantBooks.map(book => {
            const status = book.Shelf === 'read' ? 'completed' : 
                          book.Shelf === 'currently-reading' ? 'reading' : 'to-read';
            const rating = book['My Rating'] ? '★'.repeat(parseInt(book['My Rating'])) : '';
            
            // Only show action buttons in admin mode
            const actionButtons = this.isAdminMode ? `
                <div class="book-actions">
                    <button onclick="readingTracker.editBook('${book.Title}')" class="btn-secondary">Edit</button>
                    <button onclick="readingTracker.deleteBook('${book.Title}')" class="btn-danger">Delete</button>
                </div>
            ` : '';
            
            return `
                <div class="book-card ${status}">
                    <div class="book-info">
                        <h3>${book.Title || 'Unknown Title'}</h3>
                        <p class="author">by ${book.Author || 'Unknown Author'}</p>
                        <div class="book-meta">
                            <span class="status ${status}">${status}</span>
                            <span class="pages">${book['Number of Pages'] || 0} pages</span>
                            ${rating ? `<span class="rating">${rating}</span>` : ''}
                            ${book['Date Read'] ? `<span class="date">Finished: ${book['Date Read']}</span>` : ''}
                        </div>
                    </div>
                    ${actionButtons}
                </div>
            `;
        }).join('');
    }

    showAddBookModal() {
        const modal = document.getElementById('addBookModal');
        if (modal) {
            modal.style.display = 'block';
        }
    }

    showAIChatModal() {
        const modal = document.getElementById('aiChatModal');
        if (modal) {
            modal.style.display = 'block';
        }
    }

    showStatsView() {
        // For now, just scroll to stats section or could open a detailed stats modal
        const statsSection = document.querySelector('.stats-section');
        if (statsSection) {
            statsSection.scrollIntoView({ behavior: 'smooth' });
        }
    }

    showSearchModal() {
        // Placeholder for search functionality
        alert('Search functionality coming soon!');
    }

    closeModal(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'none';
        }
    }

    sendChatMessage() {
        const input = document.getElementById('chatInput');
        if (input && input.value.trim()) {
            this.chatWithAI(input.value.trim());
            input.value = '';
        }
    }

    async addBook(book) {
        this.books.push(book);
        await this.saveBooks();
        this.renderDashboard();
    }

    async editBook(title) {
        const book = this.books.find(b => b.title === title);
        if (!book) return;
        
        const newStatus = prompt(`Change status for "${title}":`, book.status);
        if (newStatus) {
            book.status = newStatus;
            if (newStatus === 'completed' && !book.completedDate) {
                book.completedDate = new Date().toISOString().split('T')[0];
            }
            await this.saveBooks();
            this.renderDashboard();
        }
    }

    async deleteBook(title) {
        if (confirm(`Delete "${title}"?`)) {
            this.books = this.books.filter(b => b.title !== title);
            await this.saveBooks();
            this.renderDashboard();
        }
    }

    async saveBooks() {
        try {
            if (this.isOnline) {
                await fetch(`${this.API_BASE}/books`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(this.books)
                });
            }
            // Also save to localStorage as backup
            localStorage.setItem('reading-tracker-books', JSON.stringify(this.books));
        } catch (error) {
            console.warn('Failed to save books:', error);
            // Save to localStorage as fallback
            localStorage.setItem('reading-tracker-books', JSON.stringify(this.books));
        }
    }

    async syncData() {
        const syncBtn = document.getElementById('sync-btn');
        if (syncBtn) {
            syncBtn.textContent = 'Syncing...';
            syncBtn.disabled = true;
        }
        
        try {
            await this.loadBooksData();
            this.renderDashboard();
            console.log('Data synced successfully');
        } catch (error) {
            console.error('Sync failed:', error);
        } finally {
            if (syncBtn) {
                syncBtn.textContent = 'Sync';
                syncBtn.disabled = false;
            }
        }
    }

    async chatWithAI(message) {
        try {
            const response = await fetch(`${this.API_BASE}/books/chat`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ 
                    message,
                    books: this.books 
                })
            });
            
            if (response.ok) {
                const data = await response.json();
                alert(`AI Assistant: ${data.reply}`);
            } else {
                alert('AI Assistant is currently unavailable');
            }
        } catch (error) {
            console.error('AI chat failed:', error);
            alert('AI Assistant is currently unavailable');
        }
    }
}

// Global function wrappers for HTML onclick handlers
function showAddBookModal() {
    if (window.readingTracker) {
        window.readingTracker.showAddBookModal();
    }
}

function showStatsView() {
    if (window.readingTracker) {
        window.readingTracker.showStatsView();
    }
}

function showSearchModal() {
    if (window.readingTracker) {
        window.readingTracker.showSearchModal();
    }
}

function toggleChatAssistant() {
    if (window.readingTracker) {
        window.readingTracker.showAIChatModal();
    }
}

function closeModal(modalId) {
    if (window.readingTracker) {
        window.readingTracker.closeModal(modalId);
    }
}

function sendChatMessage() {
    if (window.readingTracker) {
        window.readingTracker.sendChatMessage();
    }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, initializing ReadingTracker...');
    window.readingTracker = new ReadingTracker();
    
    // Also call renderProgress again after a short delay to ensure DOM is ready
    setTimeout(() => {
        if (window.readingTracker) {
            console.log('Calling renderProgress again after delay...');
            window.readingTracker.renderProgress();
        }
    }, 100);
});

// Handle online/offline status
window.addEventListener('online', () => {
    console.log('App is now online');
    window.readingTracker.isOnline = true;
});

window.addEventListener('offline', () => {
    console.log('App is now offline');
    window.readingTracker.isOnline = false;
});
