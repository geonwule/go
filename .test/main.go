package main

import (
	"fmt"
	//"time"
	"sync"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		fmt.Println("worker", id, "job", i)
		//time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, &wg)
		//time.Sleep(time.Second)
	}

	//time.Sleep(time.Second * 4)
	wg.Wait()
}