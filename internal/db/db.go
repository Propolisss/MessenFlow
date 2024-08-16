package db

import "database/sql"

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}

	createTableSQL := `
  CREATE TABLE IF NOT EXISTS users (
    login TEXT PRIMARY KEY, 
    password TEXT NOT NULL
  );
  `
	createMessagesTableSQL := `
  CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatID TEXT NOT NULL,
    user TEXT NOT NULL,
    message TEXT NOT NULL,
    time TEXT NOT NULL
  );
  `
	_, err = DB.Exec(createTableSQL)
	_, err = DB.Exec(createMessagesTableSQL)
	return err
}
