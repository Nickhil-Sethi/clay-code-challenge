package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	diff "github.com/sergi/go-diff/diffmatchpatch"
)

func getLastTwo(
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

	if len(returnRows) == 0 {
		return nil, errors.New("Not found")
	}
	return returnRows, nil
}

func sendHTTPError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Something bad happened!"))
}

func diffRequesthandler(w http.ResponseWriter, r *http.Request) {
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		sendHTTPError(w)
		return
	}

	query := r.URL.Query()
	fmt.Println(r.URL, query)
	username := query.Get("username")
	if username == "" {
		fmt.Println("username not present", username)
		sendHTTPError(w)
		return
	}

	mode := query.Get("mode")
	if mode == "" {
		fmt.Println("mode not present")
		sendHTTPError(w)
		return
	}

	// generalize this to be URL parameters
	lastTwo, queryErr := getLastTwo(
		conn, username, mode)

	if queryErr != nil {
		sendHTTPError(w)
		return
	}

	if len(lastTwo) == 0 {
		sendHTTPError(w)
		return
	}

	if len(lastTwo) < 2 {
		fmt.Fprintf(w, lastTwo[0])
		return

	}
	diffEngine := diff.New()
	diffs := diffEngine.DiffMain(lastTwo[1], lastTwo[0], false)
	diffs = diffEngine.DiffCleanupSemantic(diffs)
	html := diffEngine.DiffPrettyHtml(diffs)
	fmt.Fprintf(w, "<body>%s</body>", html)
	return
}
