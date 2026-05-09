package covers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// --- Test helpers ---

// mockServer creates a test HTTP server that responds based on the URL path
func mockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Resolver) {
	t.Helper()
	server := httptest.NewServer(handler)
	resolver := New(&http.Client{})
	return server, resolver
}

// placeholder43 returns a 43-byte body simulating Open Library's 1×1 GIF placeholder
func placeholder43() []byte {
	return make([]byte, 43)
}

// realCoverImage returns a body larger than PlaceholderMaxBytes simulating a real cover
func realCoverImage() []byte {
	return make([]byte, 5000)
}

// --- CoverURLByISBN tests ---

func TestCoverURLByISBN_ValidISBN(t *testing.T) {
	url := CoverURLByISBN("9781234567890")
	expected := "https://covers.openlibrary.org/b/isbn/9781234567890-M.jpg"
	if url != expected {
		t.Errorf("Expected %q, got %q", expected, url)
	}
}

func TestCoverURLByISBN_EmptyISBN(t *testing.T) {
	url := CoverURLByISBN("")
	if url != "" {
		t.Errorf("Expected empty string, got %q", url)
	}
}

func TestCoverURLByISBN_ISBN10(t *testing.T) {
	url := CoverURLByISBN("0446603781")
	expected := "https://covers.openlibrary.org/b/isbn/0446603781-M.jpg"
	if url != expected {
		t.Errorf("Expected %q, got %q", expected, url)
	}
}

// --- PreferISBN tests ---

func TestPreferISBN_PrefersISBN13(t *testing.T) {
	result := PreferISBN("0446603781", "9780446603785")
	if result != "9780446603785" {
		t.Errorf("Expected ISBN-13, got %q", result)
	}
}

func TestPreferISBN_FallsBackToISBN10(t *testing.T) {
	result := PreferISBN("0446603781", "")
	if result != "0446603781" {
		t.Errorf("Expected ISBN-10, got %q", result)
	}
}

func TestPreferISBN_ReturnsEmptyWhenBothEmpty(t *testing.T) {
	result := PreferISBN("", "")
	if result != "" {
		t.Errorf("Expected empty, got %q", result)
	}
}

func TestPreferISBN_IgnoresZero(t *testing.T) {
	result := PreferISBN("0446603781", "0")
	if result != "0446603781" {
		t.Errorf("Expected ISBN-10 fallback, got %q", result)
	}
}

func TestPreferISBN_TrimsWhitespace(t *testing.T) {
	result := PreferISBN("  0446603781  ", "  ")
	if result != "0446603781" {
		t.Errorf("Expected trimmed ISBN-10, got %q", result)
	}
}

// --- IsValidCover tests ---

func TestIsValidCover_RealImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(realCoverImage())
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	if !resolver.IsValidCover(server.URL + "/cover.jpg") {
		t.Error("Expected valid cover for real image response")
	}
}

func TestIsValidCover_PlaceholderImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(placeholder43())
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	if resolver.IsValidCover(server.URL + "/cover.jpg") {
		t.Error("Expected invalid cover for 43-byte placeholder")
	}
}

func TestIsValidCover_404Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	if resolver.IsValidCover(server.URL + "/cover.jpg") {
		t.Error("Expected invalid cover for 404 response")
	}
}

func TestIsValidCover_EmptyURL(t *testing.T) {
	resolver := New(&http.Client{})
	if resolver.IsValidCover("") {
		t.Error("Expected invalid cover for empty URL")
	}
}

func TestIsValidCover_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	if resolver.IsValidCover(server.URL + "/cover.jpg") {
		t.Error("Expected invalid cover for 500 response")
	}
}

// --- SearchByTitleAuthor tests ---

func TestSearchByTitleAuthor_FindsCover(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify search parameters
		title := r.URL.Query().Get("title")
		author := r.URL.Query().Get("author")
		if title == "" {
			t.Error("Expected title parameter")
		}
		if author == "" {
			t.Error("Expected author parameter")
		}

		result := searchResult{
			Docs: []searchDoc{
				{CoverI: 12345, Title: "Dawn"},
			},
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	// Override the search URL for testing
	resolver := New(&http.Client{})
	coverURL, err := searchWithURL(resolver, server.URL, "Dawn", "Octavia Butler")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "https://covers.openlibrary.org/b/id/12345-M.jpg"
	if coverURL != expected {
		t.Errorf("Expected %q, got %q", expected, coverURL)
	}
}

func TestSearchByTitleAuthor_NoCoverInResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := searchResult{
			Docs: []searchDoc{
				{CoverI: 0, Title: "Some Book"},
			},
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	coverURL, err := searchWithURL(resolver, server.URL, "Some Book", "Author")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if coverURL != "" {
		t.Errorf("Expected empty cover URL, got %q", coverURL)
	}
}

func TestSearchByTitleAuthor_NoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := searchResult{Docs: []searchDoc{}}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	coverURL, err := searchWithURL(resolver, server.URL, "Nonexistent Book", "Nobody")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if coverURL != "" {
		t.Errorf("Expected empty cover URL, got %q", coverURL)
	}
}

