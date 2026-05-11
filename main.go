package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	godotenv "github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Load environment variables and initialize local variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	connectionString := os.Getenv("DB_URL")
	scanner := bufio.NewScanner(os.Stdin)
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var totalSeats, selectedSessionId int
	var debugMsg, sessionDetail string
	var getSessionErr error

	// Connect to the database and retrieve the user's selected session.
	db, databaseConnectErr := sql.Open("pgx", connectionString)

	if databaseConnectErr != nil {
		log.Fatal(databaseConnectErr)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	selectedSessionId, sessionDetail, getSessionErr = getSession(db)

	if getSessionErr != nil {
		log.Fatal(getSessionErr)
	}

	totalSeats, _ = getAvailableSeats(db, selectedSessionId)

	nameRegex := regexp.MustCompile(`^[[:alpha:]|\s]*$`)

	// Enter the main booking loop until all seats are sold or the user exits.
	for totalSeats > 0 {
		if clearScreenErr := clearScreen(); clearScreenErr != nil {
			fmt.Println(clearScreenErr)
		}

		fmt.Printf("%v\n\n", debugMsg)
		fmt.Println("----------------- TICKET BOOKING SYSTEM -----------------")
		fmt.Printf("Session Detail:\n%v\n", sessionDetail)
		fmt.Printf("Available ticket(s): %d\n\n", totalSeats)
		fmt.Printf("GUIDE: Enter customer name and ticket amount booked, separated by space. Example: Joe 10\n\n")
		fmt.Print("Input: ")

		// Read and parse user input for customer name and ticket amount.
		if !scanner.Scan() {
			fmt.Println("ERROR:", scanner.Err())
			break
		}

		if strings.ToLower(scanner.Text()) == "exit" {
			break
		}

		if strings.ToLower(scanner.Text()) == "list" {

			selectedSessionId, sessionDetail, getSessionErr = getSession(db)

			if getSessionErr != nil {
				log.Fatal(getSessionErr)
			}

			totalSeats, _ = getAvailableSeats(db, selectedSessionId)

			continue
		}

		// Validate the customer name and process the booking if seats are available.
		ticketAmount, customerName, err := customerNameAndTicketAmountValidator(nameRegex, scanner)

		if err != "" {
			debugMsg = err
			continue
		} else {

			if totalSeats >= ticketAmount {

				trx, err := db.Begin()

				if err != nil {
					log.Fatal(err)
				}

				// Update the database with the new seat count and record the booking.
				_, execUpdateErr := db.Exec("UPDATE sessions SET available_seats = (available_seats - $1) WHERE id = $2;", ticketAmount, selectedSessionId)

				if execUpdateErr != nil {

					if rollbackErr := trx.Rollback(); rollbackErr != nil {
						log.Fatal(rollbackErr)
					}

					log.Fatal(execUpdateErr)

				}

				_, execInsertErr := db.Exec("INSERT INTO bookings (session_id, customer_name, seat_count) VALUES ($1, $2, $3);", selectedSessionId, customerName, ticketAmount)

				if execInsertErr != nil {

					if rollbackErr := trx.Rollback(); rollbackErr != nil {
						log.Fatal(rollbackErr)
					}

					log.Fatal(execInsertErr)

				}

				if commitErr := trx.Commit(); commitErr != nil {
					log.Fatal(commitErr)
				}

				totalSeats, _ = getAvailableSeats(db, selectedSessionId)

				debugMsg = fmt.Sprintf("SUCCESS: %d seat(s) booked by %s.", ticketAmount, customerName)

			} else {

				debugMsg = fmt.Sprintf("ERROR: Only %d seat(s) remaining.", totalSeats)

				continue

			}

		}

	}

	// After the loop, offer to print a summary of the bookings.
	fmt.Print("INFO: Booking process has finished. Print summary? (y/n): ")

	if scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "y" {

			if bookSumErr := printBookSummary(db, currentTime); bookSumErr != nil {

				fmt.Println(bookSumErr)

			}

		}
	} else {
		fmt.Println(scanner.Err())
	}

	fmt.Println("Program is exiting.")

	db.Close()
}
