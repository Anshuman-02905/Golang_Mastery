package handler

import (
	"go-url-shortener/database"
	"go-url-shortener/helpers"
	"go-url-shortener/internal/builder"
	"go-url-shortener/internal/strategy"
	"go-url-shortener/internal/utils"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ShortenURL(c *fiber.Ctx) error {
	body, err := utils.ParseRequst(c)

	log.Printf("Received request from IP %s for URL %s, CustomShort %s, & Expiry %s", c.IP(), body.URL, body.CustomShort, body.Expiry)

	

	ok, err, info := utils.EnforceRateLimit(c)
	if !ok {
		if err != nil && err.Error() == "rate limit exceeded" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": err.Error(),
				"info":  info,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "rate limiter internal error",
		})
	}

	ok, err = utils.ValidateUrl(body.URL)
	if !ok {
		if err != nil {
			switch err.Error() {
			case "Invalid URL":
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			case "you cannot hack the system :)":
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
	}

	//Enforce HTTPS SSL

	body.URL = helpers.EnforceHTTP(body.URL)


	//Strategy Pattern
	var strat strategy.IDStrategy
	if body.CustomShort == "" {
		strat = strategy.UUIDStrategy{}
	} else {
		strat = strategy.CustomStrategy{}
	}
	id, err := strat.Generate(body)



	rds0 := database.GetClient(0)
	defer rds0.Close()

	val, err = rds0.Get(database.Ctx, id).Result()

	if val != "" {
		log.Printf("%s provided in already use URL and the val is %v", c.IP(), val)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No Url found at Database",
		})
	} else if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error at fetching data from Databse",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = rds0.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	val, _ = rds1.Get(database.Ctx, c.IP()).Result()
	log.Printf("Value after devrement%v", val)
	ttl, _ := rds1.TTL(database.Ctx, c.IP()).Result()
	XRateRemaining, _ := strconv.Atoi(val)
	XRateLimitReset := ttl / time.Nanosecond / time.Minute
	rds1.Decr(database.Ctx, c.IP())

	//Used builder pattern to generate the response
	response := builder.NewReponseBuilder().
		SetURL(body.URL).
		SetCustomSort(os.Getenv("DOMAIN")+"/"+id).
		SetExpiry(body.Expiry).
		SetRateLimit(XRateRemaining, XRateLimitReset).
		Build()

	return c.Status(fiber.StatusOK).JSON(response)

}
