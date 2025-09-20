package main

import (
	"database/sql"
)

func initDb(dbFile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFile)
	chk(err)
	return db
}

func saveResult(db *sql.DB, result *Result) {
	_, err := db.Exec(`UPDATE task SET status = ?, headers = ?, content = ?, err = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, result.Status, result.Headers, result.Content, result.Err, result.Id)
	chk(err)
}

func loadTasks(db *sql.DB) []Task {
	rows, err := db.Query(`SELECT id, url FROM task where status is null and err is null`)
	chk(err)
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Url)
		chk(err)
		tasks = append(tasks, task)
	}
	return tasks
}

func saveUrls(db *sql.DB, urls []string) {
	for _, url := range urls {
		_, err := db.Exec(`INSERT INTO task (url, created_at) VALUES (?, CURRENT_TIMESTAMP) ON CONFLICT DO NOTHING`, url)
		chk(err)
	}
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE task (
id INTEGER PRIMARY KEY AUTOINCREMENT,
url TEXT UNIQUE,
status INTEGER,
headers TEXT,
content BLOB,
err TEXT,
created_at TIMESTAMP,
updated_at TIMESTAMP)`)
	chk(err)
}
