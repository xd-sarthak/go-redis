package core

import (
	"time"
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
	store[k] = v
}

func Get(k string) *Obj {
	obj, exists := store[k]
	if !exists {
		return nil
	}
	return obj
}