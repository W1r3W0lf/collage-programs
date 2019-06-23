package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
)

func murge(numbers []int, c chan int){
	if len(numbers) < 2 {
		c <- numbers[0]
		close(c)
		return
	}

	a := make(chan int, len(numbers)/2)
	b := make(chan int, len(numbers)/2)

	go murge(numbers[:len(numbers)/2], a)
	go murge(numbers[len(numbers)/2:], b)

	av, ac := <-a
	bv, bc := <-b

	for ac && bc {
		if av < bv{
			c <- av
			av, ac = <-a
		} else {
			c <- bv
			bv, bc = <-b
		}
	}

	for ac {
		c <- av
		av,ac = <-a
	}

	for bc {
		c <- bv
		bv,bc = <-b
	}

	close(c)
}

func insertion(numbers []int){
	i := 0
	j := 0
	var temp int
	for i < len(numbers){
		j = i
		for j > 0 && numbers[j-1] > numbers[j]{

			temp = numbers[j]
			numbers[j] = numbers[j-1]
			numbers[j-1] = temp

			j--
		}
		i++
	}
}

func hybridMurge(numbers []int, g *sync.WaitGroup){
	if len(numbers) <= 10 {
		insertion(numbers)
		g.Done()
		return
	}

	i := 0
	j := len(numbers)/2

	var wg sync.WaitGroup
	wg.Add(2)

	go hybridMurge(numbers[:j], &wg)
	go hybridMurge(numbers[j:], &wg)

	wg.Wait()

	temp := make([]int, len(numbers))
	copy(temp,numbers)

	op := 0

	for i<len(temp)/2 && j<len(temp){
		if temp[i] < temp[j]{
			numbers[op] = temp[i]
			i++
		} else {
			numbers[op] = temp[j]
			j++
		}
		op++
	}

	if i < len(temp)/2{
		copy(numbers[op:], temp[i:])
	}
	if j < len(temp){
		copy(numbers[op:],temp[j:])
	}

	g.Done()
}


var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")


func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	nl := 10000000
	numbers := make([]int,nl)
	for x:=0 ; x<len(numbers) ; x++{
		numbers[x] =  rand.Intn(20)
	}

	//fmt.Println("Slice Made")

	//fmt.Println(numbers)

	/*
	c := make(chan int, len(numbers))

	go murge(numbers,c)

	for x := 0 ; x<len(numbers) ; x++{
		numbers[x] = <- c
	}
	fmt.Println("merge Done")

	 */


	var wg sync.WaitGroup
	wg.Add(1)

	hybridMurge(numbers, &wg)
	wg.Wait()
	fmt.Println("hybrid Done")
	//fmt.Println(numbers)



	//sort.Ints(numbers)
	//fmt.Println("Defautl Done")

}
