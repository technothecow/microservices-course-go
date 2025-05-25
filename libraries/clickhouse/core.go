package clickhouse

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

var (
	ch *sql.DB
	once sync.Once
)

func getClickhouseConnection() {
	var err error

	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	dbname := os.Getenv("CLICKHOUSE_DB")

	for range(3) {
		log.Printf("Trying to connect to Clickhouse")
		ch, err = sql.Open("clickhouse", fmt.Sprintf("http://default:default@%s:%s/%s?debug=true", host, port, dbname))
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	if err = ch.Ping(); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to Clickhouse")
}

var GetClickhouseConnection = func() *sql.DB {
	once.Do(getClickhouseConnection)
	return ch
}