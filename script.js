// Reading Dashboard JavaScript - Simple & Functional

// Configuration
const CONFIG = {
    AZURE_FUNCTION_URL: 'https://your-function-app.azurewebsites.net/api',
    OPENAI_ENDPOINT: '/chat', // Azure Function endpoint
    BOOKS_ENDPOINT: '/books',  // Azure Function endpoint
    AUTO_SYNC_INTERVAL: 30000  // 30 seconds
};

// State management
let readingData = {
    books: [],
    goal: 24,
    currentYear: 2025,
    lastSync: null
};

// Initialize the dashboard
document.addEventListener('DOMContentLoaded', async function() {
    await loadReadingData();
    renderDashboard();
    setupEventListeners();
    startAutoSync();
});

// Load reading data from Azure/local storage
async function loadReadingData() {
    try {
        // Try to load from Azure first
        const response = await fetch(`${CONFIG.AZURE_FUNCTION_URL}${CONFIG.BOOKS_ENDPOINT}`);
        if (response.ok) {
            const data = await response.json();
            readingData.books = data.books || [];
            readingData.lastSync = new Date();
            updateSyncStatus('synced');
        } else {
            throw new Error('Azure fetch failed');
        }
    } catch (error) {
        console.log('Loading from local storage fallback');
        // Fallback to local storage or embedded data
        const localData = localStorage.getItem('readingData');
        if (localData) {
            readingData = { ...readingData, ...JSON.parse(localData) };
        } else {
            // Load from your existing books.json as fallback
            await loadFromBooksJson();
        }
        updateSyncStatus('offline');
    }
}

// Load from your existing books.json file
async function loadFromBooksJson() {
    try {
        const response = await fetch('./books.json');
        if (response.ok) {
            const books = await response.json();
            readingData.books = books;
        }
    } catch (error) {
        console.log('Could not load books.json');
        // Initialize with empty data
        readingData.books = [];
    }
}

