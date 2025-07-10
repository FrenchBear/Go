// g44_wc_parallel.go
// Learning go, Concurrent programming, Tests with wc using parallism
//
// 2025-07-10	PV		First version

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {
	fmt.Println("Go wc Parallel")

	test(100)
	test(200)
	test(500)
	test(1000)
	test(2000)
	test(5000)
	test(6000)
	test(7000)
	test(8000)
	test(9000)
	test(10000)
	test(20000)
	test(50000)
	test(100000)
}

func test(blocksize int) {
	start := time.Now()
	res := count(blocksize)
	duration := time.Since(start)
	fmt.Printf("Blocksize %6d: li:%6d  wo:%8d ru:%8d by:%8d, Duration: %.3fs\n", blocksize, res.lines, res.words, res.runes, res.bytes, float64(duration.Milliseconds())/1000.0)
}

type WCRes struct {
	lines int
	words int
	runes int
	bytes int
}

func count(blocksize int) WCRes {
	path := `C:\Development\TestFiles\Text\Les secrets d'Hermione.txt`

	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()

	reschan := make(chan WCRes, 100)
	sl := 0

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) >= blocksize {
			sl++
			//fmt.Println("go count_slice", sl)
			go count_slice(sl, lines, reschan)
			lines = nil
		}
	}
	if len(lines) > 0 {
		sl++
		//fmt.Println("go final count_slice", sl)
		go count_slice(sl, lines, reschan)
	}

	// fmt.Println("Wait for goroutines to end")
	// wg.Wait()

	// fmt.Println("Read and cumulate results")
	total := WCRes{}
	for i := 0; i < sl; i++ {
		res := <-reschan
		total.lines += res.lines
		total.words += res.words
		total.runes += res.runes
		total.bytes += res.bytes
	}

	return total
}

func count_slice(_ int, lines []string, reschan chan WCRes) {
	//fmt.Println("Start count_slice", sl)
	cnt := WCRes{}
	for _, line := range lines {
		cnt.lines++
		cnt.runes += utf8.RuneCountInString(line)
		cnt.bytes += len(line) + 1 // +1, arbitrary length of end-of-line, even if it's \r\n but we don't know -- Should replace that by file length

		splitFunc := func(r rune) bool {
			return r == ' ' || r == '\t'
		}
		cnt.words += len(strings.FieldsFunc(strings.Trim(line, " \t"), splitFunc))
	}
	reschan <- cnt
	//fmt.Println("End count_slice", sl)
}
