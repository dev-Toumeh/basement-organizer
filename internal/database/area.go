package database

import (
	"basement/main/internal/areas"
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

type SQLArea struct {
	SQLBasicInfo
}

// RowsToScan returns list of pointers for *sql.Rows.Scan() method.
//
//	// example usage:
//	rows.Scan(listRow.RowsToScan()...)
func (b *SQLArea) RowsToScan() []any {
	return b.SQLBasicInfo.RowsToScan()
}

// Vals returns all scanned values as strings.
func (s SQLArea) Vals() []string {
	return s.SQLBasicInfo.Vals()
}

// this function used inside of AreaByField to convert the area sql struct into normal struct
func (s SQLArea) ToArea() (areas.Area, error) {
	area := areas.Area{}
	info, err := s.ToBasicInfo()
	if err != nil {
		return area, logg.WrapErr(err)
	}

	area.BasicInfo = info
	return area, nil
}

// Create New Item Record
func (db *DB) CreateArea(newArea areas.Area) (uuid.UUID, error) {
	exists := db.AreaExists(newArea.ID)

	if exists {
		return uuid.Nil, db.ErrorExist()
	}

	id, err := db.insertNewArea(newArea)
	if err != nil {
		return uuid.Nil, logg.Errorf("error while creating new Area: %v", err)
	}
	return id, nil
}

// check if the Area Exist based on given Field
func (db *DB) AreaExists(id uuid.UUID) bool {
	exists, err := db.Exists("area", id)
	if err != nil {
		logg.Fatal(err.Error())
	}

	return exists
}

// AreaIDs returns IDs of all areas.
func (db *DB) AreaIDs() (ids []uuid.UUID, err error) {
	sqlStatement := `SELECT id FROM area`
	rows, err := db.Sql.Query(sqlStatement)
	if err != nil {
		return ids, logg.Errorf("Error while executing Area ids: %w", err)
	}
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			return ids, logg.Errorf("Error scanning Area ids: %v", err)
		}
		ids = append(ids, uuid.FromStringOrNil(idStr))
	}

	return ids, nil
}

