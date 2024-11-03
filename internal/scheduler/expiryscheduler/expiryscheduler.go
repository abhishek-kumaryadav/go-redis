package expiryscheduler

import (
	"context"
	"go-redis/internal/repository"
	"go-redis/internal/service/expire"
	"go-redis/pkg/utils/log"
	"maps"
	"math/rand"
	"slices"
	"sync"
	"time"
)

// StartScheduler
/*
Official redis implementation uses a sampling algorithm to get a set of keys to check for expiry as it cannot sweep the whole database everytime
The algorithm checks if the keys are expired and if the number of cleared keys is less than a threshold the routine stops
The official implementation also keeps a to be expired key list from the sampled keys to minimize the number of expired keys in the memory
*/
func StartScheduler(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	scheduler := time.NewTicker(5 * time.Second)

	go func() {

		for {
			select {
			case <-scheduler.C:
				keys := slices.Collect(maps.Keys(repository.KeyValueStore))
				log.InfoLog.Printf("Starting expired keys clearing routine", keys)
				percentageCleared := 100
				for percentageCleared > 25 && len(keys) > 0 {
					seedKeys := getSeedKeys(keys)
					log.InfoLog.Printf("Collected seedKeys", seedKeys)
					countExpired := 0
					for _, seedKey := range seedKeys {
						isExpired, _ := expire.CheckAndDeleteExpired(seedKey)
						if isExpired {
							countExpired += 1
						}
					}
					percentageCleared = countExpired / len(seedKeys) * 100
					log.InfoLog.Printf("Cleared %d %% of expired key", percentageCleared)
					keys = slices.Collect(maps.Keys(repository.KeyValueStore))
				}

			}
		}
	}()

	// Wait for the context to be canceled
	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		log.InfoLog.Printf("Shutting down expiry scheduler gracefully...")
		_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := scheduler.Stop
		if err != nil {
			log.ErrorLog.Printf("Expiry scheduler shutdown error: %s\n", err)
		}
	}
}

func getSeedKeys(keys []string) []string {

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	numberOfKeys := max(int(len(keys)/10), 1)
	seedSet := make(map[string]bool)
	for i := 0; i < numberOfKeys; i++ {
		randomNumber := rand.Intn(len(keys)) // Generate a random number between 0 and n-1
		seedSet[keys[randomNumber]] = true
	}
	return slices.Collect(maps.Keys(seedSet))
}
