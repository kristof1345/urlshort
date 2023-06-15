package routes

import (
	"fmt"
	"net/url"
	"os"

	"github.com/teris-io/shortid"

	"github.com/gofiber/fiber/v2"
)

type LongURL struct {
	LongURL string `json:"longUrl"`
}

func CodeHandler(c *fiber.Ctx) error {
	codeParam := c.Params("code")

	fmt.Println(codeParam)

	return c.Status(200).JSON(codeParam)
}

func ShortenUrlHandler(c *fiber.Ctx) error {
	UrlToShorten := new(LongURL)

	if err := c.BodyParser(UrlToShorten); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	baseUrl := os.Getenv("BASE_URL")
	fmt.Println(baseUrl)

	_, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return c.Status(401).JSON("Invalid base url")
	}

	urlCode, err := shortid.Generate()
	if err != nil {
		return c.Status(401).JSON("Failed to generate short id")
	}

	u, err := url.ParseRequestURI(UrlToShorten.LongURL)
	if err != nil {
		return c.Status(500).JSON("Invalid long url")
	}

	fmt.Println(u, urlCode)

	return c.Status(200).JSON(UrlToShorten)
}
