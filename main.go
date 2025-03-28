package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func prime(n int) bool {
	k := int(math.Sqrt(float64(n)))

	for iR := range k - 2 {
		i := iR + 2
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

		//for any iterations we want to run concurrently
		factorizationWait := new(sync.WaitGroup)
		//make exp the prime factorization of n!/k!
		for i := 2; i <= n; i++ {

			//n!/(k!(n-k)!)
			var runs int32 = 1
			if i <= (n - k) {
				runs = 2
			}

			if isPrimeArr[i] {
				exp[i].Add(runs)
			} else { //need to do the prime factorization of i if it's not prime
				//factorize in a coroutine
				factorizationWait.Add(1)
				go func(n int) {
					defer factorizationWait.Done()

					//n!/(k!(n-k)!)
					for range runs {
						ncurr := n
						h := 2
						for ncurr != 1 {
							if isPrimeArr[h] && ncurr%h == 0 {
								exp[h].Add(-1)
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

	}

}
