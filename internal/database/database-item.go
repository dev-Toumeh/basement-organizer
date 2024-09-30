package database

import (
	"basement/main/internal/env"
	"basement/main/internal/items"
	"basement/main/internal/logg"
	"basement/main/internal/server"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Create New Item Record
func (db *DB) CreateNewItem(newItem items.Item) error {
	if exist := db.ItemExist("label", string(newItem.Label)); exist {
		return db.ErrorExist()
	}
	err := db.insertNewItem(newItem)
	if err != nil {
		return err
	}
	return nil
}

// Get Item Record based on given Field
func (db *DB) ItemByField(field string, value string) (items.Item, error) {

	if !db.ItemExist(field, value) {
		return items.Item{}, sql.ErrNoRows
	}

	query := fmt.Sprintf("SELECT id, label, description, picture, quantity, weight, qrcode FROM item WHERE %s = ? \n", field)
	row := db.Sql.QueryRow(query, value)

	var item items.Item
	var idStr string
	err := row.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)

	if err != nil {
		return items.Item{}, fmt.Errorf("Error while checking if the Item is available: %w ", err)
	}
	item.Id = uuid.Must(uuid.FromString(idStr))
	return item, nil
}

// check if the Item exist
func (db *DB) ItemExist(field string, value string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM item WHERE %s = ?", field)
	var count int
	err := db.Sql.QueryRow(query, value).Scan(&count)
	if err != nil {
		log.Println("Error checking item existence:", err)
		return false
	}
	return count > 0
}

// Item returns new Item struct if id matches.
func (db *DB) Item(id string) (items.Item, error) {
	return db.ItemByField("id", id)
}

// ListItemById returns a single item with less information suitable for a list row.
func (db *DB) ListItemById(id uuid.UUID) (*items.VirtualItem, error) {
	item := &items.VirtualItem{}
	rows, err := db.Sql.Query(`
		SELECT 
            i.id, i.label, i.picture, i.box_id,
            b.id, b.label 
        FROM 
            item AS i
        LEFT JOIN 
            box AS b ON b.id = i.box_id 
        WHERE 
            i.id = ?;`, id.String())

	if err != nil {
		return item, logg.Errorf("%w", err)
	}
	defer rows.Close()

	var sqlItem SqlVirtualItem

	for rows.Next() {
		err = rows.Scan(&sqlItem.ItemID, &sqlItem.Label, &sqlItem.PreviewPicture, &sqlItem.OuterBoxID, &sqlItem.OuterBoxID, &sqlItem.OuterBoxLabel)
		if err != nil {
			return nil, logg.Errorf("%w", err)
		}

		if env.Development() {
			b := bytes.Buffer{}
			server.WriteJSON(&b, sqlItem)
			logg.Debugf("virtual item: %v", b.String())
		}

		if sqlItem.ItemID.Valid {
			item.Item_Id = uuid.Must(uuid.FromString(sqlItem.ItemID.String))
			item.Label = sqlItem.Label.String
			item.PreviewPicture = sqlItem.PreviewPicture.String
			item.Box_Id = uuid.Must(uuid.FromString(sqlItem.OuterBoxID.String))
			item.Box_label = sqlItem.OuterBoxLabel.String
		} else {
			return item, errors.New(fmt.Sprintf("Invalid UUID: \"%s\"", sqlItem.ItemID.String))
		}
	}

	return item, nil
}

// return items id's in array from type string
func (db *DB) ItemIDs() ([]string, error) {
	query := "SELECT id FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying item records: %v", err)
		return []string{}, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			continue
		}
		ids = append(ids, idStr)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}

	// Print all IDs
	return ids, nil
}

