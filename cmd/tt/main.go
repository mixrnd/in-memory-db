package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var sb strings.Builder
var lineCounter, lastFileNumber int

func main() {
	//f, err := os.OpenFile("data/test_2", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err) // i'm simplifying it here. you can do whatever you want.
	//}
	//sb := strings.Builder{}
	//for i := 0; i < 10; i++ {
	//	sb.WriteString("dfsfsdfs\n")
	//}
	//f.WriteString(sb.String())
	//f.Sync()
	//defer f.Close()

	files, err := os.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}

	lastFileNumber := 0
	for _, file := range files {
		cn := getFileNumber(file.Name())
		if cn > lastFileNumber {
			lastFileNumber = cn
		}
	}

	lastFileNumber++

	wg := sync.WaitGroup{}
	wg.Add(2)
	timer := time.NewTicker(5 * time.Second)
	go func() {
		//time.Sleep(time.Second)
		//defer wg.Done()
		//for {
		//	time.Sleep(time.Second)
		//	mu.Lock()
		//	lineCounter++
		//	sb.WriteString(time.Now().String() + "\n")
		//	mu.Unlock()
		//}

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("in-memory-db > ")
			query, err := reader.ReadString('\n')
			if err != nil {
				fmt.Print(err)
				return
			}
			mu.Lock()
			sb.WriteString(query)
			if lineCounter > 3 {
				writeToFile()
				timer.Reset(5 * time.Second)
			} else {
				lineCounter++
			}

			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-timer.C:
				mu.Lock()
				//можем писать в файл одновременно со счётчиком, надо как-то проверить
				//как вариант
				//if lineCounter < 5 {
				//	mu.Unlock()
				//	break
				//}
				writeToFile()
				sb.Reset()
				mu.Unlock()
			}
		}
	}()

	wg.Wait()
}

func getFileNumber(fileName string) int {
	res := strings.Split(fileName, "_")
	if len(res) != 2 {
		return 0
	}

	number, _ := strconv.Atoi(res[1])
	return number
}

func writeToFile() {
	if sb.Len() == 0 {
		return
	}
	lineCounter = 0
	f, err := os.OpenFile("data/test_"+strconv.Itoa(lastFileNumber), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	f.WriteString(sb.String())
	f.Sync()
	f.Close()
	lastFileNumber++
	sb.Reset()
}
