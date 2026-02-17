// Admin page functionality

const API_BASE = '/api';

// DOM Elements
const loginScreen = document.getElementById('login-screen');
const adminScreen = document.getElementById('admin-screen');
const loginForm = document.getElementById('login-form');
const loginError = document.getElementById('login-error');
const logoutBtn = document.getElementById('logout-btn');
const exportBtn = document.getElementById('export-btn');
const quickAddForm = document.getElementById('quick-add-form');
const formMessage = document.getElementById('form-message');
const scanBtn = document.getElementById('scan-btn');
const lookupBtn = document.getElementById('lookup-btn');
const scannerModal = document.getElementById('scanner-modal');
const cancelScanBtn = document.getElementById('cancel-scan');
const coverPreview = document.getElementById('cover-preview');
const coverImage = document.getElementById('cover-image');
const recentList = document.getElementById('recent-list');

// State
let isScanning = false;

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
    // Set default date to today
    const dateInput = document.getElementById('dateRead');
    dateInput.value = new Date().toISOString().split('T')[0];
    
    // Check authentication
    await checkAuth();
    
    // Setup event listeners
    setupEventListeners();
});

// Check if user is authenticated
async function checkAuth() {
    try {
        const response = await fetch(`${API_BASE}/auth/check`);
        const data = await response.json();
        
        if (data.authenticated) {
            showAdminScreen();
            loadRecentBooks();
        } else {
            showLoginScreen();
        }
    } catch (error) {
        console.error('Auth check failed:', error);
        showLoginScreen();
    }
}

function showLoginScreen() {
    loginScreen.classList.remove('hidden');
    adminScreen.classList.add('hidden');
}

function showAdminScreen() {
    loginScreen.classList.add('hidden');
    adminScreen.classList.remove('hidden');
    setupGoalForm();
}

// Setup goal form
async function setupGoalForm() {
    const yearSelect = document.getElementById('goal-year');
    const targetInput = document.getElementById('goal-target');
    const goalForm = document.getElementById('goal-form');
    
    if (!yearSelect || !goalForm) return;
    
    // Populate years (current year and next few years)
    const currentYear = new Date().getFullYear();
    yearSelect.innerHTML = '';
    for (let y = currentYear - 1; y <= currentYear + 2; y++) {
        const option = document.createElement('option');
        option.value = y;
        option.textContent = y;
        if (y === currentYear) option.selected = true;
        yearSelect.appendChild(option);
    }
    
    // Load current goal
    await loadGoal(currentYear);
    
    // Handle year change
    yearSelect.addEventListener('change', () => {
        loadGoal(parseInt(yearSelect.value));
    });
    
    // Handle form submit
    goalForm.addEventListener('submit', handleSetGoal);
}

async function loadGoal(year) {
    const targetInput = document.getElementById('goal-target');
    try {
        const response = await fetch(`${API_BASE}/goals/${year}`);
        const data = await response.json();
        targetInput.value = data.target || '';
    } catch (error) {
        console.error('Failed to load goal:', error);
    }
}

