package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type LongURL struct {
	LongURL string `json:"longUrl"`
}

type URL struct {
	LongURL  string `bson:"longURL"`
	ShortURL string `bson:"shortURL"`
	URLCode  string `bson:"URLCode"`
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {
		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {
	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func main() {
	err := godotenv.Load()
	fmt.Println("Connecting to db...")

	if err != nil {
		log.Fatal(".env file coundn't be loaded")
	}

	mongoUri := os.Getenv("MONGO_URI")

	app := fiber.New()

	// Get Client, Context, CancelFunc and
	// err from connect method.
	client, ctx, cancel, err := connect(mongoUri)
	if err != nil {
		panic(err)
	}

	// Release resource when the main
	// function is returned.
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	app.Get("/:code", CodeHandler)
	app.Post("/api/shorten", func(c *fiber.Ctx) error {
		return ShortenUrlHandler(c, ctx, client)
	})

	log.Fatal(app.Listen(":5000"))
}

func CodeHandler(c *fiber.Ctx) error {
	codeParam := c.Params("code")

	fmt.Println(codeParam)

	return c.Status(200).JSON(codeParam)
}

func ShortenUrlHandler(c *fiber.Ctx, ctx context.Context, client *mongo.Client) error {
	UrlToShorten := new(LongURL)

	if err := c.BodyParser(UrlToShorten); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	baseUrl := os.Getenv("BASE_URL")

	if !checkURL(baseUrl) {
		return c.Status(500).JSON("Invalid base url")
	}

	urlCode, _ := shortid.Generate()

	if checkURL(UrlToShorten.LongURL) {
		//check if there is a doc
		coll := client.Database("test").Collection("urls")
		filter := bson.D{{Key: "longURL", Value: UrlToShorten.LongURL}}

		var result URL
		err := coll.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				shortUrl := baseUrl + "/" + urlCode
				//else create short one and upload it to mongo
				// This error means your query did not match any documents.
				_, err := coll.InsertOne(ctx, bson.D{
					{Key: "longURL", Value: UrlToShorten.LongURL},
					{Key: "shortURL", Value: shortUrl},
					{Key: "URLCode", Value: urlCode},
				})

				if err != nil {
					return c.Status(500).JSON("Internal Server Error")
				}

				ret := &URL{
					LongURL:  UrlToShorten.LongURL,
					ShortURL: shortUrl,
					URLCode:  urlCode,
				}

				return c.Status(200).JSON(ret)
			} else {
				return c.Status(500).JSON("Something went wrong...")
			}
		} else {
			return c.Status(200).JSON(result)
		}
	} else {
		return c.Status(500).JSON("Invalid long url")
	}
}

func checkURL(Url string) bool {
	_, err := url.ParseRequestURI(Url)
	return err == nil
}
