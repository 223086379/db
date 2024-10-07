package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Certificate represents the certificate data structure.
type Certificate struct {
	SerialNumber string
	Signer       string
	Components   string
}

// initializeDatabase sets up the database and schema.
func initializeDatabase(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS certificates (
        serial_number TEXT PRIMARY KEY,
        signer TEXT NOT NULL,
        components TEXT NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	return err
}

// addCertificate adds a new certificate to the database.
func addCertificate(db *sql.DB, cert Certificate) error {
	insertSQL := `
    INSERT INTO certificates (serial_number, signer, components)
    VALUES (?, ?, ?);`

	_, err := db.Exec(insertSQL, cert.SerialNumber, cert.Signer, cert.Components)
	return err
}

// certificateExists checks if a certificate exists by serial number.
func certificateExists(db *sql.DB, serialNumber string) (bool, error) {
	querySQL := `
    SELECT COUNT(1) FROM certificates WHERE serial_number = ?;`

	var count int
	err := db.QueryRow(querySQL, serialNumber).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// printUsage displays the usage information.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  To add a certificate:")
	fmt.Println("    -add -serial <serial_number> -signer <signer> -components <components>")
	fmt.Println("  To check if a certificate exists:")
	fmt.Println("    -check -serial <serial_number>")
	os.Exit(1)
}

// parseArguments parses and validates command-line arguments.
func parseArguments() (string, Certificate, error) {
	addCmd := flag.Bool("add", false, "Add a new certificate")
	checkCmd := flag.Bool("check", false, "Check if a certificate exists")
	serial := flag.String("serial", "", "Serial number of the certificate")
	signer := flag.String("signer", "", "Signer of the certificate")
	components := flag.String("components", "", "Components of the certificate")

	flag.Parse()

	if *addCmd && *checkCmd {
		return "", Certificate{}, errors.New("Cannot use both -add and -check at the same time")
	}

	if !*addCmd && !*checkCmd {
		return "", Certificate{}, errors.New("Must specify either -add or -check")
	}

	if *addCmd {
		if *serial == "" || *signer == "" || *components == "" {
			return "", Certificate{}, errors.New("Missing required fields for adding a certificate")
		}
		cert := Certificate{
			SerialNumber: *serial,
			Signer:       *signer,
			Components:   *components,
		}
		return "add", cert, nil
	}

	if *checkCmd {
		if *serial == "" {
			return "", Certificate{}, errors.New("Missing serial number for checking certificate")
		}
		cert := Certificate{
			SerialNumber: *serial,
		}
		return "check", cert, nil
	}

	return "", Certificate{}, errors.New("Invalid arguments")
}

func main() {
	// Parse command-line arguments
	action, cert, err := parseArguments()
	if err != nil {
		fmt.Println("Error:", err)
		printUsage()
	}

	// Open database connection
	db, err := sql.Open("sqlite3", "./certificates.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Initialize database schema
	err = initializeDatabase(db)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	switch action {
	case "add":
		err = addCertificate(db, cert)
		if err != nil {
			log.Fatal("Failed to add certificate:", err)
		}
		fmt.Println("Certificate added successfully.")
	case "check":
		exists, err := certificateExists(db, cert.SerialNumber)
		if err != nil {
			log.Fatal("Failed to check certificate:", err)
		}
		if exists {
			fmt.Printf("Certificate with serial number '%s' exists.\n", cert.SerialNumber)
		} else {
			fmt.Printf("Certificate with serial number '%s' does not exist.\n", cert.SerialNumber)
		}
	default:
		printUsage()
	}
}
