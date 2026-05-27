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

func deleteMenu(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) (string, error) {

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "DELETE MENU"
		mr.menuSlice = tablesList()
		mr.debugMsg = debugMsg

		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return "", renderMenuErr

		}

		if !scanner.Scan() {

			return "", scanner.Err()

		}

		userInput := scanner.Text()

		switch userInput {

		case "1":

			deleteMovieErr := deleteMovie(db, scanner, mr, debugMsg)

			if deleteMovieErr != nil {

				return "", deleteMovieErr

			}

			continue

		case "2":

			if deleteTheaterErr := deleteTheater(db, scanner, mr, debugMsg); deleteTheaterErr != nil {

				return "", deleteTheaterErr

			}

			continue

		case "3":

			if deleteSessionErr := deleteSession(db, scanner, mr, debugMsg); deleteSessionErr != nil {

				return "", deleteSessionErr

			}

			continue

		case "4":

			if deleteBookingErr := deleteBooking(db, scanner, mr, debugMsg); deleteBookingErr != nil {

				return "", deleteBookingErr

			}

			continue

		case "5":

			return "back", nil

		case "6":

			return "exit", nil

		default:

			debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return "exit", nil
}

func insertMenu(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, mr *menuRenderer, debugMsg string) (string, error) {

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "INSERT MENU"
		mr.menuSlice = tablesList()
		mr.debugMsg = debugMsg

		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return "", renderMenuErr

		}

		if !scanner.Scan() {

			return "", scanner.Err()

		}

		userInput := scanner.Text()

		switch userInput {

		case "1":

			if insertMovieErr := insertMovie(db, scanner, mr); insertMovieErr != nil {

				return "", insertMovieErr

			}

			continue

		case "2":

			if insertTheaterErr := insertTheater(db, scanner, mr); insertTheaterErr != nil {

				return "", insertTheaterErr

			}

			continue

		case "3":

			if insertSessionErr := insertSession(db, scanner, mr); insertSessionErr != nil {

				return "", insertSessionErr

			}

			continue

		case "4":

			if insertBookingErr := insertBooking(db, scanner, nameRegex, mr); insertBookingErr != nil {

				return "", insertBookingErr

			}

			continue

		case "5":

			return "back", nil

		case "6":

			return "exit", nil

		default:

			debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return "exit", nil

}

func listMenu(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) (string, error) {

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "LIST MENU"
		mr.menuSlice = tablesList()
		mr.debugMsg = debugMsg

		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return "", renderMenuErr

		}

		if !scanner.Scan() {

			return "", scanner.Err()

		}

		userInput := scanner.Text()

		switch userInput {

		case "1":

			moviesRow, _, listMoviesErr := listMovies(db)

			if listMoviesErr != nil {

				return "", listMoviesErr

			}

			convertedMoviesRow := convertedListMovies(moviesRow)

			if err := clearScreen(); err != nil {

				return "", err

			}

			fmt.Print(convertedMoviesRow)

			fmt.Printf("\nPress Enter to continue...")

			if !scanner.Scan() {

				return "", scanner.Err()

			}

			continue

		case "2":

			theatersRow, _, listTheatersErr := listTheaters(db)

			if listTheatersErr != nil {

				return "", listTheatersErr

			}

			convertedTheatersRow := convertedListTheaters(theatersRow)

			if err := clearScreen(); err != nil {

				return "", err

			}

			fmt.Print(convertedTheatersRow)

			fmt.Printf("\nPress Enter to continue...")

			if !scanner.Scan() {

				return "", scanner.Err()

			}

			continue

		case "3":

			sessionsRow, _, listSessionsErr := listSessions(db)

			if listSessionsErr != nil {

				return "", listSessionsErr

			}

			convertedSessionsRow := convertedListSessions(sessionsRow)

			if err := clearScreen(); err != nil {

				return "", err

			}

			fmt.Print(convertedSessionsRow)

			fmt.Printf("\nPress Enter to continue...")

			if !scanner.Scan() {

				return "", scanner.Err()

			}

			continue

		case "4":

			bookingsRow, _, listBookingsErr := listBookings(db)

			if listBookingsErr != nil {

				return "", listBookingsErr

			}

			convertedBookingsRow := convertedListBookings(bookingsRow)

			if err := clearScreen(); err != nil {

				return "", err

			}

			fmt.Print(convertedBookingsRow)

			fmt.Printf("\nPress Enter to continue...")

			if !scanner.Scan() {

				return "", scanner.Err()

			}

			continue

		case "5":

			return "back", nil

		case "6":

			return "exit", nil

		default:

			debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return "exit", nil

}

func updateMenu(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) (string, error) {

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "UPDATE MENU"
		mr.menuSlice = tablesList()
		mr.debugMsg = debugMsg

		if renderMenuErr := mr.renderMenu(); renderMenuErr != nil {

			return "", renderMenuErr

		}

		if !scanner.Scan() {
			return "", scanner.Err()
		}

		userInput := scanner.Text()

		switch userInput {

		case "1":

			if updateMovieErr := updateMovie(db, scanner, mr, debugMsg); updateMovieErr != nil {

				return "", updateMovieErr

			}

			continue

		case "2":

			if updateTheaterErr := updateTheater(db, scanner, mr, debugMsg); updateTheaterErr != nil {

				return "", updateTheaterErr

			}

			continue

		case "3":

			if updateSessionErr := updateSession(db, scanner, mr, debugMsg); updateSessionErr != nil {

				return "", updateSessionErr

			}

			continue

		case "4":

			if updateBookingErr := updateBooking(db, scanner, mr, debugMsg); updateBookingErr != nil {

				return "", updateBookingErr

			}

			continue

		case "5":

			return "back", nil

		case "6":

			return "exit", nil

		default:

			debugMsg = "ERROR: Invalid input. Please try again."

			continue

		}

	}

	return "exit", nil

}

