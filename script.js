// Reading Tracker App - Frontend JavaScript
class ReadingTracker {
    constructor() {
        this.books = [];
        this.goal = 24; // Default annual goal
        this.API_BASE = 'http://localhost:7071/api'; // Default to local development
        this.isOnline = navigator.onLine;
        
        this.init();
    }

    async init() {
        await this.loadBooksData();
        this.setupEventListeners();
        this.renderDashboard();
        this.detectAPIEndpoint();
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
            // Try to load from Azure Functions first
            if (this.isOnline) {
                const response = await fetch(`${this.API_BASE}/books`);
                if (response.ok) {
                    this.books = await response.json();
                    return;
                }
            }
            
            // Fallback to local JSON file
            const response = await fetch('./books.json');
            if (response.ok) {
                this.books = await response.json();
            } else {
                this.books = this.getDefaultBooks();
            }
        } catch (error) {
            console.warn('Failed to load books data:', error);
            this.books = this.getDefaultBooks();
        }
    }

    getDefaultBooks() {
        return [
            {
                title: "The DevOps Handbook",
                author: "Gene Kim",
                status: "completed",
                pages: 480,
                completedDate: "2024-01-15",
                rating: 5
            },
            {
                title: "Clean Architecture",
                author: "Robert C. Martin",
                status: "reading",
                pages: 432,
                currentPage: 120,
                startDate: "2024-02-01"
            }
        ];
    }

    setupEventListeners() {
        // Goal setting
        const goalInput = document.getElementById('goal-input');
        if (goalInput) {
            goalInput.addEventListener('change', (e) => {
                this.goal = parseInt(e.target.value) || 24;
                this.renderDashboard();
            });
        }

        // Add book button
        const addBookBtn = document.getElementById('add-book-btn');
        if (addBookBtn) {
            addBookBtn.addEventListener('click', () => this.showAddBookModal());
        }

        // AI Chat button
        const aiChatBtn = document.getElementById('ai-chat-btn');
        if (aiChatBtn) {
            aiChatBtn.addEventListener('click', () => this.showAIChatModal());
        }

        // Sync button
        const syncBtn = document.getElementById('sync-btn');
        if (syncBtn) {
            syncBtn.addEventListener('click', () => this.syncData());
        }
    }

    renderDashboard() {
        this.renderStats();
        this.renderProgress();
        this.renderBooksList();
    }

    renderStats() {
        const completed = this.books.filter(book => book.status === 'completed').length;
        const reading = this.books.filter(book => book.status === 'reading').length;
        const totalPages = this.books
            .filter(book => book.status === 'completed')
            .reduce((sum, book) => sum + (book.pages || 0), 0);

        document.getElementById('books-completed').textContent = completed;
        document.getElementById('books-reading').textContent = reading;
        document.getElementById('total-pages').textContent = totalPages.toLocaleString();
        
        // Calculate pages per day this year
        const startOfYear = new Date(new Date().getFullYear(), 0, 1);
        const daysSinceStart = Math.floor((new Date() - startOfYear) / (1000 * 60 * 60 * 24));
        const pagesPerDay = daysSinceStart > 0 ? Math.round(totalPages / daysSinceStart) : 0;
        document.getElementById('pages-per-day').textContent = pagesPerDay;
    }

    renderProgress() {
        const completed = this.books.filter(book => book.status === 'completed').length;
        const progressPercentage = Math.min((completed / this.goal) * 100, 100);
        
        const progressBar = document.getElementById('progress-bar');
        const progressText = document.getElementById('progress-text');
        
        if (progressBar) {
            progressBar.style.width = `${progressPercentage}%`;
        }
        
        if (progressText) {
            progressText.textContent = `${completed} of ${this.goal} books (${Math.round(progressPercentage)}%)`;
        }
    }

    renderBooksList() {
        const booksList = document.getElementById('books-list');
        if (!booksList) return;

        booksList.innerHTML = this.books.map(book => `
            <div class="book-card ${book.status}">
                <div class="book-info">
                    <h3>${book.title}</h3>
                    <p class="author">by ${book.author}</p>
                    <div class="book-meta">
                        <span class="status ${book.status}">${book.status}</span>
                        <span class="pages">${book.pages} pages</span>
                        ${book.rating ? `<span class="rating">${'★'.repeat(book.rating)}</span>` : ''}
                    </div>
                </div>
                <div class="book-actions">
                    <button onclick="readingTracker.editBook('${book.title}')" class="btn-secondary">Edit</button>
                    <button onclick="readingTracker.deleteBook('${book.title}')" class="btn-danger">Delete</button>
                </div>
            </div>
        `).join('');
    }

    showAddBookModal() {
        // Simple prompt for now - could be enhanced with a proper modal
        const title = prompt('Book title:');
        if (!title) return;
        
        const author = prompt('Author:');
        if (!author) return;
        
        const pages = parseInt(prompt('Number of pages:')) || 0;
        
        this.addBook({
            title,
            author,
            pages,
            status: 'reading',
            startDate: new Date().toISOString().split('T')[0]
        });
    }

    showAIChatModal() {
        const message = prompt('Ask your reading assistant anything:');
        if (!message) return;
        
        this.chatWithAI(message);
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

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.readingTracker = new ReadingTracker();
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
