package utils

import (
	"errors"
	"go-url-shortener/database"
	"go-url-shortener/helpers"
	"go-url-shortener/internal/models"
	"log"
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
	val, err := getRdsIP(c)
	if err == redis.Nil {
		if err := UpdateRdsIP(c, os.Getenv("API_QUOTA"), 0); err != nil {
			return false, err, 0
		}
	} else if err != nil {
		return false, err, 0
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			//Get TTL of the key
			rds1 := database.GetClient(1)
			limit, _ := rds1.TTL(database.Ctx, c.IP()).Result()
			return false, errors.New("rate limit exceeded"), limit / time.Nanosecond / time.Minute
		}
		//Might need to decrement
	}
	return true, nil, 0

}

func getRdsIP(c *fiber.Ctx) (string, error) {
	rds1 := database.GetClient(1)
	defer rds1.Close()
	val, err := rds1.Get(database.Ctx, c.IP()).Result()
	return val, err

}

func UpdateRdsIP(c *fiber.Ctx, qouta string, duration time.Duration) error {
	if duration == 0 {
		duration = 30 * 60 * time.Second
	}
	rds1 := database.GetClient(1)
	err := rds1.Set(database.Ctx, c.IP(), qouta, duration).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func ValidateUrl(url string) (bool, error) {
	if !govalidator.IsURL(url) {
		return false, errors.New("Invalid URL")
	}
	if helpers.RemoveDomainError(url) {
		return false, errors.New("you cannot hack the system")

	}
	return true, nil
}
