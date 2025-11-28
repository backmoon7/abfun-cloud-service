package mq

import (
    "context"
    "errors"
    "log"
    "sync"

    amqp "github.com/rabbitmq/amqp091-go"
)

var (
    conn           *amqp.Connection
    publishChannel *amqp.Channel
    initOnce       sync.Once
)

// InitRabbitMQ establishes a single shared RabbitMQ connection/channel for the app.
func InitRabbitMQ(url string) error {
    var initErr error
    initOnce.Do(func() {
        var err error
        conn, err = amqp.Dial(url)
        if err != nil {
            initErr = err
            return
        }

        publishChannel, err = conn.Channel()
        if err != nil {
            initErr = err
            return
        }
    })

    if initErr != nil {
        return initErr
    }

    if conn == nil || publishChannel == nil {
        return errors.New("rabbitmq connection not initialized")
    }

    return nil
}

// Publish sends a message to the given durable queue using the default exchange.
func Publish(queue string, body []byte) error {
    if publishChannel == nil {
        return errors.New("rabbitmq channel not initialized")
    }

    if _, err := publishChannel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
        return err
    }

    return publishChannel.PublishWithContext(
        context.Background(),
        "",
        queue,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
}

// Consume starts a goroutine that passes message bodies to the handler.
func Consume(queue string, handler func([]byte)) error {
    if conn == nil {
        return errors.New("rabbitmq connection not initialized")
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    if _, err = ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
        _ = ch.Close()
        return err
    }

    deliveries, err := ch.Consume(queue, "", true, false, false, false, nil)
    if err != nil {
        _ = ch.Close()
        return err
    }

    go func() {
        for msg := range deliveries {
            func() {
                defer func() {
                    if r := recover(); r != nil {
                        log.Printf("panic in rabbitmq handler: %v", r)
                    }
                }()
                handler(msg.Body)
            }()
        }
        _ = ch.Close()
    }()

    return nil
}
