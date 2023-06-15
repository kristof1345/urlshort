package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CodeHandler(c *fiber.Ctx) error {
	codeParam := c.Params("code")

	fmt.Println(codeParam)

	return c.Status(200).JSON(codeParam)
}

func ShortenUrlHandler(c *fiber.Ctx) error {
	urlParams := c.Params("shorten")

	fmt.Println(urlParams)

	return c.Status(200).JSON(urlParams)
}
