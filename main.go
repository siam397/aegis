package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	MAX_LIMIT int64 = 5
	TIMEOUT   int64 = 10
)

func redisInitialize() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func main() {
	e := echo.New()
	e.GET("/", request)
	e.Logger.Fatal(e.Start(":8000"))
}

func request(c echo.Context) error {
	isLimited := rateLimitter(c.QueryParam("id"))
	if isLimited {
		return c.JSON(http.StatusTooManyRequests, struct{ Message string }{Message: "Too Many Requests"})
	}
	return c.JSON(http.StatusOK, struct{ Message string }{Message: "Some Data"})
}

func rateLimitter(userId string) bool {
	rdb := redisInitialize()
	return tokenBucket(userId, rdb)
}

func tokenBucket(userId string, rdb *redis.Client) bool {
	context := context.Background()
	count, err := rdb.Get(context, userId).Result()
	if err == redis.Nil {
		rdb.Set(context, userId, 1, time.Duration(TIMEOUT)*time.Second)
		return false
	}
	val, err := strconv.ParseInt(count, 10, 64)
	fmt.Println(val, userId)
	if err == nil && val < MAX_LIMIT {
		rdb.Set(context, userId, val+1, redis.KeepTTL)
		return false
	}
	return true
}

// func _slidingWindowAlgorithm(userId string) {

// }

// func _leakingBucketAlgorithm(userId string) {

// }

// func _fixedWindowCounterAlgorithm(userId string) {

// }
