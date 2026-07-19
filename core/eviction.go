package core

import (
	"github.com/xd-sarthak/go-redis/config"
)

func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		evictAllKeyRandom()
	}
}

func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

func evictAllKeyRandom() {
	evictCount := int64(config.EvictionRatio * float64(config.KeysLimit))

	// this is random in a hashmap store 
	for k := range store {
		Del(k)
		evictCount--
		if evictCount <= 0 {
			break;
	}
}
}