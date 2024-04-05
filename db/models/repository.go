package models

import (
	"database/sql"
	"fmt"
)

func getNextAutoIncrementValue(db *sql.DB, tableName string) (int64, error) {
	var nextAutoIncrementValue int64
	query := fmt.Sprintf("SELECT AUTO_INCREMENT FROM information_schema.tables WHERE table_name = '%s' AND table_schema = DATABASE()", tableName)
	row := db.QueryRow(query)
	err := row.Scan(&nextAutoIncrementValue)
	if err != nil {
		return 0, err
	}
	return nextAutoIncrementValue, nil
}
