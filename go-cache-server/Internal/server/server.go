package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"go-cache-server/Internal/cache"
	utils "go-cache-server/packageUtils/Utils"

	"github.com/gorilla/mux"
)

/*
 */
type Server struct {
	redisCache cache.CacheLibrary
	memCache   cache.CacheLibrary
	mu         sync.Mutex
}

/*
 */
func NewServer(redisCache cache.CacheLibrary, memCache cache.CacheLibrary) *Server {
	return &Server{
		redisCache: redisCache,
		memCache:   memCache,
	}
}

/*
 */
func (s *Server) determineCacheLibraryType(cacheType string) cache.CacheLibrary {
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

/*
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
		utils.RespondJSON(w, http.StatusOK, map[string]string{"key": key, "value": value})
		log.Printf("Returning from %T", cache)
		return
	}

	utils.RespondError(w, http.StatusNotFound, "Cache miss")
}

/*
 */
func (s *Server) SetCache(w http.ResponseWriter, r *http.Request) {
	CacheLibraryType := r.URL.Query().Get("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LogError("Error While decoding JSON", err)
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Set(payload.Key, payload.Value); err != nil {
		utils.LogError("Error While setting cache", err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to set cache")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/*
* Set Cache Data with Specified TTL
 */
func (s *Server) SetCacheWithTTL(w http.ResponseWriter, r *http.Request) {
	ttl := mux.Vars(r)["ttl"]
	expireDuration, err := strconv.Atoi(ttl)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	CacheLibraryType := r.URL.Query().Get("cache")
	if CacheLibraryType == "" {
		utils.RespondError(w, http.StatusBadRequest, "Cache type is missing")
		return
	}
	cache := s.determineCacheLibraryType(CacheLibraryType)
	if cache == nil {
		utils.RespondError(w, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LogError("Error decoding JSON", err)
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if err := cache.SetWithTTL(payload.Key, payload.Value, expireDuration); err != nil {
		utils.LogError("Error setting cache", err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to set cache")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/*
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
		utils.LogError("Error deleting cache", err)
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete cache")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/*
 */
func (server *Server) ClearAllCaches(writter http.ResponseWriter, request *http.Request) {
	log.Printf("Clearing all caches")
	CacheLibraryType := request.URL.Query().Get("cache")
	cache := server.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(writter, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	server.mu.Lock()
	defer server.mu.Unlock()
	if err := cache.ClearAll(); err != nil {
		utils.LogError("Error clearing cache", err)
		utils.RespondError(writter, http.StatusInternalServerError, "Failed to clear cache")
		return
	}

	utils.RespondJSON(writter, http.StatusOK, map[string]string{"status": "ok"})
}
