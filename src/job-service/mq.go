package main

import (
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

func StartRabbitMQConsumer() {
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        log.Printf("⚠️ RabbitMQ Consumer connect error: %v (will retry later)", err)
        return
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Printf("Failed to open channel: %v", err)
        return
    }
    defer ch.Close()

    q, err := ch.QueueDeclare("applications_events", true, false, false, false, nil)
    if err != nil {
        log.Printf("Failed to declare queue: %v", err)
        return
    }

    msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Printf("Failed to register consumer: %v", err)
        return
    }

    log.Println("✅ RabbitMQ Consumer started. Waiting for messages...")
    
    go func() {
        for d := range msgs {
            log.Printf("📥 Job Service received: %s", d.Body)
        }
    }()
}