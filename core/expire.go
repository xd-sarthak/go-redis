package core

import (
	"log"
	"time"
)

func expiredSample() float32 {
	var limit int = 20
	var expiredCount int = 0

	for key, obj := range store {
		if obj.ExpiresAt != -1 {
			limit--
			// if key is expired
			if obj.ExpiresAt <= time.Now().UnixMilli() {
				Del(key)
				expiredCount++
			}
		}

		if limit == 0 {
			break;
		}
	}

	return float32(expiredCount) / float32(20.0)
}


// delete all expired keys
// sampling approach: sample 20 keys, if more than 25% are expired, delete them and repeat
func DeleteExpiredKeys() {
	for {
		frac := expiredSample()
		// if sample has less than 25% expired keys, we can break the loop
		if frac < 0.25 {
			break
		}
	}
	log.Println("Expired keys deleted. total keys remaining:", len(store))
}