func deleteMovie(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var title, rating, duration, selectedRecord string
	var id int

	mr.additionalMsg = ""
	mr.debugMsg = debugMsg

	mr.columnRows, mr.columnRowsKeys, mr.listColumnsErr = listMovies(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListMovies(mr.columnRows)

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "DELETE MOVIE"

		if renderMenuTableFilterErr := mr.renderMenuFilter(); renderMenuTableFilterErr != nil {

			return renderMenuTableFilterErr

		}

		fmt.Printf("Movie ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		if strings.HasPrefix(userInput, "filter") {

			bufferUserInput := strings.Split(userInput, " ")

			if len(bufferUserInput) != 2 {

				mr.debugMsg = "ERROR: Invalid input. Refer to the given GUIDE below."

				continue

			}

			var movieRowsFilterErr error

			mr.convertedColumnRows, mr.additionalMsg, movieRowsFilterErr = movieRowsFilter(db, bufferUserInput, mr.columnRows, mr.columnRowsKeys)

			if movieRowsFilterErr != nil {

				return movieRowsFilterErr

			}

			continue

		}

		if userInput == "reset" {

			mr.additionalMsg = ""

			mr.convertedColumnRows = convertedListMovies(mr.columnRows)

			continue

		}

		mr.additionalMsg = ""

		idx, convErr := strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for movie ID. Refer to the given GUIDE below."

			continue

		}

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- DELETE CONFIRMATION ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM movies WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &title, &rating, &duration); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nTitle\t\t: %s\nRating\t\t: %s\nDuration\t: %s\n\n", id, title, rating, duration)

		fmt.Printf("%s\nDelete this record? (y/n):", selectedRecord)

		if !scanner.Scan() {

			return scanner.Err()

		}

		if scanner.Text() == "y" || scanner.Text() == "Y" {

			// here lies the trx to delete the row by the idx
			trx, beginErr := db.Begin()

			if beginErr != nil {

				return beginErr

			}

			_, execErr := trx.Exec("DELETE FROM movies WHERE id = $1;", idx)

			if execErr != nil {

				if rollbackErr := trx.Rollback(); rollbackErr != nil {
					return rollbackErr
				}

				return execErr

			}

			if commitErr := trx.Commit(); commitErr != nil {
				return commitErr
			}

			isCompleted = true

			fmt.Printf("\nSUCCESS: Record has been deleted.\n\n")

			fmt.Print("Press Enter to continue..")

			if !scanner.Scan() {

				return scanner.Err()

			}

		} else {

			continue

		}
	}

	return nil

}

func deleteTheater(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var totalCapacity, isActive, selectedRecord string
	var id int

	mr.debugMsg = debugMsg

	mr.columnRows, _, mr.listColumnsErr = listTheaters(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListTheaters(mr.columnRows)

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "DELETE THEATER"

		if renderMenuErr := mr.renderMenuNonFilter(); renderMenuErr != nil {

			return renderMenuErr

		}

		fmt.Printf("Theater ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		idx, convErr := strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for theater ID. Refer to the given GUIDE below."

			continue

		}

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- DELETE CONFIRMATION ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM theaters WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &totalCapacity, &isActive); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nTotal Capacity\t: %s\nIs Active\t: %s\n\n", id, totalCapacity, isActive)

		fmt.Printf("%s\nDelete this record? (y/n):", selectedRecord)

		if !scanner.Scan() {

			return scanner.Err()

		}

		if scanner.Text() == "y" || scanner.Text() == "Y" {

			trx, beginErr := db.Begin()

			if beginErr != nil {
				return beginErr
			}

			_, execErr := trx.Exec("DELETE FROM theaters WHERE id = $1;", idx)

			if execErr != nil {

				if rollbackErr := trx.Rollback(); rollbackErr != nil {

					return rollbackErr

				}

				return execErr

			}

			if commitErr := trx.Commit(); commitErr != nil {

				return commitErr

			}

			isCompleted = true

			fmt.Printf("SUCCESS: Record has been deleted.\n\n")

			fmt.Print("Press Enter to continue..")

			if !scanner.Scan() {

				return scanner.Err()

			}

		} else {

			return nil

		}

	}

	return nil

}

func deleteSession(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var movieId, theaterId, availableSeats, isActive, selectedRecord string
	var id int

	mr.debugMsg = debugMsg

	mr.columnRows, _, mr.listColumnsErr = listSessions(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListSessions(mr.columnRows)

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "DELETE SESSION"

		if renderMenuErr := mr.renderMenuNonFilter(); renderMenuErr != nil {

			return renderMenuErr

		}

		fmt.Printf("Session ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		idx, convErr := strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for session ID. Refer to the given GUIDE below."

			continue

		}

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- DELETE CONFIRMATION ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM sessions WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &movieId, &theaterId, &availableSeats, &isActive); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nMovie ID\t: %s\nTheater ID\t: %s\nAvailable Seats\t: %s\nIs Active\t: %s\n\n", id, movieId, theaterId, availableSeats, isActive)

		fmt.Printf("%s\nDelete this record? (y/n):", selectedRecord)

		if !scanner.Scan() {

			return scanner.Err()

		}

		if scanner.Text() == "y" || scanner.Text() == "Y" {

			// here lies the trx to delete the row by the idx
			trx, beginErr := db.Begin()

			if beginErr != nil {
				return beginErr
			}

			_, execErr := trx.Exec("DELETE FROM sessions WHERE id = $1;", idx)

			if execErr != nil {

				if rollbackErr := trx.Rollback(); rollbackErr != nil {

					return rollbackErr

				}

				return execErr

			}

			if commitErr := trx.Commit(); commitErr != nil {

				return commitErr

			}

			isCompleted = true

			fmt.Printf("SUCCESS: Record has been deleted.\n\n")

			fmt.Print("Press Enter to continue..")

			if !scanner.Scan() {

				return scanner.Err()

			}

		} else {

			return nil

		}

	}

	return nil

}

func deleteBooking(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var sessionId, customerName, seatCount, createdAt, selectedRecord string
	var id int

	mr.additionalMsg = ""
	mr.debugMsg = debugMsg

	mr.columnRows, mr.columnRowsKeys, mr.listColumnsErr = listBookings(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListBookings(mr.columnRows)

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "DELETE BOOKING"

		if renderMenuTableFilterErr := mr.renderMenuFilter(); renderMenuTableFilterErr != nil {

			return renderMenuTableFilterErr

		}

		fmt.Printf("Booking ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput := strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		if strings.HasPrefix(userInput, "filter") {

			bufferUserInput := strings.Split(userInput, " ")

			if len(bufferUserInput) != 2 {

				mr.debugMsg = "ERROR: Invalid input. Refer to the given GUIDE below."

				continue

			}

			var movieRowsFilterErr error

			mr.convertedColumnRows, mr.additionalMsg, movieRowsFilterErr = bookingRowsFilter(db, bufferUserInput, mr.columnRows, mr.columnRowsKeys)

			if movieRowsFilterErr != nil {

				return movieRowsFilterErr

			}

			continue

		}

		if userInput == "reset" {

			mr.additionalMsg = ""

			mr.convertedColumnRows = convertedListBookings(mr.columnRows)

			continue

		}

		mr.additionalMsg = ""

		idx, convErr := strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for booking ID. Refer to the given GUIDE below."

			continue

		}

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- DELETE CONFIRMATION ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM bookings WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &sessionId, &customerName, &seatCount, &createdAt); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nSession ID\t: %s\nCustomer Name\t: %s\nSeat Count\t: %s\nCreated At\t: %s\n\n", id, sessionId, customerName, seatCount, createdAt)

		fmt.Printf("%s\nDelete this record? (y/n):", selectedRecord)

		if !scanner.Scan() {

			return scanner.Err()

		}

		if scanner.Text() == "y" || scanner.Text() == "Y" {

			trx, beginErr := db.Begin()

			if beginErr != nil {

				return beginErr

			}

			_, execErr := trx.Exec("DELETE FROM bookings WHERE id = $1;", idx)

			if execErr != nil {

				if rollbackErr := trx.Rollback(); rollbackErr != nil {

					return rollbackErr

				}

				return execErr

			}

			if commitErr := trx.Commit(); commitErr != nil {

				return commitErr

			}

			isCompleted = true

			fmt.Printf("SUCCESS: Record has been deleted.\n\n")

			fmt.Print("Press Enter to continue..")

			if !scanner.Scan() {

				return scanner.Err()

			}

		} else {

			return nil

		}
	}

	return nil

}