// here we run the insert new Item query separate from the public function
// it make the code more readable
func (db *DB) insertNewItem(item items.Item) error {
	sqlStatement := `INSERT INTO item (id, label, description, picture, quantity, weight, qrcode, box_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Sql.Exec(sqlStatement, item.Id.String(), item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode, item.BoxId.String())
	if err != nil {
		log.Printf("Error while executing create new item statement: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected while executing create new item statement: %v", err)
		return err
	}
	if rowsAffected != 1 {
		log.Println("No rows affected, item not added")
		return errors.New("item not added")
	}
	return nil
}

// update the item based on the id
func (db *DB) UpdateItem(ctx context.Context, item items.Item) error {
	sqlStatement := fmt.Sprintf(`UPDATE item Set label = "%s", description = "%s", picture = "%s",
    quantity = "%d", weight = "%s", qrcode = "%s" WHERE id = ?`,
		item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
	result, err := db.Sql.ExecContext(ctx, sqlStatement, item.Id.String())
	if err != nil {
		logg.Err(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected == 0 {
		err := errors.New(fmt.Sprintf("the Record with the id: %s was not found that should not happened while updating", item.Id.String()))
		logg.Debug(err)
		return err
	} else if rowsAffected != 1 {
		err := errors.New(fmt.Sprintf("the id: %s has unexpected effected number of rows (more than one or less than 0)", item.Id.String()))
		logg.Err(err)
		return err
	}
	return nil
}

// Delete Item by Id
func (db *DB) DeleteItem(itemId uuid.UUID) error {
	sqlStatement := `DELETE FROM item WHERE id = ?;`
	result, err := db.Sql.Exec(sqlStatement, itemId.String())
	if err != nil {
		return fmt.Errorf("deleting was not succeed %W", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected == 0 {
		err := errors.New(fmt.Sprintf("the Record with the id: %s was not found that should not happened while deleting", itemId.String()))
		logg.Debug(err)
		return err
	} else if rowsAffected != 1 {
		err := errors.New(fmt.Sprintf("the id: %s has unexpected effected number of rows (more than one or less than 0)", itemId.String()))
		logg.Err(err)
		return err
	}
	return nil
}

// return all the available Items
func (db *DB) Items() ([][]string, error) {
	query := "SELECT id, label, description, picture, quantity, weight, qrcode FROM item;"
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return [][]string{}, err
	}
	defer rows.Close()

	var itemsArray [][]string
	var item items.Item
	for rows.Next() {
		var idStr string
		err := rows.Scan(&idStr, &item.Label, &item.Description, &item.Picture, &item.Quantity, &item.Weight, &item.QRcode)
		if err != nil {
			log.Printf("Error scanning item record: %v", err)
			return [][]string{}, err
		}

		formatted := fmt.Sprintf("id: %s, label: %s, description: %s, picture: %s, quantity: %d, weight: %s, qrcode: %s \n",
			idStr, item.Label, item.Description, item.Picture, item.Quantity, item.Weight, item.QRcode)
		itemsArray = append(itemsArray, []string{formatted})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
	return itemsArray, nil
}

// this is dynamic function but not ready
// am not really convinced from repeating the process every time i want to retrieve the data,
func (db *DB) ItemExperement(query string, refs []interface{}) {
	rows, err := db.Sql.Query(query)
	if err != nil {
		log.Printf("Error querying user records: %v", err)
		return
	}
	defer rows.Close()

	var results [][]interface{}

	for rows.Next() {
		err := rows.Scan(refs...)
		if err != nil {
			log.Printf("Error scanning item: %v", err)
			continue
		}

		// Copy the data from refs to a new slice to store the results
		row := make([]interface{}, len(refs))
		for i, ref := range refs {
			row[i] = *ref.(*interface{})
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
	}
	fmt.Print(results)
}

// delete one item or more
func (db *DB) DeleteItems(itemIds []uuid.UUID) error {
	if len(itemIds) == 0 {
		return nil
	}

	// Create placeholders and arguments
	placeholders := make([]string, len(itemIds))
	args := make([]interface{}, len(itemIds))
	for i, id := range itemIds {
		placeholders[i] = "?"
		args[i] = id
	}

	// Join the placeholders with commas
	sqlStatement := `DELETE FROM item WHERE id IN (` + strings.Join(placeholders, ",") + `);`

	// Execute the query with the item IDs as arguments
	result, err := db.Sql.Exec(sqlStatement, args...)
	if err != nil {
		logg.Err(err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logg.Err(err)
		return err
	}
	if rowsAffected != int64(len(itemIds)) {
		err := fmt.Errorf("unexpected number of rows affected while deleting. Expected: %d, Actual: %d", len(itemIds), rowsAffected)
		logg.Err(err)
		return err
	}

	return nil
}

func (db *DB) MoveItem(id1 uuid.UUID, id2 uuid.UUID) error {
	// updateStmt := `UPDATE item SET outerbox_id = ? WHERE Id = ?;`
	updateStmt := `UPDATE item SET box_id = ? WHERE id = ?;`
	_, err := db.Sql.Exec(updateStmt, id2, id1)
	if err != nil {
		return logg.Errorf("Placeholder function %w", err)
	}
	return nil
}
