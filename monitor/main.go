package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func writeChangeEvents(events []userChangeEvent) {
	fmt.Println(events)
}

func hitAPI() ([]userChangeEvent, error) {
	resp, err := http.Get(apiEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body = new(apiResponse)
	decodeError := json.NewDecoder(
		resp.Body).Decode(&body)

	if decodeError != nil {
		return nil, decodeError
	}

	results := body.Results
	return results, nil
}

func monitorEndpoint() {
	for {
		results, err := hitAPI()
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(2 * pollInterval)
			continue
		}
		writeChangeEvents(results)
		time.Sleep(pollInterval)
	}
}

func main() {
	runMigrations()
	monitorEndpoint()

	diff.New()

}
