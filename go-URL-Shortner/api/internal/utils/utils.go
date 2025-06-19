package utils

import (
	"errors"
	"go-url-shortener/database"
	"go-url-shortener/helpers"
	"go-url-shortener/internal/models"
	"go-url-shortener/internal/repository"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ParseRequst(c *fiber.Ctx) (*models.Request, error) {
	body := new(models.Request)

	if err := c.BodyParser(&body); err != nil {
		return body, err
	}
	return body, nil
}

func EnforceRateLimit(c *fiber.Ctx) (bool, error, time.Duration) {

	rateRepo := repository.NewRateLimitRepository(database.GetClient(1), database.Ctx)

	val, err := rateRepo.GetQouta(c.IP())
	if err == redis.Nil {
		if err := rateRepo.SetQuota(c.IP(), os.Getenv("API_QUOTA"), 0); err != nil {
			return false, err, 0
		}
	} else if err != nil {
		return false, err, 0
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			//Get TTL of the key
			limit, _ := rateRepo.GetTTL(c.IP())
			return false, errors.New("rate limit exceeded"), limit / time.Nanosecond / time.Minute
		}
		//Might need to decrement
	}
	return true, nil, 0

}

func ValidateUrl(url string) (bool, error) {
	if !govalidator.IsURL(url) {
		return false, errors.New("Invalid URL")
	}
	if !helpers.RemoveDomainError(url) {
		return false, errors.New("you cannot hack the system")

	}
	return true, nil
}
