package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/term"
)

func mainMenu(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, mr *menuRenderer, debugMsg string) error {

	isCompleted := false

	// Loop until the user chooses to exit the main menu.
	for !isCompleted {

		// Render the main menu with the current debug message.
		mr.sectionHeader = "MAIN MENU"
		mr.menuSlice = mainMenuList()
		mr.debugMsg = debugMsg

		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return renderMenuErr

		}

		// Read user input and handle different menu options.
		if !scanner.Scan() {
			return scanner.Err()
		}

		userInput := scanner.Text()

		switch userInput {
		case "1":

			// Validate admin credentials before entering the administration menu.
			isAdmin, adminValidatorErr := adminValidator()

			if adminValidatorErr != nil {

				return adminValidatorErr

			}

			if !isAdmin {

				mr.debugMsg = "ERROR: Failed to validate admin credentials."

				continue

			} else {

				// Enter the administration sub-menu and handle navigation signals.
				signal, administrationErr := administration(db, scanner, nameRegex, mr, debugMsg)

				if administrationErr != nil {

					return administrationErr

				}

				if signal == "back" {

					continue

				}

				if signal == "exit" {

					isCompleted = true

				}

			}

		case "2":

			// Enter the ticket booking flow.
			if ticketBookingErr := ticketBooking(db, scanner, nameRegex, debugMsg); ticketBookingErr != nil {

				return ticketBookingErr

			}

			continue

		case "3":

			isCompleted = true

		default:

			mr.debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return nil

}

func administration(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, mr *menuRenderer, debugMsg string) (string, error) {

	isCompleted := false

	// Loop until the user chooses to go back or exit the administration menu.
	for !isCompleted {

		mr.sectionHeader = "ADMINISTRATION MENU"
		mr.menuSlice = adminMenuList()
		mr.debugMsg = debugMsg

		// Render the administration menu.
		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return "", renderMenuErr

		}

		// Handle administrative actions like delete, insert, list, and update.
		if !scanner.Scan() {

			return "", scanner.Err()

		}

		userInput := scanner.Text()

		switch userInput {
		case "1":

			signal, deleteErr := deleteMenu(db, scanner, mr, debugMsg)

			if deleteErr != nil {

				return "", deleteErr

			}

			if signal == "back" {

				continue

			}

			if signal == "exit" {

				isCompleted = true

			}

		case "2":

			signal, insertErr := insertMenu(db, scanner, nameRegex, mr, debugMsg)

			if insertErr != nil {

				return "", insertErr

			}

			if signal == "back" {

				continue

			}

			if signal == "exit" {

				isCompleted = true

			}

		case "3":

			signal, listErr := listMenu(db, scanner, mr, debugMsg)

			if listErr != nil {

				return "", listErr

			}

			if signal == "back" {

				continue

			}

			if signal == "exit" {

				isCompleted = true

			}

		case "4":

			signal, updateErr := updateMenu(db, scanner, mr, debugMsg)

			if updateErr != nil {

				return "", updateErr

			}

			if signal == "back" {

				continue

			}

			if signal == "exit" {

				isCompleted = true

			}

		case "5":

			return "back", nil

		case "6":

			return "exit", nil

		default:

			mr.debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return "exit", nil

}

func adminValidator() (bool, error) {

	attempt := 0
	isValidated := false

	// Allow up to three attempts to enter the correct admin password.
	for attempt < 3 {

		fmt.Print("Enter password: ")

		// Read the password securely from the terminal.
		passByte, readPasswordErr := term.ReadPassword(int(os.Stdin.Fd()))

		if readPasswordErr != nil {

			return false, readPasswordErr

		}

		// Check the password against the environment variable.
		if string(passByte) != os.Getenv("ADMIN_PASSWD") {

			attempt++

			fmt.Println("ERROR: Invalid password.")

			continue

		} else {

			isValidated = true

			break

		}

	}

	return isValidated, nil

}

func clearScreen() error {

	// Execute the appropriate system command to clear the terminal screen.
	if runtime.GOOS == "windows" {

		cmd := exec.Command("cmd", "/c", "cls")

		cmd.Stdout = os.Stdout

		cmd.Run()

	} else {

		cmd := exec.Command("clear")

		cmd.Stdout = os.Stdout

		if cmdRunErr := cmd.Run(); cmdRunErr != nil {

			return cmdRunErr

		}
	}

	return nil
}

func printBookSummary(db *sql.DB, currentTime string) error {

	// Query the database for all bookings made since the program started.
	bufferRow, queryErr := db.Query("SELECT m.movie_title, s.theater_id, b.customer_name, b.seat_count, to_char(b.created_at, 'YYYY-MM-DD HH24:MI:SS') FROM bookings b JOIN sessions s ON b.session_id = s.id JOIN movies m ON m.id = s.movie_id WHERE b.created_at BETWEEN $1 AND now();", currentTime)

	if queryErr != nil {

		return queryErr

	}

	var movieTitle, theatherId, customerName, seatCount, createdAt, summaryDetail string

	// Iterate through the results to build and print a formatted summary string.
	for bufferRow.Next() {

		if rowScanErr := bufferRow.Scan(&movieTitle, &theatherId, &customerName, &seatCount, &createdAt); rowScanErr != nil {

			return rowScanErr

		}

		summaryDetail += fmt.Sprintf("Title\t\t: %s\nTheater\t\t: %s\nName\t\t: %s\nTicket(s)\t: %s\nDate\t\t: %s\n\n", movieTitle, theatherId, customerName, seatCount, createdAt)
	}

	fmt.Println("--------------- BOOKING SUMMARY ---------------")

	fmt.Print(summaryDetail)

	return nil
}

func getAvailableSeats(db *sql.DB, selectedSessionId int) (int, error) {

	var totalSeats int

	// Fetch the current available seat count for the selected session.
	bufferAvailableSeats := db.QueryRow("SELECT available_seats FROM sessions WHERE id = $1;", selectedSessionId)

	// Check for error when fetching and assign the output into a native variable
	if scanErr := bufferAvailableSeats.Scan(&totalSeats); scanErr != nil {

		return 0, scanErr

	}

	return totalSeats, nil
}

func getSession(db *sql.DB, scanner *bufio.Scanner) (int, string, error) {

	// Query and display all active movie sessions.
	var sessionId int
	var movieTitle, rating, duration, theater, sessionDetail string
	var selectedSessionIdIsValid = false

	sessionIdSlice := make([]int, 0)

	bufferMovieSessions, err := db.Query("SELECT s.id, m.movie_title, m.rating, m.duration, t.id AS theater FROM movies m JOIN sessions s ON m.id = s.movie_id JOIN theaters t ON t.id = s.theater_id WHERE t.is_active = 't' AND s.is_active = 't' ORDER BY m.movie_title;")

	if err != nil {

		return 0, "", err

	}

	if clearScreenErr := clearScreen(); clearScreenErr != nil {

		return 0, "", clearScreenErr

	}

	// Iterate through the sessions and print their details.
	for bufferMovieSessions.Next() {

		if err := bufferMovieSessions.Scan(&sessionId, &movieTitle, &rating, &duration, &theater); err != nil {

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

		if !scanner.Scan() {

			return 0, "", scanner.Err()

		}

		selectedSessionId, conversionErr := strconv.Atoi(scanner.Text())

		if conversionErr != nil {

			fmt.Println("ERROR: Expected integer. Please try again.")

			continue

		}

		// If the selection is valid, fetch and return the specific session details.
		if slices.Contains(sessionIdSlice, selectedSessionId) {

			db.QueryRow("SELECT s.id, m.movie_title, m.rating, m.duration, t.id AS theater FROM movies m JOIN sessions s ON m.id = s.movie_id JOIN theaters t ON t.id = s.theater_id WHERE s.id = $1;", selectedSessionId).Scan(&sessionId, &movieTitle, &rating, &duration, &theater)

			sessionDetail = fmt.Sprintf("Session ID: %d\nTitle: %s\nRating: %s\nDuration: %s\nTheater: %s\n\n", sessionId, movieTitle, rating, duration, theater)

			selectedSessionIdIsValid = true

		} else {

			fmt.Println("ERROR: Invalid session ID. Please try again.")

			continue

		}

	}

	return sessionId, sessionDetail, nil
}

func customerNameValidator(nameRegex *regexp.Regexp, customerName string) (string, string) {

	// Validate the customer name against the provided regex.
	if nameRegex.FindString(customerName) == "" {

		return "", "ERROR: Invalid customer name. Please try again."

	}

	return customerName, ""
}

func ticketAmountValidator(ticketAmount int) (int, string) {

	// Ensure the ticket amount is a positive integer.
	if ticketAmount <= 0 {

		return 0, "ERROR: Invalid ticket amount. Please try again."

	}

	return ticketAmount, ""

}

func scannerInitializator() (*bufio.Scanner, error) {

	// Initialize a new scanner for standard input.
	scanner := bufio.NewScanner(os.Stdin)

	if scannerErr := scanner.Err(); scannerErr != nil {

		return nil, scannerErr

	}

	return scanner, nil

}

func dbInitializator() (*sql.DB, error) {

	// Load environment variables and initialize local variables.
	if err := godotenv.Load(); err != nil {

		return nil, err

	}

	connectionString := os.Getenv("DB_URL")

	// Connect to the database and retrieve the user's selected session.
	db, databaseConnectErr := sql.Open("pgx", connectionString)

	if databaseConnectErr != nil {

		return nil, databaseConnectErr

	}

	if err := db.Ping(); err != nil {

		return nil, err

	}

	return db, nil

}

type menuRenderer struct {
	sectionHeader string
	menuSlice     []string
	debugMsg      string

	additionalMsg       string
	columnRows          map[int]map[string]string
	columnRowsKeys      []int
	listColumnsErr      error
	convertedColumnRows string
}

func (mr *menuRenderer) renderMenu() error {

	// Clear the screen and render the standard menu header and options.
	if clearScreenErr := clearScreen(); clearScreenErr != nil {

		return clearScreenErr

	}

	fmt.Printf("------------------------------------------- %s -------------------------------------------\n\n", mr.sectionHeader)

	for idx, menu := range mr.menuSlice {

		fmt.Printf("%d. %s\n", idx+1, menu)

	}

	// Print the debug message and input prompt, then reset the debug message.
	fmt.Printf("\n%s\n\n", mr.debugMsg)
	fmt.Printf("Input: ")

	mr.debugMsg = ""

	return nil

}

func (mr *menuRenderer) renderMenuFilter() error {

	// Clear the screen and render the menu header with filter-specific content and guides.
	if clearScreenErr := clearScreen(); clearScreenErr != nil {

		return clearScreenErr

	}

	fmt.Printf("------------------------------------------- %s -------------------------------------------\n\n", mr.sectionHeader)
	fmt.Print(mr.additionalMsg)
	fmt.Print(mr.convertedColumnRows)
	fmt.Printf("\n%s\n\n", mr.debugMsg)
	fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n")
	fmt.Printf("GUIDE: Type \"Filter <keyword>\" to filter the movie title out by <keyword>. Example: Filter avangers\n")
	fmt.Printf("GUIDE: Type \"Reset\" to reset the filter result.\n\n")

	mr.debugMsg = ""

	return nil

}

func (mr *menuRenderer) renderMenuNonFilter() error {

	// Clear the screen and render the menu header with non-filter content and a back guide.
	if clearScreenErr := clearScreen(); clearScreenErr != nil {

		return clearScreenErr

	}

	fmt.Printf("------------------------------------------- %s -------------------------------------------\n\n", mr.sectionHeader)
	fmt.Print(mr.convertedColumnRows)
	fmt.Printf("\n%s\n\n", mr.debugMsg)
	fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n")

	mr.debugMsg = ""

	return nil

}

func (mr *menuRenderer) renderMenuInsert() error {

	// Clear the screen and render the header for the insertion menu.
	if clearScreenErr := clearScreen(); clearScreenErr != nil {

		return clearScreenErr

	}

	fmt.Printf("------------------------------------------- %s -------------------------------------------\n\n", mr.sectionHeader)
	fmt.Printf("GUIDE: Type \"debug back\" to go back to the previous menu.\n")

	return nil

}
