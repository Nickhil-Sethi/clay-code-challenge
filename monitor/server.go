package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/sergi/go-diff/diffmatchpatch"
	diff "github.com/sergi/go-diff/diffmatchpatch"
)

const errorHTML = "<body>Oops! Something's wrong. Try again later.</body>"
const badRequestHTML = "<body>Bad request. Check that you have both username and mode parameters</body>"
const notFoundHTML = "<body>No events found for that user and mode.</body>"

func getLastTwoEvents(
	conn *sql.DB,
	username string,
	mode string) ([]string, error) {

	rows, queryErr := conn.Query(`
		SELECT content 
		FROM user_change_events 
		WHERE 
			username = $1 AND 
			type = $2
		ORDER BY 
			CREATED DESC 
		LIMIT 2`,
		username, mode)

	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var returnRows []string
	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		returnRows = append(returnRows, c)
	}
	return returnRows, nil
}

func sendHTTPResponse(
	w http.ResponseWriter, status int, body template.HTML) {
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func getDiffs(before string, after string) []diffmatchpatch.Diff {
	diffEngine := diff.New()

	// compute the diff
	diffs := diffEngine.DiffMain(
		before, after, false)

	// clean up the diff
	diffs = diffEngine.DiffCleanupSemantic(
		diffs)

	return diffs
}

func getDiffHTML(before string, after string) string {
	diffEngine := diff.New()

	// compute the diff
	diffs := diffEngine.DiffMain(
		before, after, false)

	// clean up the diff
	diffs = diffEngine.DiffCleanupSemantic(
		diffs)

	// Stock function returns pretty HTML
	html := diffEngine.DiffPrettyHtml(diffs)
	return fmt.Sprintf("<body>%s</body>", html)
}

func validateQuery(query url.Values) bool {
	username := query.Get("username")
	if username == "" {
		return false
	}

	mode := query.Get("mode")
	if mode == "" {
		return false
	}
	return true
}

func diffRequestHandlerJSON(w http.ResponseWriter, r *http.Request) {
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		sendHTTPResponse(
			w,
			http.StatusInternalServerError,
			errorHTML)
		return
	}

	query := r.URL.Query()
	if !validateQuery(query) {
		sendHTTPResponse(
			w,
			http.StatusBadRequest,
			badRequestHTML)
		return
	}

	// extract query parameters
	username := query.Get("username")
	mode := query.Get("mode")

	// get the last two events of the
	// mode requested for this user
	events, queryErr := getLastTwoEvents(
		conn, username, mode)

	if queryErr != nil {
		sendHTTPResponse(
			w,
			http.StatusInternalServerError,
			errorHTML)
		return
	}

	// compute json
	// write json
	diffs := getDiffs(events[1], events[0])
	bytes, err := json.Marshal(diffs)
	if err != nil {
		panic(err)
	}
	// set Header content-type: application/json
	w.Header().Set(
		"Content-Type",
		"application/json")
	fmt.Fprintf(w, string(bytes))

}

func diffRequesthandler(w http.ResponseWriter, r *http.Request) {
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		sendHTTPResponse(
			w,
			http.StatusInternalServerError,
			errorHTML)
		return
	}

	query := r.URL.Query()
	if !validateQuery(query) {
		sendHTTPResponse(
			w,
			http.StatusBadRequest,
			badRequestHTML)
		return
	}

	// extract query parameters
	username := query.Get("username")
	mode := query.Get("mode")

	// get the last two events of the
	// mode requested for this user
	events, queryErr := getLastTwoEvents(
		conn, username, mode)

	if queryErr != nil {
		sendHTTPResponse(
			w,
			http.StatusInternalServerError,
			errorHTML)
		return
	}

	switch len(events) {
	case 0:
		sendHTTPResponse(
			w,
			http.StatusNotFound,
			notFoundHTML)
	case 1:
		// pretty print HTML
		html := getDiffHTML(events[1], "")
		fmt.Fprintf(w, html)
	case 2:
		html := getDiffHTML(events[1], events[0])
		fmt.Fprintf(w, html)
	}

	return
}
