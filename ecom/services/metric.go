package services

import (
	"log"

	"ecom.com/repository"

	"ecom.com/common"
	"ecom.com/constants"
)

type Metric struct {
	Repo repository.MetricRepositoryI
}

func NewMetricService(repo repository.MetricRepositoryI) *Metric {
	return &Metric{
		Repo: repo,
	}
}

func (m *Metric) GetMetrics() (*common.Metrics, error) {
	totalOrderReceived, err := m.Repo.GetMetricCount()
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	averageProcessingTime, err := m.Repo.GetAverageProcessingTime()
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	ordersPending, err := m.Repo.GetCountByStatus(string(constants.PENDING))
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	ordersProcessing, err := m.Repo.GetCountByStatus(string(constants.PROCESSING))
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	ordersCompleted, err := m.Repo.GetCountByStatus(string(constants.COMPELETED))
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	if err != nil {
		log.Printf("failed to get data from repository %v", err)
	}
	metrics := common.Metrics{
		TotalOrdersReceived:   int64(*totalOrderReceived),
		AverageProcessingTime: *averageProcessingTime,
		OrdersPending:         *ordersPending,
		OrdersProcessing:      *ordersProcessing,
		OrdersCompleted:       *ordersCompleted,
	}

	return &metrics, nil
}
