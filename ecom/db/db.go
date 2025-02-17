package db

import (
	"database/sql"
	"log"
)

var (
	// Orders DB holds orders data.
	DB *sql.DB
	// MetricsDB holds metrics data.
	MetricsDB *sql.DB
)

// InitDB initializes the orders database.
func InitDB(driver, dsn string) {
	var err error
	DB, err = sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("Error opening orders database: %v", err)
	}

	ordersQuery := `CREATE TABLE IF NOT EXISTS orders (
		order_id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		item_ids TEXT NOT NULL,
		total_amount DECIMAL(10,2) NOT NULL,
		status TEXT CHECK (status IN ('Pending', 'Processing', 'Completed')) NOT NULL
	);`
	_, err = DB.Exec(ordersQuery)
	if err != nil {
		log.Fatalf("Error creating orders table: %v", err)
	}

	log.Println("Orders database initialized successfully")
}

// InitMetricsDB initializes the metrics database.
func InitMetricsDB(driver, dsn string) {
	var err error
	MetricsDB, err = sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("Error opening metrics database: %v", err)
	}

	metricsQuery := `CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id TEXT,
		processing_time REAL
	);`
	_, err = MetricsDB.Exec(metricsQuery)
	if err != nil {
		log.Fatalf("Error creating metrics table: %v", err)
	}

	log.Println("Metrics database initialized successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func CloseMetricsDB() {
	if MetricsDB != nil {
		MetricsDB.Close()
	}
}
