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

	nextnextPage := nextPage + 1
	if nextnextPage < 1 {
		nextnextPage = 1
	}
	if nextnextPage > totalPages {
		nextnextPage = totalPages
	}

	prevPage := currPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	if prevPage > totalPages {
		prevPage = totalPages
	}

	prevprevPage := prevPage - 1
	if prevprevPage < 1 {
		prevprevPage = 1
	}
	if prevprevPage > totalPages {
		prevprevPage = totalPages
	}

	logg.Debugf("currentPage %d", currPage)

	pages := make([]PaginationButton, 0)

	numberFmt := "%d"
	if totalPages >= 10 {
		numberFmt = "%02d"
	}
	if totalPages >= 100 {
		numberFmt = "%03d"
	}
	if totalPages >= 1000 {
		numberFmt = "%04d"
	}

	// more pagination
	firstPage := 1
	disablePrev := false
	disableNext := false
	disableFirst := false
	disableLast := false
	disablePrevprev := false
	disableNextnext := false
	nextPageText := fmt.Sprintf(numberFmt, nextPage)
	prevPageText := fmt.Sprintf(numberFmt, prevPage)
	nextnextPageText := fmt.Sprintf(numberFmt, nextnextPage)
	prevprevPageText := fmt.Sprintf(numberFmt, prevprevPage)
	firstPageText := fmt.Sprintf(numberFmt, firstPage)
	currPageText := fmt.Sprintf(numberFmt, currPage)
	lastPageText := fmt.Sprintf(numberFmt, totalPages)
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
	if nextnextPage == nextPage {
		disableNextnext = true
	}
	if prevprevPage == prevPage {
		disablePrevprev = true
	}

	pages = append(pages, PaginationButton{PageNumber: firstPage, Disabled: disableFirst, Text: firstPageText})
	pages = append(pages, PaginationButton{PageNumber: prevprevPage, Disabled: disablePrevprev, Text: prevprevPageText})
	pages = append(pages, PaginationButton{PageNumber: prevPage, Disabled: disablePrev, Text: prevPageText})
	pages = append(pages, PaginationButton{PageNumber: currPage, Selected: true, Text: currPageText})
	pages = append(pages, PaginationButton{PageNumber: nextPage, Disabled: disableNext, Text: nextPageText})
	pages = append(pages, PaginationButton{PageNumber: nextnextPage, Disabled: disableNextnext, Text: nextnextPageText})
	pages = append(pages, PaginationButton{PageNumber: totalPages, Disabled: disableLast, Text: lastPageText})

	// Putting required data for templates together.
	data["Pages"] = pages
	data["Limit"] = limit
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

	paginationData := Pagination(data.GetTypeMap(), count, limit, pageNr)

	// Putting required data for templates together.
	data.SetPages(paginationData["Pages"].([]PaginationButton))
	data.SetNextPage(paginationData["NextPage"].(int))
	data.SetPrevPage(paginationData["PrevPage"].(int))
	data.SetPageNumber(paginationData["PageNumber"].(int))
	data.SetPaginationButtons(paginationData["Pages"].([]PaginationButton))
	data.SetMove(paginationData["Move"].(bool))
	data.SetPagination(true)

	return data
}
