package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// check solutions (try to use mutex)
// why does my solutions takes same time as with 1 goroutine
const (
	defaultGoroutineNumber = 20
	defaultFilename        = "./test/asd.json"
)

type Objects []Object

type Object struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {

	args := os.Args[1:]

	var (
		goroutineNumber int
		fileName        string
	)

	switch len(args) {
	case 0:
		goroutineNumber = defaultGoroutineNumber
		fileName = defaultFilename
	case 1:
		inputGoroutineNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error while getting goroutine number: %v\n", err)
			return
		}
		if inputGoroutineNumber <= 0 {
			fmt.Printf("goroutine number is less than required 0\n")
			return
		}
		goroutineNumber = inputGoroutineNumber
		fileName = defaultFilename
	default:
		inputGoroutineNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error while getting goroutine number: %v\n", err)
			return
		}
		goroutineNumber = inputGoroutineNumber
		fileName = args[1]
	}

	input, err := getInputsFromFile(fileName)
	if err != nil {
		fmt.Printf("error while getting inputs from file: %v\n", err)
	}

	start := time.Now()

	res := solution2(goroutineNumber, input)

	resTime := time.Since(start)

	fmt.Printf("%v\n", resTime)

	fmt.Printf("%d\n", res)
}

func getInputsFromFile(file string) (Objects, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var input Objects
	if err = json.Unmarshal(b, &input); err != nil {
		return nil, err
	}

	return input, nil
}

func splitObjects(objects Objects, goroutineNum int) []Objects {
	size := len(objects) / goroutineNum
	extra := len(objects) % goroutineNum

	res := []Objects{}
	start := 0
	for i := 0; i < goroutineNum; i++ {
		end := start + size
		if extra > 0 {
			end++
			extra--
		}
		res = append(res, objects[start:end])
		start = end
	}

	return res
}

// just compare with 1 goroutine solution
func notSolution(goroutineNumber int, objects Objects) int {
	s := 0
	for _, obj := range objects {
		s += obj.A + obj.B
	}
	return s
}

// first solution that I came up with
// the idea is to split big object into small ones, and sum small parts of it
func solution1(goroutineNum int, input Objects) int {
	sum := 0

	// we do not need to create extra goroutines to get sum of objects
	// so why cant we just count everything in main goroutine?
	if goroutineNum == 1 {
		for _, obj := range input {
			sum += obj.A + obj.B
		}

		return sum
	}

	var (
		wg          = sync.WaitGroup{}
		ch          = make(chan int, goroutineNum)
		splitObject = splitObjects(input, goroutineNum)
	)

	for _, obj := range splitObject {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := 0
			for _, o := range obj {
				s += o.A + o.B
			}
			ch <- s
		}()
	}

	wg.Wait()
	close(ch)

	for s := range ch {
		sum += s
	}

	return sum
}

// same solution, but I decided to not iterate twice over 1 list of objects
func solution2(goroutineNum int, input Objects) int {
	sum := 0

	if goroutineNum == 1 {
		for _, obj := range input {
			sum += obj.A + obj.B
		}

		return sum
	}

	var (
		wg    = sync.WaitGroup{}
		ch    = make(chan int, goroutineNum)
		size  = len(input) / goroutineNum
		extra = len(input) % goroutineNum
		start = 0
	)

	for i := 0; i < goroutineNum; i++ {
		end := start + size
		if extra > 0 {
			end++
			extra--
		}
		tmp := input[start:end]
		start = end
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := 0
			for _, o := range tmp {
				s += o.A + o.B
			}
			ch <- s
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for s := range ch {
		sum += s
	}

	return sum
}

func solution3(goroutineNum int, input Objects) int {
	sum := 0

	if goroutineNum == 1 {
		for _, obj := range input {
			sum += obj.A + obj.B
		}

		return sum
	}

	var (
		wg    = sync.WaitGroup{}
		mux   = sync.Mutex{}
		size  = len(input) / goroutineNum
		extra = len(input) % goroutineNum
		start = 0
	)

	for i := 0; i < goroutineNum; i++ {
		end := start + size
		if extra > 0 {
			end++
			extra--
		}
		tmp := input[start:end]
		start = end
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := 0
			for _, o := range tmp {
				s += o.A + o.B
			}
			mux.Lock()
			sum += s
			mux.Unlock()
		}()
	}

	wg.Wait()

	return sum
}

func solution4(goroutineNum int, input Objects) int {
	sum := 0

	if goroutineNum == 1 {
		for _, obj := range input {
			sum += obj.A + obj.B
		}

		return sum
	}

	var (
		wg    = sync.WaitGroup{}
		size  = len(input) / goroutineNum
		extra = len(input) % goroutineNum
		start = 0
	)

	tmpSum := int64(sum)

	for i := 0; i < goroutineNum; i++ {
		end := start + size
		if extra > 0 {
			end++
			extra--
		}
		tmp := input[start:end]
		start = end
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := 0
			for _, o := range tmp {
				s += o.A + o.B
			}
			atomic.AddInt64(&tmpSum, int64(s))
		}()
	}
	wg.Wait()

	sum = int(tmpSum)

	return sum
}
