package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"

	godotenv "github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Load environment variables and initialize local variables.
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	connectionString := os.Getenv("DB_URL")
	scanner := bufio.NewScanner(os.Stdin)
	bookRecord := map[string]int{}
	var totalSeats int
	var debugMsg string

	// Connect to the database and retrieve the user's selected session.
	db, databaseConnectErr := sql.Open("pgx", connectionString)

	selectedSessionId, sessionDetail, getSessionErr := getSession(db)

	if getSessionErr != nil {
		log.Fatal(getSessionErr)
	}

	if databaseConnectErr != nil {
		log.Fatal(databaseConnectErr)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	// Fetch the current available seat count for the selected session.
	bufferAvailableSeats := db.QueryRow("SELECT available_seats FROM sessions WHERE id = $1;", selectedSessionId)

	// Check for error when fetching and assign the output into a native variable
	scanErr := bufferAvailableSeats.Scan(&totalSeats)

	if scanErr != nil {
		log.Fatal(scanErr)
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

		// Read and parse user input for customer name and ticket amount.
		scanner.Scan()

		if scanner.Text() == "EXIT" {
			break
		}

		bufferInput := strings.Split(scanner.Text(), " ")
		ticketAmount, err := strconv.Atoi(bufferInput[len(bufferInput)-1])

		if err != nil {
			debugMsg = "ERROR: Invalid input. Refer to the given example below."
			continue
		}

		customerName := strings.ToLower(strings.Join(bufferInput[:len(bufferInput)-1], " "))

		// Validate the customer name and process the booking if seats are available.
		nameRegex := regexp.MustCompile(`^[[:alpha:]|\s]*$`)
		nameCheck := nameRegex.FindString(customerName)

		if nameCheck == "" {
			debugMsg = "ERROR: Invalid customer name. Please try again."
			continue
		} else {

			if totalSeats >= ticketAmount {

				// Update the database with the new seat count and record the booking.
				_, execUpdateErr := db.Exec("UPDATE sessions SET available_seats = (available_seats - $1) WHERE id = $2;", ticketAmount, selectedSessionId)

				if execUpdateErr != nil {

					log.Fatal(execUpdateErr)

				} else {

					_, execInsertErr := db.Exec("INSERT INTO bookings (session_id, customer_name, seat_count) VALUES ($1, $2, $3);", selectedSessionId, customerName, ticketAmount)

					if execInsertErr != nil {

						log.Fatal(execInsertErr)

					} else {

						totalSeats -= ticketAmount

					}
				}

				// Update the local booking record for the final summary.
				keys := make([]string, 0, len(bookRecord))

				for k := range bookRecord {

					keys = append(keys, k)

				}

				if slices.Contains(keys, customerName) {

					bookRecord[customerName] += ticketAmount

				} else {

					bookRecord[customerName] = ticketAmount

				}

				debugMsg = fmt.Sprintf("SUCCESS: %d seat(s) booked by %s.", ticketAmount, customerName)

			} else {

				debugMsg = fmt.Sprintf("ERROR: Only %d seat(s) remaining.", totalSeats)

				continue

			}

		}

	}

	// Print the final booking summary and close the database connection.
	keys := make([]string, 0, len(bookRecord))

	for k := range bookRecord {

		keys = append(keys, k)

	}

	fmt.Printf("\n----------------- BOOKING SUMMARY -----------------\n")

	for key := range keys {

		fmt.Printf("%s: %d ticket(s).\n", keys[key], bookRecord[keys[key]])

	}

	db.Close()
}

func getSession(db *sql.DB) (int, string, error) {

	// Query and display all active movie sessions.
	var sessionId int
	var movieTitle, rating, duration, theater, sessionDetail string
	var selectedSessionIdIsValid = false
	sessionIdSlice := make([]int, 0)
	var getSessionScanner = bufio.NewScanner(os.Stdin)

	bufferMovieSessions, err := db.Query("select s.id, m.movie_title, m.rating, m.duration, t.id as theater from movies m join sessions s on m.id = s.movie_id join theaters t on t.id = s.theater_id where t.is_active = 't' and s.is_active = 't' order by m.movie_title;")

	if err != nil {
		return 0, "", err
	}

	clearScreen()

	for !selectedSessionIdIsValid {

		for bufferMovieSessions.Next() {
			err := bufferMovieSessions.Scan(&sessionId, &movieTitle, &rating, &duration, &theater)

			if err != nil {

				fmt.Println("ERROR:", err)
				continue

			}

			sessionIdSlice = append(sessionIdSlice, sessionId)

			fmt.Printf("Session ID: %d\nTitle: %s\nRating: %s\nDuration: %s\nTheater: %s\n\n", sessionId, movieTitle, rating, duration, theater)

		}

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

			err := db.QueryRow("select s.id, m.movie_title, m.rating, m.duration, t.id as theater from movies m join sessions s on m.id = s.movie_id join theaters t on t.id = s.theater_id where s.id = $1;", selectedSessionId).Scan(&sessionId, &movieTitle, &rating, &duration, &theater)

			if err != nil {

				fmt.Println("ERROR:", err)
				continue
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

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
