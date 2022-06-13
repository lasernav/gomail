package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/BinodKafle/gomail/gomail"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	params := os.Args
	paramsLength := len(params)
	if paramsLength < 2 {
		log.Println("Please add SMTP or OAUTH along with go run main.go command")
		log.Println("Eg: go run main.go SMTP or go run main.go OAUTH")
		os.Exit(1)
	}

	inputMethod := os.Args[1]

	valid := IsValidInputMethod(inputMethod)

	emailTo := os.Getenv("EMAIL_TO")

	if valid {
		data := struct {
			ReceiverName string
			SenderName   string
		}{
			ReceiverName: "David Gilmour",
			SenderName:   "Binod Kafle",
		}

		if inputMethod == "SMTP" {
			status, err := gomail.SendEmailSMTP([]string{emailTo}, data, "sample_template.txt")
			if err != nil {
				log.Println(err)
			}
			if status {
				log.Println("Email sent successfully using SMTP")
			}
		}

		if inputMethod == "OAUTH" {
			gomail.OAuthGmailService()
			status, err := gomail.SendEmailOAUTH2(emailTo, data, "sample_template.txt")
			if err != nil {
				log.Println(err)
			}
			if status {
				log.Println("Email sent successfully using OAUTH")
			}
		}

		if inputMethod == "LSR" {
			gomail.OAuthGmailService()

			status, err := gomail.SendEmail("TARGA-MEZZO", "https://rp.lasernavigation.it:6014", "Camera Recorder START", emailTo)
			if err != nil {
				log.Println(err)
			}
			if status {
				log.Println("Email sent successfully using LSR")
			}
		}

		if inputMethod == "LSRATT" {
			gomail.OAuthGmailService()

			status, err := gomail.SendEmailAttachment("TARGA-MEZZO", "https://rp.lasernavigation.it:6014", "Invasione GEOFENCE", emailTo, "/home/laser/gomail/", "2022_06_09_08_19_43.jpg")
			if err != nil {
				log.Println(err)
			}
			if status {
				log.Println("Email sent successfully using LSRATT")
			}
		}
	} else {
		log.Println("Please add SMTP or OAUTH along with go run main.go command")
		log.Println("Eg: go run main.go SMTP or go run main.go OAUTH")
		os.Exit(1)
	}
}

func IsValidInputMethod(method string) bool {
	switch method {
	case
		"SMTP",
		"OAUTH",
		"LSR",
		"LSRATT":
		return true
	}
	return false
}
