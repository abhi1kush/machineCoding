package constants

type OrderStates string

const (
	PENDING    OrderStates = "Pending"
	PROCESSING OrderStates = "Processing"
	COMPELETED OrderStates = "Completed"
)

type MetricName string

const (
	PROCESSING_TIME MetricName = "processing_time"
	CREATION_TIME   MetricName = "creation_time"
)