async function handleSetGoal(e) {
    e.preventDefault();
    
    const year = parseInt(document.getElementById('goal-year').value);
    const target = parseInt(document.getElementById('goal-target').value);
    const message = document.getElementById('goal-message');
    
    if (!target || target < 1) {
        showGoalMessage('Please enter a valid goal', 'error');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/goals`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ year, target })
        });
        
        if (response.ok) {
            showGoalMessage(`Goal set: ${target} books for ${year}`, 'success');
        } else {
            showGoalMessage('Failed to save goal', 'error');
        }
    } catch (error) {
        showGoalMessage('Failed to save goal', 'error');
    }
}

function showGoalMessage(text, type) {
    const message = document.getElementById('goal-message');
    message.textContent = text;
    message.className = type;
    message.classList.remove('hidden');
    setTimeout(() => message.classList.add('hidden'), 3000);
}

// Setup all event listeners
function setupEventListeners() {
    // Login form
    loginForm.addEventListener('submit', handleLogin);
    
    // Logout
    logoutBtn.addEventListener('click', handleLogout);
    
    // Export
    exportBtn.addEventListener('click', handleExport);
    
    // Book form
    quickAddForm.addEventListener('submit', handleAddBook);
    
    // ISBN lookup
    lookupBtn.addEventListener('click', handleIsbnLookup);
    
    // Barcode scanner
    scanBtn.addEventListener('click', startScanner);
    cancelScanBtn.addEventListener('click', stopScanner);
    
    // ISBN field change - show cover preview
    document.getElementById('isbn').addEventListener('change', updateCoverPreview);
}

// Handle login
async function handleLogin(e) {
    e.preventDefault();
    
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password })
        });
        
        if (response.ok) {
            showAdminScreen();
            loadRecentBooks();
            loginError.classList.add('hidden');
        } else {
            loginError.textContent = 'Invalid password';
            loginError.classList.remove('hidden');
        }
    } catch (error) {
        loginError.textContent = 'Login failed. Please try again.';
        loginError.classList.remove('hidden');
    }
}

// Handle logout
async function handleLogout() {

// Handle export - download books.json backup
async function handleExport() {
    try {
        const response = await fetch(`${API_BASE}/export`);
        if (!response.ok) throw new Error('Export failed');
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'books.json';
        a.click();
        URL.revokeObjectURL(url);
    } catch (error) {
        console.error('Export failed:', error);
        showMessage('Failed to export books. Are you logged in?', 'error');
    }
}

    try {
        await fetch(`${API_BASE}/auth/logout`, { method: 'POST' });
    } catch (error) {
        console.error('Logout error:', error);
    }
    showLoginScreen();
}

// Handle add book form submission
async function handleAddBook(e) {
    e.preventDefault();
    
    const formData = new FormData(quickAddForm);
    const book = {
        title: formData.get('title'),
        author: formData.get('author'),
        dateRead: formData.get('dateRead')?.replace(/-/g, '/') || '',
        pages: parseInt(formData.get('pages')) || 0,
        yearPublished: parseInt(formData.get('yearPublished')) || 0,
        isbn: formData.get('isbn') || '',
        shelf: formData.get('shelf') || 'read',
        review: formData.get('review') || '',
        coverUrl: getCoverUrl(formData.get('isbn'))
    };
    
    try {
        const response = await fetch(`${API_BASE}/books`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(book)
        });
        
        if (response.ok) {
            showMessage('Book added successfully!', 'success');
            quickAddForm.reset();
            document.getElementById('dateRead').value = new Date().toISOString().split('T')[0];
            hideCoverPreview();
            loadRecentBooks();
        } else {
            const error = await response.text();
            showMessage(`Failed to add book: ${error}`, 'error');
        }
    } catch (error) {
        showMessage('Failed to add book. Please try again.', 'error');
    }
}

// ISBN lookup via Open Library
async function handleIsbnLookup() {
    const isbn = document.getElementById('isbn').value.trim();
    if (!isbn) {
        showMessage('Please enter an ISBN first', 'error');
        return;
    }
    
    lookupBtn.disabled = true;
    lookupBtn.textContent = '⏳';
    
    try {
        const response = await fetch(`https://openlibrary.org/isbn/${isbn}.json`);
        if (!response.ok) {
            throw new Error('Book not found');
        }
        
        const data = await response.json();
        
        // Fill in form fields
        if (data.title) {
            document.getElementById('title').value = data.title;
        }
        if (data.number_of_pages) {
            document.getElementById('pages').value = data.number_of_pages;
        }
        if (data.publish_date) {
            const year = extractYear(data.publish_date);
            if (year) {
                document.getElementById('yearPublished').value = year;
            }
        }
        
        // Fetch author info if available
        if (data.authors && data.authors.length > 0) {
            const authorKey = data.authors[0].key;
            const authorResponse = await fetch(`https://openlibrary.org${authorKey}.json`);
            if (authorResponse.ok) {
                const authorData = await authorResponse.json();
                document.getElementById('author').value = authorData.name || '';
            }
        }
        
        updateCoverPreview();
        showMessage('Book info loaded!', 'success');
        
    } catch (error) {
        showMessage('Could not find book. Try entering details manually.', 'error');
    } finally {
        lookupBtn.disabled = false;
        lookupBtn.textContent = '🔍';
    }
}

