package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//linear(10)
	parallelFirst(10, 5)
}

func createArr(n int) ([]int, []int) {
	arr := make([]int, n)
	arrFind := make([]int, n, n)

	for i := 0; i < len(arr); i++ {
		arr[i] = i + 2
	}
	return arr, arrFind
}

func linear(n int) {
	arr, arrFind := createArr(n)
	t := time.Now()
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[j]%arr[i] == 0 {
				arrFind[j] = 1
			} else if arr[j]%arr[i] != 0 {
				if arrFind[j] == 1 {
					continue
				}
				arrFind[j] = 0
			}
		}
		//time.Sleep(time.Second)
	}
	fmt.Println(arrFind)
	for idx, value := range arrFind {
		if value == 0 {
			fmt.Print(arr[idx], " ")
		}
	}
	fmt.Println()
	elapsedTime := time.Since(t)
	fmt.Println("Прошло времени:", elapsedTime)
}

func parallelFirst(n, m int) {

	task := make(chan int, n)
	results := make(chan int, n)
	var wg sync.WaitGroup
	for i := 0; i < m; i++ {
		wg.Add(1)
		go isPrime(i+1, task, results, &wg)
	}

	for i := 0; i < n; i++ {
		value := i
		task <- value
	}
	close(task)

	t := time.Now()
	go func() {
		wg.Wait()
		close(results)
	}()
	for value := range results {
		fmt.Printf("значение = %d\n", value)
	}
	fmt.Println("Время начала", t)

	fmt.Println("Прошло времени:", time.Since(t))
}

func isPrime(id int, task <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range task {

		if value <= 1 {
			fmt.Printf("Поток %d закончил\n", id)
			continue
		}
		if value <= 3 {
			results <- value
			fmt.Printf("Поток %d закончил\n", id)
			//continue
		} else if value%2 == 0 || value%3 == 0 {
			fmt.Printf("Поток %d закончил\n", id)
			continue
		}
		results <- value
	}
}
