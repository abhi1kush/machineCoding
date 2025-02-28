package repository

import (
	"ecom.com/models"
)

type MetricRepositoryI interface {
	CreateMetric(metric *models.Metric) error
	GetMetricByID(id int) (*models.Metric, error)
	GetMetricCount() (*int, error)
	GetAverageProcessingTime() (*float64, error)
	GetCountByStatus(status string) (*int, error)
}
