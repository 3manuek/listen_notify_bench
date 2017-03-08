package main

import (
	"net/http"
	"os"

	qebench "github.com/3manuek/listen_notify_bench"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	var err error
	pgxpool, qc, err = qebench.Setup(dbURL)
	if err != nil {
		log.WithField("DATABASE_URL", dbURL).Fatal("Unable to setup queue / database")
	}
	defer pgxpool.Close()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/index", handleIndexRequest)
	log.Println(http.ListenAndServe(":"+port, nil))
}
