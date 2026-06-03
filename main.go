package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	mr := &menuRenderer{}

	nameRegex := regexp.MustCompile(`^[[:alpha:]|\s]*$`)

	// Initialize the database connection and ensure it closes on exit.
	db, dbInitErr := dbInitializator()

	if dbInitErr != nil {
		log.Fatal(dbInitErr)
	}

	defer db.Close()

	// Initialize the input scanner.
	scanner, scannerInitErr := scannerInitializator()

	if scannerInitErr != nil {
		log.Fatal(scannerInitErr)
	}

	// Run the main menu loop and handle any fatal errors.
	if mainMenuErr := mainMenu(db, scanner, nameRegex, mr); mainMenuErr != nil {
		log.Fatal(mainMenuErr)
	}

	// Print termination message and exit the program.
	fmt.Println("Program is terminated.")

	os.Exit(0)
}