func insertMovie(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer) error {

	var movieTitle, rating, duration string
	var convErr error

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "INSERT MOVIE"

		if renderMenuInsertErr := mr.renderMenuInsert(); renderMenuInsertErr != nil {

			return renderMenuInsertErr

		}

		fmt.Print("Movie Title: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		movieTitle = strings.TrimSpace(scanner.Text())

		if strings.ToLower(movieTitle) == "debug back" {

			return nil

		}

		isRatingValid := false

		for !isRatingValid {

			fmt.Printf("\nGUIDE: Type \"debug back\" to go back to the previous menu.\n\n")

			fmt.Print("Rating: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			rating = strings.TrimSpace(scanner.Text())

			if strings.ToLower(rating) == "debug back" {

				return nil

			}

			bufferRating, convErr := strconv.ParseFloat(rating, 64)

			if convErr != nil {

				fmt.Printf("\n\nERROR: Invalid input. Expected numeric in double precision float format for rating (1 - 10 scale).\n\n")

				continue

			}

			rating = fmt.Sprintf("%.1f", bufferRating)

			isRatingValid = true

		}

		isDurationValid := false

		for !isDurationValid {

			fmt.Printf("\nGUIDE: Type \"debug back\" to go back to the previous menu.\n\n")

			fmt.Print("Duration: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			duration = strings.TrimSpace(scanner.Text())

			if strings.ToLower(duration) == "debug back" {

				return nil

			}

			if _, convErr = strconv.Atoi(duration); convErr != nil {

				fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for duration (in minutes).\n\n")

				continue

			}

			isDurationValid = true

		}

		trx, beginErr := db.Begin()

		if beginErr != nil {

			return beginErr

		}

		_, execErr := trx.Exec("INSERT INTO movies (movie_title, rating, duration) VALUES ($1, $2, $3);", movieTitle, rating, duration)

		if execErr != nil {

			if rollbackErr := trx.Rollback(); rollbackErr != nil {

				return rollbackErr

			}

			return execErr

		}

		if commitErr := trx.Commit(); commitErr != nil {

			return commitErr

		}

		fmt.Printf("\nSUCCESS: Record has been inserted.\n\n")

		fmt.Println("Press Enter to continue..")

		if !scanner.Scan() {

			return scanner.Err()

		}

		isCompleted = true

	}

	return nil

}

func insertTheater(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer) error {

	var totalCapacity, isActive string
	var convErr error

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "INSERT THEATER"

		if renderMenuInsertErr := mr.renderMenuInsert(); renderMenuInsertErr != nil {

			return renderMenuInsertErr

		}

		fmt.Print("Total Capacity: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		totalCapacity = strings.TrimSpace(scanner.Text())

		if strings.ToLower(totalCapacity) == "debug back" {

			return nil

		}

		if _, convErr = strconv.Atoi(totalCapacity); convErr != nil {

			fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for total capacity.\n\n")

			continue

		}

		isIsActiveValid := false

		for !isIsActiveValid {

			fmt.Print("Is currently active? (y/n): ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			isActive = strings.TrimSpace(strings.ToLower(scanner.Text()))

			if isActive == "debug back" {

				return nil

			}

			if isActive != "y" && isActive != "n" {

				fmt.Printf("\n\nERROR: Invalid input. Expected \"y\" or \"n\".\n\n")

				continue

			}

		}

		trx, beginErr := db.Begin()

		if beginErr != nil {

			return beginErr

		}

		_, execErr := trx.Exec("INSERT INTO theaters (total_capacity, is_active) VALUES ($1, $2);", totalCapacity, isActive)

		if execErr != nil {

			if rollbackErr := trx.Rollback(); rollbackErr != nil {

				return rollbackErr

			}

			return execErr

		}

		if commitErr := trx.Commit(); commitErr != nil {

			return commitErr

		}

		fmt.Printf("\nSUCCESS: Record has been inserted.\n\n")

		fmt.Println("Press Enter to continue..")

		if !scanner.Scan() {

			return scanner.Err()

		}

		isCompleted = true

	}

	return nil

}

func insertSession(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer) error {

	var movieId, theaterId, availableSeats, isActive string
	var convErr error
	isCompleted := false

	moviesIdSlice, getMoviesIdErr := getMoviesId(db)

	if getMoviesIdErr != nil {

		return getMoviesIdErr

	}

	theatersIdSlice, getTheatersIdErr := getTheatersId(db)

	if getTheatersIdErr != nil {

		return getTheatersIdErr

	}

	for !isCompleted {

		mr.sectionHeader = "INSERT SESSION"

		if renderMenuInsertErr := mr.renderMenuInsert(); renderMenuInsertErr != nil {

			return renderMenuInsertErr

		}
		fmt.Print("Movie ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		movieId = strings.TrimSpace(scanner.Text())

		if strings.ToLower(movieId) == "debug back" {

			return nil

		}

		var tempInput int

		tempInput, convErr = strconv.Atoi(movieId)

		if convErr != nil {

			fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for movie ID.\n\n")

			continue

		}

		if !slices.Contains(moviesIdSlice, tempInput) {

			fmt.Printf("\n\nERROR: Non-existing movie ID. Please try again.\n\n")

			continue

		}

		isTheaterIdValid := false

		for !isTheaterIdValid {

			fmt.Print("Theater ID: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			theaterId = strings.TrimSpace(scanner.Text())

			if strings.ToLower(theaterId) == "debug back" {

				return nil

			}

			var tempInput int

			tempInput, convErr = strconv.Atoi(theaterId)

			if convErr != nil {

				fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for theater ID.\n\n")

				continue

			}

			if !slices.Contains(theatersIdSlice, tempInput) {

				fmt.Printf("\n\nERROR: Non-existing theater ID. Please try again.\n\n")

				continue

			}

			isTheaterIdValid = true

		}

		isSeatCountValid := false

		for !isSeatCountValid {

			fmt.Print("Seat(s) Count: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			availableSeats = strings.TrimSpace(scanner.Text())

			if strings.ToLower(availableSeats) == "debug back" {

				return nil

			}

			if _, convErr = strconv.Atoi(availableSeats); convErr != nil {

				fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for total seats.\n\n")

				continue

			}

			isSeatCountValid = true

		}

		isActiveValid := false

		for !isActiveValid {

			fmt.Print("Is currently active? (y/n): ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			isActive = strings.TrimSpace(strings.ToLower(scanner.Text()))

			if isActive == "debug back" {

				return nil

			}

			if isActive != "y" && isActive != "n" {

				fmt.Printf("\n\nERROR: Invalid input. Expected \"y\" or \"n\".\n\n")

				continue

			}

			isActiveValid = true

		}

		trx, beginErr := db.Begin()

		if beginErr != nil {

			return beginErr

		}

		_, execErr := trx.Exec("INSERT INTO sessions (movie_id, theater_id, available_seats, is_active) VALUES ($1, $2, $3, $4);", movieId, theaterId, availableSeats, isActive)

		if execErr != nil {

			if rollbackErr := trx.Rollback(); rollbackErr != nil {

				return rollbackErr

			}

			return execErr

		}

		if commitErr := trx.Commit(); commitErr != nil {

			return commitErr

		}

		fmt.Printf("\nSUCCESS: Record has been inserted.\n\n")

		fmt.Println("Press Enter to continue..")

		if !scanner.Scan() {

			return scanner.Err()

		}

		isCompleted = true

	}

	return nil

}

func insertBooking(db *sql.DB, scanner *bufio.Scanner, nameRegex *regexp.Regexp, mr *menuRenderer) error {

	var sessionId, customerName, seatCount string
	var convErr error

	sessionsIdSlice, getSessionsIdErr := getSessionsId(db)

	if getSessionsIdErr != nil {

		return getSessionsIdErr

	}

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "INSERT BOOKING"

		if renderMenuInsertErr := mr.renderMenuInsert(); renderMenuInsertErr != nil {

			return renderMenuInsertErr

		}

		fmt.Print("Session ID: ")

		if !scanner.Scan() {

			return scanner.Err()

		}

		sessionId = strings.TrimSpace(scanner.Text())

		if strings.ToLower(sessionId) == "debug back" {

			return nil

		}

		var tempInput int

		tempInput, convErr = strconv.Atoi(sessionId)

		if convErr != nil {

			fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for session ID.\n\n")

			continue

		}

		if !slices.Contains(sessionsIdSlice, tempInput) {

			fmt.Printf("\n\nERROR: Non-existing session ID. Please try again.\n\n")

			continue

		}

		isCustomerNameValid := false

		for !isCustomerNameValid {

			fmt.Print("Customer Name: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			customerName = strings.TrimSpace(scanner.Text())

			if strings.ToLower(customerName) == "debug back" {

				return nil

			}

			var customerNameErr string

			customerName, customerNameErr = customerNameValidator(nameRegex, customerName)

			if customerNameErr != "" {

				fmt.Printf("\n\n%s\n\n", customerNameErr)

				continue

			}

			isCustomerNameValid = true

		}

		isSeatCountValid := false

		for !isSeatCountValid {

			fmt.Print("Seat(s) Count: ")

			if !scanner.Scan() {

				return scanner.Err()

			}

			seatCount = strings.TrimSpace(scanner.Text())

			if strings.ToLower(seatCount) == "debug back" {

				return nil

			}

			var tempSeatCount int

			tempSeatCount, convErr = strconv.Atoi(seatCount)

			if convErr != nil {

				fmt.Printf("\n\nERROR: Invalid input. Expected numeric input for seat(s) count.\n\n")

				continue

			}

			if tempSeatCount <= 0 {

				fmt.Printf("\n\nERROR: Invalid seat(s) count.\n\n")

				continue

			}

			isSeatCountValid = true

		}

		trx, beginErr := db.Begin()

		if beginErr != nil {

			return beginErr

		}

		_, execErr := trx.Exec("INSERT INTO bookings (session_id, customer_name, seat_count, created_at) VALUES ($1, $2, $3, $4);", sessionId, customerName, seatCount, time.Now().Format("2006-01-02 15:04:05"))

		if execErr != nil {

			if rollbackErr := trx.Rollback(); rollbackErr != nil {

				return rollbackErr

			}

			return execErr

		}

		if commitErr := trx.Commit(); commitErr != nil {

			return commitErr

		}

		fmt.Printf("SUCCESS: Record has been inserted.\n\n")

		fmt.Println("Press Enter to continue..")

		if !scanner.Scan() {

			return scanner.Err()

		}

		isCompleted = true

	}

	return nil

}

func listMovies(db *sql.DB) (map[int]map[string]string, []int, error) {

	keys := make([]int, 0)

	bufferRows, queryErr := db.Query("SELECT * FROM movies;")

	if queryErr != nil {

		return nil, nil, queryErr

	}

	var id int
	var title, rating, duration string
	resultMap := make(map[int]map[string]string)

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id, &title, &rating, &duration); scanErr != nil {
			return nil, nil, scanErr
		}

		if _, exists := resultMap[id]; !exists {
			resultMap[id] = make(map[string]string)
		}

		resultMap[id]["title"] = title
		resultMap[id]["rating"] = rating
		resultMap[id]["duration"] = duration

		keys = append(keys, id)

	}

	return resultMap, keys, nil

}

func listTheaters(db *sql.DB) (map[int]map[string]string, []int, error) {

	keys := make([]int, 0)

	bufferRows, queryErr := db.Query("SELECT * FROM theaters;")

	if queryErr != nil {

		return nil, nil, queryErr

	}

	var id int
	var totalCapacity, isActive string
	resultMap := make(map[int]map[string]string)

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id, &totalCapacity, &isActive); scanErr != nil {
			return nil, nil, scanErr
		}

		if _, exists := resultMap[id]; !exists {
			resultMap[id] = make(map[string]string)
		}

		resultMap[id]["total_capacity"] = totalCapacity
		resultMap[id]["is_active"] = isActive

		keys = append(keys, id)

	}

	return resultMap, keys, nil

}

func listSessions(db *sql.DB) (map[int]map[string]string, []int, error) {

	keys := make([]int, 0)

	bufferRows, queryErr := db.Query("SELECT * FROM sessions;")

	if queryErr != nil {

		return nil, nil, queryErr

	}

	var id int
	var movieId, theaterId, availableSeats, isActive string
	resultMap := make(map[int]map[string]string)

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id, &movieId, &theaterId, &availableSeats, &isActive); scanErr != nil {
			return nil, nil, scanErr
		}

		if _, exists := resultMap[id]; !exists {
			resultMap[id] = make(map[string]string)
		}

		resultMap[id]["movie_id"] = movieId
		resultMap[id]["theater_id"] = theaterId
		resultMap[id]["available_seats"] = availableSeats
		resultMap[id]["is_active"] = isActive

		keys = append(keys, id)

	}

	return resultMap, keys, nil

}

func listBookings(db *sql.DB) (map[int]map[string]string, []int, error) {

	keys := make([]int, 0)

	bufferRows, queryErr := db.Query("SELECT * FROM bookings;")

	if queryErr != nil {

		return nil, nil, queryErr

	}

	var id int
	var sessionId, customerName, seatCount, createdAt string
	resultMap := make(map[int]map[string]string)

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id, &sessionId, &customerName, &seatCount, &createdAt); scanErr != nil {
			return nil, nil, scanErr
		}

		if _, exists := resultMap[id]; !exists {
			resultMap[id] = make(map[string]string)
		}

		resultMap[id]["session_id"] = sessionId
		resultMap[id]["customer_name"] = customerName
		resultMap[id]["seat_count"] = seatCount
		resultMap[id]["created_at"] = createdAt

		keys = append(keys, id)

	}

	return resultMap, keys, nil

}

