package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var channelName = os.Getenv("PEER_CHANNEL_NAME")
var mtrScriptPath = os.Getenv("MTR_SCRIPT_PATH")

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func exec_mtr(wg *sync.WaitGroup, host string, port string) {
	defer wg.Done()

	cmd, err := exec.Command(
		"/bin/sh",
		mtrScriptPath,
		host, port).Output()

	if err != nil {
		fmt.Printf("error %s", err)
	}

	output := string(cmd)
	fmt.Println(output)
}

func main() {
	var wg sync.WaitGroup

	pubsub := redisClient.Subscribe(ctx, channelName)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	// Schedule to exit
	go func() {
		time.Sleep(24 * time.Hour)

		err := pubsub.Unsubscribe(ctx, channelName)
		if err != nil {
			panic(err)
		}

		wg.Wait() // Wait for any mtr executions to finish

		os.Exit(0)
	} ()

	for msg := range ch {
		var dat map[string]interface{}

		json.Unmarshal([]byte(msg.Payload), &dat)
		if err != nil {
			fmt.Printf("error %s", err)
		}

		peerDest := dat["peer"].(string)
		split := strings.Split(peerDest, ":")
		host := split[0]
		port := split[1]

		fmt.Println(host)
		fmt.Println(port)


		// Execute once for now
		wg.Add(1)
		go exec_mtr(&wg, host, port)


		// Periodically run mtr on endpoints
		// gocron.Every(1).Hour().Do(exec_mtr, host, port)
		// gocron.Start()
	}
}
