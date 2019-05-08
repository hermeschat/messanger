package auth

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"strconv"
)

func ConnectRedis() (*redis.Client, error) {
	dbName, err := strconv.Atoi(config.RedisDBName)
	if err != nil {
		return nil, err
	}
	Addr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: config.RedisPassword, // no password set
		DB:       dbName,               // use default DB
	})
	return client, nil

}

// helpers function for checking card is blocked or not
func CheckIsBlock(h echo.HandlerFunc) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sourceUserid := c.Param("source")
			redisCon, err := ConnectRedis()
			if err != nil {
				return utils.GetValidationError("تعداد رمز اشتباه بیش از حد مجاز می باشد")
			}
			key := "Block:" + sourceUserid
			_, err = redisCon.Get(key).Result()
			if err == redis.Nil {
				return h(c)
			} else if err != nil {
				return err
			} else {
				return utils.GetValidationError("تعداد رمز اشتباه بیش از حد مجاز می باشد")
			}
			return h(c)
		}
	}
}