func updateMovie(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var title, rating, duration, selectedRecord, userInput string

	mr.additionalMsg = ""
	mr.debugMsg = debugMsg

	var id, idx int
	var convErr error

	mr.columnRows, mr.columnRowsKeys, mr.listColumnsErr = listMovies(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListMovies(mr.columnRows)

	columns, columnsErr := listMoviesColumns(db)

	if columnsErr != nil {

		return columnsErr

	}

	keys := make([]string, 0, len(columns))

	for k := range columns {

		keys = append(keys, k)

	}

	slices.SortFunc(keys, func(a, b string) int {

		return cmp.Compare(a, b)

	})

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "UPDATE MOVIE"

		if renderMenuNonFilterErr := mr.renderMenuNonFilter(); renderMenuNonFilterErr != nil {

			return renderMenuNonFilterErr

		}

		fmt.Print("Movie ID: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		if strings.HasPrefix(userInput, "filter") {

			bufferUserInput := strings.Split(userInput, " ")

			if len(bufferUserInput) != 2 {

				mr.debugMsg = "ERROR: Invalid input. Refer to the given GUIDE below."

				continue

			}

			var movieRowsFilterErr error

			mr.convertedColumnRows, mr.additionalMsg, movieRowsFilterErr = movieRowsFilter(db, bufferUserInput, mr.columnRows, mr.columnRowsKeys)

			if movieRowsFilterErr != nil {

				return movieRowsFilterErr

			}

			continue

		}

		if userInput == "reset" {

			mr.additionalMsg = ""

			mr.convertedColumnRows = convertedListMovies(mr.columnRows)

			continue

		}

		mr.additionalMsg = ""

		idx, convErr = strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for movie ID. Refer to the given GUIDE below."

			continue

		}

		isCompleted = true

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

	}

	isColumnPicking := true

	mr.debugMsg = ""

	for isColumnPicking {

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- PICK MOVIE COLUMNS ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM movies WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &title, &rating, &duration); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		fmt.Printf("Selected record:\n\n")

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nTitle\t\t: %s\nRating\t\t: %s\nDuration\t: %s\n\n", id, title, rating, duration)

		fmt.Printf("%s\n\n", selectedRecord)

		fmt.Printf("Available columns:\n\n")

		for i, key := range keys {

			fmt.Printf("%d. %s\n", i+1, key)

		}

		fmt.Printf("\n%s\n\n", debugMsg)
		fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n\n")
		fmt.Print("Input column number: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		switch userInput {

		case "1":

			if updateMovieErr := updateMovieColumn(db, scanner, columns, "duration", idx); updateMovieErr != "" {

				mr.debugMsg = updateMovieErr

				continue

			}

		case "2":

			if updateMovieErr := updateMovieColumn(db, scanner, columns, "movie_title", idx); updateMovieErr != "" {

				mr.debugMsg = updateMovieErr

				continue

			}

		case "3":

			if updateMovieErr := updateMovieColumn(db, scanner, columns, "rating", idx); updateMovieErr != "" {

				mr.debugMsg = updateMovieErr

				continue

			}

		default:

			mr.debugMsg = "ERROR: Invalid input. Refer to the listed columns above."

			continue

		}

		isColumnPicking = false

	}

	mr.debugMsg = ""

	return nil

}

