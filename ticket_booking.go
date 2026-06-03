package main

import (
	"bufio"
	"cmp"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

func ticketBooking(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, mr *menuRenderer) error {

	// Retrieve active sessions and initialize the menu renderer state.
	sessionsMap, listSessionsErr := listActiveSessions(db)

	if listSessionsErr != nil {

		return listSessionsErr

	}

	mr.debugMsg = ""
	mr.sessionsMap = sessionsMap
	mr.sessionsMapKeys = make([]int, 0, len(sessionsMap))
	mr.sectionHeader = "TICKET BOOKING"

	for k := range sessionsMap {

		mr.sessionsMapKeys = append(mr.sessionsMapKeys, k)

	}

	// Sort session IDs to ensure a consistent display order.
	slices.SortFunc(mr.sessionsMapKeys, func(a, b int) int {

		return cmp.Compare(a, b)

	})

	// Enter the main booking loop.
	for {

		// Render the session selection menu and get user input.
		if renderMenuErr := mr.renderMenuTicket(); renderMenuErr != nil {

			return renderMenuErr

		}

		if !scanner.Scan() {

			return scanner.Err()

		}

		// Validate the session ID input.
		userInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		sessionId, convErr := strconv.Atoi(userInput)

		if convErr != nil {

			switch userInput {

			case "back":

				return nil

			default:

				mr.debugMsg = "ERROR: Invalid input. Expected numeric input for session ID."

				continue

			}

		}

		if !slices.Contains(mr.sessionsMapKeys, sessionId) {

			mr.debugMsg = "ERROR: Invalid session ID. Refer to the list of available sessions above."

			continue

		}

		// Display the seating map for the chosen session and prompt for a seat.
		if seatsMapErr := renderSeatsMap(db, sessionId); seatsMapErr != nil {

			return seatsMapErr

		}

		fmt.Printf("%s\nPick a seat: ", mr.backGuideMsg)

		if !scanner.Scan() {

			return scanner.Err()

		}

		seatInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if seatInput == "back" {

			return nil

		}

		// Parse the seat code and verify the seat is available.
		rowLetter, seatNumber, getSeatCodeErr := getSeatCode(seatInput)

		if getSeatCodeErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for seat number."

			continue

		}

		sessionSeatId, getUnbookedSeatErr := getUnbookedSeat(db, sessionId, rowLetter, seatNumber)

		if getUnbookedSeatErr != nil {

			mr.debugMsg = "ERROR: Invalid seat or the seat has already been booked. Please try again."

			continue

		}

		// Prompt for and validate the customer's name.
		fmt.Printf("\n%s\nEnter customer name: ", mr.backGuideMsg)

		if !scanner.Scan() {

			return scanner.Err()

		}

		customerName := strings.TrimSpace(scanner.Text())

		if strings.ToLower(customerName) == "back" {

			return nil

		}

		var validatorErr string

		customerName, validatorErr = customerNameValidator(nameRegex, customerName)

		if validatorErr != "" {

			mr.debugMsg = validatorErr

			continue

		}

		// Execute the booking transaction and update the UI message.
		ticketTransErr := ticketTrans(db, sessionId, customerName, sessionSeatId)

		if ticketTransErr != nil {

			return ticketTransErr

		}

		mr.debugMsg = fmt.Sprintf("SUCCESS: Seat %s%d booked by %s.", rowLetter, seatNumber, customerName)

	}

}

func listActiveSessions(db *sql.DB) (map[int]map[string]string, error) {

	sessionsMap := make(map[int]map[string]string)

	// Query the database for active sessions with available seats.
	bufferQuery, querryErr := db.Query("SELECT s.id AS session_id, t.id AS theater_id, m.movie_title, m.rating, m.duration, s.available_seats FROM sessions s JOIN movies m ON s.movie_id = m.id JOIN theaters t ON s.theater_id = t.id WHERE t.is_active = 't' AND s.is_active = 't' AND s.available_seats > 0;")

	if querryErr != nil {

		return nil, querryErr

	}

	var theaterId, movieTitle, rating, duration, availableSeats string
	var sessionId int

	// Iterate through the results and populate the sessions map.
	for bufferQuery.Next() {

		if scanErr := bufferQuery.Scan(&sessionId, &theaterId, &movieTitle, &rating, &duration, &availableSeats); scanErr != nil {

			return nil, scanErr

		}

		if _, ok := sessionsMap[sessionId]; !ok {

			sessionsMap[sessionId] = make(map[string]string)

		}

		sessionsMap[sessionId]["theater_id"] = theaterId
		sessionsMap[sessionId]["movie_title"] = movieTitle
		sessionsMap[sessionId]["rating"] = rating
		sessionsMap[sessionId]["duration"] = duration
		sessionsMap[sessionId]["available_seats"] = availableSeats

	}

	return sessionsMap, nil

}

func ticketTrans(db *sql.DB, sessionId int, customerName string, sessionSeatId int) error {

	// Start a new database transaction.
	trx, beginErr := db.Begin()

	if beginErr != nil {

		return beginErr

	}

	// Insert the new booking record.
	_, execErr := trx.Exec("INSERT INTO bookings (session_id, customer_name, created_at, session_seat_id) VALUES ($1, $2, $3, $4);", sessionId, customerName, time.Now().Format("2006-01-02 15:04:05"), sessionSeatId)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return rollbackErr

		}

		return execErr

	}

	// Decrement the available seats for the session.
	_, execErr = trx.Exec("UPDATE sessions SET available_seats = available_seats - 1 WHERE id = $1;", sessionId)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return rollbackErr

		}

		return execErr

	}

	// Mark the specific seat as booked.
	_, execErr = trx.Exec("UPDATE session_seats SET is_booked = 't' WHERE id = $1;", sessionSeatId)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return rollbackErr

		}

		return execErr

	}

	// Commit the transaction.
	if commitErr := trx.Commit(); commitErr != nil {

		return commitErr

	}

	return nil

}

func getUnbookedSeat(db *sql.DB, sessionId int, rowLetter string, seatNumber int) (int, error) {

	var sessionSeatId int

	// Query for an unbooked seat ID matching the session and physical seat coordinates.
	bufferQueryRow := db.QueryRow("SELECT ss.id FROM session_seats ss JOIN physical_seats ps ON ss.physical_seat_id = ps.id WHERE ss.session_id = $1 AND ps.row_letter = $2 AND ps.seat_num = $3 AND ss.is_booked = 'f';", sessionId, rowLetter, seatNumber)

	if scanErr := bufferQueryRow.Scan(&sessionSeatId); scanErr != nil {

		return 0, scanErr

	}

	return sessionSeatId, nil

}

func getSeatCode(seatInput string) (string, int, error) {

	// Extract the row letter and convert the remaining input to a seat number.
	rowLetter := strings.ToUpper(seatInput[0:1])
	seatNumber, convErr := strconv.Atoi(seatInput[1:])

	if convErr != nil {

		return "", 0, convErr

	}

	return rowLetter, seatNumber, nil

}
