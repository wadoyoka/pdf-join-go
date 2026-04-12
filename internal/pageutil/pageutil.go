package pageutil

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ParsePages parses a comma-separated string of page numbers (e.g. "2,5,8")
// into a sorted slice of positive integers with no duplicates.
func ParsePages(s string) ([]int, error) {
	if s == "" {
		return nil, fmt.Errorf("page numbers cannot be empty")
	}

	parts := strings.Split(s, ",")
	pages := make([]int, 0, len(parts))
	seen := make(map[int]bool)

	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid page number %q: %w", strings.TrimSpace(p), err)
		}
		if n < 1 {
			return nil, fmt.Errorf("page number must be positive, got %d", n)
		}
		if seen[n] {
			return nil, fmt.Errorf("duplicate page number: %d", n)
		}
		seen[n] = true
		pages = append(pages, n)
	}

	sort.Ints(pages)
	return pages, nil
}