// Checkpoint

func updateTheater(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var totalCapacity, isActive, selectedRecord, userInput string
	var id, idx int
	var convErr error

	mr.debugMsg = debugMsg

	mr.columnRows, _, mr.listColumnsErr = listTheaters(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListTheaters(mr.columnRows)

	columns, columnsErr := listTheatersColumns(db)

	if columnsErr != nil {

		return columnsErr

	}

	keys := make([]string, 0, len(columns))

	for k := range columns {

		keys = append(keys, k)

	}

	slices.SortFunc(keys, func(a, b string) int {

		return cmp.Compare(a, b)

	})

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "UPDATE THEATER"

		if renderMenuNonFilterErr := mr.renderMenuNonFilter(); renderMenuNonFilterErr != nil {

			return renderMenuNonFilterErr

		}

		fmt.Print("Theater ID: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		idx, convErr = strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for theater ID. Refer to the given GUIDE below."

			continue

		}

		isCompleted = true

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

	}

	isColumnPicking := true

	mr.debugMsg = ""

	for isColumnPicking {

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- PICK THEATER COLUMNS ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM theaters WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &totalCapacity, &isActive); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		fmt.Printf("Selected record:\n\n")

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nTotal Capacity\t: %s\nIs Active\t: %s\n\n", id, totalCapacity, isActive)

		fmt.Printf("%s\n\n", selectedRecord)

		fmt.Printf("Available columns:\n\n")

		for i, key := range keys {

			fmt.Printf("%d. %s\n", i+1, key)

		}

		fmt.Printf("\n%s\n\n", debugMsg)
		fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n\n")
		fmt.Print("Input column number: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		switch userInput {

		case "1":

			if updateTheaterErr := updateTheaterColumn(db, scanner, columns, "is_active", idx); updateTheaterErr != "" {

				mr.debugMsg = updateTheaterErr

				continue

			}

		case "2":

			if updateTheaterErr := updateTheaterColumn(db, scanner, columns, "total_capacity", idx); updateTheaterErr != "" {

				mr.debugMsg = updateTheaterErr

				continue

			}

		default:

			fmt.Println("ERROR: Invalid input. Refer to the listed columns above.")

			continue

		}

		isColumnPicking = false

	}

	mr.debugMsg = ""

	return nil

}

