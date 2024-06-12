package server

import (
	"go-cache-server/Internal/cache"
	utils "go-cache-server/packageUtils/Utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

func (s *Server) GetCacheHandler(c *gin.Context) {
	key := c.Param("key")
	CacheLibraryType := c.Query("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	value, err := cache.Get(key)
	if err != nil {
		utils.LogError("Error while getting cache", err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to get cache")
		return
	}

	utils.RespondJSON(c.Writer, http.StatusOK, value)
}

// func (s *Server) GetCacheWithTTLHandler(c *gin.Context) {
// 	key := c.Param("key")
// 	CacheLibraryType := c.Query("cache")
// 	cache := s.determineCacheLibraryType(CacheLibraryType)

// 	if cache == nil {
// 		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
// 		return
// 	}

// 	value, ttl, err := cache.GetWithTTL(key)
// 	if err != nil {
// 		utils.LogError("Error while getting cache with TTL", err)
// 		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to get cache with TTL")
// 		return
// 	}

// 	utils.RespondJSON(c.Writer, http.StatusOK, map[string]interface{}{
// 		"value": value,
// 		"ttl":   ttl.Seconds(),
// 	})
// }

func (s *Server) SetCacheHandler(c *gin.Context) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	CacheLibraryType := c.Query("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)
	if cache == nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	ttl := time.Duration(payload.TTL) * time.Second

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Set(payload.Key, payload.Value, ttl); err != nil {
		utils.LogError("Error while setting cache", err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to set cache")
		return
	}

	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) DeleteCacheHandler(c *gin.Context) {
	key := c.Param("key")
	CacheLibraryType := c.Query("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Delete(key); err != nil {
		utils.LogError("Error while deleting cache", err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Cache not Found - Failed to delete cache")
		return
	}

	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ClearAllCachesHandler(c *gin.Context) {
	CacheLibraryType := c.Query("cache")
	cache := s.determineCacheLibraryType(CacheLibraryType)

	if cache == nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.ClearAll(); err != nil {
		utils.LogError("Error while clearing cache", err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to clear cache")
		return
	}
	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}
