// g44_wc_parallel.go
// Learning go, Concurrent programming, Tests with wc using parallism
//
// 2025-07-10	PV		First version

/*
2025-07-10 on WOTAN, repeats=10, blocksize around 6000 is constantly the best

Go wc Parallel
Blocksize    100: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.017s
Blocksize    200: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.017s
Blocksize    500: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.017s
Blocksize   1000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.015s
Blocksize   2000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.011s
Blocksize   3000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.007s
Blocksize   4000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.007s
Blocksize   5000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize   6000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.006s
Blocksize   7000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize   8000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.008s
Blocksize   9000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.007s
Blocksize  10000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize  11000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.008s
Blocksize  12000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.008s
Blocksize  13000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize  14000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize  15000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.009s
Blocksize  16000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.010s
Blocksize  17000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.010s
Blocksize  18000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.011s
Blocksize  19000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.011s
Blocksize  20000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.011s
Blocksize  50000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.021s
Blocksize 100000: li: 56337  wo: 1363732 ru: 7833526 by: 8434125, Duration: 0.025s
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {
	fmt.Println("Go wc Parallel")

	test(100)
	test(200)
	test(500)
	for bs:=1000;bs<=20000;bs+=1000 {
		test(bs)
	}
	test(50000)
	test(100000)
}

func test(blocksize int) {
	var res WCRes
	const REPEATS=10
	times := make([]float64, REPEATS)
	for repeat := 0; repeat < REPEATS; repeat++ {
		start := time.Now()
		res = count(blocksize)
		duration := time.Since(start)
		times = append(times, float64(duration.Milliseconds())/1000.0)
	}
	d := Median(times)
	fmt.Printf("Blocksize %6d: li:%6d  wo:%8d ru:%8d by:%8d, Duration: %.3fs\n", blocksize, res.lines, res.words, res.runes, res.bytes, d)
}

func Median(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)
	sort.Float64s(dataCopy)
	n := len(dataCopy)
	mid := n / 2
	if n%2 != 0 {
		return dataCopy[mid]
	}
	return (dataCopy[mid-1] + dataCopy[mid]) / 2.0
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
