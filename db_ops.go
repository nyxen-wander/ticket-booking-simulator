package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func printBookSummary(db *sql.DB, currentTime string) error {

	trx, err := db.Begin()

	if err != nil {
		return err
	}

	// Query the database for all bookings made since the program started.
	bufferRow, queryErr := db.Query("SELECT m.movie_title, s.theater_id, b.customer_name, b.seat_count, to_char(b.created_at, 'YYYY-MM-DD HH24:MI:SS') FROM bookings b JOIN sessions s ON b.session_id = s.id JOIN movies m ON m.id = s.movie_id WHERE b.created_at BETWEEN $1 AND now();", currentTime)

	if queryErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return queryErr
	}

	if commitErr := trx.Commit(); commitErr != nil {
		return commitErr
	}

	var movieTitle, theatherId, customerName, seatCount, createdAt, summaryDetail string

	// Iterate through the results to build and print a formatted summary string.
	for bufferRow.Next() {

		rowScanErr := bufferRow.Scan(&movieTitle, &theatherId, &customerName, &seatCount, &createdAt)

		if rowScanErr != nil {
			return rowScanErr
		}

		summaryDetail += fmt.Sprintf("Title\t\t: %s\nTheater\t\t: %s\nName\t\t: %s\nTicket(s)\t: %s\nDate\t\t: %s\n\n", movieTitle, theatherId, customerName, seatCount, createdAt)
	}

	fmt.Printf("\n--------------- BOOKING SUMMARY ---------------\n")
	fmt.Print(summaryDetail)

	return nil
}

// Fetch the current available seat count for the selected session from the database.
func getAvailableSeats(db *sql.DB, selectedSessionId int) (int, error) {

	var totalSeats int

	// Start a transaction and fetch the current available seat count for the selected session.
	trx, err := db.Begin()

	if err != nil {
		return 0, err
	}

	// Fetch the current available seat count for the selected session.
	bufferAvailableSeats := db.QueryRow("SELECT available_seats FROM sessions WHERE id = $1;", selectedSessionId)

	if commitErr := trx.Commit(); commitErr != nil {
		return 0, commitErr
	}

	// Check for error when fetching and assign the output into a native variable
	scanErr := bufferAvailableSeats.Scan(&totalSeats)

	if scanErr != nil {
		return 0, scanErr
	}

	return totalSeats, nil
}

func getSession(db *sql.DB) (int, string, error) {

	// Query and display all active movie sessions.
	var sessionId int
	var movieTitle, rating, duration, theater, sessionDetail string
	var selectedSessionIdIsValid = false
	sessionIdSlice := make([]int, 0)
	var getSessionScanner = bufio.NewScanner(os.Stdin)

	trx, err := db.Begin()

	if err != nil {
		return 0, "", err
	}

	bufferMovieSessions, err := db.Query("SELECT s.id, m.movie_title, m.rating, m.duration, t.id AS theater FROM movies m JOIN sessions s ON m.id = s.movie_id JOIN theaters t ON t.id = s.theater_id WHERE t.is_active = 't' AND s.is_active = 't' ORDER BY m.movie_title;")

	if err != nil {
		if rollbackErr := trx.Rollback(); rollbackErr != nil {
			return 0, "", rollbackErr
		}
		return 0, "", err
	}

	if commitErr := trx.Commit(); commitErr != nil {
		return 0, "", commitErr
	}

	if clearScreenErr := clearScreen(); clearScreenErr != nil {
		return 0, "", clearScreenErr
	}

	for bufferMovieSessions.Next() {
		err := bufferMovieSessions.Scan(&sessionId, &movieTitle, &rating, &duration, &theater)

		if err != nil {

			fmt.Println("ERROR:", err)
			continue

		}

		sessionIdSlice = append(sessionIdSlice, sessionId)

		fmt.Printf("Session ID: %d\nTitle: %s\nRating: %s\nDuration: %s\nTheater: %s\n\n", sessionId, movieTitle, rating, duration, theater)

	}

	// Prompt the user to select a valid session ID and fetch its details.
	for !selectedSessionIdIsValid {

		// Prompt the user to select a session ID.
		fmt.Print("Enter session ID: ")

		getSessionScanner.Scan()
		selectedSessionId, conversionErr := strconv.Atoi(getSessionScanner.Text())

		if conversionErr != nil {

			fmt.Println("ERROR: Expected integer. Please try again.")
			continue

		}

		// If the selection is valid, fetch and return the specific session details.
		if slices.Contains(sessionIdSlice, selectedSessionId) {

			trx, err := db.Begin()

			if err != nil {
				return 0, "", err
			}

			queryRowErr := db.QueryRow("SELECT s.id, m.movie_title, m.rating, m.duration, t.id AS theater FROM movies m JOIN sessions s ON m.id = s.movie_id JOIN theaters t ON t.id = s.theater_id WHERE s.id = $1;", selectedSessionId).Scan(&sessionId, &movieTitle, &rating, &duration, &theater)

			if queryRowErr != nil {

				if rollbackErr := trx.Rollback(); rollbackErr != nil {
					return 0, "", rollbackErr
				}

				return 0, "", queryRowErr
			}

			if commitErr := trx.Commit(); commitErr != nil {
				return 0, "", commitErr
			}

			sessionDetail = fmt.Sprintf("Session ID: %d\nTitle: %s\nRating: %s\nDuration: %s\nTheater: %s\n\n", sessionId, movieTitle, rating, duration, theater)

			selectedSessionIdIsValid = true

		} else {

			fmt.Println("ERROR: Invalid session ID. Please try again.")
			continue

		}

	}

	return sessionId, sessionDetail, nil
}
