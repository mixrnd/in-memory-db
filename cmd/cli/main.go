package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"in-memory-db/internal/config"
	"in-memory-db/internal/network"
)

var address = flag.String("address", "127.0.0.1:3030", "server address")
var configPath = flag.String("config", "config.yml", "Path to config file")

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	maxMessageSize, err := cfg.Network.MessageSizeToSizeInBytes()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := network.NewClient(*address, cfg.Network.IdleTimeout, maxMessageSize)
	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("in-memory-db > ")
		query, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}

		res, err := client.Send(query)
		if err != nil {
			fmt.Printf("error > %s\n", err.Error())
		} else {
			fmt.Printf("result > %s\n", res)
		}
	}
}
