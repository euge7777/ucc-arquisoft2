package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// PaymentEventConsumer - Consumidor de eventos de pago
type PaymentEventConsumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
}

// PaymentEvent - Estructura del evento de pago (debe coincidir con el publicado por payments-api)
type PaymentEvent struct {
	Action         string                 `json:"action"`          // "payment.completed", "payment.failed", etc.
	Type           string                 `json:"type"`            // "payment"
	PaymentID      string                 `json:"payment_id"`      // ID del pago
	EntityType     string                 `json:"entity_type"`     // "subscription", "inscription"
	EntityID       string                 `json:"entity_id"`       // ID de la suscripción
	UserID         string                 `json:"user_id"`         // ID del usuario
	Amount         float64                `json:"amount"`          // Monto pagado
	Currency       string                 `json:"currency"`        // ARS, USD, etc.
	TransactionID  string                 `json:"transaction_id"`  // ID en gateway
	PaymentGateway string                 `json:"payment_gateway"` // mercadopago, stripe
	Timestamp      time.Time              `json:"timestamp"`       // Fecha del evento
	Metadata       map[string]interface{} `json:"metadata"`        // Datos adicionales
}

// PaymentEventHandler - Interface para manejar eventos de pago
type PaymentEventHandler interface {
	HandlePaymentCompleted(ctx context.Context, event PaymentEvent) error
}

// NewPaymentEventConsumer - Constructor
func NewPaymentEventConsumer(url, exchange, queueName string) (*PaymentEventConsumer, error) {
	// 1. Conectar a RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %w", err)
	}

	// 2. Crear canal
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error creando canal: %w", err)
	}

	// 3. Declarar exchange (debe ser el mismo que usa payments-api)
	err = channel.ExchangeDeclare(
		exchange, // name: "gym.events"
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando exchange: %w", err)
	}

	// 4. Declarar cola
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable (sobrevive a reinicios)
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando cola: %w", err)
	}

	// 5. Bind queue al exchange con routing key para escuchar eventos de pago completado de suscripciones
	// Routing key: "payment.completed.subscription"
	err = channel.QueueBind(
		queue.Name,                      // queue name
		"payment.completed.subscription", // routing key
		exchange,                        // exchange
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("error binding queue: %w", err)
	}

	log.Printf("✅ Consumidor RabbitMQ conectado (Exchange: %s, Queue: %s, RoutingKey: payment.completed.subscription)\n",
		exchange, queueName)

	return &PaymentEventConsumer{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		queue:    queueName,
	}, nil
}

// StartConsuming - Inicia el consumo de eventos
func (c *PaymentEventConsumer) StartConsuming(ctx context.Context, handler PaymentEventHandler) error {
	// Configurar QoS (prefetch)
	err := c.channel.Qos(
		1,     // prefetch count (procesar 1 mensaje a la vez)
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("error configurando QoS: %w", err)
	}

	// Consumir mensajes
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer tag (auto-generado)
		false,   // auto-ack (false para ack manual)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("error iniciando consumo: %w", err)
	}

	log.Printf("📥 Escuchando eventos de pago...")

	// Procesar mensajes en goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 Deteniendo consumidor de eventos...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("⚠️ Canal de mensajes cerrado")
					return
				}

				// Procesar mensaje
				if err := c.processMessage(ctx, msg, handler); err != nil {
					log.Printf("❌ Error procesando mensaje: %v\n", err)
					// NACK (no requeue) para evitar loops infinitos
					msg.Nack(false, false)
				} else {
					// ACK mensaje procesado correctamente
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// processMessage - Procesa un mensaje individual
func (c *PaymentEventConsumer) processMessage(ctx context.Context, msg amqp.Delivery, handler PaymentEventHandler) error {
	// 1. Deserializar evento
	var event PaymentEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("error deserializando evento: %w", err)
	}

	log.Printf("📨 Evento recibido: %s | PaymentID: %s | EntityType: %s | EntityID: %s\n",
		event.Action, event.PaymentID, event.EntityType, event.EntityID)

	// 2. Validar que sea un evento de pago completado para una suscripción
	if event.Action != "payment.completed" {
		log.Printf("⏭️ Ignorando evento (action: %s)\n", event.Action)
		return nil
	}

	if event.EntityType != "subscription" {
		log.Printf("⏭️ Ignorando evento (entity_type: %s)\n", event.EntityType)
		return nil
	}

	// 3. Delegar al handler
	if err := handler.HandlePaymentCompleted(ctx, event); err != nil {
		return fmt.Errorf("error manejando evento: %w", err)
	}

	log.Printf("✅ Evento procesado exitosamente: PaymentID %s activó Subscription %s\n",
		event.PaymentID, event.EntityID)

	return nil
}

// Close - Cierra la conexión
func (c *PaymentEventConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
