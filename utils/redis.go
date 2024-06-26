package utils

import (
	"time"

	redisdb "github.com/go-redis/redis"
)

type RedisMutexLock struct {
	client *redisdb.Client
	prefix string
	timeout time.Duration
}

// 加锁
func (rml *RedisMutexLock) Lock(key string) bool {
	var try int
	for {
		if rml.client.SetNX(rml.prefix + key, time.Now().UnixNano(), rml.timeout).Val() {
			return true
		} else {
			s := Atoi64(rml.client.Get(rml.prefix + key).Val()) + int64(rml.timeout) - time.Now().UnixNano()
			if s > 0 {
				time.Sleep(100 * time.Millisecond)
			}
			try++
			if try > 10 {
				break
			}
		}
	}
	return false
}

// 解锁
func (rml *RedisMutexLock) Unlock(key string) {
	rml.client.Del(rml.prefix + key)
}

// 检查是否存在锁, 给读操作使用
func (rml *RedisMutexLock) IsExist(key string) bool {
	var try int
	for {
		s := Atoi64(rml.client.Get(rml.prefix + key).Val()) + int64(rml.timeout) - time.Now().UnixNano()
		if s > 0 {
			time.Sleep(100 * time.Millisecond)
			try++
			if try > 10 {
				return true
			}
		} else {
			return false
		}
	}
}

func NewRedisMutexLock(client *redisdb.Client, prefix string, timeout time.Duration) *RedisMutexLock {
	if timeout == 0 {
		timeout = time.Second * 5
	}
	return &RedisMutexLock{
		client: client,
		prefix: prefix,
		timeout: timeout,
	}
}