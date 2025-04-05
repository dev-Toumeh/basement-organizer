package common

import (
	"basement/main/internal/env"
	"basement/main/internal/logg"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
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
//
// Example buttons layout for 99 total pages with current page = 15:
//
//	[firstPage] [prevFive] [prevPage] [currPage] [nextPage] [nextFive] [lastPage]
//	   [1]         [10]       [14]       [15]      [16]        [20]       [99]
func Pagination(data map[string]any, count int, limit int, pageNr int) map[string]any {
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	logg.Debugf("limit: %d, totalPages: %d, results: %d", limit, totalPages, count)

	currPage := 1
	if pageNr > totalPages {
		currPage = totalPages
	} else {
		currPage = pageNr
	}

	if totalPages == 0 {
		totalPages = 1
	}

	nextPage := currPage + 1
	if nextPage < 1 {
		nextPage = 1
	}
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	if prevPage > totalPages {
		prevPage = totalPages
	}

	logg.Debugf("currentPage %d", currPage)

	pages := make([]PaginationButton, 0)

	// more pagination
	disablePrev := false
	disableNext := false
	disableFirst := false
	disableLast := false
	if currPage == nextPage {
		disableNext = true
	}
	if currPage == totalPages {
		disableLast = true
	}
	if currPage == prevPage {
		disablePrev = true
	}
	if currPage == 1 {
		disableFirst = true
	}

	pages = append(pages, PaginationButton{PageNumber: 1, Disabled: disableFirst})

	if totalPages >= 10 {
		disabled := false
		prevFive := currPage - 5
		if prevFive < 1 {
			prevFive = 1
		}
		if currPage == prevFive {
			disabled = true
		}

		pages = append(pages, PaginationButton{PageNumber: prevFive, Disabled: disabled})
	}
	pages = append(pages, PaginationButton{PageNumber: prevPage, Disabled: disablePrev})
	pages = append(pages, PaginationButton{PageNumber: currPage, Selected: true})
	pages = append(pages, PaginationButton{PageNumber: nextPage, Disabled: disableNext})

	if totalPages >= 10 {
		disabled := false
		nextFive := currPage + 5
		if nextFive > totalPages {
			nextFive = totalPages
		}
		if currPage == nextFive {
			disabled = true
		}
		pages = append(pages, PaginationButton{PageNumber: nextFive, Disabled: disabled})
	}
	pages = append(pages, PaginationButton{PageNumber: totalPages, Disabled: disableLast})

	// Putting required data for templates together.
	data["Pages"] = pages
	data["Limit"] = fmt.Sprint(limit)
	data["NextPage"] = nextPage
	data["PrevPage"] = prevPage
	data["PageNumber"] = currPage

	move := false
	data["Move"] = move

	return data
}

func ParsePageNumber(r *http.Request) int {
	pageNr, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || pageNr < 1 {
		pageNr = 1
	}
	return pageNr
}

func ParseLimit(r *http.Request) int {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = env.CurrentConfig().DefaultTableSize()
	}
	if limit == 0 {
		limit = env.CurrentConfig().DefaultTableSize()
	}
	return limit
}

func ParseOrigin(r *http.Request) (origin string) {
	p := r.URL.Path
	switch {
	case strings.Contains(p, "/items"):
		origin = "Items"
		break
	case strings.Contains(p, "/boxes"):
		origin = "Boxes"
		break
	case strings.Contains(p, "/shelves"):
		origin = "Shelves"
		break
	case strings.Contains(p, "/areas"):
		origin = "Areas"
		break
	}
	return origin
}

func Pagination2(data Data) Data {
	count := data.GetCount()
	limit := data.GetLimit()
	pageNr := data.GetPageNumber()
	totalPages := int(math.Ceil(float64(count) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	logg.Debugf("limit: %d, totalPages: %d, results: %d", limit, totalPages, count)

	currPage := 1
	if pageNr > totalPages {
		currPage = totalPages
	} else {
		currPage = pageNr
	}

	if totalPages == 0 {
		totalPages = 1
	}

	nextPage := currPage + 1
	if nextPage < 1 {
		nextPage = 1
	}
	if nextPage > totalPages {
		nextPage = totalPages
	}

	prevPage := currPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	if prevPage > totalPages {
		prevPage = totalPages
	}

	logg.Debugf("currentPage %d", currPage)

	pages := make([]PaginationButton, 0)

	// more pagination
	disablePrev := false
	disableNext := false
	disableFirst := false
	disableLast := false
	if currPage == nextPage {
		disableNext = true
	}
	if currPage == totalPages {
		disableLast = true
	}
	if currPage == prevPage {
		disablePrev = true
	}
	if currPage == 1 {
		disableFirst = true
	}

	pages = append(pages, PaginationButton{PageNumber: 1, Disabled: disableFirst})

	if totalPages >= 10 {
		disabled := false
		prevFive := currPage - 5
		if prevFive < 1 {
			prevFive = 1
		}
		if currPage == prevFive {
			disabled = true
		}

		pages = append(pages, PaginationButton{PageNumber: prevFive, Disabled: disabled})
	}
	pages = append(pages, PaginationButton{PageNumber: prevPage, Disabled: disablePrev})
	pages = append(pages, PaginationButton{PageNumber: currPage, Selected: true})
	pages = append(pages, PaginationButton{PageNumber: nextPage, Disabled: disableNext})

	if totalPages >= 10 {
		disabled := false
		nextFive := currPage + 5
		if nextFive > totalPages {
			nextFive = totalPages
		}
		if currPage == nextFive {
			disabled = true
		}
		pages = append(pages, PaginationButton{PageNumber: nextFive, Disabled: disabled})
	}
	pages = append(pages, PaginationButton{PageNumber: totalPages, Disabled: disableLast})
	move := false

	// Putting required data for templates together.
	data.SetPages(pages)
	data.SetNextPage(nextPage)
	data.SetPrevPage(prevPage)
	data.SetPageNumber(currPage)
	data.SetPaginationButtons(pages)
	data.SetMove(move)
	data.SetPagination(true)

	return data
}
