package core

import (
	"time"
	"github.com/xd-sarthak/go-redis/config"
)

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64,objType uint8, objEncoding uint8) *Obj {
	var expiresAt int64 = -1
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	
	return &Obj{
		Value:     value,
		TypeEncoding: objType | objEncoding,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, v *Obj) {
	if len(store) >= config.KeysLimit {
		// Handle keys limit exceeded (e.g., remove oldest key)
		evict();
	}
	store[k] = v
	if KeySpaceStat[0] == nil {
		KeySpaceStat[0] = make(map[string]int)
	}
	KeySpaceStat[0]["keys"]++
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil {
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli(){
			delete(store, k)
			return nil
		}
	}
	return v
}

func Del(k string) bool {
	if _,ok := store[k]; ok {
		delete(store, k)
		KeySpaceStat[0]["keys"]--
		return true
	}
	return false
}