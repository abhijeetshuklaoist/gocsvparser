package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main11() {
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	wg.Add(2)
	go printOdds(&wg)
	go printEvens(&wg)
	wg.Wait()
	fmt.Println("Done")
}


func printOdds(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i < 102; i = i + 2 {
		time.Sleep(100)
		fmt.Printf("\nOdd %d", i)
	}
}


func printEvens(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 102; i = i + 2 {
		time.Sleep(100)
		fmt.Printf("\nEven %d", i)
	}
}
