package repository

import (
	"ecom.com/models"
)

type MetricRepositoryI interface {
	CreateMetric(metric *models.Metric) error
	GetMetricByID(id int, name string) (*models.Metric, error)
	GetMetricCount() (*int, error)
	GetAverageTime(metricname string) (*float64, error)
}
