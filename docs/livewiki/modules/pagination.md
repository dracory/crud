---
path: modules/pagination.md
page-type: module
summary: Documentation for the server-side pagination system including configuration, rendering, and integration with FuncRows.
tags: [module, pagination, page-size, server-side]
created: 2025-02-11
updated: 2025-02-11
version: 2.0.0
---

# Module: Pagination

## Purpose

The pagination module renders Bootstrap pagination controls for the entity manager when server-side pagination is enabled. It calculates total pages from the row count and page size, and generates navigation links.

**File:** `v2/pagination.go`

## Configuration

Server-side pagination requires two config fields:

| Field | Type | Description |
|-------|------|-------------|
| `PageSize` | `int` | Number of rows per page. `0` = disabled (default). |
| `FuncRowsCount` | `func(r *http.Request) (int64, error)` | Returns total row count. Required when `PageSize > 0`. |

## How It Works

1. The entity manager controller checks if `PageSize > 0` and `FuncRowsCount != nil`
2. Reads the current page from the `page` query parameter (defaults to 1)
3. Calls `FuncRowsCount(r)` to get the total number of rows
4. Calls `renderPagination(currentPage, totalRows, baseURL)` to generate HTML
5. Appends the pagination HTML to the entity manager container

## renderPagination Method

```go
func (crud *Crud) renderPagination(currentPage int, totalRows int64, baseURL string) string
```

**Behavior:**
- Returns empty string if `pageSize <= 0`
- Returns empty string if total pages <= 1
- Calculates total pages: `(totalRows + pageSize - 1) / pageSize`
- Generates Bootstrap pagination with:
  - Previous button (`«`) - disabled on first page
  - Numbered page links - active class on current page
  - Next button (`»`) - disabled on last page

**Generated HTML structure:**

```html
<nav aria-label="Page navigation">
  <ul class="pagination justify-content-center mt-3">
    <li class="page-item disabled">
      <a class="page-link" href="...&page=1">&laquo;</a>
    </li>
    <li class="page-item active">
      <a class="page-link" href="...&page=1">1</a>
    </li>
    <li class="page-item">
      <a class="page-link" href="...&page=2">2</a>
    </li>
    <li class="page-item">
      <a class="page-link" href="...&page=2">&raquo;</a>
    </li>
  </ul>
</nav>
```

## Integration with FuncRows

Your `FuncRows` implementation must handle pagination when `PageSize > 0`:

```go
FuncRows: func(r *http.Request) ([]crud.Row, error) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    pageSize := 25
    offset := (page - 1) * pageSize

    // Query with LIMIT and OFFSET
    products, err := db.Query(
        "SELECT id, name, status FROM products LIMIT ? OFFSET ?",
        pageSize, offset,
    )
    // Convert to []crud.Row and return
},

FuncRowsCount: func(r *http.Request) (int64, error) {
    return db.Count("SELECT COUNT(*) FROM products")
},
```

## Client-Side vs Server-Side Pagination

| Feature | Client-Side (default) | Server-Side |
|---------|----------------------|-------------|
| `PageSize` | `0` | `> 0` |
| `FuncRowsCount` | Not needed | Required |
| Data loaded | All rows at once | One page at a time |
| Pagination UI | jQuery DataTables | Bootstrap pagination |
| Best for | Small datasets (< 1000 rows) | Large datasets |

## See Also

- [Modules: Controllers](controllers.md) - Entity manager controller
- [Configuration](../configuration.md) - PageSize and FuncRowsCount config
- [Cheatsheet](../cheatsheet.md) - Pagination setup example
