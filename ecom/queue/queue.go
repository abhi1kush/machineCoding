package queue

type QueueI interface {
	StartOrderProcessor() error
	StopOrderProcessor()
	Enqueue(item Item)
}
