package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//linear(10000)
	//parallelFirst(10000, 2)
	//parallelSecond(1000, 4)
	parallelThird(10, 5)
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
		time.Sleep(time.Microsecond)
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
		go isPrime(task, results, &wg)
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

func isPrime(task <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range task {
		time.Sleep(time.Microsecond)
		if value <= 1 {
			continue
		}
		if value <= 3 {
			results <- value
		} else if value%2 == 0 || value%3 == 0 {
			continue
		} else {
			results <- value
		}
	}
}

func parallelSecond(n, m int) {
	arr, _ := createArr(n)
	task := make(chan []int, n)
	results := make(chan int)
	partSize := n / m
	var wg sync.WaitGroup

	for i := 0; i < m; i++ {
		task <- arr[partSize*i : partSize*(i+1)]
	}
	close(task)

	for i := 0; i < m; i++ {
		wg.Add(1)
		go isPrimeSecond(task, results, &wg)
	}

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

func isPrimeSecond(task <-chan []int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for arr := range task {
		for _, value := range arr {
			time.Sleep(time.Microsecond)
			if value <= 1 {
				continue
			}
			if value <= 3 {
				results <- value
			} else if value%2 == 0 || value%3 == 0 {
				continue
			} else {
				results <- value
			}
		}
	}
}

func parallelThird(n, m int) {
	arr, _ := createArr(n)
	primeArr := [10]int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	task := make(chan int, 10)
	results := make(chan int)
	var wg sync.WaitGroup

	t := time.Now()
	fmt.Println("Время начала", t)

	for i := 0; i < m; i++ {
		wg.Add(1)
		go isPrimeThird(arr, results, &wg, task)
	}

	// Отправляем все значения в канал task.
	for i := 0; i < len(primeArr); i++ {
		task <- primeArr[i]
	}

	close(task) // Закрываем канал task после отправки всех значений.

	resultMap := make(map[int]int)

	go func() {
		wg.Wait() // Ожидаем завершения всех горутин.
		close(results)
	}()

	for value := range results {
		// Увеличиваем счетчик встречаемости числа в мапе.
		resultMap[value]++
	}

	for num, count := range resultMap {
		if count == 10 {
			fmt.Printf("Число %d \n", num)
		}
	}

	fmt.Println("Прошло времени:", time.Since(t))
}

func isPrimeThird(arr []int, results chan<- int, wg *sync.WaitGroup, task <-chan int) {
	defer wg.Done()

	for primeArr := range task {
		for _, value := range arr {
			if value <= 1 {
				continue
			} else if value <= 3 {
				results <- value
			} else if value == primeArr {
				results <- value
			} else if value%primeArr == 0 && value != primeArr {
				continue
			} else if value%primeArr != 0 {
				results <- value
			}
		}
	}
}
