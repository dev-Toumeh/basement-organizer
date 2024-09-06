package database

import "fmt"

func (db *DB) SearchItemsByLabel(query string) ([]struct {
	Id    string
	Label string
}, error) {
	rows, err := db.Sql.Query(`
        SELECT item_id, label
        FROM item_fts
        WHERE label LIKE ? 
        ORDER BY item_id;
    `, query+"%")
	if err != nil {
		return nil, fmt.Errorf("error while fetching the items for search engine: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Id    string
		Label string
	}
	for rows.Next() {
		var result struct {
			Id    string
			Label string
		}
		if err := rows.Scan(&result.Id, &result.Label); err != nil {
			return nil, fmt.Errorf("error while fetching the items for search engine: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}
