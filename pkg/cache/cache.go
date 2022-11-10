package cache

import (
	"errors"
	"fmt"
	"sync"
)

type cache struct {
	Data map[string]int `json:"data"`
}

var (
	judgeCache *cache
	judgeQueue *cache
	cacheLock  = new(sync.RWMutex)
)

func GetCacheKey(taskID string, ruleID int) string {
	return fmt.Sprintf("%s-%v", taskID, ruleID)
}

// Get 获取Key值
func Get(k string) (value int, err error) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if v, ok := judgeCache.Data[k]; ok {
		value = v
	} else {
		err = errors.New("non-existent")
	}
	return
}

// Set 设置Key值
func Set(k string, val int) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	judgeCache.Data[k] = val
}

// Delete 删除Key
func Delete(k string) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if _, ok := judgeCache.Data[k]; ok {
		delete(judgeCache.Data, k)
	}
}

// Pop 删除Key
func Pop(k string) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if _, ok := judgeQueue.Data[k]; ok {
		delete(judgeQueue.Data, k)
	}
}

// Pull 获取Key值
func Pull(k string) (value int, err error) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if v, ok := judgeQueue.Data[k]; ok {
		value = v
	} else {
		err = errors.New("non-existent")
	}
	return
}

// Push 设置Key值
func Push(k string, val int) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	judgeQueue.Data[k] = val
}

// NewInit 初始化
func NewInit() {
	data1 := make(map[string]int, 0)
	data2 := make(map[string]int, 0)
	judgeCache = &cache{Data: data1}
	judgeQueue = &cache{Data: data2}
}