func updateSession(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var movieId, theaterId, availableSeats, isActive, selectedRecord, userInput string
	var id, idx int
	var convErr error

	mr.debugMsg = debugMsg

	mr.columnRows, _, mr.listColumnsErr = listSessions(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListSessions(mr.columnRows)

	columns, columnsErr := listSessionsColumns(db)

	if columnsErr != nil {

		return columnsErr

	}

	keys := make([]string, 0, len(columns))

	for k := range columns {

		keys = append(keys, k)

	}

	slices.SortFunc(keys, func(a, b string) int {

		return cmp.Compare(a, b)

	})

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "UPDATE SESSION"

		if renderMenuNonFilterErr := mr.renderMenuNonFilter(); renderMenuNonFilterErr != nil {

			return renderMenuNonFilterErr

		}

		fmt.Print("Session ID: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		idx, convErr = strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for session ID. Refer to the given GUIDE below."

			continue

		}

		isCompleted = true

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

	}

	isColumnPicking := true

	mr.debugMsg = ""

	for isColumnPicking {

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- PICK SESSION COLUMNS ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM theaters WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &movieId, &theaterId, &availableSeats, &isActive); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		fmt.Printf("Selected record:\n\n")

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nMovie ID\t: %s\nTheater ID\t: %s\nAvailable Seats\t: %s\nIs Active\t: %s\n\n", id, movieId, theaterId, availableSeats, isActive)

		fmt.Printf("%s\n\n", selectedRecord)

		fmt.Printf("Available columns:\n\n")

		for i, key := range keys {

			fmt.Printf("%d. %s\n", i+1, key)

		}

		fmt.Printf("\n%s\n\n", debugMsg)
		fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n\n")
		fmt.Print("Input column number: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		switch userInput {

		case "1":

			if updateSessionErr := updateSessionColumn(db, scanner, columns, "available_seats", idx); updateSessionErr != "" {

				mr.debugMsg = updateSessionErr

				continue

			}

		case "2":

			if updateSessionErr := updateSessionColumn(db, scanner, columns, "movie_id", idx); updateSessionErr != "" {

				mr.debugMsg = updateSessionErr

				continue

			}

		case "3":

			if updateSessionErr := updateSessionColumn(db, scanner, columns, "is_active", idx); updateSessionErr != "" {

				mr.debugMsg = updateSessionErr

				continue

			}

		case "4":

			if updateSessionErr := updateSessionColumn(db, scanner, columns, "theater_id", idx); updateSessionErr != "" {

				mr.debugMsg = updateSessionErr

				continue

			}

		default:

			mr.debugMsg = "ERROR: Invalid input. Refer to the listed columns above."

			continue

		}

		isColumnPicking = false

	}

	mr.debugMsg = ""

	return nil

}

func updateBooking(db *sql.DB, scanner *bufio.Scanner, mr *menuRenderer, debugMsg string) error {

	var sessionId, customerName, seatCount, createdAt, selectedRecord, userInput string
	var id, idx int
	var convErr error

	mr.additionalMsg = ""
	mr.debugMsg = debugMsg

	mr.columnRows, mr.columnRowsKeys, mr.listColumnsErr = listBookings(db)

	if mr.listColumnsErr != nil {

		return mr.listColumnsErr

	}

	mr.convertedColumnRows = convertedListBookings(mr.columnRows)

	columns, columnsErr := listBookingsColumns(db)

	if columnsErr != nil {

		return columnsErr

	}

	keys := make([]string, 0, len(columns))

	for k := range columns {

		keys = append(keys, k)

	}

	slices.SortFunc(keys, func(a, b string) int {

		return cmp.Compare(a, b)

	})

	isCompleted := false

	for !isCompleted {

		mr.sectionHeader = "UPDATE BOOKING"

		if renderMenuFilterErr := mr.renderMenuFilter(); renderMenuFilterErr != nil {

			return renderMenuFilterErr

		}

		fmt.Print("Booking ID: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		if strings.HasPrefix(userInput, "filter") {

			bufferUserInput := strings.Split(userInput, " ")

			if len(bufferUserInput) != 2 {

				mr.debugMsg = "ERROR: Invalid input. Refer to the given GUIDE below."

				continue

			}

			var bookingRowsFilterErr error

			mr.convertedColumnRows, mr.additionalMsg, bookingRowsFilterErr = bookingRowsFilter(db, bufferUserInput, mr.columnRows, mr.columnRowsKeys)

			if bookingRowsFilterErr != nil {

				return bookingRowsFilterErr

			}

			continue

		}

		if userInput == "reset" {

			mr.additionalMsg = ""

			mr.convertedColumnRows = convertedListBookings(mr.columnRows)

			continue

		}

		mr.additionalMsg = ""

		idx, convErr = strconv.Atoi(userInput)

		if convErr != nil {

			mr.debugMsg = "ERROR: Invalid input. Expected numeric input for booking ID. Refer to the given GUIDE below."

			continue

		}

		isCompleted = true

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

	}

	isColumnPicking := true

	mr.debugMsg = ""

	for isColumnPicking {

		if clearScreenErr := clearScreen(); clearScreenErr != nil {

			return clearScreenErr

		}

		fmt.Printf("---------------------------------------- PICK BOOKING COLUMNS ---------------------------------------\n\n")

		bufferQueryRow := db.QueryRow("SELECT * FROM bookings WHERE id = $1;", idx)

		if scanErr := bufferQueryRow.Scan(&id, &sessionId, &customerName, &seatCount, &createdAt); scanErr != nil {

			mr.debugMsg = "ERROR: Invalid ID. Please try again."

			continue

		}

		fmt.Printf("Selected record:\n\n")

		selectedRecord = fmt.Sprintf("ID\t\t: %d\nSession ID\t: %s\nCustomer Name\t: %s\nSeat Count\t: %s\nCreated At\t: %s\n\n", id, sessionId, customerName, seatCount, createdAt)

		fmt.Printf("%s\n\n", selectedRecord)

		fmt.Printf("Available columns:\n\n")

		for i, key := range keys {

			fmt.Printf("%d. %s\n", i+1, key)

		}

		fmt.Printf("\n%s\n\n", debugMsg)
		fmt.Printf("GUIDE: Type \"Back\" to go back to the previous menu.\n\n")
		fmt.Print("Input column number: ")

		mr.debugMsg = ""

		if !scanner.Scan() {

			return scanner.Err()

		}

		userInput = strings.TrimSpace(strings.ToLower(scanner.Text()))

		if userInput == "back" {

			return nil

		}

		switch userInput {

		case "1":

			if updateBookingErr := updateBookingColumn(db, scanner, columns, "created_at", idx); updateBookingErr != "" {

				mr.debugMsg = updateBookingErr

				continue

			}

		case "2":

			if updateBookingErr := updateBookingColumn(db, scanner, columns, "customer_name", idx); updateBookingErr != "" {

				mr.debugMsg = updateBookingErr

				continue

			}

		case "3":

			if updateBookingErr := updateBookingColumn(db, scanner, columns, "seat_count", idx); updateBookingErr != "" {

				mr.debugMsg = updateBookingErr

				continue

			}

		case "4":

			if updateBookingErr := updateBookingColumn(db, scanner, columns, "session_id", idx); updateBookingErr != "" {

				mr.debugMsg = updateBookingErr

				continue

			}

		default:

			mr.debugMsg = "ERROR: Invalid input. Refer to the listed columns above."

			continue

		}

		isColumnPicking = false

	}

	mr.debugMsg = ""

	return nil

}

func convertedListMovies(listMovies map[int]map[string]string) string {

	var result string

	for id, movie := range listMovies {

		result += fmt.Sprintf("ID\t\t: %d\nTitle\t\t: %s\nRating\t\t: %s\nDuration\t: %s\n\n", id, movie["title"], movie["rating"], movie["duration"])

	}

	return result

}

func convertedListTheaters(listTheaters map[int]map[string]string) string {

	var result string

	for id, theater := range listTheaters {

		result += fmt.Sprintf("ID\t\t: %d\nTotal Capacity\t: %s\nIs Active\t: %s\n\n", id, theater["total_capacity"], theater["is_active"])

	}

	return result

}

