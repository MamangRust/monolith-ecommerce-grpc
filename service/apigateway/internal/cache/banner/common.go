package banner_cache

import "time"

const (
	bannerAllCacheKey     = "banner:all:page:%d:pageSize:%d:search:%s"
	bannerByIdCacheKey    = "banner:id:%d"
	bannerActiveCacheKey  = "banner:active:page:%d:pageSize:%d:search:%s"
	bannerTrashedCacheKey = "banner:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)