// update area data
func (db *DB) UpdateArea(area areas.Area, ignorePicture bool) error {
	exist := db.AreaExists(area.ID)
	if !exist {
		return logg.Errorf("the area does not exist")
	}

	var err error
	var stmt string
	var result sql.Result
	if ignorePicture {
		stmt = "UPDATE area SET label = ?, description = ?, qrcode = ? WHERE id = ?"
		result, err = db.Sql.Exec(stmt, area.Label, area.Description, area.QRCode, area.ID)
	} else {
		area.PreviewPicture, err = ResizePNG(area.Picture, 50)
		if err != nil {
			return logg.Errorf("Error while resizing picture of area '%s' to create a preview picture %w", area.Label, err)
		}
		stmt = "UPDATE area SET label = ?, description = ?, picture = ?, preview_picture = ?, qrcode = ? WHERE id = ?"
		result, err = db.Sql.Exec(stmt, area.Label, area.Description, area.Picture, area.PreviewPicture, area.QRCode, area.ID)
	}

	if err != nil {
		return logg.Errorf("something wrong happened while running the area update query:\n\"%s\" %w", stmt, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.Errorf("error while finding the area %w", err)
	}
	if rowsAffected == 0 {
		return logg.Errorf("the Record with the id: %s was not found; this should not have happened while updating", area.ID.String())
	} else if rowsAffected != 1 {
		return logg.Errorf("the id: %s has an unexpected number of rows affected (more than one or less than 0)", area.ID.String())
	}
	return nil
}

// delete Area
func (db *DB) DeleteArea(areaId uuid.UUID) error {
	id := areaId.String()

	areaExist := db.AreaExists(areaId)
	if !areaExist {
		return logg.Errorf(`the area with id="` + id + `" doesn't exist`)
	}

	err := db.deleteFrom("area", areaId)
	if err != nil {
		return logg.WrapErr(err)
	}
	return nil
}

// Get Area based on his ID
// Wrapper function for AreaByField
func (db *DB) AreaById(id uuid.UUID) (areas.Area, error) {
	area := areas.Area{}
	if !db.AreaExists(id) {
		return area, logg.Errorf("area is not exist \n")
	}
	area, err := db.areaByField("id", id.String())
	if err != nil {
		return area, logg.WrapErr(err)
	}
	return area, err
}

// Get Area based on given Field
func (db *DB) areaByField(field string, value string) (areas.Area, error) {
	var sqlArea SQLArea
	stmt := "SELECT " + ALL_AREA_COLS + " FROM area WHERE " + field + " = ?;"

	err := db.Sql.QueryRow(stmt, value).Scan(sqlArea.RowsToScan()...)
	if err != nil {
		return areas.Area{}, logg.WrapErr(err)
	}

	area, err := sqlArea.ToArea()
	if err != nil {
		return areas.Area{}, logg.WrapErr(err)
	}

	return area, nil
}

// insert new Area record in the Database
func (db *DB) insertNewArea(area areas.Area) (uuid.UUID, error) {
	if db.AreaExists(area.ID) {
		return uuid.Nil, db.ErrorExist()
	}

	sqlStatement := "INSERT INTO area (" + ALL_AREA_COLS + ") VALUES (?,?,?,?,?,?)"

	updatePicture(&area.Picture, &area.PreviewPicture)

	result, err := db.Sql.Exec(sqlStatement, area.ID.String(), area.Label, area.Description, area.Picture, area.PreviewPicture, area.QRCode)
	if err != nil {
		return uuid.Nil, logg.Errorf("Error while executing create new area statement: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, logg.Errorf("Error while executing create new area statement: %w", err)
	}
	if rowsAffected != 1 {
		return uuid.Nil, logg.Errorf("unexpected number of effected rows, check insirtNewArea")
	}

	return area.ID, nil
}

// AreaListRows retrieves virtual areas by label.
// If the query is empty or contains only spaces, it returns default results.
func (db *DB) AreaListRows(searchQuery string, limit int, page int) (listRows []common.ListRow, err error) {
	listRows, err = db.listRowsPaginatedFrom("area_fts", searchQuery, limit, page)
	if err != nil {
		return listRows, logg.WrapErr(err)
	}
	return listRows, nil
}

// Get the virtual Area based on his ID
func (db *DB) AreaListRowByID(id uuid.UUID) (listRow common.ListRow, err error) {
	exists := db.AreaExists(id)
	if !exists {
		return listRow, logg.NewError("the Area Id does not exsist in the virtual table")
	}

	query := "SELECT " + ALL_FTS_COLS + " FROM area_fts WHERE id = ?"
	row, err := db.Sql.Query(query, id.String())
	if err != nil {
		return listRow, fmt.Errorf("error while fetching the virtual area: %w", err)
	}

	var sqlListRow SQLListRow
	for row.Next() {
		err := row.Scan(sqlListRow.RowsToScan()...)
		if err != nil {
			return common.ListRow{}, fmt.Errorf("error while assigning the Data to the Virtualarea struct : %w", err)
		}
	}

	vArea, err := sqlListRow.ToListRow()
	if err != nil {
		return common.ListRow{}, err
	}
	return *vArea, nil
}

// returns the count of rows in the area_fts table that match the specified searchString.
// If queryString is empty, it returns the count of all rows in the table.
func (db *DB) AreaListCounter(searchString string) (count int, err error) {
	countQuery := `SELECT COUNT(*) FROM area_fts;`

	if searchString != "" {
		countQuery = ` SELECT COUNT(*) FROM area_fts WHERE label MATCH '` + searchString + `*'`
	}

	err = db.Sql.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while fetching the number of area from the database: %v", err)
	}
	return count, nil
}

// check if the area row  exist
func (db *DB) AreaRowExist(id uuid.UUID) (bool, error) {
	return db.Exists("area_fts", id)
}
