package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func ConnectDB(driver string, dsn string) *sql.DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create table if not exists
	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );`

	if _, err = db.Exec(query); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Create orders if not exists
	ordersQuery := `CREATE TABLE IF NOT EXISTS orders (
		order_id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		total_amount DECIMAL(10,2) NOT NULL,
		status TEXT CHECK (status IN ('Pending', 'Processing', 'Completed')) NOT NULL
	);`
	_, err = db.Exec(ordersQuery)
	if err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}

	// Create items if not exists
	itemQuery := `CREATE TABLE IF NOT EXISTS items (
		item_id TEXT,
		amount DECIMAL(10,2) NOT NULL,
		order_id TEXT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE ON UPDATE CASCADE
		PRIMARY KEY (item_id, order_id)
	);`
	_, err = db.Exec(itemQuery)
	if err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}

	return db
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func ConnectMetricsDB(driver string, dsn string) *sql.DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	metricsQuery := `CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		order_id TEXT NOT NULL,
		metric_name TEXT NOT NULL,
		duration REAL
	);`
	_, err = db.Exec(metricsQuery)
	if err != nil {
		log.Fatalf("Error creating metrics table: %v", err)
	}
	return db
}
