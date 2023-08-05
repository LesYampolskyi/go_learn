package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"web_dev/models"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	es, err := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	if err != nil {
		panic(err)
	}

	err = es.ForgotPassword("test@mail.com", "https://some.site.com/reset-pw?token=abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent")

	// msg := mail.NewMessage()
	// msg.SetHeader("To", to)
	// msg.SetHeader("From", from)
	// msg.SetHeader("Subject", subject)
	// msg.SetBody("text/plain", plaintext)
	// msg.AddAlternative("text/html", html)
	// msg.WriteTo(os.Stdout)

	// dialer := mail.NewDialer(host, port, username, password)
	// err := dialer.DialAndSend(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("message sent")
}
