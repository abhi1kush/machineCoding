package repository

import (
	"database/sql"

	"ecom.com/models"
)

type PostgeSqlMetricRepository struct {
	DB *sql.DB
}

func NewPostgeSqlMetricRepository(db *sql.DB) MetricRepositoryI {
	return &PostgeSqlMetricRepository{DB: db}
}

func (r *PostgeSqlMetricRepository) CreateMetric(m *models.Metric) error {
	query := `INSERT INTO metrics (order_id, duration, metric_name) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, m.OrderId, m.Duration, m.MetricName)
	return err
}

func (r *PostgeSqlMetricRepository) GetMetricByID(id int, name string) (*models.Metric, error) {
	query := `SELECT order_id, duration FROM metrics WHERE order_id = $1 AND metric_name = $2`
	row := r.DB.QueryRow(query, id, name)

	var metric models.Metric
	err := row.Scan(&metric.OrderId, &metric.Duration)
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *PostgeSqlMetricRepository) GetMetricCount() (*int, error) {
	var TotalOrdersReceived int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&TotalOrdersReceived)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &TotalOrdersReceived, nil
}

func (r *PostgeSqlMetricRepository) GetAverageTime(metricName string) (*float64, error) {
	var AverageDuration float64
	err := r.DB.QueryRow("SELECT COALESCE(AVG(duration), 0) FROM metrics WHERE metric_name = $1", metricName).Scan(&AverageDuration)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &AverageDuration, nil
}

func (r *PostgeSqlMetricRepository) GetCountByStatus(status string) (*int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = $1", string(status)).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &count, nil
}
