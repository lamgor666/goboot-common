package cachex

import (
	"github.com/gomodule/redigo/redis"
	"github.com/lamgor666/goboot-common/util/fsx"
	"github.com/lamgor666/goboot-common/util/stringx"
	"github.com/patrickmn/go-cache"
	"os"
	"strings"
	"time"
)

var cacheDir string
var cacheKeyPrefix string
var cacheKeyRedismqNormal = "redismq.normal"
var cacheKeyRedismqDelayable = "redismq.delayable"
var defaultCacheStore string
var gocache *cache.Cache
var cacheStores = map[string]ICache{}

func CacheDir(dir ...string) string {
	if len(dir) > 0 {
		_dir := dir[0]

		if _dir != "" {
			_dir = fsx.GetRealpath(_dir)
			_dir = strings.ReplaceAll(_dir, "\\", "/")
			_dir = strings.TrimRight(_dir, "/")

			if stat, err := os.Stat(_dir); err != nil || !stat.IsDir() {
				os.MkdirAll(_dir, 0755)
			}

			if stat, err := os.Stat(_dir); err != nil || !stat.IsDir() {
				_dir = ""
			}
		}

		cacheDir = _dir
	}

	_dir := cacheDir

	if _dir != "" {
		return _dir
	}

	_dir = fsx.GetRealpath("datadir:cache")
	_dir = strings.ReplaceAll(_dir, "\\", "/")
	_dir = strings.TrimRight(_dir, "/")

	if stat, err := os.Stat(_dir); err != nil || !stat.IsDir() {
		os.MkdirAll(_dir, 0755)
	}

	if stat, err := os.Stat(_dir); err != nil || !stat.IsDir() {
		return ""
	}

	return _dir
}

func CacheKeyPrefix(prefix ...string) string {
	if len(prefix) > 0 {
		_prefix := prefix[0]

		if _prefix != "" {
			_prefix = strings.TrimRight(_prefix, ".")
		}

		if _prefix != "" {
			cacheKeyPrefix = _prefix
		}
	}

	return cacheKeyPrefix
}

func BuildCacheKey(cacheKey string) string {
	if cacheKeyPrefix == "" {
		return cacheKey
	}

	return cacheKeyPrefix + stringx.EnsureLeft(cacheKey, ".")
}

func CacheKeyRedismqNormal(cacheKey ...string) string {
	if len(cacheKey) > 0 {
		key := cacheKey[0]

		if key != "" {
			cacheKeyRedismqNormal = key
		}
	}

	return BuildCacheKey(cacheKeyRedismqNormal)
}

func CacheKeyRedismqDelayable(cacheKey ...string) string {
	if len(cacheKey) > 0 {
		key := cacheKey[0]

		if key != "" {
			cacheKeyRedismqDelayable = key
		}
	}

	return BuildCacheKey(cacheKeyRedismqDelayable)
}

func DefaultStore(name ...string) string {
	if len(name) > 0 {
		if name[0] != "" {
			defaultCacheStore = name[0]
		}
	}

	_name := defaultCacheStore

	if _name == "" {
		_name = "memory"
	}

	return _name
}

func WithMemoryCache(defaultTtl time.Duration, cleanupInterval time.Duration) {
	if gocache == nil {
		gocache = cache.New(defaultTtl, cleanupInterval)
	}

	cacheStores["memory"] = &memoryCache{}
}

func WithRedisCache(pool *redis.Pool) {
	cacheStores["redis"] = &redisCache{pool: pool}
}

func WithFileCache() {
	cacheStores["file"] = newFileCache()
}

func Store(name string) ICache {
	if c, ok := cacheStores[name]; ok {
		return c
	}

	return &noopCache{}
}
