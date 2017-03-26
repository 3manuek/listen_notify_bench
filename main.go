package main

/*
	Simple Postgres async stresser 

	It's built to compare 2 instances in different versions.

	Inspired in:
	
	Probably this is the one I should follow => https://godoc.org/github.com/lib/pq/listen_example
	https://github.com/jackc/pgx/blob/master/examples/chat/main.go


	N listeners per instance (or notifiers/2 as a rule of the thumb)
	N notifiers per instance

*/

import (
	//"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	pgx   "github.com/jackc/pgx"
//	toml  "github.com/BurntSushi/toml"
)

var(
  pool1 *pgx.ConnPool
  pool2 *pgx.ConnPool
  dbURL1 string
  dbURL2 string 
)

type ConfFile struct {
	connectionURI map[string]url
}

type url struct{
	url 	string 
}

func main() {
	var err error
	// var conf ConfFile

/* In the future:
	if _, err := toml.Decode(tomlData, &conf); err != nil {
	// handle error
	}

	for i, tt := range conf.connectionURI {
			connParams, err := pgx.ParseURI(tt.url)
			if err != nil {
				//t.Errorf("%d. Unexpected error from pgx.ParseURL(%q) => %v", i, tt.url, err)
				continue
			}
			// map[] = tt.url
		}
*/

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(2)
	}()


	dbURL1 = os.Getenv("DATABASE_URL1")
	dbURL2 = os.Getenv("DATABASE_URL2")

	fmt.Fprintln(os.Stderr,"Print conn 1", dbURL1)

	connParams1, err := pgx.ParseURI(dbURL1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "The URL1 is unparseable:", err)
		//continue
		os.Exit(4)
	}
	connParams2, err := pgx.ParseURI(dbURL2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "The URL2 is unparseable:", err)
		//continue
		os.Exit(4)
	}
	fmt.Fprintln(os.Stderr,"Print conn 1", connParams1)
	fmt.Fprintln(os.Stderr,"Print conn 2", connParams1)	

	connPoolParams1 := pgx.ConnPoolConfig{
		ConnConfig: connParams1,
		MaxConnections: 4,
		//AfterConnect: afterConnect,
	}
	connPoolParams2 := pgx.ConnPoolConfig{
		ConnConfig: connParams2,
		MaxConnections: 4,
		//AfterConnect: afterConnect,
	}
	

	pool1, err = pgx.NewConnPool(connPoolParams1)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Print conn 1", connPoolParams1)
		fmt.Fprintln(os.Stderr, "Unable to connect to database 1:", err)
		os.Exit(1)
	}

	pool2, err = pgx.NewConnPool(connPoolParams2)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Print conn 2", connPoolParams2)
		fmt.Fprintln(os.Stderr, "Unable to connect to database 2:", err)
		os.Exit(1)
	}

	defer pool2.Close()
	defer pool1.Close()


	fmt.Println(`Starting the tests.
Please sit tight. 
`)

	// range url
	go listen(pool1,"olakase")
	go listen(pool2,"olakase")

	// range notify
	for {
		go notify(pool1,"olakase")
		go notify(pool2,"olakase")
		// notify
	}
}

/*

*/
func generateRandomMsg() string {
	return "a"
}


/*
Receives: pointer to connection,
		  channel,
		  payload
*/
func notify(pool *pgx.ConnPool, channel string) {
	var msg string 
	var notifyQuery string
	var err error  
	msg = generateRandomMsg()
	//notifyQuery = fmt.Sprintf("select pg_notify('",channel, "', '", msg,"')")
	notifyQuery = fmt.Sprintf("notify %s, %s",channel, msg)
	//fmt.Fprintln(os.Stderr, "NOtify query" , notifyQuery)
	// notifyQuery = "select pg_notify('chat', $1)"
	_, err = pool.Exec(notifyQuery, msg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error sending notification:", err)
			os.Exit(1)
		}
}


/*
Receives: 	pointer to connection,
			channel
*/
func listen(pool *pgx.ConnPool, channel string ) {
	var err error 
	conn, err := pool.Acquire()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		os.Exit(1)
	}
	defer pool.Release(conn)

	conn.Listen(channel)
	defer conn.Unlisten(channel)

	for {
		notification, err := conn.WaitForNotification(time.Second)
		if err == pgx.ErrNotificationTimeout {
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
			os.Exit(1)
		}
		// counter process and back to the channel 
		fmt.Println("PID:", notification.Pid, "Channel:", notification.Channel, "Payload:", notification.Payload)
	}
}

