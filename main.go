package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func prime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	//precache primes
	isPrimeArr := make([]bool, 432)

	//how many threads we'll need to wait for
	primeWait := new(sync.WaitGroup)
	primeWait.Add(432)
	//concurrently check prime numbers
	for iterationNumber := range 432 {
		//inline concurrent function
		go func(routineInput int) {
			//when thread completes, mark thread as completed
			defer primeWait.Done()
			//check primeness
			isPrimeArr[routineInput] = prime(routineInput)
			//fmt.Println(isPrimeArr[iterationNumber])
		}(iterationNumber)
	}

	//wait for all concurrent threads to complete
	primeWait.Wait()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		//parse in text
		in := scanner.Text()
		nStr, kStr, _ := strings.Cut(in, " ")
		n, _ := strconv.Atoi(nStr)
		k, _ := strconv.Atoi(kStr)

		exp := make([]atomic.Int32, 432)
		for i := range exp {
			exp[i].Store(0) // Explicitly initialize to 0
		}

		//for any iterations we want to run concurrently
		factorizationWait := new(sync.WaitGroup)
		//make exp the prime factorization of n!/k!
		factorizationWait.Add(n - 1)
		for i := 2; i <= n; i++ {
			fmt.Println("i", i)

			//n!/(k!(n-k)!)
			var runs int32 = 1
			if i <= (n - k) {
				runs = 2
			}

			if isPrimeArr[i] {
				exp[i].Add(runs)
				fmt.Println("exp", exp[i].Load())
				factorizationWait.Done()
			} else { //need to do the prime factorization of i if it's not prime
				//factorize in a coroutine

				go func(n int) {
					defer factorizationWait.Done()

					//first run is n!/k!
					//second run is n!/(k!(n-k)!)
					for range runs {
						ncurr := n
						h := 2
						for ncurr != 1 {
							if isPrimeArr[h] && ncurr%h == 0 {
								exp[h].Add(-1)
								fmt.Println("exp", exp[i].Load())
								ncurr /= h
							} else {
								h++
							}
						}
					}

				}(i)
			}
		}
		//wait for all processing to finish
		factorizationWait.Wait()

		var res int32 = 1
		for i := range 432 {
			res *= exp[i].Load()
		}
		fmt.Println(res)
	}

}
