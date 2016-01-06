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
import "math"
import "bufio"
import "strings"
import "strconv"
import "fmt"

type DiagnoseIndices struct {
	avgWidth   int            /* average number of features in one row */
	rowCnt     int            /* number of sample rows */
	maxWidth   int            /* max number of features in one row */
	minWidth   int            /* min number of reatures in one row */
	positive   int            /* number of positive samples */
	negative   int            /* nmuber of negative samples */
	features   int            /* not in use */
	coverage   map[string]int /* storage of feature and coresponding coverage */
	featureSum int            /* tmp value of feautre statistic */
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func basic_statistics(cfg Configuration) {
	/*
		Performing a diagnose on sample rows.

		input:
			cfg        a Configuration object which contains all thresholds which loaded
			           from configuration file or specified by command line options.

		output:
			Write diagnose results into different files.

			    sample.sumit.txt           :  contains all summarys(rowCnt, avg of features, etc.)
			    feature.coverage_more.txt  :  contains all features whose coverage are more than the coverage upper bound.
			    feature.coverage_less.txt  :  contains all features whose coverage are less than the coverage lower bound.
	*/

	var dig = DiagnoseIndices{0, 0, MinInt, MaxInt, 0, 0, 0, make(map[string]int, 1), 0}
	var group = make(map[string]int, 1)

	if len(cfg.groupTags) > 0 {
		for _, tag := range cfg.groupTags {
			group[tag] = 0
		}
	}

	var columes = make(map[string][]int, 1)
	var labels = make(map[int]int, 1)

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

		statisticRowIndex(fields, &dig)
		statisticColumeIndex(fields, &dig, group)
		storageByColumes(fields, labels, columes, dig.rowCnt)
	}

	dumpIndices(dig, group, cfg)
	calulateFeatureMutalInformation(cfg, columes, labels)
}

func statisticRowIndex(fields []string, dig *DiagnoseIndices) {

	dig.rowCnt += 1 /* 总行数 */
	var width int = len(fields) - 1
	dig.featureSum += width

	if fields[0] == "-1" {
		dig.negative += 1
	} else {
		dig.positive += 1
	}

	if width < dig.minWidth {
		dig.minWidth = width /* 最小特征数 */
	}

	if width > dig.maxWidth {
		dig.maxWidth = width /* 最大特征数 */
	}
}

func dumpIndices(dig DiagnoseIndices, group map[string]int, cfg Configuration) {

	/* Write results into file(sample.summary.txt).  */

	file, err := os.Create("sample.summary.txt")
	check(err)

	file.WriteString(fmt.Sprintf("positive    = %d\n", dig.positive))
	file.WriteString(fmt.Sprintf("negative    = %d\n", dig.negative))
	file.WriteString(fmt.Sprintf("totalcnt    = %d (%d + %d)\n", dig.rowCnt, dig.positive, dig.negative))

	file.WriteString(fmt.Sprintf("feature.max = %d\n", dig.maxWidth))
	file.WriteString(fmt.Sprintf("feature.min = %d\n", dig.minWidth))
	file.WriteString(fmt.Sprintf("feature.avg = %.3f\n", float32(dig.featureSum)/float32(dig.rowCnt)))

	for tag, cnt := range group {
		file.WriteString(fmt.Sprintf("%s.coverage = %.3f\n", tag, float32(cnt)/float32(dig.rowCnt)))
	}

	file.Close()

	fileMore, err := os.Create("feature.coverage_more.txt")
	check(err)

	fileLess, err := os.Create("feature.coverage_less.txt")
	check(err)

	for k, v := range dig.coverage {
		ratio := float32(v) / float32(dig.rowCnt)
		if ratio > cfg.thresholds["cover_max"] {
			_, err = fileMore.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
			check(err)
		} else if ratio < cfg.thresholds["cover_min"] {
			_, err = fileLess.WriteString(fmt.Sprintf("%s\t%d\n", k, v))
			check(err)
		}
	}

	fileMore.Close()
	fileLess.Close()
}

func statisticColumeIndex(fields []string, dig *DiagnoseIndices, group map[string]int) {
	/*
		Statistic indices by sample columes

		input:
			fields        a list of sample columes(label and features) from a single sample row.
			dig           storage of indices
			group         a hash stroage of feature tags and sample numbers

		output:
			fulfilled <dig> and <group>
	*/

	var mathcedTags = make(map[string]interface{})
	for _, featureName := range fields {
		for tag, _ := range group {
			_, exist := mathcedTags[tag]
			if strings.Contains(featureName, tag) && exist == false {
				group[tag] += 1
				mathcedTags[tag] = nil
			}
		}
	}

	for _, field := range fields[1:] {
		var feature = strings.Split(field, ":")[0]
		dig.coverage[feature] += 1
	}
}

