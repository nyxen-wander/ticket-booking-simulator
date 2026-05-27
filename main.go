package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	var debugMsg string
	mr := &menuRenderer{}

	nameRegex := regexp.MustCompile(`^[[:alpha:]|\s]*$`)

	db, dbInitErr := dbInitializator()

	if dbInitErr != nil {
		log.Fatal(dbInitErr)
	}

	defer db.Close()

	scanner, scannerInitErr := scannerInitializator()

	if scannerInitErr != nil {
		log.Fatal(scannerInitErr)
	}

	if mainMenuErr := mainMenu(db, scanner, nameRegex, mr, debugMsg); mainMenuErr != nil {
		log.Fatal(mainMenuErr)
	}

	fmt.Println("Program is terminated.")

	os.Exit(0)
}
