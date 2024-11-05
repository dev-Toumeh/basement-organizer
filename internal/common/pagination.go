package common

import (
	"basement/main/internal/logg"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
)

// retrieve the query Value from the Request and filter, sanitize it
func SearchString(r *http.Request) string {
	searchTrimmed := strings.TrimSpace(r.FormValue("query"))
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(searchTrimmed, "*")
}

// Generates pagination data for a given dataset.
// @param data - The main data structure used to render the template.
// @param count - The total number of records found from the search query.
// @param limit - The maximum number of records displayed per page.
// @param pageNr - The current page number.
// @return - A map containing the pagination information to be used in the template.
func Pagination(data map[string]any, count int, limit int, pageNr int) map[string]any {
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	logg.Debugf("limit: %d, totalPages: %d, results: %d", limit, totalPages, count)

	currentPage := 1
	if pageNr > totalPages {
		currentPage = totalPages
	} else {
		currentPage = pageNr
	}

	if totalPages == 0 {
		totalPages = 1
	}

	nextPage := currentPage + 1
	if nextPage < 1 {
		nextPage = 1
	}
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	if prevPage > totalPages {
		prevPage = totalPages
	}

	logg.Debugf("currentPage %d", currentPage)

	pages := make([]map[string]any, 0)

	// more pagination
	disablePrev := false
	disableNext := false
	disableFirst := false
	disableLast := false
	if currentPage == nextPage {
		disableNext = true
	}
	if currentPage == totalPages {
		disableLast = true
	}
	if currentPage == prevPage {
		disablePrev = true
	}
	if currentPage == 1 {
		disableFirst = true
	}

	pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", 1), "Limit": fmt.Sprint(limit),
		"ID": fmt.Sprintf("pagination-%d", 1), "Disabled": disableFirst})

	if totalPages >= 10 {
		disabled := false
		prevFive := currentPage - 5
		if prevFive < 1 {
			prevFive = 1
		}
		if currentPage == prevFive {
			disabled = true
		}

		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", prevFive), "Limit": fmt.Sprint(limit),
			"ID": fmt.Sprintf("pagination-%d", prevFive), "Disabled": disabled})
	}
	pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", prevPage), "Limit": fmt.Sprint(limit),
		"ID": fmt.Sprintf("pagination-%d", prevPage), "Disabled": disablePrev})
	pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", currentPage), "Limit": fmt.Sprint(limit),
		"Selected": true, "ID": fmt.Sprintf("pagination-%d", currentPage)})
	pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", nextPage), "Limit": fmt.Sprint(limit),
		"ID": fmt.Sprintf("pagination-%d", nextPage), "Disabled": disableNext})

	if totalPages >= 10 {
		disabled := false
		nextFive := currentPage + 5
		if nextFive > totalPages {
			nextFive = totalPages
		}
		if currentPage == nextFive {
			disabled = true
		}
		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", nextFive), "Limit": fmt.Sprint(limit),
			"ID": fmt.Sprintf("pagination-%d", nextFive), "Disabled": disabled})
	}
	pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", totalPages), "Limit": fmt.Sprint(limit),
		"ID": fmt.Sprintf("pagination-%d", totalPages), "Disabled": disableLast})

	// Putting required data for templates together.
	data["Pages"] = pages
	data["Limit"] = fmt.Sprint(limit)
	data["NextPage"] = nextPage
	data["PrevPage"] = prevPage
	data["PageNumber"] = currentPage

	move := false
	data["Move"] = move

	return data
}