// Extract year from various date formats
function extractYear(dateStr) {
    const match = dateStr.match(/\d{4}/);
    return match ? parseInt(match[0]) : null;
}

// Barcode scanner
async function startScanner() {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        showMessage('Camera not supported on this device', 'error');
        return;
    }
    
    scannerModal.classList.remove('hidden');
    const video = document.getElementById('scanner-video');
    const status = document.getElementById('scanner-status');
    
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ 
            video: { facingMode: 'environment' } 
        });
        video.srcObject = stream;
        await video.play();
        
        isScanning = true;
        
        // Initialize Quagga
        if (typeof Quagga !== 'undefined') {
            Quagga.init({
                inputStream: {
                    name: 'Live',
                    type: 'LiveStream',
                    target: video,
                    constraints: {
                        facingMode: 'environment'
                    }
                },
                decoder: {
                    readers: ['ean_reader', 'ean_8_reader', 'upc_reader', 'upc_e_reader']
                }
            }, (err) => {
                if (err) {
                    console.error('Quagga init error:', err);
                    status.textContent = 'Scanner initialization failed';
                    return;
                }
                Quagga.start();
            });
            
            Quagga.onDetected((result) => {
                if (result.codeResult) {
                    const code = result.codeResult.code;
                    document.getElementById('isbn').value = code;
                    status.textContent = `Found: ${code}`;
                    stopScanner();
                    handleIsbnLookup();
                }
            });
        } else {
            status.textContent = 'Barcode library not loaded. Enter ISBN manually.';
        }
        
    } catch (error) {
        console.error('Camera error:', error);
        status.textContent = 'Could not access camera';
    }
}

function stopScanner() {
    isScanning = false;
    scannerModal.classList.add('hidden');
    
    const video = document.getElementById('scanner-video');
    if (video.srcObject) {
        video.srcObject.getTracks().forEach(track => track.stop());
        video.srcObject = null;
    }
    
    if (typeof Quagga !== 'undefined') {
        Quagga.stop();
    }
}

// Cover preview
function getCoverUrl(isbn) {
    if (!isbn) return '';
    return `https://covers.openlibrary.org/b/isbn/${isbn}-M.jpg`;
}

function updateCoverPreview() {
    const isbn = document.getElementById('isbn').value.trim();
    if (isbn) {
        const url = getCoverUrl(isbn);
        coverImage.src = url;
        coverImage.onerror = () => {
            coverPreview.classList.add('hidden');
        };
        coverImage.onload = () => {
            coverPreview.classList.remove('hidden');
        };
    } else {
        hideCoverPreview();
    }
}

function hideCoverPreview() {
    coverPreview.classList.add('hidden');
}

// Load recent books
async function loadRecentBooks() {
    try {
        // Get current year's books as a proxy for recent
        const year = new Date().getFullYear();
        const response = await fetch(`${API_BASE}/books?year=${year}`);
        const data = await response.json();
        
        if (data.books && data.books.length > 0) {
            // Show last 5 books
            const recent = data.books.slice(0, 5);
            recentList.innerHTML = recent.map(book => `
                <div class="recent-book">
                    <div class="recent-book-info">
                        <div class="recent-book-title">${escapeHtml(book.title)}</div>
                        <div class="recent-book-author">${escapeHtml(book.author)}</div>
                        <div class="recent-book-date">${book.dateRead || 'No date'}</div>
                    </div>
                </div>
            `).join('');
        } else {
            recentList.innerHTML = '<p style="color: #7f8c8d;">No books added yet this year.</p>';
        }
    } catch (error) {
        console.error('Failed to load recent books:', error);
        recentList.innerHTML = '<p style="color: #e74c3c;">Failed to load recent books.</p>';
    }
}

// Show message
function showMessage(text, type) {
    formMessage.textContent = text;
    formMessage.className = type;
    formMessage.classList.remove('hidden');
    
    setTimeout(() => {
        formMessage.classList.add('hidden');
    }, 5000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
