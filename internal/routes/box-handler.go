package routes

import (
	"basement/main/internal/auth"
	"basement/main/internal/common"
	"basement/main/internal/database"
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"basement/main/internal/templates"
	"fmt"
	"maps"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type BoxDatabase interface {
	CreateBox(newBox *items.Box) (uuid.UUID, error)
	MoveBox(id1 uuid.UUID, id2 uuid.UUID) error
	UpdateBox(box items.Box) error
	DeleteBox(boxId uuid.UUID) error
	BoxById(id uuid.UUID) (items.Box, error)
	BoxIDs() ([]string, error) // @TODO: Change string to uuid.UUID
	BoxFuzzyFinder(query string, limit int, page int) ([]items.BoxListRow, error)
	BoxListRowByID(id uuid.UUID) (items.BoxListRow, error)
}

func registerBoxRoutes(db *database.DB) {
	// Box templates
	http.HandleFunc("/box", boxHandler(db))
	http.HandleFunc("/box/{id}/move", boxPageMove(db))
	// Box api
	http.HandleFunc("/api/v1/box", boxHandler(db))
	http.HandleFunc("/api/v1/box/{id}", boxHandler(db))
	http.HandleFunc("/api/v1/box/{id}/move", moveBox(db))
	// Boxes templates
	http.HandleFunc("/boxes", boxesPage(db))
	http.HandleFunc("/boxes/{id}", boxDetailsPage(db))
	http.HandleFunc("/boxes/move", boxesPageMove(db))
	http.HandleFunc("/boxes-list", boxesHandler(db))
	// Boxes api
	http.HandleFunc("/api/v1/boxes", boxesHandler(db))
	http.HandleFunc("/api/v1/boxes/move", moveBoxes(db))
}

// boxHandler handles read, create, update and delete for single box.
func boxHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			const errMsgForUser = "Can't find box"

			id := server.ValidID(w, r, errMsgForUser)
			if id.IsNil() {
				return
			}

			box, err := db.BoxById(id)
			if err != nil {
				server.WriteNotFoundError(errMsgForUser, err, w, r)
				return
			}

			// Use API data writer
			if !wantsTemplateData(r) {
				server.WriteJSON(w, box)
				return
			}

			// Template writer
			renderBoxTemplate(&box, w, r)
			break

		case http.MethodPost:
			createBox(w, r, db)
			break

		case http.MethodDelete:
			deleteBox(w, r, db)
			return

		case http.MethodPut:
			updateBox(w, r, db)
			break

		default:
			// Other methods are not allowed.
			w.Header().Add("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method:'", r.Method, "' not allowed")
		}
	}
}

// boxesHandler handles read and delete for multiple boxes.
func boxesHandler(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			if !wantsTemplateData(r) {
				boxes, err := db.BoxFuzzyFinder("", 5, 1)
				if err != nil {
					server.WriteInternalServerError("cant query boxes", logg.Errorf("%w", err), w, r)
					return
				}
				server.WriteJSON(w, boxes)
				return
			}

			boxes, err := db.BoxFuzzyFinder("", 2, 1)
			if err != nil {
				server.WriteNotFoundError("Can't find boxes", err, w, r)
				return
			}
			if wantsTemplateData(r) {
				err = renderBoxesListTemplate(w, r, db, boxes, "")
				if err != nil {
					server.WriteInternalServerError("Can't render box list", err, w, r)
				}
			} else {
				server.WriteJSON(w, boxes)
			}
			break

		case http.MethodPost:
			query := r.PostFormValue("query")
			logg.Debugf("search query: %s", query)
			boxes, err := db.BoxFuzzyFinder(query, 5, 1)
			if err != nil {
				server.WriteInternalServerError("cant query boxes", err, w, r)
				return
			}
			err = renderBoxesListTemplate(w, r, db, boxes, query)
			if err != nil {
				server.WriteInternalServerError("cant render boxlist", err, w, r)
				return
			}

		case http.MethodPut:
			server.WriteNotImplementedWarning("Multiple boxes edit?", w, r)
			break

		case http.MethodDelete:
			deleteBoxes(w, r, db)
			break

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			break
		}
	}
}

