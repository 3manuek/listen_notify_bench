package shared

// https://github.com/heroku-examples/go-queue-example/blob/master/cmd/queue-example-worker/main.go

import (
	que "github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

// IndexRequest container.
// The URL is the url to index content from.
//type IndexRequest struct {
//	URL string `json:"url"`
//}

const (
	QueueName = "benchmarkQueue"
)

type QueueMessage struct {
	Message string 
}

func Setup(dbURL string) (*pgx.ConnPool, *que.Client, error) {
	pgxpool, err := GetPgxPool(dbURL)
	if err != nil {
		return nil, nil, err
	}

	qc := que.NewClient(pgxpool)

	return pgxpool, qc, err
}

// GetPgxPool based on the provided database URL
func GetPgxPool(dbURL string) (*pgx.ConnPool, error) {
	pgxcfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})

	if err != nil {
		return nil, err
	}

	return pgxpool, nil
}
