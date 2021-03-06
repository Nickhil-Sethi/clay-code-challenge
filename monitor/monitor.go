package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type endpointMonitor struct {
	Endpoint     string
	PollInterval time.Duration
	connString   string
}

func (m *endpointMonitor) parseTimestamp(timestring string) (time.Time, error) {
	layout := "2006-01-02 15:04:05.000000+00:00"
	return time.Parse(layout, timestring)
}

func (m *endpointMonitor) writeChangeEvents(events []userChangeEvent) {
	conn, err := sql.Open(
		"postgres",
		m.connString)

	if err != nil {
		log.Println(err)
	}

	defer conn.Close()

	for _, evt := range events {
		timestamp, _ := m.parseTimestamp(evt.Created)
		_, err := conn.Exec(
			`INSERT INTO user_change_events (
				username,
				created,
				content,
				type
			) VALUES (
				$1,
				$2,
				$3,
				$4
			) ON CONFLICT DO NOTHING;`,
			evt.Username, timestamp, evt.Content, evt.Type)

		if err != nil {
			log.Println(err)
		}
	}
}

func (m *endpointMonitor) hitEndpoint() ([]userChangeEvent, error) {
	resp, err := http.Get(m.Endpoint)
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

func (m *endpointMonitor) Run() {
	for {
		results, err := m.hitEndpoint()
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(2 * pollInterval)
			continue
		}
		m.writeChangeEvents(results)
		time.Sleep(pollInterval)
	}
}