// boxesPage shows a page with a box list.
func boxesPage(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		// page template
		page := templates.NewPageTemplate()
		page.Title = "Boxes"
		page.Authenticated = authenticated
		page.User = user
		data := page.Map()

		// search-input template
		query := r.FormValue("query")
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputValue = query
		logg.Debugf("searchInput %v", searchInput.Map())
		maps.Copy(data, searchInput.Map())

		// @TODO: Implement move page
		move := false
		urlQuery := r.URL.Query()
		logg.Debugf("query values: %v", urlQuery)

		pageNr, err := strconv.Atoi(r.FormValue("page"))
		if err != nil || pageNr < 1 {
			pageNr = 1
		}
		limit, err := strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			limit = env.DefaultTableSize()
		}
		if limit == 0 {
			limit = env.DefaultTableSize()
		}

		// box-list-row to fill box-list template
		logg.Info("has query: ", urlQuery.Has("query"))
		var boxes []items.BoxListRow
		err = nil
		usedSearch := false
		if urlQuery.Has("query") && query != "" {
			boxes, err = db.BoxFuzzyFinder(query, 10000, 1)
			usedSearch = true
		} else {
			boxes, err = db.BoxFuzzyFinder(query, limit, pageNr)
		}
		if err != nil {
			server.WriteInternalServerError("cant query boxes", err, w, r)
			return
		}

		// pagination
		results := 0
		totalPages := 1
		if usedSearch {
			results = len(boxes)

		} else {
			ids, _ := db.BoxIDs()
			results = len(ids)
		}
		logg.Debugf("limit: %d, boxes: %d, totalPages: %d, results: %d", limit, results, totalPages, results)

		totalPagesF := float64(results) / float64(limit)
		totalPagesCeil := math.Ceil(float64(totalPagesF))
		totalPages = int(totalPagesCeil)
		if totalPages < 1 {
			totalPages = 1
		}

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
		// Search is not paginated and returns all results.
		// Limit items per page manually
		if usedSearch {
			fromOffset := (currentPage - 1) * limit
			toOffset := currentPage * limit
			if toOffset > results {
				toOffset = results
			}
			if toOffset < 0 {
				toOffset = 0
			}
			if fromOffset < 0 {
				fromOffset = 0
			}
			logg.Debugf("fromOffset: %d, toOffset: %d", fromOffset, toOffset)
			boxes = boxes[fromOffset:toOffset]
		}

		// fill Boxes field for box-list template
		boxesMaps := make([]map[string]any, len(boxes))
		for i := range boxes {
			boxesMaps[i] = boxes[i].Map()
			boxesMaps[i]["Move"] = move
		}
		fillEmpty := limit - len(boxes)
		for range fillEmpty {
			boxesMaps = append(boxesMaps, map[string]any{})
		}
		maps.Copy(data, map[string]any{"Boxes": boxesMaps})
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

		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", 1), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", 1), "Disabled": disableFirst})

		if totalPages >= 10 {
			disabled := false
			prevFive := currentPage - 5
			if prevFive < 1 {
				prevFive = 1
			}
			if currentPage == prevFive {
				disabled = true
			}
			pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", prevFive), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", prevFive), "Disabled": disabled})
		}
		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", prevPage), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", prevPage), "Disabled": disablePrev})
		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", currentPage), "Limit": fmt.Sprint(limit), "Selected": true, "ID": fmt.Sprintf("pagination-%d", currentPage)})
		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", nextPage), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", nextPage), "Disabled": disableNext})
		if totalPages >= 10 {
			disabled := false
			nextFive := currentPage + 5
			if nextFive > totalPages {
				nextFive = totalPages
			}
			if currentPage == nextFive {
				disabled = true
			}
			pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", nextFive), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", nextFive), "Disabled": disabled})
		}

		pages = append(pages, map[string]any{"PageNumber": fmt.Sprintf("%d", totalPages), "Limit": fmt.Sprint(limit), "ID": fmt.Sprintf("pagination-%d", totalPages), "Disabled": disableLast})

		// Putting required data for templates together.
		data["Pages"] = pages
		data["Limit"] = fmt.Sprint(limit)
		data["NextPage"] = nextPage
		data["PrevPage"] = prevPage
		data["PageNumber"] = currentPage
		data["Move"] = move

		server.MustRender(w, r, templates.TEMPLATE_BOXES_PAGE, data)
	}
}

// boxDetailsPage shows a page with details of a specific box.
func boxDetailsPage(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, _ := auth.Authenticated(r)
		user, _ := auth.UserSessionData(r)

		id := server.ValidID(w, r, "no box")
		if id.IsNil() {
			return
		}
		logg.Debug(id)

		notFound := false
		box, err := db.BoxById(id)
		if err != nil {
			notFound = true
		}
		box.ID = id
		data := items.BoxPageTemplateData()
		data.Box = &box

		data.Title = fmt.Sprintf("Box - %s", box.Label)
		data.Authenticated = authenticated
		data.User = user
		data.NotFound = notFound
		nd := data.Map()
		maps.Copy(nd, map[string]any{"Boxes": &box.InnerBoxes})
		searchInput := items.NewSearchInputTemplate()
		searchInput.SearchInputLabel = "Search boxes"
		searchInput.SearchInputHxTarget = "#box-list"
		searchInput.SearchInputHxPost = "/boxes"
		maps.Copy(nd, searchInput.Map())
		itemsMap := make([]map[string]any, len(box.Items))
		for i, v := range box.Items {
			itemsMap[i] = v.Map()
		}
		nd["ItemListItems"] = itemsMap

		server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS_PAGE, nd)
	}
}

func createBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	box := items.NewBox()
	logg.Debug("create box: ", box)
	id, err := db.CreateBox(&box)
	if err != nil {
		server.WriteNotFoundError("error while creating the box", err, w, r)
		return
	}
	if wantsTemplateData(r) {
		box, err := db.BoxListRowByID(id)
		logg.Debug(box)
		if err != nil {
			server.WriteNotFoundError("error while fetching the box based on Id", err, w, r)
			return
		}
		server.MustRender(w, r, templates.TEMPLATE_BOX_LIST_ROW, box.Map())
	} else {
		server.WriteJSON(w, id)
	}
}

func updateBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't update box."
	id := server.ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}

	box := boxFromPostFormValue(id, r)
	err := db.UpdateBox(box)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	if wantsTemplateData(r) {
		boxTemplate := items.BoxTemplateData{Box: &box, Edit: false}
		err := server.RenderWithSuccessNotification(w, r, templates.TEMPLATE_BOX_DETAILS, boxTemplate.Map(), fmt.Sprintf("Updated box: %v", boxTemplate.Label))
		if err != nil {
			server.WriteInternalServerError(errMsgForUser, err, w, r)
			return
		}
	} else {
		server.WriteJSON(w, box)
	}
	logg.Debug("Updated Box: ", box)
}

func deleteBoxes(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't delete boxes"
	r.ParseForm()
	toDelete := make([]uuid.UUID, 0)
	for k, v := range r.Form {
		logg.Debugf("k: %v, v:%v", k, v)
		if strings.Contains(k, "delete:") {
			ids := strings.Split(k, "delete:")
			if len(ids) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, errMsgForUser)
				logg.Debugf("Wrong delete key value pair: '%v'\n\t%s", k, errMsgForUser)
				return
			}
			id, err := uuid.FromString(ids[1])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, errMsgForUser)
				logg.Errorf("%s: Malformed uuid \"%s\". %w", errMsgForUser, k, err)
				return
			}
			toDelete = append(toDelete, id)
		}
	}
	deleteErrorIds := []string{}
	var err error
	errOccurred := false
	for _, deleteId := range toDelete {
		err = nil
		err = db.DeleteBox(deleteId)
		if err != nil {
			errOccurred = true
			deleteErrorIds = append(deleteErrorIds, deleteId.String())
			logg.Errorf("%v: %v. %w", errMsgForUser, deleteId, err)
		} else {
			logg.Debug("Box deleted: ", deleteId)
		}
	}
	if errOccurred {
		errIds := strings.Join(deleteErrorIds, ",")
		server.TriggerErrorNotification(w, errMsgForUser+errIds)
		// @TODO: Update partial table, even if error happens.
		return
	}

	if wantsTemplateData(r) {
		boxesPage(db).ServeHTTP(w, r)
		for _, id := range toDelete {
			templates.RenderSuccessNotification(w, "Box deleted: "+id.String())
		}
		return
	}
	fmt.Fprint(w, nil)
	w.WriteHeader(http.StatusOK)
}

// deleteBox deletes a single box.
func deleteBox(w http.ResponseWriter, r *http.Request, db BoxDatabase) {
	errMsgForUser := "Can't delete box"
	id := server.ValidID(w, r, errMsgForUser)
	if id.IsNil() {
		return
	}
	err := db.DeleteBox(id)
	if err != nil {
		server.WriteNotFoundError(errMsgForUser, err, w, r)
		return
	}
	server.RedirectWithSuccessNotification(w, "/boxes", fmt.Sprintf("%s deleted", id))
}

// @TODO: Implement move box.
// boxPageMove shows the page where client can choose where to move the box.
func boxPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box page", w, r)
	}
}

// boxesPageMove shows the page where client can choose where to move the selected boxes from the boxes page.
func boxesPageMove(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move multiple boxes page", w, r)
	}
}

// @TODO: Implement moveBox.
func moveBox(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move single box", w, r)
	}
}

// @TODO: Implement move boxes.
func moveBoxes(db BoxDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.WriteNotImplementedWarning("Move multiple boxes", w, r)
	}
}

