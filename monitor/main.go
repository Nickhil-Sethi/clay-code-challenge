package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type userChangeEvent struct {
	Username string `json:username`
	Type     string `json:type`
	Created  string `json:created`
	Content  string `json:content`
}

type apiResponse struct {
	Results []userChangeEvent
}

const apiEndpoint = "https://api.clay.earth/api/v1/network/test/twitter/updates"
const pollInterval = time.Second * 30

var (
	host     = os.Getenv("DB_ENDPOINT")
	port     = 5432
	user     = "postgres"
	password = os.Getenv("DB_PASS")
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s sslmode=disable",
	host, port, user, password)

func main() {

	// initialize the database
	_, error := runMigrations()
	if error != nil {
		panic(error)
	}

	// monitor the endpoint
	// in the background
	monitor := endpointMonitor{
		Endpoint:     apiEndpoint,
		PollInterval: pollInterval,
		connString:   psqlInfo,
	}
	go monitor.Run()

	// serve our diff request endpoint
	http.HandleFunc("/", diffRequesthandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
