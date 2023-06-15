package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kristof1345/urlshort/db"
	"github.com/kristof1345/urlshort/routes"
)

func main() {
	go db.Connect()

	app := fiber.New()

	app.Get("/:code", routes.CodeHandler)
	app.Post("/api/shorten", routes.ShortenUrlHandler)

	log.Fatal(app.Listen(":5000"))
}
