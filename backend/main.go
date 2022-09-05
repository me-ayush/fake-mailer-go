package main

import (
	"fake-mailer/models"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Post("/", createMail)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).Render("index", fiber.Map{
			// "Msg": "404 Not Found",
		})
	})
}

func main() {
	err := godotenv.Load()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = ":3000"
	}

	if err != nil {
		fmt.Println(err)
	}

	engine := html.New(".", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	setupRoutes(app)
	log.Fatal(app.Listen(PORT))
}

func createMail(c *fiber.Ctx) error {
	var mail models.Mailer

	if err := c.BodyParser(&mail); err != nil {
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON in login"})
		return c.Status(fiber.StatusOK).Render("index", fiber.Map{
			"Success": true,
			"Message": "Cannot Parse JSON",
		})
	}

	if mail.Sender == "" || mail.To == nil {
		// return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "Please Fill All Fields"})
		return c.Status(fiber.StatusOK).Render("index", fiber.Map{
			"Success": true,
			"Message": "Please Fill All Fields",
		})
	}

	msg := sendMail(mail)

	return c.Status(fiber.StatusOK).Render("index", fiber.Map{
		"Success": true,
		"Message": msg,
	})
}

func sendMail(request models.Mailer) string {
	_ = godotenv.Load()

	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASS_WORD")
	smtpHost := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", username, password, smtpHost)

	sender := request.Sender
	sent_by := request.SenderName
	to := request.To

	message := buildMessage(request, sent_by)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, []byte(message))
	if err != nil {
		return err.Error()
	}

	return "Message Was Sent"
}

func buildMessage(mail models.Mailer, name string) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s %s\r\n", name, mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}
