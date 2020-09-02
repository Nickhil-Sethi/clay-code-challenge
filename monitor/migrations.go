package main

import "database/sql"

const migrationStr = `
	CREATE TABLE IF NOT EXISTS user_change_events (
		username TEXT,
		created TIMESTAMP,
		content TEXT,
		type TEXT,
		PRIMARY KEY (username, created)
);`

func runMigrations() (sql.Result, error) {
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Exec(migrationStr)
}
