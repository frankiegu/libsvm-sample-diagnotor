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

type DiagnoseIndecis struct {
	avgWidth int            /* average number of features in one row */
	rowCnt   int            /* number of sample rows */
	maxWidth int            /* max number of features in one row */
	minWidth int            /* min number of reatures in one row */
	positive int            /* number of positive samples */
	negative int            /* nmuber of negative samples */
	features int            /* not in use */
	coverage map[string]int /* storage of feature and coresponding coverage */
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func Diagnose(cfg map[string]float32) interface{} {
	/*
		Performing a diagnose on sample rows.

		input:
			cfg        a map[string]float32 contains all thresholds which loaded
			           from configuration file or specified by command line options.

		output:
			Write diagnose results into different files.

			    sample.sumit.txt           :  contains all summarys(rowCnt, avg of features, etc.)
			    feature.coverage_more.txt  :  contains all features whose coverage are more than the coverage upper bound.
			    feature.coverage_less.txt  :  contains all features whose coverage are less than the coverage lower bound.
			    feature.coverage_full.txt  :  contains all features coverage(the raw count number, not in precent)
	*/

	var dig = DiagnoseIndecis{0, 0, MinInt, MaxInt, 0, 0, 0, make(map[string]int, 1)}
	var featureSum int = 0

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF || len(line) == 0 {
			break
		}
		check(err)

		fields := strings.Split(strings.Trim(line, " \t\n"), "\t")
		if len(fields) < 2 {
			fmt.Println("Expected samples from stdin; Abort!")
			os.Exit(-1)
		}

		dig.rowCnt += 1 /* 总行数 */
		width := len(fields) - 1
		featureSum += width

		if fields[0] == "-1" {
			dig.negative += 1
		} else {
			dig.positive += 1
		}

		for _, field := range fields[1:] {
			var feature = strings.Split(field, ":")[0]
			dig.coverage[feature] += 1
		}

		if width < dig.minWidth {
			dig.minWidth = width /* 最小特征数 */
		}

		if width > dig.maxWidth {
			dig.maxWidth = width /* 最大特征数 */
		}

	}

	file, err := os.Create("sample.sumit.txt")
	check(err)

	file.WriteString(fmt.Sprintf("feature.max = %d\n", dig.maxWidth))
	file.WriteString(fmt.Sprintf("feature.min = %d\n", dig.minWidth))
	file.WriteString(fmt.Sprintf("feature.avg = %.3f\n", float32(featureSum)/float32(dig.rowCnt)))
	file.WriteString(fmt.Sprintf("positive    = %d\n", dig.positive))
	file.WriteString(fmt.Sprintf("negative    = %d\n", dig.negative))
	file.WriteString(fmt.Sprintf("totalcnt    = %d (%d + %d)\n", dig.rowCnt, dig.positive, dig.negative))
	file.Close()

	fileMore, err := os.Create("feature.coverage_more.txt")
	check(err)

	fileLess, err := os.Create("feature.coverage_less.txt")
	check(err)

	fileFull, err := os.Create("feature.coverage_full.txt")
	check(err)

	for k, v := range dig.coverage {
		ratio := float32(v) / float32(dig.rowCnt)
		if ratio > cfg["cover_max"] {
			_, err = fileMore.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
			check(err)
		} else if ratio < cfg["cover_min"] {
			_, err = fileLess.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
			check(err)
		}
		fileFull.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
	}

	fileMore.Close()
	fileLess.Close()
	fileFull.Close()

	return nil
}
