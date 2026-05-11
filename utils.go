package main

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func customerNameAndTicketAmountValidator(nameRegex *regexp.Regexp, scanner *bufio.Scanner) (int, string, string) {

	// Split the input and attempt to parse the ticket amount from the last part.
	bufferInput := strings.Split(scanner.Text(), " ")

	ticketAmount, err := strconv.Atoi(bufferInput[len(bufferInput)-1])

	if err != nil {
		return 0, "", "ERROR: Invalid input. Refer to the given GUIDE below."
	}

	if ticketAmount <= 0 {
		return 0, "", "ERROR: Invalid ticket amount. Please try again."
	}

	// Join the remaining parts to form the customer name and validate it against the regex.
	customerName := strings.ToLower(strings.Join(bufferInput[:len(bufferInput)-1], " "))

	if nameRegex.FindString(customerName) == "" {
		return 0, "", "ERROR: Invalid customer name. Please try again."
	}

	return ticketAmount, customerName, ""
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
