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

var rdb *redis.Client
var ctx context.Context

func redisInitialize() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	redisInitialize()
	ctx = context.Background()
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
	// return tokenBucket(userId)
	return leakingBucketAlgorithm(userId)
}

func tokenBucket(userId string) bool {

	val := getApiCallCount(userId, 1)
	if val < MAX_LIMIT {
		rdb.Set(ctx, userId, val+1, redis.KeepTTL)
		return false
	}
	return true
}

// func _slidingWindowAlgorithm(userId string) {

// }

func leakingBucketAlgorithm(userId string) bool {
	val := getApiCallCount(userId, MAX_LIMIT)
	if val == 0 {
		return true
	}
	rdb.Set(ctx, userId, val-1, redis.KeepTTL)
	return false
}

// func _fixedWindowCounterAlgorithm(userId string) {

// }

func getApiCallCount(userId string, defaultValue int64) int64 {
	count, err := rdb.Get(ctx, userId).Result()
	if err == redis.Nil {
		rdb.Set(ctx, userId, defaultValue, time.Duration(TIMEOUT)*time.Second)
		return defaultValue
	}
	val, err := strconv.ParseInt(count, 10, 64)
	fmt.Println(userId, val)
	if err != nil {
		panic("How the f!!!")
	}
	return val
}
