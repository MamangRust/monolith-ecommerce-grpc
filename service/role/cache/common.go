package cache

import "time"

const (
	roleAllCacheKey      = "role:all:page:%d:pageSize:%d:search:%s"
	roleByIdCacheKey     = "role:id:%d"
	roleByNameCacheKey   = "role:name:%s"
	roleByUserIdCacheKey = "role:user_id:%d"
	roleActiveCacheKey   = "role:active:page:%d:pageSize:%d:search:%s"
	roleTrashedCacheKey  = "role:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)
