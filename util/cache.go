package util

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"
)

type UtiCache struct {
	Data      map[string]Uti
	Timestamp time.Time
}

type RecommendedAppsCache struct {
	Data      map[string][]string // key is suffix, value is app list
	Timestamp time.Time
}

func getCacheFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".cache", "dutis")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "uti_cache.gob"), nil
}

func getRecommendedAppsCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".cache", "dutis")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "recommended_apps_cache.gob"), nil
}

func LoadUtiCache() (map[string]Uti, bool) {
	cachePath, err := getCacheFilePath()
	if err != nil {
		return nil, false
	}

	file, err := os.Open(cachePath)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	var cache UtiCache
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, false
	}

	// Cache valid for 24 hours
	if time.Since(cache.Timestamp) > 24*time.Hour {
		return nil, false
	}

	return cache.Data, true
}

func SaveUtiCache(data map[string]Uti) error {
	cachePath, err := getCacheFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	cache := UtiCache{
		Data:      data,
		Timestamp: time.Now(),
	}

	encoder := gob.NewEncoder(file)
	return encoder.Encode(cache)
}

func GetCachedUtiMap() map[string]Uti {
	if cached, ok := LoadUtiCache(); ok {
		return cached
	}

	// Cache miss, build and save
	utiMap := ListApplicationsUti()
	_ = SaveUtiCache(utiMap)
	return utiMap
}

func LoadRecommendedAppsCache(suffix string) ([]string, bool) {
	cachePath, err := getRecommendedAppsCachePath()
	if err != nil {
		return nil, false
	}

	file, err := os.Open(cachePath)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	var cache RecommendedAppsCache
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, false
	}

	// Cache valid for 24 hours
	if time.Since(cache.Timestamp) > 24*time.Hour {
		return nil, false
	}

	apps, ok := cache.Data[suffix]
	return apps, ok
}

func SaveRecommendedAppsCache(suffix string, apps []string) error {
	cachePath, err := getRecommendedAppsCachePath()
	if err != nil {
		return err
	}

	// Load existing cache or create new
	var cache RecommendedAppsCache
	file, err := os.Open(cachePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		_ = decoder.Decode(&cache)
		file.Close()
	}

	if cache.Data == nil {
		cache.Data = make(map[string][]string)
	}

	// Update with new data
	cache.Data[suffix] = apps
	cache.Timestamp = time.Now()

	// Save
	outFile, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	return encoder.Encode(cache)
}
