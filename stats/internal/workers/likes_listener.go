package workers

import (
	"encoding/json"
	"log"

	"sn/libraries/clickhouse"
	"sn/libraries/kafka"

	"github.com/IBM/sarama"
)

type LikesListener struct{}

func (k *LikesListener) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (k *LikesListener) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (k *LikesListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ch := clickhouse.GetClickhouseConnection()
	for {
		select {
		case msg := <-claim.Messages():
			log.Printf("[likes] Received message: %s", msg.Value)
			body := make(map[string]any)
			json.Unmarshal(msg.Value, &body)
			_, err := ch.Exec("INSERT INTO stats.likes (user_id, post_id) VALUES (?, ?)", body["user_id"], body["post_id"])
			if err != nil {
				log.Printf("Failed to insert into clickhouse: %v", err)
				return err
			}
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func StartLikesListener() {
	log.Printf("Starting likes listener")
	ch := clickhouse.GetClickhouseConnection()

	query := `CREATE TABLE IF NOT EXISTS stats.likes
(
    user_id UUID,
    post_id UUID,
    created_at DateTime DEFAULT CAST(now(), 'DateTime')
) ENGINE = MergeTree()
ORDER BY (post_id, user_id);`
	ch.Exec(query)

	log.Printf("Table created")
	consumer, err := kafka.SetupConsumer("posts_likes", new(LikesListener))
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error from likes consumer: %v", err)
		}
	}()
}
