package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kristof1345/urlshort/db"
)

func main() {
	go db.Connect()

	app := fiber.New()

	app.Get("/:code")
	app.Post("/api/:shorten")

	log.Fatal(app.Listen(":5000"))
}
