package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)


func openWords(path string)[]string{

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words
}

/*////////////////////////////
	Hash Table and helper functions
////////////////////////////*/

type HashTable struct{
	values [][]string
	p uint64
	k uint64
	Bcol int
	col int
}

func newHashTable(size int, p, k uint64) *HashTable{
	var ht HashTable
	ht.values = make([][]string, size, size)
	ht.p = p
	ht.k = k
	ht.col = 0
	ht.Bcol = 0

	for i := 0 ; i < len(ht.values) ; i++{
		ht.values[i] = make([]string, 0)
	}

	return &ht
}

func (h HashTable) add(str string) {
	index := h.hash(str)
	h.values[index] = append(h.values[index], str)
}

func (h *HashTable) calcStats(){
	if h.col != 0 {
		return
	}

	for x := 0 ; x < len(h.values) ; x++{
		if len(h.values[x]) > 1{
			h.Bcol++
			h.col += len(h.values[x])-1
		}
	}

}

func (h HashTable) printStats(){

	h.calcStats()

	fmt.Println("Total Buckets with colitions", h.col)
	fmt.Println("p: ", h.p)
	fmt.Println("k: ", h.k)
	fmt.Print("\n")
}

func wordTotal(word string)uint64 {
	
	var total uint64 = 1
	for i:=0 ; i< len(word); i++{
		total *= uint64(int(word[i])/(i+1))
	}
	total *= uint64(len(word))
	return total
}

func (h HashTable) hash(word string)uint64{
	var x = wordTotal(word)
	//var p uint64 = 100123456789
	var n = uint64(len(h.values))
	//var k uint64 = 887
	return uint64(((h.k*x)%h.p)%n)
}

/*/////////////////////
	Prime Stuff
/////////////////////*/


func isPrime(number uint64) bool{
	if number % 2 == 0 && number > 2{
		return false
	}
	for t := uint64(3) ; t < uint64(math.Sqrt(float64(number))) ; t+=2{
		if number % t == 0{
			return false
		}
	}
	return true
}

func primeList(start uint64, numPrimes int) []uint64{
	primes := make([]uint64, numPrimes, numPrimes)

	if start % 2 == 0{
		start += 1
	}

	var test = start

	for i:= 0 ; i < numPrimes ; i++{
		for ! isPrime(test){
			test +=2
		}
		primes[i] = test
		test +=2
	}

	return primes
}

func primeRange(start, end uint64) []uint64{
	primes := make([]uint64,0)

	if start %2 == 0{
		start += 1
	}

	var test = start

	for test < end{
		if isPrime(test){
			primes = append(primes, test)
		}
		test += 2
	}

	return primes
}

func pRThread(subRange []uint64, primePipe chan []uint64){
	miniPrimeList := make([]uint64,0)

	for _, number := range subRange{
		if isPrime(number){
			miniPrimeList = append(miniPrimeList, number)
		}
	}
	primePipe <- miniPrimeList
}

func primeRangeT(start, end uint64, threads int) []uint64{
	primes := make([]uint64,0)
	numberRange := make([]uint64, (end-start)/2, (end-start)/2)

	if start %2 == 0{
		start += 1
	}

	i := 0
	for x := uint64(start) ; x < end ; x+=2{
		numberRange[i] = x
		i++
	}

	var primePipe = make(chan []uint64, threads)

	for t:=0 ; t < threads ; t++{

		tStart := int(float64(t)/float64(threads)*float64(len(numberRange)))
		tEnd := int(float64(t+1)/float64(threads)*float64(len(numberRange)))
		go pRThread(numberRange[tStart:tEnd], primePipe)
	}


	for t:=0 ;t < threads ; t++{
		primes = append(primes, <-primePipe...)
	}


	sort.Slice(primes, func(i, j int) bool {return primes[i] < primes[j]})
	return primes
}



func hashTableTesting(){
	ht := newHashTable(2999, 100123456789, 887)

	words := openWords("/home/wire_wolf/programs/go/awesomeProject/common-words.txt")

	for _, word := range words{
		ht.add(word)
	}

	/*
		tableSizes := primeRangeT(6,30000,8)
		//tableSizes := primeRangeT(2000,3000,2)
		//tableSizes := primeList(2000, 129)
		fmt.Println(tableSizes[0])
		fmt.Println(tableSizes[len(tableSizes)-1])

	*/


	//pl := primeList(500, 100000)
	pl := primeRangeT(500,1000000,8)

	least := 20000
	var beatVal uint64

	for x := 0 ; x < len(pl) ; x++{
		test := newHashTable(2999, 100123456789, pl[x])

		for _, word := range words{
			test.add(word)
		}
		test.calcStats()

		if test.col < least{
			least = test.col
			beatVal = test.k
			test.printStats()
		}

	}

	fmt.Println("Best # ,",least, " val",beatVal)

	//ht.stats()
}




func main() {

	hashTableTesting()

}
