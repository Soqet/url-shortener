package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDb(filePath string) (db *DB, err error) {
	db = new(DB)
	db.db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS urls ( 
			url TEXT,
			id TEXT UNIQUE,
			clicks NUMBER,
			lastTimeClicked TEXT
		)`,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) AddLink(url string, short string) error {
	_, err := db.db.Exec("REPLACE INTO urls ( url, id ) VALUES ( $1, $2 )", url, short)
	return err
}

func (db *DB) GetInitialLink(id string) (string, error) {
	rows, err := db.db.Query("SELECT url FROM urls WHERE id = $1", id)
	if err != nil {
		return "", nil
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		rows.Scan(&url)
		return url, nil
	}
	return "", nil
}

func (db *DB) UpdateLink(id string, clicks int, currentTime string) error {
	_, err := db.db.Exec("UPDATE urls SET clicks = clicks + $2, lastTimeClicked = $3 WHERE id = $1", id, clicks, currentTime)
	return err
}
