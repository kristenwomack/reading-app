package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Book represents a book in the database
type Book struct {
	ID                      int64
	Title                   string
	Author                  string
	AdditionalAuthors       string
	ISBN                    string
	ISBN13                  string
	Publisher               string
	Pages                   int
	YearPublished           int
	OriginalPublicationYear int
	DateRead                string
	DateAdded               string
	Shelf                   string
	Review                  string
	CoverURL                string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// Store handles database operations
type Store struct {
	db *sql.DB
}

// New creates a new Store with the given database path
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// migrate runs database migrations
func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		additional_authors TEXT DEFAULT '',
		isbn TEXT DEFAULT '',
		isbn13 TEXT DEFAULT '',
		publisher TEXT DEFAULT '',
		pages INTEGER DEFAULT 0,
		year_published INTEGER DEFAULT 0,
		original_publication_year INTEGER DEFAULT 0,
		date_read TEXT DEFAULT '',
		date_added TEXT DEFAULT '',
		shelf TEXT DEFAULT 'read',
		review TEXT DEFAULT '',
		cover_url TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_books_date_read ON books(date_read);
	CREATE INDEX IF NOT EXISTS idx_books_shelf ON books(shelf);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS goals (
		year INTEGER PRIMARY KEY,
		book_target INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(schema)
	return err
}

// GetAllBooks returns all books from the database
func (s *Store) GetAllBooks() ([]Book, error) {
	rows, err := s.db.Query(`
		SELECT id, title, author, additional_authors, isbn, isbn13, publisher,
		       pages, year_published, original_publication_year, date_read,
		       date_added, shelf, review, cover_url, created_at, updated_at
		FROM books
		ORDER BY date_read DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.AdditionalAuthors, &b.ISBN, &b.ISBN13,
			&b.Publisher, &b.Pages, &b.YearPublished, &b.OriginalPublicationYear,
			&b.DateRead, &b.DateAdded, &b.Shelf, &b.Review, &b.CoverURL,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// GetBook returns a single book by ID
func (s *Store) GetBook(id int64) (*Book, error) {
	var b Book
	err := s.db.QueryRow(`
		SELECT id, title, author, additional_authors, isbn, isbn13, publisher,
		       pages, year_published, original_publication_year, date_read,
		       date_added, shelf, review, cover_url, created_at, updated_at
		FROM books WHERE id = ?
	`, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.AdditionalAuthors, &b.ISBN, &b.ISBN13,
		&b.Publisher, &b.Pages, &b.YearPublished, &b.OriginalPublicationYear,
		&b.DateRead, &b.DateAdded, &b.Shelf, &b.Review, &b.CoverURL,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// CreateBook inserts a new book and returns its ID
func (s *Store) CreateBook(b *Book) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO books (title, author, additional_authors, isbn, isbn13, publisher,
		                   pages, year_published, original_publication_year, date_read,
		                   date_added, shelf, review, cover_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, b.Title, b.Author, b.AdditionalAuthors, b.ISBN, b.ISBN13, b.Publisher,
		b.Pages, b.YearPublished, b.OriginalPublicationYear, b.DateRead,
		b.DateAdded, b.Shelf, b.Review, b.CoverURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateBook updates an existing book
func (s *Store) UpdateBook(b *Book) error {
	_, err := s.db.Exec(`
		UPDATE books SET
			title = ?, author = ?, additional_authors = ?, isbn = ?, isbn13 = ?,
			publisher = ?, pages = ?, year_published = ?, original_publication_year = ?,
			date_read = ?, date_added = ?, shelf = ?, review = ?, cover_url = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, b.Title, b.Author, b.AdditionalAuthors, b.ISBN, b.ISBN13,
		b.Publisher, b.Pages, b.YearPublished, b.OriginalPublicationYear,
		b.DateRead, b.DateAdded, b.Shelf, b.Review, b.CoverURL, b.ID)
	return err
}

// DeleteBook removes a book by ID
func (s *Store) DeleteBook(id int64) error {
	_, err := s.db.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

// GetSetting retrieves a setting value
func (s *Store) GetSetting(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SetSetting stores a setting value
func (s *Store) SetSetting(key, value string) error {
	_, err := s.db.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}

// BookCount returns the total number of books
func (s *Store) BookCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	return count, err
}

// Goal represents a yearly reading goal
type Goal struct {
	Year       int
	BookTarget int
}

// GetGoal returns the goal for a specific year
func (s *Store) GetGoal(year int) (*Goal, error) {
	var g Goal
	err := s.db.QueryRow("SELECT year, book_target FROM goals WHERE year = ?", year).Scan(&g.Year, &g.BookTarget)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// SetGoal creates or updates a goal for a year
func (s *Store) SetGoal(year, bookTarget int) error {
	_, err := s.db.Exec(`
		INSERT INTO goals (year, book_target) VALUES (?, ?)
		ON CONFLICT(year) DO UPDATE SET book_target = excluded.book_target
	`, year, bookTarget)
	return err
}
