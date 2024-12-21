package database

import (
	"basement/main/internal/common"
	"basement/main/internal/logg"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Exists checks existence of entity (item, box, shelf, area).
// Returns error only if other internal errors happen.
func (db *DB) Exists(entityType string, id uuid.UUID) (bool, error) {
	err1 := ValidTable(entityType)
	err2 := ValidVirtualTable(entityType)
	// not a vaid table and not valid virtual table
	if err1 != nil && err2 != nil {
		return false, logg.NewError(fmt.Sprintf("no entity with type: \"%s\"", entityType))
	}

	var itemExists int
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)`, entityType)
	err := db.Sql.QueryRow(query, id.String()).Scan(&itemExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, logg.WrapErr(err)
		}
	}
	if itemExists != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *DB) deleteFrom(table string, id uuid.UUID) error {
	err := ValidTable(table)
	if err != nil {
		return logg.WrapErr(err)
	}

	stmt := fmt.Sprintf(`DELETE FROM %s WHERE id = ?;`, table)
	result, err := db.Sql.Exec(stmt, id.String())
	if err != nil {
		return logg.Errorf(`can't delete "%s" from "%s" %w`, id, table, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return logg.WrapErr(err)
	}
	if rowsAffected == 0 {
		return logg.NewError(fmt.Sprintf(`id "%s" not found in "%s"`, id, table))
	} else if rowsAffected != 1 {
		return logg.NewError(fmt.Sprintf(`unexpected number of rows affected (%d) while deleting "%s" from "%s"`, rowsAffected, id, table))
	}
	return nil
}

func ValidTable(table string) error {
	_, ok := (*mainTables)[table]
	if !ok {
		return logg.NewError(fmt.Sprintf(`"%s" is not a valid table`, table))
	}
	return nil
}

func ValidVirtualTable(table string) error {
	_, ok := (*virtualTables)[table]
	if !ok {
		return logg.NewError(fmt.Sprintf(`"%s" is not a valid virtual table`, table))
	}
	return nil
}

