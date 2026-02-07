package env

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/s-hammon/p"
)

type CacheFile struct {
	ETags map[string]string `json:"eTags"`
}

func GetCachePath(dir string) string {
	return path.Join(dir, "cache.json")
}

func LoadCacheFile() CacheFile {
	data, _ := os.ReadFile(homeDirCacheFile)

	var cache CacheFile
	json.Unmarshal(data, &cache)
	return cache
}

func WriteCacheFile(key string, val map[string]string) {
	var cache CacheFile

	switch key {
	default:
	case "eTags":
		cache.ETags = val
	}

	data, _ := json.MarshalIndent(cache, "", "  ")
	if err := os.WriteFile(homeDirCacheFile, data, 0600); err != nil {
		msg := p.Format("couldn't write to %s: %v", homeDirCacheFile, err)
		log.Println(msg)
	}
}
