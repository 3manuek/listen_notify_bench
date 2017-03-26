package main

/*
	This is the queuer. We generate random data and push it to a central queue
	to be processed by the worker.

	Worker have 2 channels, both LISTENing. Enqueue action NOTIFies.

	This is intended for benchmarking CPU usage on Postgres versions. Not for
	benchmarking the que-go nor neither any of the involved drivers.  

*/

import (
	// "net/http"
	"os"
	"pgx"

	"github.com/Sirupsen/logrus"
	que "github.com/bgentry/que-go"
	// qebench "github.com/3manuek/listen_notify_bench"
	quebench "../main_lib.go"
)

var (
	// log     = logrus.WithField("cmd", "queue-example-web")
	qc      *que.Client
	pgxpool *pgx.ConnPool
)

func queueRequest(ir quebench.QueueMessage) error {
	//enc, err := json.Marshal(ir)
	//if err != nil {
	//	return errors.Wrap(err, "Marshalling the IndexRequest")
	//}

	j := que.Job{
		Type: quebench.QueueName,
		Args: ir,
	}

	return errors.Wrap(qc.Enqueue(&j), "Enqueueing Job")
}

func main() {
	// postgres://postgres@localhost/queue
	// default='postgres://user:pass@localhost/dbname'
	dbURL := os.Getenv("DATABASE_URL")
	var err error

	pgxpool, qc, err = qebench.Setup(dbURL)

	if err != nil {
		//log.WithField("DATABASE_URL", dbURL).Fatal("Unable to setup queue / database")

	}
	defer pgxpool.Close()

	// http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/index", handleIndexRequest)
	// log.Println(http.ListenAndServe(":"+port, nil))


	if err := queueRequest(ir); err != nil {
		l.Println(err.Error())
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
