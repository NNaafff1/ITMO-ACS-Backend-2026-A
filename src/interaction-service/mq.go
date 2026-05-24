package main

import (
    "encoding/json"
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitConn *amqp.Connection

func InitRabbitMQProducer() {
    var err error
    rabbitConn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    log.Println("✅ RabbitMQ Producer connected")
}

func PublishApplicationEvent(vacancyID, resumeID int) {
    if rabbitConn == nil || rabbitConn.IsClosed() {
        log.Println("⚠️ RabbitMQ connection closed, skipping publish")
        return
    }

    ch, err := rabbitConn.Channel()
    if err != nil {
        log.Printf("Failed to open channel: %v", err)
        return
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "applications_events", // имя очереди
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        nil,
    )
    if err != nil {
        log.Printf("Failed to declare queue: %v", err)
        return
    }

    msg := map[string]interface{}{
        "vacancy_id": vacancyID,
        "resume_id":  resumeID,
        "status":     "pending",
    }
    body, _ := json.Marshal(msg)

    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        log.Printf("Failed to publish message: %v", err)
        return
    }
    log.Printf("📤 Sent to RabbitMQ: %s", body)
}