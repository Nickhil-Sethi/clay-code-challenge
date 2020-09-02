package main

import (
	"fmt"
	"time"

	diff "github.com/sergi/go-diff/diffmatchpatch"
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

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		host, port, user, password)

	monitor := endpointMonitor{
		Endpoint:     apiEndpoint,
		PollInterval: pollInterval,
		connString:   psqlInfo,
	}

	monitor.Run()
	diff.New()
}
