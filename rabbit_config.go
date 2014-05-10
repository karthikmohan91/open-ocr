package ocrworker

import (
	"flag"
)

type RabbitConfig struct {
	AmqpURI            string
	Exchange           string
	ExchangeType       string
	RoutingKey         string
	CallbackRoutingKey string
	Reliable           bool
	QueueName          string
	CallbackQueueName  string
}

func DefaultTestConfig() RabbitConfig {

	// Reliable: false due to major issues that would completely
	// wedge the rpc worker.  Setting the buffered channels length
	// higher would delay the problem, but then it would still happen later.

	rabbitConfig := RabbitConfig{
		AmqpURI:            "amqp://guest:guest@localhost:5672/",
		Exchange:           "test-exchange",
		ExchangeType:       "direct",
		RoutingKey:         "test-key",
		CallbackRoutingKey: "callback-key",
		Reliable:           false, // setting to false because of observed issues
		QueueName:          "test-queue",
		CallbackQueueName:  "callback-queue",
	}
	return rabbitConfig

}

func DefaultConfigFlagsOverride() RabbitConfig {
	rabbitConfig := DefaultTestConfig()

	var AmqpURI string
	flag.StringVar(
		&AmqpURI,
		"amqp_uri",
		"",
		"The Amqp URI, eg: amqp://guest:guest@localhost:5672/",
	)

	flag.Parse()
	if len(AmqpURI) > 0 {
		rabbitConfig.AmqpURI = AmqpURI
	}

	return rabbitConfig

}
