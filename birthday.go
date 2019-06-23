package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)


func bday(itter int, randFunc func(int)int)float64{
	days := make([]int, 365, 365)
	grandSum := 0
	n := 0
	counter := 0

	for i := 0 ; i < itter ; i++ {
		for ; days[n] < 2 ; counter++{
			n = randFunc(365)
			days[n]++
		}
		grandSum += counter
		counter = 0
		days = make([]int, 365, 365)
	}
	//result := math.Ceil(float64(grandSum)/float64(itter))
	//return int(result)
	return float64(grandSum)/float64(itter)
}

func bdayQ(itter, qSize int, randFunc func(int, chan int, chan bool))float64{
	days := make([]int, 365, 365)
	grandSum := 0
	n := 0
	counter := 0

	randomIn := make(chan int, qSize)
	killer := make(chan bool)

	go randFunc(365, randomIn, killer)

	for i := 0 ; i < itter ; i++ {
		for ; days[n] < 2 ; counter++{
			n = <- randomIn
			days[n]++
		}
		grandSum += counter
		counter = 0
		days = make([]int, 365, 365)
	}
	//result := math.Ceil(float64(grandSum)/float64(itter))
	//return int(result)
	return float64(grandSum)/float64(itter)
}


func randQ(max int, out chan int, kill chan bool){
	for{
		select {
		case out <- rand.Intn(max):

		case <- kill:
			return
		}
	}
}

func devrand(max int)int{
	file, err := os.Open("/dev/random")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes := make([]byte, 16)

	file.Read(bytes)

	output := int(binary.BigEndian.Uint16(bytes)) % max

	return output
}

func orgrand(max int)int{
	url := fmt.Sprintf("https://www.random.org/integers/?num=1&min=0&base=10&format=plain&rnd=new&col=1&max=%d",max-1)
	resp, err := http.Get(url)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	var number []byte

	number, err = ioutil.ReadAll(resp.Body)

	var out int

	power := 0
	for i := len(number)-2 ; i >= 0 ; i--{
		out += (int(number[i]) - 48) * int(math.Pow10(power))
		power++
	}

	return out
}

func orgrandQ(max int, out chan int, kill chan bool){
	bufferSize := 1000
	url := fmt.Sprintf("https://www.random.org/integers/?num=%d&min=0&base=10&format=plain&rnd=new&col=1&max=%d", bufferSize, max-1)
	numbers := make([]int,0)
	for{
		if len(numbers) == 0 {
			resp, err := http.Get(url)
			if err != nil{
				panic(err)
			}

			var rawBytes []byte

			rawBytes, err = ioutil.ReadAll(resp.Body)

			processedNumbers := make([]int,bufferSize)
			processedIndex := 0

			power := 0
			randomBench := 0
			for  i := len(rawBytes)-2 ; i >= 0 ; i--{

				if rawBytes[i] == 10{
					processedNumbers[processedIndex] = randomBench
					processedIndex++
					randomBench = 0
					power = 0
				}else{
					randomBench += (int(rawBytes[i])-48) * int(math.Pow10(power))
					power++
				}
			}
			// Adding the last value that is left on the bench
			processedNumbers[len(processedNumbers)-1] = randomBench

			resp.Body.Close()

			numbers = append(numbers, processedNumbers...)
		}
		select {
		case <- kill:
			return

		case out <- numbers[len(numbers)-1]:
			numbers = numbers[:len(numbers)-1]
			// This pops the last element of the slice.

		default:

		}
	}
}

func main() {

	its := 1000


	rand.Seed(time.Now().UTC().UnixNano())


	a := bday(its, rand.Intn)
	fmt.Println("Stock ", a)

	b := bday(its, devrand)
	fmt.Println("DEV ", b)

	/*
	// This is too slow for practical use
	c := bday(its, orgrand)
	fmt.Println("ORG ", c)
	*/

	//qSize := 1000


	d := bdayQ(its, 10, randQ)
	fmt.Println("StockQ ", d)

	/*
	e := bdayQ(its, 10, devrandQ)
	fmt.Println("DEV ", e)
	*/

	f := bdayQ(its, 1000, orgrandQ)
	fmt.Println("OrgQ ", f)



	/*
	for x := 0 ; x < 10 ; x++{
		//fmt.Println("DEV: ",devrand(100))
		fmt.Println("ORG: ",orgrand(100))
	}
	*/

}
