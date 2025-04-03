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
	isPrimeArr := make([]bool, 433)

	//how many threads we'll need to wait for
	primeWait := new(sync.WaitGroup)
	primeWait.Add(433)
	//concurrently check prime numbers
	for iterationNumber := range 433 {
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
		n_int, _ := strconv.Atoi(nStr)
		n := int64(n_int)
		k_int, _ := strconv.Atoi(kStr)
		k := int64(k_int)

		exp := make([]atomic.Int64, 433)
		for i := range exp {
			exp[i].Store(0) // Explicitly initialize to 0
		}

		//for any iterations we want to run concurrently
		factorizationWait := new(sync.WaitGroup)
		//make exp the prime factorization of n!/k!
		factorizationWait.Add(int(n - 1))
		for i := int64(2); i <= n; i++ {

			//just coroutine shit
			go func(passedN int64) {
				defer factorizationWait.Done()

				if isPrimeArr[passedN] {
					exp[passedN].Add(1)
				} else { //need to do the prime factorization of i if it's not prime
					//first run is passedN!/k!
					ncurr := passedN
					var h int64 = 2
					for ncurr != 1 {
						if isPrimeArr[h] && ncurr%h == 0 {
							exp[h].Add(1)
							ncurr /= h
						} else {
							h++
						}
					}

				}

				divis := func() {
					if isPrimeArr[passedN] {
						exp[passedN].Add(-1)
					} else { //need to do the prime factorization of i if it's not prime
						//second run is passedN!/(k!(passedN-k)!)
						ncurr := passedN
						var h int64 = 2
						for ncurr != 1 {
							if isPrimeArr[h] && ncurr%h == 0 {
								exp[h].Add(-1)
								//fmt.Println("exp false", exp[i].Load())
								ncurr /= h
							} else {
								h++
							}
						}

					}
				}

				if passedN <= (k) {
					divis()
				}
				if passedN <= (n - k) {
					divis()
				}
			}(i)
		}
		//wait for all processing to finish
		factorizationWait.Wait()

		var res int64 = 1
		for i := range 433 {
			res *= exp[i].Load() + 1
		}
		fmt.Println(res)
	}

}