// Save reading data to Azure
async function saveReadingData() {
    try {
        const response = await fetch(`${CONFIG.AZURE_FUNCTION_URL}${CONFIG.BOOKS_ENDPOINT}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(readingData)
        });
        
        if (response.ok) {
            readingData.lastSync = new Date();
            updateSyncStatus('synced');
        } else {
            throw new Error('Save failed');
        }
    } catch (error) {
        console.error('Save to Azure failed:', error);
        // Save to local storage as fallback
        localStorage.setItem('readingData', JSON.stringify(readingData));
        updateSyncStatus('local');
    }
}

// Render the entire dashboard
function renderDashboard() {
    renderGoalProgress();
    renderCurrentlyReading();
    renderRecentActivity();
}

// Render goal progress section
function renderGoalProgress() {
    const currentYear = new Date().getFullYear();
    const booksRead = readingData.books.filter(book => 
        book['Date Read'] && book['Date Read'].startsWith(currentYear.toString())
    ).length;
    
    const progressPercent = (booksRead / readingData.goal) * 100;
    
    document.getElementById('booksRead').textContent = booksRead;
    document.getElementById('booksGoal').textContent = readingData.goal;
    document.getElementById('goalProgress').style.width = `${Math.min(progressPercent, 100)}%`;
    
    // Calculate stats
    const thisMonth = new Date().getMonth() + 1;
    const booksThisMonth = readingData.books.filter(book => {
        if (!book['Date Read']) return false;
        const dateRead = new Date(book['Date Read'].replace(/\//g, '-'));
        return dateRead.getFullYear() === currentYear && dateRead.getMonth() + 1 === thisMonth;
    }).length;
    
    document.getElementById('booksThisMonth').textContent = booksThisMonth;
    document.getElementById('currentStreak').textContent = calculateReadingStreak();
    document.getElementById('avgPagesPerDay').textContent = calculateAvgPagesPerDay();
}

// Render currently reading books
function renderCurrentlyReading() {
    const currentlyReading = readingData.books.filter(book => 
        book.Bookshelves === 'currently-reading' || book.Shelf === 'currently-reading'
    );
    
    document.getElementById('currentlyReadingCount').textContent = `(${currentlyReading.length})`;
    
    const grid = document.getElementById('currentlyReadingGrid');
    grid.innerHTML = '';
    
    if (currentlyReading.length === 0) {
        grid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-secondary);">
                📚 No books currently being read
                <br><br>
                <button class="action-btn primary" onclick="showAddBookModal()" style="margin-top: 16px;">
                    <span class="action-icon">➕</span>
                    <span class="action-text">Add Your First Book</span>
                </button>
            </div>
        `;
        return;
    }
    
    currentlyReading.forEach(book => {
        const bookElement = createBookCard(book);
        grid.appendChild(bookElement);
    });
}

// Create a book card element
function createBookCard(book) {
    const div = document.createElement('div');
    div.className = 'book-card';
    div.onclick = () => showBookDetails(book);
    
    const progress = book.currentPage ? (book.currentPage / book['Number of Pages']) * 100 : 0;
    const progressText = book.currentPage ? 
        `${book.currentPage} / ${book['Number of Pages']} pages` : 
        `${book['Number of Pages']} pages`;
    
    div.innerHTML = `
        <div class="book-cover">📖</div>
        <div class="book-title">${book.Title}</div>
        <div class="book-author">by ${book.Author}</div>
        <div class="book-progress">
            <div class="progress-bar-small">
                <div class="progress-fill-small" style="width: ${progress}%"></div>
            </div>
            <div class="progress-label">${Math.round(progress)}% • ${progressText}</div>
        </div>
    `;
    
    return div;
}

// Render recent activity
function renderRecentActivity() {
    const recentBooks = readingData.books
        .filter(book => book['Date Read'])
        .sort((a, b) => new Date(b['Date Read']) - new Date(a['Date Read']))
        .slice(0, 5);
    
    const activityList = document.getElementById('recentActivity');
    activityList.innerHTML = '';
    
    if (recentBooks.length === 0) {
        activityList.innerHTML = `
            <div class="activity-item">
                <span class="activity-icon">📚</span>
                <span class="activity-text">No recent reading activity</span>
                <span class="activity-time">Start your first book!</span>
            </div>
        `;
        return;
    }
    
    recentBooks.forEach(book => {
        const activityItem = document.createElement('div');
        activityItem.className = 'activity-item';
        
        const dateRead = new Date(book['Date Read'].replace(/\//g, '-'));
        const timeAgo = getTimeAgo(dateRead);
        
        activityItem.innerHTML = `
            <span class="activity-icon">✅</span>
            <span class="activity-text">Finished reading "${book.Title}" by ${book.Author}</span>
            <span class="activity-time">${timeAgo}</span>
        `;
        
        activityList.appendChild(activityItem);
    });
}

// Calculate reading streak
function calculateReadingStreak() {
    const today = new Date();
    let streak = 0;
    let currentDate = new Date(today);
    
    while (true) {
        const dateStr = currentDate.toISOString().slice(0, 10).replace(/-/g, '/');
        const hasBookOnDate = readingData.books.some(book => 
            book['Date Read'] === dateStr
        );
        
        if (hasBookOnDate) {
            streak++;
            currentDate.setDate(currentDate.getDate() - 1);
        } else {
            break;
        }
    }
    
    return streak;
}

// Calculate average pages per day
function calculateAvgPagesPerDay() {
    const thisMonth = new Date().getMonth() + 1;
    const currentYear = new Date().getFullYear();
    
    const booksThisMonth = readingData.books.filter(book => {
        if (!book['Date Read']) return false;
        const dateRead = new Date(book['Date Read'].replace(/\//g, '-'));
        return dateRead.getFullYear() === currentYear && dateRead.getMonth() + 1 === thisMonth;
    });
    
    const totalPages = booksThisMonth.reduce((sum, book) => sum + (book['Number of Pages'] || 0), 0);
    const daysInMonth = new Date().getDate();
    
    return Math.round(totalPages / daysInMonth);
}

// Get time ago string
function getTimeAgo(date) {
    const now = new Date();
    const diffMs = now - date;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
    return `${Math.floor(diffDays / 30)} months ago`;
}

// Modal functions
function showAddBookModal() {
    document.getElementById('addBookModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

function showBookDetails(book) {
    alert(`Book: ${book.Title}\nAuthor: ${book.Author}\nPages: ${book['Number of Pages']}\nStatus: ${book.Shelf}`);
}

function showStatsView() {
    alert('Stats view coming soon!');
}

function showSearchModal() {
    alert('Search functionality coming soon!');
}

// Chat assistant functions
function toggleChatAssistant() {
    const chatPanel = document.getElementById('chatPanel');
    chatPanel.classList.toggle('active');
}

async function sendChatMessage() {
    const input = document.getElementById('chatInput');
    const message = input.value.trim();
    
    if (!message) return;
    
    // Add user message to chat
    addChatMessage(message, 'user');
    input.value = '';
    
    try {
        // Send to Azure Function with OpenAI
        const response = await fetch(`${CONFIG.AZURE_FUNCTION_URL}${CONFIG.OPENAI_ENDPOINT}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: message,
                books: readingData.books,
                context: 'reading_assistant'
            })
        });
        
        if (response.ok) {
            const data = await response.json();
            addChatMessage(data.response, 'assistant');
            
            // Handle any book updates
            if (data.bookUpdates) {
                await processBookUpdates(data.bookUpdates);
            }
        } else {
            throw new Error('AI response failed');
        }
    } catch (error) {
        console.error('Chat error:', error);
        addChatMessage('Sorry, I'm having trouble connecting to the AI service. Please try again later.', 'assistant');
    }
}

function addChatMessage(message, sender) {
    const messagesContainer = document.getElementById('chatMessages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `chat-message ${sender}`;
    messageDiv.innerHTML = `<div class="message-content">${message}</div>`;
    messagesContainer.appendChild(messageDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

function handleChatKeyPress(event) {
    if (event.key === 'Enter') {
        sendChatMessage();
    }
}

// Process book updates from AI
async function processBookUpdates(updates) {
    for (const update of updates) {
        if (update.action === 'add') {
            readingData.books.push(update.book);
        } else if (update.action === 'update') {
            const bookIndex = readingData.books.findIndex(book => book.Title === update.title);
            if (bookIndex !== -1) {
                readingData.books[bookIndex] = { ...readingData.books[bookIndex], ...update.changes };
            }
        }
    }
    
    await saveReadingData();
    renderDashboard();
}

// Setup event listeners
function setupEventListeners() {
    // Add book form submission
    document.getElementById('addBookForm').addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData(e.target);
        const book = {
            Title: formData.get('title'),
            Author: formData.get('author'),
            'Additional Authors': '',
            ISBN: '',
            ISBN13: '',
            Publisher: '',
            'Number of Pages': 0,
            'Year Published': new Date().getFullYear(),
            'Original Publication Year': new Date().getFullYear(),
            'Date Read': '',
            'Date Added': new Date().toISOString().slice(0, 10).replace(/-/g, '/'),
            Bookshelves: formData.get('shelf'),
            'Bookshelves with positions': '',
            Shelf: formData.get('shelf'),
            'My Review': ''
        };
        
        readingData.books.push(book);
        await saveReadingData();
        renderDashboard();
        closeModal('addBookModal');
        e.target.reset();
    });
    
    // Close modals when clicking outside
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('modal')) {
            e.target.classList.remove('active');
        }
    });
}

// Update sync status indicator
function updateSyncStatus(status) {
    const statusElement = document.getElementById('syncStatus');
    const dot = statusElement.querySelector('.status-dot');
    const text = statusElement.querySelector('.status-text');
    
    switch (status) {
        case 'synced':
            dot.style.background = 'var(--accent-green)';
            text.textContent = 'Synced';
            break;
        case 'syncing':
            dot.style.background = 'var(--accent-blue)';
            text.textContent = 'Syncing...';
            break;
        case 'offline':
            dot.style.background = 'var(--accent-orange)';
            text.textContent = 'Offline';
            break;
        case 'local':
            dot.style.background = 'var(--text-secondary)';
            text.textContent = 'Local only';
            break;
    }
}

// Auto-sync functionality
function startAutoSync() {
    setInterval(async () => {
        if (readingData.lastSync && Date.now() - readingData.lastSync > CONFIG.AUTO_SYNC_INTERVAL) {
            updateSyncStatus('syncing');
            await saveReadingData();
        }
    }, CONFIG.AUTO_SYNC_INTERVAL);
}

// Export functions for global access
window.showAddBookModal = showAddBookModal;
window.closeModal = closeModal;
window.showStatsView = showStatsView;
window.showSearchModal = showSearchModal;
window.toggleChatAssistant = toggleChatAssistant;
window.sendChatMessage = sendChatMessage;
window.handleChatKeyPress = handleChatKeyPress;