func convertedListSessions(listSessions map[int]map[string]string) string {

	var result string

	for id, session := range listSessions {

		result += fmt.Sprintf("ID\t\t: %d\nMovie ID\t: %s\nTheater ID\t: %s\nAvailable Seats\t: %s\nIs Active\t: %s\n\n", id, session["movie_id"], session["theater_id"], session["available_seats"], session["is_active"])

	}

	return result

}

func convertedListBookings(listBookings map[int]map[string]string) string {

	var result string

	for id, booking := range listBookings {

		result += fmt.Sprintf("ID\t\t: %d\nSession ID\t: %s\nCustomer Name\t: %s\nSeat Count\t: %s\nCreated At\t: %s\n\n", id, booking["session_id"], booking["customer_name"], booking["seat_count"], booking["created_at"])

	}

	return result

}

func mainMenuList() []string {

	listSlice := []string{
		"Administration",
		"Ticket Booking",
		"Exit",
	}

	return listSlice

}

func adminMenuList() []string {

	listSlice := []string{
		"Delete",
		"Insert",
		"List",
		"Update",
		"Back",
		"Exit",
	}

	return listSlice

}

func tablesList() []string {

	listSlice := []string{
		"Movie",
		"Theater",
		"Session",
		"Booking",
		"Back",
		"Exit",
	}

	return listSlice

}

func movieRowsFilter(db *sql.DB, bufferUserInput []string, movieRows map[int]map[string]string, movieRowsKeys []int) (string, string, error) {

	var id, title, rating, duration string
	keyword := strings.ToLower(bufferUserInput[1])
	filteredMovieRows := make(map[int]map[string]string)

	for _, idx := range movieRowsKeys {

		titleTemp := strings.ToLower(movieRows[idx]["title"])

		if strings.Contains(titleTemp, keyword) {

			bufferMovieRows := db.QueryRow("SELECT * FROM movies WHERE id = $1;", idx)

			if scanErr := bufferMovieRows.Scan(&id, &title, &rating, &duration); scanErr != nil {

				return "", "", scanErr

			}

			filteredMovieRows[idx] = map[string]string{
				"title":    title,
				"rating":   rating,
				"duration": duration,
			}

		}

	}

	convertedMovieRows := convertedListMovies(filteredMovieRows)

	additionalMsg := fmt.Sprintf("Filter result for \"%s\":\n\n", keyword)

	return convertedMovieRows, additionalMsg, nil

}

func bookingRowsFilter(db *sql.DB, bufferUserInput []string, bookingRows map[int]map[string]string, bookingRowsKeys []int) (string, string, error) {

	// Continue from here

	var id, sessionId, customerName, seatCount, createdAt string
	keyword := strings.ToLower(bufferUserInput[1])
	filteredBookingRows := make(map[int]map[string]string)

	for _, idx := range bookingRowsKeys {

		customerTemp := strings.ToLower(bookingRows[idx]["customer_name"])

		if strings.Contains(customerTemp, keyword) {

			bufferBookingRows := db.QueryRow("SELECT * FROM bookings WHERE id = $1;", idx)

			if scanErr := bufferBookingRows.Scan(&id, &sessionId, &customerName, &seatCount, &createdAt); scanErr != nil {

				return "", "", scanErr

			}

			filteredBookingRows[idx] = map[string]string{
				"session_id":    sessionId,
				"customer_name": customerName,
				"seat_count":    seatCount,
				"created_at":    createdAt,
			}

		}

	}

	convertedBookingRows := convertedListMovies(filteredBookingRows)

	additionalMsg := fmt.Sprintf("Filter result for \"%s\":\n\n", keyword)

	return convertedBookingRows, additionalMsg, nil

}

func getMoviesId(db *sql.DB) ([]int, error) {

	bufferRows, queryErr := db.Query("SELECT id FROM movies;")

	if queryErr != nil {

		return nil, queryErr

	}

	var id int
	var result []int

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id); scanErr != nil {
			return nil, scanErr
		}

		result = append(result, id)

	}

	return result, nil

}

func getTheatersId(db *sql.DB) ([]int, error) {

	bufferRows, queryErr := db.Query("SELECT id FROM theaters;")

	if queryErr != nil {

		return nil, queryErr

	}

	var id int
	var result []int

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id); scanErr != nil {
			return nil, scanErr
		}

		result = append(result, id)

	}

	return result, nil

}

func getSessionsId(db *sql.DB) ([]int, error) {

	bufferRows, queryErr := db.Query("SELECT id FROM sessions;")

	if queryErr != nil {

		return nil, queryErr

	}

	var id int
	var result []int

	for bufferRows.Next() {

		if scanErr := bufferRows.Scan(&id); scanErr != nil {
			return nil, scanErr
		}

		result = append(result, id)

	}

	return result, nil

}

func listMoviesColumns(db *sql.DB) (map[string]string, error) {

	columns := make(map[string]string, 0)

	bufferRows, queryErr := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'movies' AND column_name != 'id';")

	if queryErr != nil {

		return nil, queryErr

	}

	for bufferRows.Next() {

		var columnName, columnType string

		if scanErr := bufferRows.Scan(&columnName, &columnType); scanErr != nil {

			return nil, scanErr

		}

		columns[columnName] = columnType

	}

	return columns, nil

}

func listTheatersColumns(db *sql.DB) (map[string]string, error) {

	columns := make(map[string]string, 0)

	bufferRows, queryErr := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'theaters' AND column_name != 'id';")

	if queryErr != nil {

		return nil, queryErr

	}

	for bufferRows.Next() {

		var columnName, columnType string

		if scanErr := bufferRows.Scan(&columnName, &columnType); scanErr != nil {

			return nil, scanErr

		}

		columns[columnName] = columnType

	}

	return columns, nil

}

func listSessionsColumns(db *sql.DB) (map[string]string, error) {

	columns := make(map[string]string, 0)

	bufferRows, queryErr := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'sessions' AND column_name != 'id';")

	if queryErr != nil {

		return nil, queryErr

	}

	for bufferRows.Next() {

		var columnName, columnType string

		if scanErr := bufferRows.Scan(&columnName, &columnType); scanErr != nil {

			return nil, scanErr

		}

		columns[columnName] = columnType

	}

	return columns, nil

}

func listBookingsColumns(db *sql.DB) (map[string]string, error) {

	columns := make(map[string]string, 0)

	bufferRows, queryErr := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'bookings' AND column_name != 'id';")

	if queryErr != nil {

		return nil, queryErr

	}

	for bufferRows.Next() {

		var columnName, columnType string

		if scanErr := bufferRows.Scan(&columnName, &columnType); scanErr != nil {

			return nil, scanErr

		}

		columns[columnName] = columnType

	}

	return columns, nil

}

