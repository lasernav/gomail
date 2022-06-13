package gomail

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"io/ioutil"
	"net/http"
	"math/rand"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized \n")
	}
}

func SendEmailOAUTH2(to string, data interface{}, template string) (bool, error) {

	emailBody, err := parseTemplate(template, data)
	if err != nil {
		return false, errors.New("unable to parse email template")
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Test Email form Gmail API using OAuth" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

func SendEmail(mezzoTarga string, mezzoUrl string, emailSubject string, emailTo string) (bool, error) {

	var message gmail.Message

	now := time.Now().UTC().Add(2 * time.Hour)

	send_emailTo := "To: " + emailTo + "\r\n"
	send_emailSubject := "Subject: " + "[" + mezzoTarga + "] " + emailSubject + "\n"
	send_mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	send_emailBody := "<ul><li>"+ emailSubject +"<li>Targa mezzo: " + mezzoTarga +"</li><li>Timestamp: " + now.Format("2006-01-02 3:4:5 pm") + "</li></ul><h1><a href=" + mezzoUrl + " target='_blank'>Portale mezzo</a></h1>"

	msg := []byte(send_emailTo + send_emailSubject + send_mime + "\n" + send_emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

func SendEmailAttachment(mezzoTarga string, mezzoUrl string, emailSubject string, emailTo string, fileDir string, fileName string) (bool, error) {

	var message gmail.Message

	now := time.Now().UTC().Add(2 * time.Hour)

	fileBytes, err := ioutil.ReadFile(fileDir + fileName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fileMIMEType := http.DetectContentType(fileBytes)
	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32, "alphanum")

	messageBody := []byte(
		"Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"To: " + emailTo + "\n" +
		"Subject: " + "[" + mezzoTarga + "] " + emailSubject + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: text/html; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: quoted-printable\n\n" +
		"<ul><li>"+ emailSubject +"<li>Targa mezzo: " + mezzoTarga +"</li><li>Timestamp: " + now.Format("2006-01-02 3:4:5 pm") + "</li></ul><h1><a href=" + mezzoUrl + " target='_blank'>Portale mezzo</a></h1>" + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
	return true, nil
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}
