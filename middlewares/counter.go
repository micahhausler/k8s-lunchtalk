package middlewares

import (
	"net/http"
	"fmt"
	"os"
	"log"

	redis "gopkg.in/redis.v3"
)

var client *redis.Client


func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}


// RequestCounter counts the requests by path
func RequestCounter() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(){
				if err := client.Incr(r.URL.Path).Err(); err != nil {
					log.Print(err.Error())
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}

