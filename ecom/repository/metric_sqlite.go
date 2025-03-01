package repository

import (
	"database/sql"

	"ecom.com/models"
)

type SQLiteMetricRepository struct {
	DB *sql.DB
}

func NewSQLiteMetricRepository(db *sql.DB) MetricRepositoryI {
	return &SQLiteMetricRepository{DB: db}
}

func (r *SQLiteMetricRepository) CreateMetric(m *models.Metric) error {
	query := `INSERT INTO metrics (order_id, duration, metric_name) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, m.OrderId, m.Duration, m.MetricName)
	return err
}

func (r *SQLiteMetricRepository) GetMetricByID(id int, name string) (*models.Metric, error) {
	query := `SELECT order_id, duration FROM metrics WHERE order_id = ? AND metric_name = ?`
	row := r.DB.QueryRow(query, id, name)

	var metric models.Metric
	err := row.Scan(&metric.OrderId, &metric.Duration)
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

func (r *SQLiteMetricRepository) GetAverageTime(metricName string) (*float64, error) {
	var AverageDuration float64
	err := r.DB.QueryRow("SELECT COALESCE(AVG(duration), 0) FROM metrics WHERE metric_name = ?", metricName).Scan(&AverageDuration)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &AverageDuration, nil
}

func (r *SQLiteMetricRepository) GetCountByStatus(status string) (*int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = ?", string(status)).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &count, nil
}
