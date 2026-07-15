package core

import (
	"time"
	"github.com/xd-sarthak/go-redis/config"
)

var store map[string]*Obj

type Obj struct {
	Value     interface{}
	ExpiresAt int64 // Unix timestamp in milliseconds
}

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64) *Obj {
	var expiresAt int64 = -1
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	
	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, v *Obj) {
	if len(store) >= config.KeysLimit {
		// Handle keys limit exceeded (e.g., remove oldest key)
		evict();
	}
	store[k] = v
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil {
		if v.ExpiresAt <= time.Now().UnixMilli(){
			delete(store, k)
			return nil
		}
	}
	return v
}

func Del(k string) bool {
	if _,ok := store[k]; ok {
		delete(store, k)
		return true
	}
	return false
}