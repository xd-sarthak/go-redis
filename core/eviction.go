package core

import (
	"github.com/xd-sarthak/go-redis/config"
)

func evict() {
	switch config.EvictionPolicy {
	case "simple-first":
		evictFirst()
	}
}

func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}