package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

type SidraRequest struct {
	Headers map[string]string `json:"Headers"`
	Url     string            `json:"Url"`
	Method  string            `json:"Method"`
	Body    string            `json:"Body"`
}

type SidraResponse struct {
	StatusCode int    `json:"StatusCode"`
	Body       string `json:"Body"`
}

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

func rateLimitHandler(req SidraRequest) SidraResponse {
	clientIP := req.Headers["X-Real-Ip"]
	if clientIP == "" {
		log.Println("Missing X-Real-IP header")
		return SidraResponse{
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
		return SidraResponse{
			StatusCode: 429,
			Body:       "Rate limit exceeded",
		}
	}

	log.Printf("Rate limit OK for IP: %s", clientIP)
	return SidraResponse{
		StatusCode: 200,
		Body:       "Request allowed",
	}
}

func main() {
	go resetRateLimit()

	listener, err := net.Listen("unix", "/tmp/rate-limit.sock")
	if err != nil {
		log.Fatalf("Error setting up Unix domain socket: %v", err)
	}
	defer listener.Close()
	log.Println("Unix domain socket server listening on /tmp/rate-limit.sock")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			decoder := json.NewDecoder(conn)
			var req SidraRequest
			if err := decoder.Decode(&req); err != nil {
				log.Printf("Error decoding request: %v", err)
				return
			}
			log.Printf("Received request: %+v\n", req)

			resp := rateLimitHandler(req)

			encoder := json.NewEncoder(conn)
			if err := encoder.Encode(resp); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}(conn)
	}
}
