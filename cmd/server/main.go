package main

import (
	"go-redis/internal/repository"
	"go-redis/internal/service/datastructure"
	"go-redis/internal/service/expire"
	"go-redis/internal/service/hashmap"
	"go-redis/pkg/utils/log"
	"go-redis/pkg/utils/tcp"
	"maps"
	"math/rand"
	"net"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	HOST = "localhost"
	PORT = "7369"
	TYPE = "tcp4"
)

func main() {
	log.InitLog("build/logs/server.log")
	repository.InitRepositories()

	args := os.Args

	var wg sync.WaitGroup
	wg.Add(2)
	go startHttpServer(args, &wg)
	go startScheduler(&wg)
	wg.Wait()

}

/*
Official redis implementation uses a sampling algorithm to get a set of keys to check for expiry as it cannot sweep the whole database everytime
The algorithm checks if the keys are expired and if the number of cleared keys is less than a threshold the routine stops
The official implementation also keeps a to be expired key list from the sampled keys to minimize the number of expired keys in the memory
*/
func startScheduler(wg *sync.WaitGroup) {
	defer wg.Done()
	scheduler := time.NewTicker(5 * time.Second)
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
func startHttpServer(args []string, s *sync.WaitGroup) {
	port := PORT
	if len(args) == 2 && args[1] != "" {
		port = args[1]
	}

	listener, err := net.Listen(TYPE, HOST+":"+port)
	if err != nil {
		log.InfoLog.Fatal("Error: ", err)
		return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.InfoLog.Fatal("Error: ", err)
			return
		}
		go handleConnection(connection)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.InfoLog.Printf("Serving %s\n", c.RemoteAddr().String())

	packet := tcp.ReadFromConn(c)
	commands := strings.Split(string(packet), " ")
	primaryCommand := strings.TrimSpace(commands[0])

	var response string
	var ok bool

	result, ok := datastructure.GetDataStructureFromCommand(primaryCommand)
	if !ok {
		response = result
	} else {
		switch result {
		case datastructure.HASHMAP:
			response, ok = hashmap.Execute(commands)
		case datastructure.EXPIRE:
			response, ok = expire.Execute(commands)

		}
		if !ok {
			response = "Error running command: " + response
		}
	}

	num, _ := c.Write([]byte(response))
	log.InfoLog.Printf("Wrote back %d bytes, the payload is %s\n", num, response)
}
