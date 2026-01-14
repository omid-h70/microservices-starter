package messaging

type QueueConsumer struct {
	rb          *RabbitMQ
	connManager *ConnManager
	queueName   string
}

func NewQueueConsumer(rb *RabbitMQ, connManager *ConnManager, qName string) *QueueConsumer {
	return &QueueConsumer{
		rb:          rb,
		connManager: connManager,
		queueName:   qName,
	}
}

func (qc *QueueConsumer) Start() error {
	return nil
}