// boxFromPostFormValue returns items.Box without references to inner boxes, outer box and items.
func boxFromPostFormValue(id uuid.UUID, r *http.Request) items.Box {
	box := items.Box{}
	box.ID = id
	box.Label = r.PostFormValue("label")
	box.Description = r.PostFormValue("description")
	box.Picture = common.ParsePicture(r)
	box.QRcode = r.PostFormValue("qrcode")
	return box
}

// wantsTemplateData checks if current request requires template data.
// Helps deciding how to write the data.
func wantsTemplateData(r *http.Request) bool {
	return !strings.Contains(r.URL.Path, "/api/")
}

func renderBoxTemplate(box *items.Box, w http.ResponseWriter, r *http.Request) {
	editParam := r.FormValue("edit")
	edit := false
	if editParam == "true" {
		edit = true
	}
	b := items.BoxTemplateData{Box: box, Edit: edit}
	server.MustRender(w, r, templates.TEMPLATE_BOX_DETAILS, b.Map())
}

func renderBoxesListTemplate(w http.ResponseWriter, r *http.Request, db BoxDatabase, boxes []items.BoxListRow, query string) error {
	searchInput := items.NewSearchInputTemplate()
	searchInput.SearchInputLabel = "Search boxes"
	searchInput.SearchInputHxTarget = "#box-list"
	searchInput.SearchInputHxPost = "/boxes-list"
	searchInput.SearchInputValue = query
	logg.Debugf("searchInput %v", searchInput.Map())
	boxesMaps := make([]map[string]any, len(boxes))
	for i := range boxes {
		boxesMaps[i] = boxes[i].Map()
	}
	data := map[string]any{"Boxes": boxesMaps}
	maps.Copy(data, searchInput.Map())
	logg.Debug(data)

	err := templates.SafeRender(w, templates.TEMPLATE_BOX_LIST, data)
	if err != nil {
		return logg.Errorf("%w", err)
	}
	return nil
}

// // sampleBoxDB to test request handler.
// type sampleBoxDB struct {
// 	Boxes map[string]*items.Box
// }
//
// func newSampleBoxDB() *sampleBoxDB {
// 	db := sampleBoxDB{Boxes: make(map[string]*items.Box, 100)}
// 	for i := range 10 {
// 		box := items.NewBox()
// 		box.Label = fmt.Sprintf("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa %d", i)
// 		db.CreateBox(&box)
// 	}
// 	return &db
// }
//
// // func (db *sampleBoxDB) CreateBox() (uuid.UUID, error) {
// // 	box := items.NewBox()
// // 	db.Boxes[box.Id.String()] = &box
// // 	logg.Debug(db.Boxes)
// // 	return box.Id, nil
// // }
//
// func (db *sampleBoxDB) CreateBox(box *items.Box) (uuid.UUID, error) {
// 	// box := items.NewBox()
// 	db.Boxes[box.Id.String()] = box
// 	logg.Debug(db.Boxes)
// 	return box.Id, nil
// }
//
// func (db *sampleBoxDB) BoxById(id uuid.UUID) (items.Box, error) {
// 	box, ok := db.Boxes[id.String()]
// 	if !ok {
// 		// logg.Debug("BoxByID: ",db.Boxes)
// 		return items.Box{}, errors.New("ID " + id.String() + " doesn't exist")
// 	}
// 	return *box, nil
// }
//
// // func (db *sampleBoxDB) BoxIDs() ([]uuid.UUID, error) {
// func (db *sampleBoxDB) BoxIDs() ([]string, error) {
// 	// ids := make([]uuid.UUID, len(db.Boxes))
// 	ids := make([]string, len(db.Boxes))
// 	i := 0
// 	for _, v := range db.Boxes {
// 		// ids[i] = v.Id
// 		ids[i] = v.Id.String()
// 		i++
// 	}
// 	slices.Sort(ids)
// 	return ids, nil
// }
//
// func (db *sampleBoxDB) UpdateBox(box items.Box) error {
// 	oldBox, err := db.BoxById(box.Id)
// 	if err != nil {
// 		return logg.Errorf("UpdateBox(): %w", err)
// 	}
// 	db.Boxes[oldBox.Id.String()] = &box
// 	return nil
// }
//
// func (db *sampleBoxDB) DeleteBox(id uuid.UUID) error {
// 	_, err := db.BoxById(id)
// 	if err != nil {
// 		return logg.Errorf("DeleteBox(): %w", err)
// 	}
// 	delete(db.Boxes, id.String())
// 	return nil
// }
