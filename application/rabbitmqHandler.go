package application

import (
	"context"
	"encoding/json"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/dulguundd/logError-lib/logger"
	"github.com/streadway/amqp"
	"log"
)

type RabbitmqVar struct {
	ch                *amqp.Channel
	conn              *amqp.Connection
	rabbitmqCloseFunc func()
}

func connectRabbitmq() RabbitmqVar {
	conn, err := amqp.Dial("amqp://guest:guest@172.30.52.239:5672/")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel")
	}

	return RabbitmqVar{ch, conn, func() {
		defer conn.Close()
		defer ch.Close()
	}}
}

func (h Handlers) rpcCreateInstance() {
	r := connectRabbitmq()
	q, err := r.ch.QueueDeclare(
		"rpcCreateInstanceQueue", // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		logger.Error("Failed to declare a queue")
	}
	err = r.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logger.Error("Failed to set QoS")
	}
	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Error("Failed to register a consumer")
	}
	_ = make(chan bool)

	go func() {
		for d := range msgs {
			var req DeployInstanceRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				logger.Error("Failed to convert body to request body")
			} else {
				zbClient, err := zbc.NewClient(h.zeebeClient)
				if err != nil {
					logger.Error("Cant connect zeebe")
				} else {
					variables := make(map[string]interface{})
					for variableCount, variablesValue := range req.Variables {
						variables[req.Variables[variableCount].Name] = variablesValue.Value
					}
					ctx := context.Background()
					request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(req.BpmnProcessId).LatestVersion().VariablesFromMap(variables)
					_, err = request.Send(ctx)
					if err != nil {
						logger.Error("Deploy error")
					} else {
						logger.Info("Create Instance successful")
					}
				}
			}

			//err = r.ch.Publish(
			//	"",        // exchange
			//	d.ReplyTo, // routing key
			//	false,     // mandatory
			//	false,     // immediate
			//	amqp.Publishing{
			//		ContentType:   "text/plain",
			//		CorrelationId: d.CorrelationId,
			//		Body:          []byte(strconv.Itoa(response)),
			//	})

			d.Ack(false)
		}
	}()
	log.Printf(" [*] Awaiting RPC requests")
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
