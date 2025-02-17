package services

import (
	"database/sql"

	"ecom.com/constants"
	"ecom.com/db"
	"ecom.com/models"
)

func GetMetrics() (*models.Metrics, error) {
	var metrics models.Metrics
	var err error

	err = db.MetricsDB.QueryRow("SELECT COUNT(*) FROM metrics").Scan(&metrics.TotalOrdersReceived)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = db.MetricsDB.QueryRow("SELECT COALESCE(AVG(processing_time), 0) FROM metrics").Scan(&metrics.AverageProcessingTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	err = db.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = ?", string(constants.PENDING)).Scan(&metrics.OrdersPending)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	err = db.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = ?", string(constants.PROCESSING)).Scan(&metrics.OrdersProcessing)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	err = db.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE status = ?", string(constants.COMPELETED)).Scan(&metrics.OrdersCompleted)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &metrics, nil
}
