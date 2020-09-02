package main

import (
	"fmt"
	"log"
	"net/http"
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
const (
	host     = ""
	port     = 5432
	user     = "postgres"
	password = ""
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s sslmode=disable",
	host, port, user, password)

func main() {
	_, error := runMigrations()
	if error != nil {
		panic(error)
	}

	monitor := endpointMonitor{
		Endpoint:     apiEndpoint,
		PollInterval: pollInterval,
		connString:   psqlInfo,
	}
	// monitor the endpoint
	// in the background
	go monitor.Run()

	http.HandleFunc("/", diffRequesthandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