func updateMovieColumn(db *sql.DB, scanner *bufio.Scanner, column map[string]string, columnName string, idx int) string {

	columnType := strings.ToLower(column[columnName])

	fmt.Printf("\nEnter new value for column %s: ", columnName)

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	newValue := strings.TrimSpace(scanner.Text())

	if strings.Contains(columnType, "int") {

		if _, convErr := strconv.Atoi(newValue); convErr != nil || newValue == "" {

			return "ERROR: Invalid input. Expected non-empty numeric input."

		}

		if updateTransErr := updateMovieColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	} else if strings.Contains(columnType, "char") || strings.Contains(columnType, "text") {

		if len(newValue) == 0 {

			return "ERROR: Invalid input. Expected non-empty alphanumeric input."

		}

		if updateTransErr := updateMovieColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	} else if strings.Contains(columnType, "float") || strings.Contains(columnType, "double") {

		bufferNewValue, convErr := strconv.ParseFloat(newValue, 64)

		if convErr != nil || newValue == "" {

			return "ERROR: Invalid input. Expected non-empty numeric input in float format (2 decimal places)."

		}

		newValue = fmt.Sprintf("%.1f", bufferNewValue)

		if updateTransErr := updateMovieColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	}

	return ""

}

func updateTheaterColumn(db *sql.DB, scanner *bufio.Scanner, column map[string]string, columnName string, idx int) string {

	columnType := strings.ToLower(column[columnName])

	fmt.Printf("\nEnter new value for column %s: ", columnName)

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	newValue := strings.TrimSpace(scanner.Text())

	if strings.Contains(columnType, "int") {

		if _, convErr := strconv.Atoi(newValue); convErr != nil || newValue == "" {

			return "ERROR: Invalid input. Expected non-empty numeric input."

		}

		if updateTransErr := updateTheaterColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	} else if strings.Contains(columnType, "bool") {

		if _, convErr := strconv.ParseBool(newValue); convErr != nil {

			return "ERROR: Invalid input. Expected \"true\" or \"false\"."

		}

		if updateTransErr := updateTheaterColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	}

	return ""

}

func updateSessionColumn(db *sql.DB, scanner *bufio.Scanner, column map[string]string, columnName string, idx int) string {

	columnType := strings.ToLower(column[columnName])

	fmt.Printf("\nEnter new value for column %s: ", columnName)

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	newValue := strings.TrimSpace(scanner.Text())

	if strings.Contains(columnType, "int") {

		if _, convErr := strconv.Atoi(newValue); convErr != nil || newValue == "" {

			return "ERROR: Invalid input. Expected non-empty numeric input."

		}

		if updateTransErr := updateSessionColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	} else if strings.Contains(columnType, "bool") {

		if _, convErr := strconv.ParseBool(newValue); convErr != nil {

			return "ERROR: Invalid input. Expected \"true\" or \"false\"."

		}

		if updateTransErr := updateSessionColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	}

	return ""

}

func updateBookingColumn(db *sql.DB, scanner *bufio.Scanner, column map[string]string, columnName string, idx int) string {

	columnType := strings.ToLower(column[columnName])

	fmt.Printf("\nEnter new value for column %s: ", columnName)

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	newValue := strings.TrimSpace(scanner.Text())

	if strings.Contains(columnType, "int") {

		if _, convErr := strconv.Atoi(newValue); convErr != nil || newValue == "" {

			return "ERROR: Invalid input. Expected non-empty numeric input."

		}

		if updateTransErr := updateBookingColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	} else if strings.Contains(columnType, "char") || strings.Contains(columnType, "text") || strings.Contains(columnType, "time") {

		if len(newValue) == 0 {

			return "ERROR: Invalid input. Expected non-empty alphanumeric input."

		}

		if updateTransErr := updateBookingColumnTrans(db, scanner, columnName, newValue, idx); updateTransErr != "" {

			return updateTransErr

		}

	}

	return ""

}

func updateMovieColumnTrans(db *sql.DB, scanner *bufio.Scanner, columnName string, newValue string, idx int) string {

	trx, beginErr := db.Begin()

	if beginErr != nil {

		return "ERROR: Something wrong with the transaction begin."

	}

	query := fmt.Sprintf("UPDATE movies SET %s = $1 WHERE id = $2;", columnName)

	_, execErr := trx.Exec(query, newValue, idx)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return "ERROR: Something wrong with the transaction rollback."

		}

		return "ERROR: Something wrong with the query."

	}

	if commitErr := trx.Commit(); commitErr != nil {

		return "ERROR: Something wrong with the transaction commit."

	}

	fmt.Printf("SUCCESS: Record has been updated.\n\n")

	fmt.Print("Press Enter to continue...")

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	return ""

}

func updateTheaterColumnTrans(db *sql.DB, scanner *bufio.Scanner, columnName string, newValue string, idx int) string {

	trx, beginErr := db.Begin()

	if beginErr != nil {

		return "ERROR: Something wrong with the transaction begin."

	}

	query := fmt.Sprintf("UPDATE theaters SET %s = $1 WHERE id = $2;", columnName)

	_, execErr := trx.Exec(query, newValue, idx)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return "ERROR: Something wrong with the transaction rollback."

		}

		return "ERROR: Something wrong with the query."

	}

	if commitErr := trx.Commit(); commitErr != nil {

		return "ERROR: Something wrong with the transaction commit."

	}

	fmt.Printf("SUCCESS: Record has been updated.\n\n")

	fmt.Print("Press Enter to continue...")

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	return ""

}

func updateSessionColumnTrans(db *sql.DB, scanner *bufio.Scanner, columnName string, newValue string, idx int) string {

	trx, beginErr := db.Begin()

	if beginErr != nil {

		return "ERROR: Something wrong with the transaction begin."

	}

	query := fmt.Sprintf("UPDATE sessions SET %s = $1 WHERE id = $2;", columnName)

	_, execErr := trx.Exec(query, newValue, idx)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return "ERROR: Something wrong with the transaction rollback."

		}

		return "ERROR: Something wrong with the query."

	}

	if commitErr := trx.Commit(); commitErr != nil {

		return "ERROR: Something wrong with the transaction commit."

	}

	fmt.Printf("SUCCESS: Record has been updated.\n\n")

	fmt.Print("Press Enter to continue...")

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	return ""

}

func updateBookingColumnTrans(db *sql.DB, scanner *bufio.Scanner, columnName string, newValue string, idx int) string {

	trx, beginErr := db.Begin()

	if beginErr != nil {

		return "ERROR: Something wrong with the transaction begin."

	}

	query := fmt.Sprintf("UPDATE bookings SET %s = $1 WHERE id = $2;", columnName)

	_, execErr := trx.Exec(query, newValue, idx)

	if execErr != nil {

		if rollbackErr := trx.Rollback(); rollbackErr != nil {

			return "ERROR: Something wrong with the transaction rollback."

		}

		return "ERROR: Something wrong with the query.\nGUIDE: Follow \"YYYY-MM-DD HH:MM:SS\" format for timestamp value.\n\n"

	}

	if commitErr := trx.Commit(); commitErr != nil {

		return "ERROR: Something wrong with the transaction commit."

	}

	fmt.Printf("SUCCESS: Record has been updated.\n\n")

	fmt.Print("Press Enter to continue...")

	if !scanner.Scan() {

		return "ERROR: Something wrong with the scanner."

	}

	return ""

}
