module pizza

go 1.18

require common v0.0.0-00010101000000-000000000000

require github.com/rabbitmq/amqp091-go v1.8.1 // indirect

replace common => ../common
