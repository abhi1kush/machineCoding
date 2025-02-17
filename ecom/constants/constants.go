package constants

type OrderStates string

const (
	PENDING    OrderStates = "Pending"
	PROCESSING OrderStates = "Processing"
	COMPELETED OrderStates = "Completed"
)