func storageByColumes(fields []string, labels map[int]int, columes map[string][]int, lineno int) {
	/*
		Stoage sample labels and features by colume.

		input:
			fields        a list of sample columes(label and features) from a single sample row.
			labels        a hash storage of sample labels
			columes       a hash storage of feature name and the feature colume
			lineno        line number of the single samle(which had been split into <fields>)

		outpu:
			fulfilled <labels> and <columes>.
	*/
	for i, s := range fields {
		if i == 0 {
			sample_class_label, err := strconv.ParseInt(s, 10, 32)
			check(err)
			labels[lineno] = int(sample_class_label)
		}
		if v, ok := columes[s]; ok == false {
			columes[s] = []int{lineno}
		} else {
			columes[s] = append(v, lineno)
		}
	}
}

func calulateFeatureMutalInformation(cfg Configuration, columes map[string][]int, labels map[int]int) {
	/*
		Calulate the mutal infermation between given feature and sample labels

		input:
			labels        a hash of sample-id(lineno) and its label
			columes       a hash of feature name and coresponding sample-id(lineno) vector
			cfg           configurations

		output:
			a file named "feature.mutal.info" which contains feature-names and
			its mutal information writen in the same line.

	*/

	if cfg.enable_mi == false {
		return
	}

	var positive int = 0
	for _, v := range labels {
		if v == 1 {
			positive += 1
		}
	}

	if positive == len(labels) || positive == 0 {
		fmt.Println("Bad Samples: Both negative and positive samples are needed!")
		os.Exit(-1)
	}

	mutal, err := os.Create("feature.mutal.info")
	check(err)

	var total = float64(len(labels))

	var plabel1 int = positive
	var plabel0 int = len(labels) - positive

	for k, v := range columes {

		var pfeatu1 int = len(v)
		var pfeatu0 int = len(labels) - len(v)

		var pxy00 = probability(labels, v, "00")
		var pxy10 = probability(labels, v, "10")
		var pxy01 = probability(labels, v, "01")
		var pxy11 = probability(labels, v, "11")

		var mi float64 = 0.0
		if pxy00 != 0 {
			mi += float64(pxy00) / total * math.Log(total*float64(pxy00)/float64(pfeatu0*plabel0))
		}
		if pxy10 != 0 {
			mi += float64(pxy10) / total * math.Log(total*float64(pxy10)/float64(pfeatu1*plabel0))
		}
		if pxy01 != 0 {
			mi += float64(pxy01) / total * math.Log(total*float64(pxy01)/float64(pfeatu0*plabel1))
		}
		if pxy11 != 0 {
			mi += float64(pxy11) / total * math.Log(total*float64(pxy11)/float64(pfeatu1*plabel1))
		}

		mutal.WriteString(fmt.Sprintf("%s\t%.3f\n", k, mi*float64(10000)))
	}

	mutal.Close()
}

func probability(labels map[int]int, feature []int, mode string) int {
	/*
		calculate probabily of feautre (in cnt)

		input:
			labels        a hash of sample-id(lineno) and its label
			feature       a hash of feature name and sample-id(lineno)
			mode          a string which value should be one of {"00", "10", "01", "11"}.
			              different cnt value will be returned according to "mode".

			ex:
			    mode="00" means counting the number of negative sample without given feature.
			    mode="10" means counting the number of negative sample with the given feature.
				(vice versa)


		output:
			return the cnt(type int) according to different "mode"
	*/

	var sub int = 0

	switch mode {
	case "00": /* (feature=0, label=negative) */
		var curr int = 0
		for i, v := range labels {
			if i+1 == feature[curr] {
				if curr+1 < len(feature) {
					curr++
				}
				continue
			}
			if v == -1 {
				sub += 1
			}
		}
		return sub
	case "01": /* (feature=0, label=positive) */
		var curr int = 0
		for i, v := range labels {
			if i+1 == feature[curr] {
				if curr+1 < len(feature) {
					curr++
				}
				continue
			}
			if v == 1 {
				sub += 1
			}
		}
		return sub
	case "10": /* (feature=1, label=negative) */
		for _, k := range feature {
			if labels[k] == -1 {
				sub += 1
			}
		}
		return sub
	case "11": /* (feature=1, label=positive) */
		for _, k := range feature {
			if labels[k] == 1 {
				sub += 1
			}
		}
		return sub
	default:
		panic("Unknown mode, Abort!")
	}
}
