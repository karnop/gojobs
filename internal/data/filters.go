package data

import (
	"strings"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // Allowed sort fields (security)
}

// CalculateMetadata calculates limit and offset for SQL
func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// sortColumn checks if the sort field is safe to use in SQL.
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
		// Handle descending order (e.g. "-salary")
		if f.Sort == "-"+safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	// Default fallback: panic-safe
	return "id" 
}

// sortDirection returns "ASC" or "DESC"
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}