package routes

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/kristof1345/urlshort/db"
	"github.com/teris-io/shortid"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type LongURL struct {
	LongURL string `json:"longUrl"`
}

type URL struct {
	LongURL  *url.URL `bson:"longUrl"`
	ShortURL string   `bson:"shortUrl"`
	URLCode  string   `bson:"urlCode"`
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

	fmt.Println(u)

	filter := bson.D{{"longUrl", u}}

	var result URL

	err = db.UrlCollection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println("its ok1")

	if err != nil {
		shortenedUrl := baseUrl + "/" + urlCode
		fmt.Println("its ok")

		urlResults, err := db.UrlCollection.InsertOne(context.TODO(), bson.D{
			{Key: "longUrl", Value: u},
			{Key: "shortUrl", Value: shortenedUrl},
			{Key: "urlCode", Value: urlCode},
		})

		fmt.Println("its ok here")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("its ok here tpoo")

		c.Status(200).JSON(urlResults)
	}

	// fmt.Println(u, urlCode)

	return c.Status(201).JSON(result)
}
