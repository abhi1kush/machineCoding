package repository

import (
	"database/sql"

	"ecom.com/models"
)

type SQLiteMetricRepository struct {
	DB *sql.DB
}

func NewSQLiteMetricRepository(db *sql.DB) *SQLiteMetricRepository {
	return &SQLiteMetricRepository{DB: db}
}

func (r *SQLiteMetricRepository) CreateMetric(m *models.Metric) error {
	query := `INSERT INTO metrics (order_id, processing_time) VALUES (?, ?)`
	_, err := r.DB.Exec(query, m.OrderId, m.ProcessingTime)
	return err
}

func (r *SQLiteMetricRepository) GetMetricByID(id int) (*models.Metric, error) {
	query := `SELECT order_id, processing_time FROM metrics WHERE order_id = ?`
	row := r.DB.QueryRow(query, id)

	var metric models.Metric
	err := row.Scan(&metric.OrderId, &metric.ProcessingTime)
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *SQLiteMetricRepository) GetMetricCount() (*int, error) {
	var TotalOrdersReceived int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&TotalOrdersReceived)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &TotalOrdersReceived, nil
}

func (r *SQLiteMetricRepository) GetAverageProcessingTime() (*float64, error) {
	var AverageProcessingTime float64
	err := r.DB.QueryRow("SELECT COALESCE(AVG(processing_time), 0) FROM metrics").Scan(&AverageProcessingTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &AverageProcessingTime, nil
}

func (r *SQLiteMetricRepository) GetCountByStatus(status string) (*int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = ?", string(status)).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &count, nil
}
