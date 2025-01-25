package main

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"
	"in-memory-db/internal"
	"in-memory-db/internal/compute"
	inmemory "in-memory-db/internal/storage/in-memory"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	logger, _ := zap.NewProduction()

	e := inmemory.NewEngine()
	p := compute.NewParser()
	db := internal.NewDatabase(e, p, logger)

	for {
		fmt.Print("in-memory-db > ")
		query, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}

		res, err := db.RunQuery(query)
		if err != nil {
			fmt.Printf("error > %s\n", err.Error())
		} else {
			fmt.Printf("result > %s\n", res)
		}
	}
}
