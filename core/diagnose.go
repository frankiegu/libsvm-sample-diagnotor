/*
The MIT License (MIT)

Copyright (c) [2015] [liangchengming]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import "os"
import "io"
import "bufio"
import "strings"
import "fmt"

type Diagnose struct {
	avg_width int
	row_cnt   int
	max_width int
	min_width int
	positive  int
	negative  int
	features  int
	coverage  map[string]int
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func diagnose(cfg map[string]interface{}) interface{} {

	var dig Diagnose = Diagnose{0, 0, MinInt, MaxInt, 0, 0, 0, make(map[string]int, 1)}
	var feature_sum int = 0

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF || len(line) == 0 {
			break
		}
		check(err)

		fields := strings.Split(strings.Trim(line, " \t\n"), "\t")
		if len(fields) < 2 {
			fmt.Println("Bad sample:%v", line)
			panic(fmt.Sprintf("Bad sample:%v", line))
		}

		dig.row_cnt += 1 /* 总行数 */
		width := len(fields) - 1
		feature_sum += width

		if fields[0] == "-1" {
			dig.negative += 1
		} else {
			dig.positive += 1
		}

		for _, field := range fields[1:] {
			var feature = strings.Split(field, ":")[0]
			dig.coverage[feature] += 1
		}

		if width < dig.min_width {
			dig.min_width = width /* 最小特征数 */
		}

		if width > dig.max_width {
			dig.max_width = width /* 最大特征数 */
		}

	}

	file, err := os.Create("sample.sumit.txt")
	check(err)

	file.WriteString(fmt.Sprintf("feature.max = %d\n", dig.max_width))
	file.WriteString(fmt.Sprintf("feature.min = %d\n", dig.min_width))
	file.WriteString(fmt.Sprintf("feature.avg = %.3f\n", float32(feature_sum)/float32(dig.row_cnt)))
	file.WriteString(fmt.Sprintf("positive    = %d\n", dig.positive))
	file.WriteString(fmt.Sprintf("negative    = %d\n", dig.negative))
	file.WriteString(fmt.Sprintf("totalcnt    = %d (%d + %d)\n", dig.row_cnt, dig.positive, dig.negative))
	file.Close()

	file, err = os.Create("feature.coverage.txt")
	check(err)

	for k, v := range dig.coverage {
		_, err = file.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
		check(err)
	}

	file.Close()

	return nil
}
