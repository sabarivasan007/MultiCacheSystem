package main

import (
	"go-cache-server/Internal/cache"
	"go-cache-server/Internal/server"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	redisCache := cache.NewRedisCache("localhost:6379", "", 0, 1*time.Minute)
	memCache := cache.NewMemCache("localhost:11211", 60)

	cacheSystemType := server.NewServer(redisCache, memCache)

	router := gin.Default()

	// Cache System routes
	router.GET("/cache/:key", cacheSystemType.GetCacheHandler)
	//r.GET("/cache/TTL/:key", cacheSystemType.GetCacheWithTTLHandler)
	router.POST("/cache", cacheSystemType.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystemType.DeleteCacheHandler)
	router.PUT("/cache/clear", cacheSystemType.ClearAllCachesHandler)

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server started at %s\n", addr)
	log.Fatal(router.Run(addr))
}
