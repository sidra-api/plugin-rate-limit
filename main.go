package main

import (
	"log"
	"sync"
	"time"
	"github.com/sidra-gateway/go-pdk/server"
)

var rateLimitMap = make(map[string]int)
var rateLimitMutex = &sync.Mutex{}
const rateLimitPerMinute = 5

func resetRateLimit() {
	for range time.Tick(time.Minute) {
		rateLimitMutex.Lock()
		rateLimitMap = make(map[string]int)
		rateLimitMutex.Unlock()
	}
}

func rateLimitHandler(req server.Request) server.Response {
	clientIP := req.Headers["X-Real-Ip"]
	if clientIP == "" {
		log.Println("Missing X-Real-IP header")
		return server.Response{
			StatusCode: 400,
			Body:       "Missing X-Real-IP header",
		}
	}

	rateLimitMutex.Lock()
	rateLimitMap[clientIP]++
	currentCount := rateLimitMap[clientIP]
	rateLimitMutex.Unlock()

	if currentCount > rateLimitPerMinute {
		log.Printf("Rate limit exceeded for IP: %s", clientIP)
		return server.Response{
			StatusCode: 429,
			Body:       "Rate limit exceeded",
		}
	}

	log.Printf("Rate limit OK for IP: %s", clientIP)
	return server.Response{
		StatusCode: 200,
		Body:       "Request allowed",
	}
}

func main() {
	go resetRateLimit()
	server.NewServer("rate-limit", rateLimitHandler).Start()
}
