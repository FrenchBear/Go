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

const path string = `C:\Development\TestFiles\Text\Les secrets d'Hermione.txt`

func main() {
	fmt.Println("Go wc Parallel")

	// test_block(100)
	// test_block(200)
	// test_block(500)
	// for bs := 1000; bs <= 20000; bs += 1000 {
	// 	test_block(bs)
	// }
	// test_block(50000)
	// test_block(100000)

	test_block(6000)

	test(count_linear_1, "Linear 1")
	test(count_linear_2, "Linear 2")
	test(count_linear_3, "Linear 3")
	test(count_linear_4, "Linear 4")
	test(count_parallel_readall, "B6000 readall")
}

func test(count_linear_func func() WCRes, name string) {
	var res WCRes
	const REPEATS = 10
	times := make([]float64, REPEATS)
	for repeat := 0; repeat < REPEATS; repeat++ {
		start := time.Now()
		res = count_linear_func()
		duration := time.Since(start)
		times = append(times, float64(duration.Milliseconds())/1000.0)
	}
	d := Median(times)
	fmt.Printf("%-16s  li:%6d  wo:%8d ru:%8d by:%8d, Duration: %.3fs\n", name, res.lines, res.words, res.runes, res.bytes, d)
}

func test_block(blocksize int) {
	var res WCRes
	const REPEATS = 10
	times := make([]float64, REPEATS)
	for repeat := 0; repeat < REPEATS; repeat++ {
		start := time.Now()
		res = count_parallel(blocksize)
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

func count_linear_1() WCRes {
	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return count_slice_core(lines)
}

func count_linear_2() WCRes {
	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()
	res := WCRes{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count_line(&res, scanner.Text())
	}
	return res
}

func count_linear_3() WCRes {
	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()
	res := WCRes{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		go count_line(&res, scanner.Text())
	}
	return res
}

func count_linear_4() WCRes {
	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()

	data := make(chan string)
	result := make(chan WCRes)

	go func() {
		res := WCRes{}
		for s := range data {
			count_line(&res, s)
		}
		result <- res
		close(result)
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data <- scanner.Text()
	}
	close(data)
	return <-result
}

func count_parallel_readall() WCRes {
	file, err := os.Open(path)
	if err != nil {
		return WCRes{}
	}
	defer file.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		return WCRes{}
	}

	reschan := make(chan WCRes, 100)

	lines := strings.Split(strings.ReplaceAll(strings.ReplaceAll(string(data), "\r\n", "\n"), "\r", "\n"), "\n")
	SLICESIZE := 6000
	sl := 0
	for i := 0; i < len(lines); i += SLICESIZE {
		end := i + SLICESIZE
		if end > len(lines) {
			end = len(lines)
		}
		sl++
		go count_slice_to_reschan(lines[i:end], reschan)
	}

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

func count_parallel(blocksize int) WCRes {
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
			go count_slice_to_reschan(lines, reschan)
			lines = nil
		}
	}
	if len(lines) > 0 {
		sl++
		go count_slice_to_reschan(lines, reschan)
	}

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

func count_slice_to_reschan(lines []string, reschan chan WCRes) {
	reschan <- count_slice_core(lines)
}

func count_slice_core(lines []string) WCRes {
	cnt := WCRes{}
	for _, line := range lines {
		count_line(&cnt, line)
	}
	return cnt
}

func count_line(cnt *WCRes, line string) {
	cnt.lines++
	cnt.runes += utf8.RuneCountInString(line)
	cnt.bytes += len(line) + 1 // +1, arbitrary length of end-of-line, even if it's \r\n but we don't know -- Should replace that by file length

	splitFunc := func(r rune) bool {
		return r == ' ' || r == '\t'
	}
	cnt.words += len(strings.FieldsFunc(strings.Trim(line, " \t"), splitFunc))
}
