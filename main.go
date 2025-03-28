package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
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
	primeWait := new(sync.WaitGroup)
	isPrimeArr := make([]bool, 432)

	//concurrently check prime numbers
	for iterationNumber := range 432 {
		primeWait.Add(1)
		go func(routineInput int) {
			defer primeWait.Done()
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

		exp := make([]int, 432)

	}

}
