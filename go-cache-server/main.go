// main.go
package main

import (
	"fmt"
	"go-cache-server/Internal/cache"
	"go-cache-server/Internal/server"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	redisCache := cache.NewRedisCache("localhost:6379", "", 0, 1*time.Minute)
	memCache := cache.NewMemCache("localhost:11211", 60)

	srv := server.NewServer(redisCache, memCache)

	r := mux.NewRouter()

	//
	r.HandleFunc("/cache/{key}", srv.GetCache).Methods("GET")
	r.HandleFunc("/cache", srv.SetCache).Methods("POST")
	r.HandleFunc("/cache/{ttl}", srv.SetCacheWithTTL).Methods("POST")
	r.HandleFunc("/cache/{key}", srv.DeleteCache).Methods("DELETE")
	r.HandleFunc("/cache/clear", srv.ClearAllCaches).Methods("PUT")

	// For Redis
	// r.HandleFunc("/redis/cache/{key}", srv.GetCache).Methods("GET")
	// r.HandleFunc("/redis/cache/{key}", srv.SetCache).Methods("POST")
	// r.HandleFunc("/redis/cache/{ttl}", srv.SetCacheWithTTL).Methods("POST")
	// r.HandleFunc("/redis/cache/{key}", srv.DeleteCache).Methods("DELETE")
	// r.HandleFunc("/redis/cache/clearAll", srv.ClearAllCaches).Methods("PUT")

	//For MemCache
	// r.HandleFunc("/memcache/cache/{key}", srv.GetCache).Methods("GET")
	// r.HandleFunc("/memcache/cache/{key}", srv.SetCache).Methods("POST")
	// r.HandleFunc("/memcache/cache/{ttl}", srv.SetCacheWithTTL).Methods("POST")
	// r.HandleFunc("/memcache/cache/{key}", srv.DeleteCache).Methods("DELETE")
	// r.HandleFunc("/memcache/cache/clearAll", srv.ClearAllCaches).Methods("PUT")

	// Use your custom router instance with middleware or additional configuration if needed
	router := Router{r}

	// Start the HTTP server
	addr := "127.0.0.1:8080"
	fmt.Printf("Server started at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

// Router is a custom type for mux.Router that can be used to add additional methods if needed
type Router struct {
	*mux.Router
}