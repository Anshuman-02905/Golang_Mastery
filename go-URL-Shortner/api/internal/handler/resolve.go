package handler

import (
	"go-url-shortener/database"
	"go-url-shortener/internal/repository"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)


func ResolveUrl(c *fiber.Ctx)error{

	url:= c.Params("url")
	r:=database.GetClient(0)
	urlRepo:= repository.NewUrlRepository(database.GetClient(0),database.Ctx)

	defer r.Close()

	
	_,value,err:= urlRepo.Exists(url,) 

	if err!=redis.Nil{

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"Short Not Found",
		})
	}else if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"cannot connect to DB",
	})}

		rInr:= database.GetClient(1)
		defer rInr.Close()

		_=rInr.Incr(database.Ctx,"counter")
		return c.Redirect(value,301 )
	

}
