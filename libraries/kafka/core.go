package kafka

import (
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

var (
	producer  sarama.SyncProducer
	once      sync.Once
	retryOnce sync.Once
	// Channel for handling failed messages that need retry
	retryChannel      = make(chan *messageRetry, 100)
	isTesting    bool = false
)

// SetTestingMode sets the testing mode for the Kafka producer
func SetTestingMode(testing bool) {
	isTesting = testing
}

type messageRetry struct {
	message     *sarama.ProducerMessage
	shouldRetry bool
	attempts    int
}

// GetProducer returns a singleton Kafka producer and initializes the retry processor
func GetProducer() (sarama.SyncProducer, error) {
	var err error
	once.Do(func() {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		producer, err = sarama.NewSyncProducer([]string{"kafka:29092"}, config)

		// Initialize the retry processor only once when producer is created
		retryOnce.Do(func() {
			go retryProcessor()
		})
	})
	return producer, err
}

// CloseProducer closes the Kafka producer connection
func CloseProducer() error {
	if producer != nil {
		return producer.Close()
	}
	return nil
}

// SendMessageSync sends a message to Kafka synchronously
// Returns the partition and offset where the message was stored
// If shouldRetry is true and sending fails, it will queue the message for retries
func SendMessageSync(topic string, key, value []byte, shouldRetry bool) (partition int32, offset int64, err error) {
	if isTesting {
		return 0, 0, nil
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	prod, err := GetProducer()
	if err != nil {
		log.Printf("Failed to get Kafka producer: %v", err)
		if shouldRetry {
			retryChannel <- &messageRetry{
				message:     msg,
				shouldRetry: true,
				attempts:    1,
			}
		}
		return 0, 0, err
	}

	partition, offset, err = prod.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		if shouldRetry {
			retryChannel <- &messageRetry{
				message:     msg,
				shouldRetry: true,
				attempts:    1,
			}
		}
		return 0, 0, err
	}

	return partition, offset, nil
}

// SendMessageAsync sends a message to Kafka asynchronously
// If shouldRetry is true and sending fails, it will queue the message for retries
func SendMessageAsync(topic string, key, value []byte, shouldRetry bool) {
	if isTesting {
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	go func() {
		prod, err := GetProducer()
		if err != nil {
			log.Printf("Failed to get Kafka producer: %v", err)
			if shouldRetry {
				retryChannel <- &messageRetry{
					message:     msg,
					shouldRetry: true,
					attempts:    1,
				}
			}
			return
		}

		partition, offset, err := prod.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			if shouldRetry {
				retryChannel <- &messageRetry{
					message:     msg,
					shouldRetry: true,
					attempts:    1,
				}
			}
			return
		}

		log.Printf("Message sent successfully to partition %d at offset %d", partition, offset)
	}()
}

// retryProcessor handles message retries in a separate goroutine
func retryProcessor() {
	log.Println("Starting Kafka message retry processor")

	for retry := range retryChannel {
		if retry.shouldRetry && retry.attempts <= 3 {
			// Exponential backoff: 500ms, 1s, 2s
			backoff := time.Duration(500*(1<<(retry.attempts-1))) * time.Millisecond
			time.Sleep(backoff)

			go func(r *messageRetry) {
				log.Printf("Retrying message, attempt %d", r.attempts)

				prod, err := GetProducer()
				if err != nil {
					log.Printf("Retry failed to get Kafka producer: %v", err)
					if r.attempts < 3 {
						r.attempts++
						retryChannel <- r
					} else {
						log.Printf("Failed to send message after 3 attempts")
					}
					return
				}

				partition, offset, err := prod.SendMessage(r.message)
				if err != nil {
					log.Printf("Retry failed to send message: %v", err)
					if r.attempts < 3 {
						r.attempts++
						retryChannel <- r
					} else {
						log.Printf("Failed to send message after 3 attempts")
					}
					return
				}

				log.Printf("Message sent successfully on retry %d to partition %d at offset %d",
					r.attempts, partition, offset)
			}(retry)
		}
	}
}
