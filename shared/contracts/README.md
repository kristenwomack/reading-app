# API Contracts

This directory contains the API contract specifications for the Reading Tracker application.

## OpenAPI Specification

See `openapi.yaml` for the complete API specification.

## Endpoints

### GET /api/years

Returns all years that have at least one book with a "Date Read" entry.

**Response**: `200 OK`
```json
{
  "years": [
    {"year": 2025, "count": 12},
    {"year": 2024, "count": 45},
    {"year": 2023, "count": 38}
  ]
}
```

### GET /api/books?year={year}

Returns all books read in the specified year with "read" shelf status.

**Parameters**:
- `year` (required): Year to filter books by (e.g., 2025)

**Response**: `200 OK`
```json
{
  "books": [
    {
      "title": "Co-Intelligence",
      "author": "Ethan Mollick",
      "dateRead": "2025/09/19",
      "pages": 0,
      "month": 9
    }
  ]
}
```

**Error Responses**:
- `400 Bad Request`: Invalid or missing year parameter
- `500 Internal Server Error`: Server error

### GET /api/stats?year={year}

Returns comprehensive reading statistics for the specified year.

**Parameters**:
- `year` (required): Year to calculate statistics for (e.g., 2025)

**Response**: `200 OK`
```json
{
  "year": 2025,
  "totalBooks": 12,
  "totalPages": 3240,
  "averagePerMonth": 1.0,
  "monthlyBreakdown": [
    {"month": 1, "monthName": "Jan", "count": 2},
    {"month": 2, "monthName": "Feb", "count": 1},
    ...
  ]
}
```

**Error Responses**:
- `400 Bad Request`: Invalid or missing year parameter
- `500 Internal Server Error`: Server error

## Contract Testing

Both backend and frontend have contract tests:

**Backend**: Tests verify that responses match the schema
**Frontend**: Tests verify that requests are formatted correctly

## Viewing the Specification

You can view the OpenAPI specification using:

- Swagger UI: https://editor.swagger.io/
- Redoc: https://redocly.github.io/redoc/
- VS Code: OpenAPI (Swagger) Editor extension

Simply paste the contents of `openapi.yaml` into any of these tools.