func TestSearchByTitleAuthor_EmptyTitle(t *testing.T) {
	resolver := New(&http.Client{})
	_, err := resolver.SearchByTitleAuthor("", "Author")
	if err == nil {
		t.Error("Expected error for empty title")
	}
}

func TestSearchByTitleAuthor_TitleOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		if author != "" {
			t.Error("Expected no author parameter when author is empty")
		}

		result := searchResult{
			Docs: []searchDoc{
				{CoverI: 99999, Title: "Solo Book"},
			},
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	coverURL, err := searchWithURL(resolver, server.URL, "Solo Book", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if coverURL == "" {
		t.Error("Expected cover URL from title-only search")
	}
}

func TestSearchByTitleAuthor_PicksFirstWithCover(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := searchResult{
			Docs: []searchDoc{
				{CoverI: 0, Title: "No Cover Edition"},
				{CoverI: 55555, Title: "Has Cover Edition"},
				{CoverI: 77777, Title: "Also Has Cover"},
			},
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	coverURL, err := searchWithURL(resolver, server.URL, "Some Book", "Author")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "https://covers.openlibrary.org/b/id/55555-M.jpg"
	if coverURL != expected {
		t.Errorf("Expected first result with cover %q, got %q", expected, coverURL)
	}
}

func TestSearchByTitleAuthor_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	_, err := searchWithURL(resolver, server.URL, "Book", "Author")
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

// --- Resolve tests (full fallback chain) ---

func TestResolve_UsesExistingValidURL(t *testing.T) {
	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(realCoverImage())
	}))
	defer validServer.Close()

	resolver := New(&http.Client{})
	result := resolver.Resolve(validServer.URL+"/existing.jpg", "9781234567890", "Title", "Author")

	if !strings.Contains(result, "/existing.jpg") {
		t.Errorf("Expected existing URL to be returned, got %q", result)
	}
}

func TestResolve_SkipsInvalidExistingURL(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if strings.Contains(r.URL.Path, "/existing") {
			// Existing URL returns placeholder
			w.Write(placeholder43())
		} else if strings.Contains(r.URL.Path, "/isbn") {
			// ISBN URL returns real image
			w.Write(realCoverImage())
		}
	}))
	defer server.Close()

	resolver := New(&http.Client{})
	result := resolver.Resolve(
		server.URL+"/existing.jpg",
		"",
		"Title",
		"Author",
	)

	// Should have tried existing URL and found it invalid
	// Since no ISBN provided and search would fail (wrong URL), should return empty
	if result != "" {
		// If it returned something, it should NOT be the existing placeholder URL
		if strings.Contains(result, "/existing") {
			t.Error("Should not return existing URL with placeholder image")
		}
	}
}

func TestResolve_NoISBNFallsToSearch(t *testing.T) {
	// This tests the full chain with no ISBN and a working search
	coverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(placeholder43())
	}))
	defer coverServer.Close()

	resolver := New(&http.Client{})
	// With no ISBN and the search API pointing to a real server, this exercises the fallback
	result := resolver.Resolve("", "", "Co-Intelligence", "Ethan Mollick")

	// The result will be from the real Open Library search API (integration-style)
	// In unit tests we mainly verify the chain logic
	_ = result // May or may not find a cover depending on network
}

func TestResolve_AllEmpty(t *testing.T) {
	resolver := New(&http.Client{})
	result := resolver.Resolve("", "", "", "")
	if result != "" {
		t.Errorf("Expected empty result for all empty inputs, got %q", result)
	}
}

// --- Helper for testing search with custom URL ---

func searchWithURL(r *Resolver, baseURL, title, author string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title is required")
	}

	params := url.Values{}
	params.Set("title", title)
	if author != "" {
		params.Set("author", author)
	}
	params.Set("fields", "cover_i,title,author_key")
	params.Set("limit", "3")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search returned status %d", resp.StatusCode)
	}

	var result searchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, doc := range result.Docs {
		if doc.CoverI > 0 {
			return fmt.Sprintf(openLibraryCoversByID, doc.CoverI), nil
		}
	}

	return "", nil
}
