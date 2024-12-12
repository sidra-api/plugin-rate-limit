package main

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/sidra-gateway/go-pdk/server"
)

var rateLimitMap = make(map[string]int)
var rateLimitMutex = &sync.Mutex{}
var rateLimitPerMinute  int
var pluginName  string 

// Fungsi untuk mereset rate limit setiap menit
func resetRateLimit() {
	for range time.Tick(time.Minute) {
		rateLimitMutex.Lock()
		rateLimitMap = make(map[string]int)
		rateLimitMutex.Unlock()
	}
}

// Handler utama untuk memproses setiap permintaan yang masuk
func rateLimitHandler(req server.Request) server.Response {
	clientIP := req.Headers["X-Real-Ip"]
	if clientIP == "" {
		log.Println("Missing X-Real-IP header")
		return server.Response{
			StatusCode: 400,
			Body:       "Missing X-Real-IP header",
		}
	}

	// Menggunakan mutex untuk memastikan operasi thread-safe pada map rateLimitMap
	rateLimitMutex.Lock()
	rateLimitMap[clientIP]++
	currentCount := rateLimitMap[clientIP]
	rateLimitMutex.Unlock()

	// Jika jumlah permintaan melebihi batas, kembalikan respons 429 Too Many Requests
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
	// Ambil nama plugin dari variabel lingkungan PLUGIN_NAME, gunakan default jika kosong
	pluginName = os.Getenv("PLUGIN_NAME")
	if pluginName == "" {
		pluginName = "rate-limit" 
	}

	// Ambil nilai rate limit dari variabel lingkungan RATE_LIMIT
	rateLimitEnv := os.Getenv("RATE_LIMIT")
	if rateLimitEnv != "" {
		// Parse nilai rate limit menjadi integer
		parsedRate, err := strconv.Atoi(rateLimitEnv)
		if err == nil {
			rateLimitPerMinute = parsedRate
		} else {
			// Jika parsing gagal, gunakan nilai default
			log.Println("Invalid RATE_LIMIT, using default: 5")
			rateLimitPerMinute = 5
		}
	} else {
		// Jika variabel lingkungan tidak diset, gunakan nilai default
		rateLimitPerMinute = 5
	}

	// Log informasi awal tentang plugin yang dijalankan
	log.Printf("Starting plugin %s with rate limit: %d requests per minute", pluginName, rateLimitPerMinute)

	// Jalankan fungsi resetRateLimit sebagai goroutine untuk mereset map setiap menit
	go resetRateLimit()

	server.NewServer(pluginName, rateLimitHandler).Start()
}
