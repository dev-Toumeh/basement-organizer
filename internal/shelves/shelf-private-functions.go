package shelves

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// retrieve the Page Number Value from the Request and filter, sanitize it
func pageNumber(r *http.Request) int {
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pageNumber = 1
	}

	return pageNumber
}

// retrieve the desired type from the Request (add, move or search) and validate it.
func typeFromRequest(r *http.Request) (string, error) {
	t := strings.TrimSpace(r.URL.Query().Get("type"))

	if t == "" {
		return "", nil
	}

	// Validate that 'type' is one of the allowed values
	allowedTypes := map[string]bool{"add": true, "move": true, "search": true}
	if _, ok := allowedTypes[t]; !ok {
		return "", fmt.Errorf("unexpected type: %s, while preparing the search Template", t)
	}

	return strings.ToUpper(t[:1]) + t[1:], nil
}
