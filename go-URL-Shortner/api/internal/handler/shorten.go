package handler

import (
	"go-url-shortener/database"
	"go-url-shortener/helpers"
	"go-url-shortener/internal/builder"
	"go-url-shortener/internal/logx"
	"go-url-shortener/internal/repository"
	"go-url-shortener/internal/strategy"
	"go-url-shortener/internal/utils"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ShortenURL(c *fiber.Ctx) error {

	body, err := utils.ParseRequst(c)

	if err != nil {
		logx.Error("Error at Parsing the request")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errot":    err.Error(),
			"location": "At Parsing Request",
		})
	}
	logx.Info("Requst Parsing is done")
	logx.Info("Starting EnforceRateLimit")
	ok, err, info := utils.EnforceRateLimit(c)
	if !ok {
		logx.Info("Enforce Limit Returned not ok")
		if err != nil && err.Error() == "rate limit exceeded" {
			logx.Error("Error RATE LIMIT EXCEEDED")
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": err.Error(),
				"info":  info,
			})
		}
		logx.Error("Some Error in Rate limit")
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
	//CHECK URL
	id, _ := strat.Generate(body)

	urlRepo := repository.NewUrlRepository(database.GetClient(0), database.Ctx)

	ok, val, err := urlRepo.Exists(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error at fetching data from Database",
		})
	}
	if ok {
		log.Printf("%s provided an already-used short URL. Value: %v", c.IP(), val)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Short URL already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	if err := urlRepo.Save(id, body.URL, body.Expiry*3600*time.Second); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to save to database",
		})
	}

	rateRepo := repository.NewRateLimitRepository(database.GetClient(1), database.Ctx)
	val, _ = rateRepo.GetQouta(c.IP())
	ttl, _ := rateRepo.GetTTL(c.IP())
	XRateRemaining, _ := strconv.Atoi(val)
	XRateLimitReset := ttl / time.Nanosecond / time.Minute
	_ = rateRepo.DecrementQuota(c.IP())

	//Used builder pattern to generate the response
	response := builder.NewReponseBuilder().
		SetURL(body.URL).
		SetCustomSort(os.Getenv("DOMAIN")+"/"+id).
		SetExpiry(body.Expiry).
		SetRateLimit(XRateRemaining, XRateLimitReset).
		Build()

	return c.Status(fiber.StatusOK).JSON(response)

}
