package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"go-cache-server/Internal/cache"
	utils "go-cache-server/packageUtils/Utils"

	"github.com/gorilla/mux"
)

/* Structure for multiple Cache System
 */
type Server struct {
	redisCache cache.CacheSystem
	memCache   cache.CacheSystem
	mu         sync.Mutex
}

/* Creating a New server
 */
func NewServer(redisCache cache.CacheSystem, memCache cache.CacheSystem) *Server {
	return &Server{
		redisCache: redisCache,
		memCache:   memCache,
	}
}

/* Determine the cache Library Type based on URI Param.
 */
func (s *Server) determineCacheLibraryType(cacheType string) cache.CacheSystem {
	//cacheType := mux.Vars(r)["cacheType"]
	switch cacheType {
	case "redis":
		return s.redisCache
	case "memcache":
		return s.memCache
	default:
		return nil
	}
}

/* Get Cache data from specfied Cache System.
 */
func (s *Server) GetCache(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	CacheLibraryType := r.URL.Query().Get("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	value, err := cache.Get(key)
	if err == nil {
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"key": key, "value": value})
		log.Printf("Returning from %T", cache)
		return
	}

	utils.RespondError(w, http.StatusNotFound, "Cache miss")
}

/* Get Cache data from specfied Cache System.
 */
func (s *Server) GetCacheWithTTL(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	CacheLibraryType := r.URL.Query().Get("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	value, ttl, err := cache.GetWithTTL(key)
	log.Println("values:", value)
	log.Println("values:", ttl)
	log.Println("values:", err)
	//str := strconv.Itoa(ttl)
	if err == nil {
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"key": key, "value": value, "ttl": ttl.String()})
		log.Printf("Returning from %T", cache)
		return
	}

	utils.RespondError(w, http.StatusNotFound, "Cache miss")
}

/* Set the Cache data to specific Cache System
 */
func (s *Server) SetCache(w http.ResponseWriter, r *http.Request) {
	CacheLibraryType := r.URL.Query().Get("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	var payload struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
		TTL   int64       `json:"ttl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LogError("Error While decoding JSON", err)
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	log.Println("Step1: ", payload.Key, payload.Value, payload.TTL)

	ttl := time.Duration(payload.TTL) * time.Second

	log.Println("Step2: ", payload.Key, payload.Value, &ttl)

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Set(payload.Key, payload.Value, ttl); err != nil {
		utils.LogError("Error While setting cache", err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to set cache")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/* Delete the given key from specified Cache System.
 */
func (s *Server) DeleteCache(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	CacheLibraryType := r.URL.Query().Get("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	log.Printf("Deleting cache for key: %s", key)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Delete(key); err != nil {
		utils.LogError("Error while deleting cache", err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete cache")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/* Clear All data from specified cache System.
 */
func (server *Server) ClearAllCaches(writter http.ResponseWriter, request *http.Request) {
	log.Print("Started Clearing all caches")
	CacheLibraryType := request.URL.Query().Get("cache")
	cache := server.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(writter, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	server.mu.Lock()
	defer server.mu.Unlock()
	if err := cache.ClearAll(); err != nil {
		utils.LogError("Error while clearing cache System", err)
		utils.RespondError(writter, http.StatusInternalServerError, "Failed to clear cache")
		return
	}

	utils.RespondJSON(writter, http.StatusOK, map[string]string{"status": "ok"})
}
