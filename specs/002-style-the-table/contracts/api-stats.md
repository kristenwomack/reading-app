# API Contracts: Summary Statistics

**Feature**: 002-style-the-table  
**Contract Status**: REUSED (no changes to existing API)

---

## Existing Contract: GET /api/stats

**Purpose**: This feature reuses the existing `/api/stats` endpoint without modifications. This document serves as reference for the contract being consumed.

**Endpoint**: `GET /api/stats?year={year}`

**Description**: Returns aggregated reading statistics for the specified year.

---

### Request

**Method**: `GET`

**URL Pattern**: `/api/stats`

**Query Parameters**:
| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `year` | integer | Yes | 4-digit year | `2025` |

**Headers**: None required

**Example Request**:
```http
GET /api/stats?year=2025 HTTP/1.1
Host: localhost:8080
```

---

### Response

**Success (200 OK)**:

**Content-Type**: `application/json`

**Schema**:
```json
{
  "year": <integer>,
  "total_books": <integer>,
  "total_pages": <integer>,
  "avg_pages_per_book": <integer>
}
```

**Field Definitions**:
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| `year` | integer | 4-digit year | Year for which stats were calculated |
| `total_books` | integer | >= 0 | Total number of books finished in the year |
| `total_pages` | integer | >= 0 | Sum of page counts across all books |
| `avg_pages_per_book` | integer | >= 0 | Average pages per book (floor division) |

**Example Response (Normal Year)**:
```json
{
  "year": 2025,
  "total_books": 42,
  "total_pages": 12450,
  "avg_pages_per_book": 296
}
```

**Example Response (Empty Year)**:
```json
{
  "year": 2024,
  "total_books": 0,
  "total_pages": 0,
  "avg_pages_per_book": 0
}
```

---

**Error Responses**:

**400 Bad Request** (Invalid year parameter):
```json
{
  "error": "Invalid year parameter"
}
```

**500 Internal Server Error** (Server error):
```json
{
  "error": "Internal server error"
}
```

---

## Frontend Usage

**Consumed By**: `frontend/src/api-client.js` (existing implementation)

**Current Implementation** (reference only, no changes):
```javascript
// In api-client.js
async function getStats(year) {
  const response = await fetch(`/api/stats?year=${year}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch stats: ${response.statusText}`);
  }
  return response.json(); // Returns ReadingStatistics object
}
```

**UI Binding**: Data from this endpoint is consumed by `ui.js` to populate the summary card.

---

## Contract Validation

**Existing Tests**: Backend contract tests exist in `backend/internal/handlers/*_test.go` (no changes needed)

**Frontend Validation**: Frontend tests will mock this contract response for UI testing (new tests added in this feature)

**Test Data Examples**:
```javascript
// In frontend/tests/ui.test.js (to be created)
const mockStatsResponse = {
  year: 2025,
  total_books: 42,
  total_pages: 12450,
  avg_pages_per_book: 296
};

const mockEmptyStatsResponse = {
  year: 2024,
  total_books: 0,
  total_pages: 0,
  avg_pages_per_book: 0
};
```

---

## Breaking Change Policy

**This feature does NOT introduce breaking changes.**

If future features require changes to `/api/stats`:
1. Maintain backward compatibility (add optional fields only)
2. OR: Version the endpoint (`/api/v2/stats`)
3. Update contract tests to verify both old and new formats

---

## Related Contracts

**Other endpoints used by the application** (not modified by this feature):
- `GET /api/years` - List available years
- `GET /api/books?year={year}` - Get books for year (used for detailed chart data)

These contracts remain unchanged and are not impacted by the summary table layout feature.

---

**Last Updated**: 2025-01-16  
**Contract Version**: Existing (no version change)  
**Backend Implementation**: `backend/internal/handlers/stats_handler.go` (no changes)
