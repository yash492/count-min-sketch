package main

import (
	"fmt"
	"hash"
	"slices"
	"sync"
	"time"

	"github.com/twmb/murmur3"
)

type CountMinSketch struct {
	depth   int
	width   int
	Arr     [][]int
	HashArr []hash.Hash32
	mu      *sync.RWMutex
}

func NewCountMinSketch(hashArr []hash.Hash32, depth, width int) *CountMinSketch {
	arr := make([][]int, depth)
	for i := range arr {
		arr[i] = make([]int, width)
	}
	return &CountMinSketch{
		depth:   depth,
		width:   width,
		Arr:     arr,
		HashArr: hashArr,
		mu:      &sync.RWMutex{},
	}
}

func main() {

	arr := []string{
		"chintya",
		"ashwin",
		"sahil",
		"pintya",
		"mintya",
		"john",
		"chintya",
		"sahil",
		"ashwin",
		"yash",
		"chintya",
		"doe",
		"ashwin",
		"yash",
		"john",
		"ashwin",
		"chintya",
		"mintya",
		"mintya",
		"john",
		"mintya",
		"deepli",
		"manish",
		"deepli",
		"john",
		"deepli",
		"doe",
		"pintya",
		"yash",
		"manish",
		"doe",
		"deepli",
		"john",
		"ashwin",
		"john",
		"chintya",
		"chintya",
		"mintya",
		"manish",
		"sahil",
		"yash",
		"manish",
		"mintya",
		"yash",
		"sahil",
		"pintya",
		"mintya",
		"chintya",
		"doe",
		"sahil",
		"john",
		"manish",
		"doe",
		"doe",
		"mintya",
		"manish",
		"ashwin",
		"ashwin",
		"doe",
		"manish",
		"yash",
		"doe",
		"ashwin",
		"pintya",
		"john",
		"yash",
		"manish",
		"chintya",
		"yash",
		"sahil",
		"john",
		"doe",
		"pintya",
		"yash",
		"ashwin",
		"deepli",
	}

	hasharr := []hash.Hash32{}

	depth := 5
	width := 2716

	for i := range depth {
		hasharr = append(hasharr, murmur3.SeedNew32(uint32(i)))
	}

	realData := map[string]int{}

	for _, a := range arr {
		realData[a] += 1
	}

	fmt.Printf("Real Data count: %v\n", realData)

	countMinSketchInstance := NewCountMinSketch(hasharr, depth, width)

	valCh := make(chan string, 5)

	go func() {
		wg := &sync.WaitGroup{}

		for _, a := range arr {
			wg.Go(func() {
				valCh <- a
			})
		}

		wg.Wait()
		close(valCh)
	}()

	wg := &sync.WaitGroup{}

	i := 0
	for ch := range valCh {
		fmt.Println(ch)
		wg.Go(func() {
			countMinSketchInstance.Insert(ch)
		})
		if i > 0 && i%10 == 0 {
			wg.Wait()
		}
		i++

		fmt.Println(i)
	}

	wg.Wait()

	query := "yash"
	count := countMinSketchInstance.GetCount(query)
	fmt.Printf("count for %v: %v\n", query, count)

}

func (c *CountMinSketch) Insert(val string) {
	time.Sleep(5 * time.Second)
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.depth {
		hasher := c.HashArr[i]
		hasher.Write([]byte(val))
		columnIdx := hasher.Sum32() % uint32(c.width)
		hasher.Reset()

		c.Arr[i][columnIdx] += 1
	}

}

func (c *CountMinSketch) GetCount(val string) int {
	countArr := []int{}
	for i := range c.depth {
		hasher := c.HashArr[i]
		hasher.Write([]byte(val))
		columnIdx := hasher.Sum32() % uint32(c.width)
		hasher.Reset()

		countArr = append(countArr, c.Arr[i][columnIdx])
	}

	return slices.Min(countArr)
}
