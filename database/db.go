package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"project-root/db"
	"project-root/models"
	"time"
)

func main() {
	// Check if the database file already exists
	dbFile := "./platform_certificates.db"
	if _, err := os.Stat(dbFile); err == nil {
		fmt.Println("Existing database found. Resuming operations...")
	} else {
		fmt.Println("No existing database found. Creating a new database...")
	}

	// Initialize the SQLiteRepository
	repo, err := db.NewSQLiteRepository(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Database setup complete. Waiting for instructions...")

	// Enter command loop to wait for instructions
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command (insert, check, get, exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Process commands
		switch {
		case strings.HasPrefix(input, "insert"):
			handleInsert(repo, input)

		case strings.HasPrefix(input, "check"):
			handleCheck(repo, input)

		case strings.HasPrefix(input, "get"):
			handleGet(repo)

		case input == "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command. Valid commands are: insert, check, get, exit.")
		}
	}
}

// handleInsert processes the insert command
func handleInsert(repo *db.SQLiteRepository, input string) {
	parts := strings.Fields(input)
	if len(parts) < 4 {
		fmt.Println("Usage: insert <serialNumber> <signer> <component1,component2,...>")
		return
	}

	serialNumber := parts[1]
	signer := parts[2]
	components := strings.Split(parts[3], ",")

	cert := models.Certificate{
		SerialNumber: serialNumber,
		Signer:       signer,
		Components:   components,
		IssueDate:    time.Now(),
		ExpiryDate:   time.Now().AddDate(1, 0, 0),
	}

	err := repo.InsertCertificate(cert)
	if err != nil {
		fmt.Printf("Failed to insert certificate: %v\n", err)
	} else {
		fmt.Println("Certificate inserted successfully.")
	}
}

// handleCheck processes the check command
func handleCheck(repo *db.SQLiteRepository, input string) {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		fmt.Println("Usage: check <serialNumber>")
		return
	}

	serialNumber := parts[1]
	exists, err := repo.CheckCertificateExists(serialNumber)
	if err != nil {
		fmt.Printf("Error checking certificate existence: %v\n", err)
	} else if exists {
		fmt.Printf("Certificate with serial number %s exists.\n", serialNumber)
	} else {
		fmt.Printf("Certificate with serial number %s does not exist.\n", serialNumber)
	}
}

// handleGet processes the get command
func handleGet(repo *db.SQLiteRepository) {
	certificates, err := repo.GetCertificates()
	if err != nil {
		fmt.Printf("Failed to retrieve certificates: %v\n", err)
		return
	}

	for _, cert := range certificates {
		fmt.Printf("Certificate: Serial=%s, Signer=%s, Components=%v\n", cert.SerialNumber, cert.Signer, cert.Components)
	}
}