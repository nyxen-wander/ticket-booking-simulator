package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ticketBooking(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, debugMsg string) error {

	clearScreen()
	debugMsg = ""

	listActiveSessionsErr := error(nil)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	selectedSessionId, sessionDetail, totalSeats, listActiveSessionsErr := listActiveSessions(db, scanner)

	if listActiveSessionsErr != nil {
		return listActiveSessionsErr
	}

	// Enter the main booking loop until all seats are sold or the user exits.
	for totalSeats > 0 {
		clearScreen()

		fmt.Printf("%v\n\n", debugMsg)
		fmt.Println("----------------- TICKET BOOKING SYSTEM -----------------")
		fmt.Printf("Session Detail:\n%v\n", sessionDetail)
		fmt.Printf("Available ticket(s): %d\n\n", totalSeats)
		fmt.Printf("GUIDE: Enter customer name and ticket amount booked, separated by space. Example: Joe 10\n\n")
		fmt.Print("Input: ")
		debugMsg = ""

		// Read and parse user input for customer name and ticket amount.
		if !scanner.Scan() {
			return scanner.Err()
		}

		userInput := strings.ToLower(scanner.Text())

		switch userInput {
		case "exit":
			// After the loop, offer to print a summary of the bookings.
			fmt.Print("INFO: Booking process has finished. Print summary? (y/n): ")

			if scanner.Scan() {
				if strings.ToLower(scanner.Text()) == "y" {

					bookSumErr := printBookSummary(db, currentTime)

					if bookSumErr != nil {

						return bookSumErr

					}

				}
			} else {
				return scanner.Err()
			}

			fmt.Printf("\nPress enter to continue...")

			if scanner.Scan() {
				return nil
			} else {
				return scanner.Err()
			}

		case "list":
			selectedSessionId, sessionDetail, totalSeats, listActiveSessionsErr = listActiveSessions(db, scanner)

			if listActiveSessionsErr != nil {
				return listActiveSessionsErr

			}

			continue

		default:
			// Validate the customer name and process the booking if seats are available.

			bufferInput := strings.Split(userInput, " ")

			if len(bufferInput) < 2 {

				debugMsg = "ERROR: Invalid input. Refer to the given GUIDE below."

				continue

			}

			ticketAmount, convErr := strconv.Atoi(bufferInput[1])

			if convErr != nil {

				debugMsg = "ERROR: Invalid input. Expected numeric input for ticket amount. Refer to the given GUIDE below."

				continue

			}

			var ticketAmountErr string

			ticketAmount, ticketAmountErr = ticketAmountValidator(ticketAmount)

			if ticketAmountErr != "" {

				debugMsg = ticketAmountErr

				continue

			}

			var customerNameErr string

			customerName := strings.Join(bufferInput[0:len(bufferInput)-1], " ")

			customerName, customerNameErr = customerNameValidator(nameRegex, customerName)

			if customerNameErr != "" {

				debugMsg = customerNameErr

				continue

			}

			if totalSeats >= ticketAmount {

				trx, err := db.Begin()

				if err != nil {

					return err

				}

				// Update the database with the new seat count and record the booking.
				_, execUpdateErr := trx.Exec("UPDATE sessions SET available_seats = (available_seats - $1) WHERE id = $2;", ticketAmount, selectedSessionId)

				if execUpdateErr != nil {

					if rollbackErr := trx.Rollback(); rollbackErr != nil {
						return rollbackErr
					}

					return execUpdateErr

				}

				_, execInsertErr := trx.Exec("INSERT INTO bookings (session_id, customer_name, seat_count) VALUES ($1, $2, $3);", selectedSessionId, customerName, ticketAmount)

				if execInsertErr != nil {

					if rollbackErr := trx.Rollback(); rollbackErr != nil {
						return rollbackErr
					}

					return execInsertErr

				}

				if commitErr := trx.Commit(); commitErr != nil {
					return commitErr
				}

				totalSeats, _ = getAvailableSeats(db, selectedSessionId)

				debugMsg = fmt.Sprintf("SUCCESS: %d seat(s) booked by %s.", ticketAmount, customerName)

			} else {

				debugMsg = fmt.Sprintf("ERROR: Only %d seat(s) remaining.", totalSeats)

				continue

			}

		}
	}

	return nil

}

func listActiveSessions(db *sql.DB, scanner *bufio.Scanner) (int, string, int, error) {

	selectedSessionId, sessionDetail, getSessionErr := getSession(db, scanner)

	if getSessionErr != nil {

		return 0, "", 0, getSessionErr

	}

	totalSeats, getSeatsErr := getAvailableSeats(db, selectedSessionId)

	if getSeatsErr != nil {

		return 0, "", 0, getSeatsErr

	}

	return selectedSessionId, sessionDetail, totalSeats, nil
}
