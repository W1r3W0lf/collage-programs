package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

func mandlebrot(z complex128 ,maxiter int)int {
	c := z
	for n :=0 ; n < maxiter; n++{
		if real(z)*real(z)+imag(z)*imag(z) > 4{
			return n
		}
		z = z*z + c
	}
	return 0
}

func mandelbrotThread( xslice [][]complex128, output [][]int, maxiter int, wg *sync.WaitGroup){

	for x := 0 ; x < len(xslice) ; x++ {
		for y := 0 ; y < len(xslice[0]) ; y++ {
			output[x][y] = mandlebrot(xslice[x][y], maxiter)
		}
	}
	wg.Done()
}

func mandelbrotSet(xmin, xmax, ymin, ymax float64, width, height, maxiter, threads int)[][]int{

	var complexPlane = make([][]complex128, width)
	var outputImage = make([][]int, width)

	for i := range complexPlane {
		complexPlane[i] = make([]complex128, height)
		outputImage[i] = make([]int, height)
	}

	dx := (xmax-xmin)/float64(width)
	dy := (ymax-ymin)/float64(width)

	for x := 0 ; x < width ; x++ {
		for y := 0 ; y < height ; y++ {
			r := xmin + float64(x) * dx
			i := ymin + float64(y) * dy
			complexPlane[x][y] = complex(r,i)
		}
	}


	var wg sync.WaitGroup
	wg.Add(threads)

	for t := 0 ; t < threads ; t++ {
		start := int(float64(t) / float64(threads) * float64(width))
		end := int(float64(t+1) / float64(threads) * float64(width))

		go mandelbrotThread( complexPlane[ start : end ], outputImage[ start : end ], maxiter, &wg)
	}

	wg.Wait()

	return outputImage
}


func writeImageSlice(image *image.RGBA, start, end int, rawValues [][]int, wg *sync.WaitGroup){

	for y:= start ; y < end ; y++{
		for x:=0 ; x < image.Stride ; x+=4{
			image.Pix[x+y*image.Stride+0] = uint8(rawValues[x/4][y])*255
			image.Pix[x+y*image.Stride+1] = uint8(rawValues[x/4][y])*255
			image.Pix[x+y*image.Stride+2] = uint8(rawValues[x/4][y])*255
			image.Pix[x+y*image.Stride+3] = 255
		}
	}
	wg.Done()
}

func writeImage(image *image.RGBA, height int, rawValues [][]int, threads int){


	var wg sync.WaitGroup
	wg.Add(threads)

	for t := 0 ; t < threads ; t++ {
		start := int(float64(t) / float64(threads)*float64(height))
		end := int(float64(t+1) / float64(threads)*float64(height))
		go writeImageSlice(image, start, end, rawValues, &wg)
	}

	wg.Wait()
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

	fmt.Println("Start")
	width := 10000
	height := 10000
	threads := 8

	var scale float64 = 2
	var x float64 = 0
	var y float64 = 0


	var mandelbrotImage = image.NewRGBA(image.Rect(0, 0, width, height))

	fmt.Println("Calculating Values")
	var rawValues = mandelbrotSet(-scale+x, scale+x, -scale+y, scale+y, width, height, 100, threads)

	fmt.Println("Converting 2DSlice to RGBA")
	writeImage(mandelbrotImage, height, rawValues, threads)

	fmt.Println("Writing Image")
	var imageOut, err= os.Create("mand.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(imageOut, mandelbrotImage)

	imageOut.Close()



}