// MoveTo moves item/box/shelf to a box/shelf/area.
//
// Example move item to a box:
//
//	MoveTo("item", itemID, "box", boxID)
//
// To move things out set
//
//	toTableID = uuid.Nil
func (db *DB) MoveTo(table string, id uuid.UUID, toTable string, toTableID uuid.UUID) error {
	if id == toTableID {
		return logg.NewError(fmt.Sprintf(`can't move "%s" to itself. ID=%s`, table, id.String()))
	}
	err := ValidTable(table)
	if err != nil {
		return logg.Errorf(`table: "%s" %w`, table, err)
	}
	err = ValidTable(toTable)
	if err != nil {
		return logg.Errorf(`toTable: "%s" %w`, toTable, err)
	}

	errMsg := fmt.Sprintf(`moving "%s" "%s" to "%s" "%s"`, table, id.String(), toTable, toTableID.String())

	exists, err := db.Exists(table, id)
	if err != nil {
		return logg.WrapErr(err)
	}
	if exists == false {
		return logg.Errorf("%s %w", errMsg, ErrNotExist)
	}

	// check if the table where the item is being moved to exists
	if toTableID != uuid.Nil {
		exists, err := db.Exists(toTable, toTableID)
		if err != nil {
			return logg.WrapErr(err)
		}
		if !exists {
			return logg.Errorf("%s %w", errMsg, ErrNotExist)
		}
	}

	// Update the item's shelf_id
	stmt := fmt.Sprintf(`UPDATE %s SET %s_id = ? WHERE id = ?`, table, toTable)
	result, err := db.Sql.Exec(stmt, toTableID.String(), id.String())
	if err != nil {
		return logg.WrapErr(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return logg.WrapErr(err)
	}
	if rows != 1 {
		return logg.NewError(fmt.Sprintf("rows should be != 1 but is %d", rows))
	}
	logg.Debugf("moved %s to %s", id, toTableID)
	return nil
}

// allListRowsFrom returns all items/boxes/shelves/etc from FTS tables item_fts, box_fts, shelf_fts, area_fts.
func (db *DB) allListRowsFrom(listRowsTable string) (listRows []common.ListRow, err error) {
	err = ValidVirtualTable(listRowsTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	stmt := " SELECT " + ALL_FTS_COLS + " FROM " + listRowsTable + ";"

	rows, err := db.Sql.Query(stmt)
	if err != nil {
		return nil, logg.Errorf("%s %w", stmt, err)
	}
	defer rows.Close()

	for rows.Next() {
		var sqlrow SQLListRow

		err = rows.Scan(sqlrow.RowsToScan()...)
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		lrow, err := sqlrow.ToListRow()
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		listRows = append(listRows, *lrow)
	}
	return listRows, nil
}

// listRowsPaginatedFrom Handles paginated retrieval of rows from a virtual table with optional search filtering.
//
// listRowsTable must be valid fts table (item_fts, box_fts, shelf_fts, area_fts).
//
// Empty searchQuery will return all rows.
//
// Panics if page or limit is zero, both must be at least 1.
func (db *DB) listRowsPaginatedFrom(listRowsTable string, searchQuery string, limit int, page int) (listRows []common.ListRow, err error) {
	if page == 0 {
		panic("offset starts at 1, can't be 0")
	}
	if limit == 0 {
		panic("limit starts at 1, can't be 0")
	}

	err = ValidVirtualTable(listRowsTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	offset := (page - 1) * limit

	var stmt string
	var rows *sql.Rows

	if strings.TrimSpace(searchQuery) != "" {
		stmt = "" +
			"SELECT " + ALL_FTS_COLS + " " +
			"FROM " + listRowsTable + " " +
			"WHERE label MATCH ?" + " " +
			"LIMIT ? OFFSET ?;"
		rows, err = db.Sql.Query(stmt, searchQuery+"*", limit, offset)
	} else {
		stmt = "" +
			"SELECT " + ALL_FTS_COLS + " " +
			"FROM " + listRowsTable + " " +
			"LIMIT ? OFFSET ?;"
		rows, err = db.Sql.Query(stmt, limit, offset)
	}
	defer rows.Close()

	if err != nil {
		return []common.ListRow{}, fmt.Errorf("error while fetching rows from %s: %w", listRowsTable, err)
	}

	var sqlListRow SQLListRow

	for rows.Next() {
		err := rows.Scan(sqlListRow.RowsToScan()...)
		if err != nil {
			return []common.ListRow{}, fmt.Errorf("error while scanning %s row: %w", listRowsTable, err)
		}
		shelfRow, err := sqlListRow.ToListRow()
		if err != nil {
			return []common.ListRow{}, fmt.Errorf("error while converting a %s row to ListRow: %w", listRowsTable, err)
		}
		listRows = append(listRows, *shelfRow)
	}

	return listRows, nil
}

// innerListRowsFrom returns all items/boxes/shelves/etc belonging to another box, shelf or area.
//
// Example:
//
//	// get all items that belongs to a shelf.
//	innerListRowsFrom("shelf", shelf.ID, "item_fts")
//
// listRowsTable:
//
//	FROM "item_fts"
//	FROM "box_ftx"
//	...
//
// belongsToTable:
//
//	"item", "box", "shelf", ...
//
// belongsToTableID:
//
//	WHERE "item"_id = ID
//	WHERE "box"_id = ID
//	...
func (db *DB) innerListRowsFrom(belongsToTable string, belongsToTableID uuid.UUID, listRowsTable string) ([]*common.ListRow, error) {
	err := ValidVirtualTable(listRowsTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	err = ValidTable(belongsToTable)
	if err != nil {
		return nil, logg.WrapErr(err)
	}

	var listRows []*common.ListRow

	stmt := "SELECT " + ALL_FTS_COLS + " FROM " + listRowsTable + "	WHERE " + belongsToTable + "_id = ?;"

	rows, err := db.Sql.Query(stmt, belongsToTableID.String())
	if err != nil {
		return nil, logg.Errorf("%s %w", stmt, err)
	}
	defer rows.Close()
	for rows.Next() {
		var sqlrow SQLListRow

		err = rows.Scan(sqlrow.RowsToScan()...)
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		lrow, err := sqlrow.ToListRow()
		if err != nil {
			return nil, logg.WrapErr(err)
		}
		listRows = append(listRows, lrow)
	}
	return listRows, nil
}

// Helper function to check for null strings and return empty if null
func ifNullString(sqlStr sql.NullString) string {
	if sqlStr.Valid {
		return sqlStr.String
	}
	return ""
}

// Helper function to check for null strings and return empty if null
func ifNullFloat64(sqlFloat sql.NullFloat64) float64 {
	if sqlFloat.Valid {
		return sqlFloat.Float64
	}
	return float64(0)
}

// Helper function to check for null strings and return empty if null
func ifNullInt(sqlInt sql.NullInt64) int {
	if sqlInt.Valid {
		return int(sqlInt.Int64)
	}
	return 0
}

// Helper function to check for null UUIDs and return uuid.Nil if null
func ifNullUUID(sqlUUID sql.NullString) uuid.UUID {
	if sqlUUID.Valid {
		return uuid.FromStringOrNil(sqlUUID.String)
	}
	return uuid.Nil
}

func UUIDFromSqlString(boxID sql.NullString) (uuid.UUID, error) {
	if boxID.Valid {
		id, err := uuid.FromString(boxID.String)
		if err != nil {
			return uuid.Nil, logg.Errorf("error while converting the string id into uuid: %w", err)
		}
		return id, nil
	}
	return uuid.Nil, logg.Errorf("invalid Virtual Id string")
}
