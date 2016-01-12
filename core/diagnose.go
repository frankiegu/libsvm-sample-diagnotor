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
import "sync"
import "runtime"
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

	var featureAppearance = make(map[string]int, 1)
	var featurePositiveAppearance = make(map[string]int, 1)

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
		storageByColumes(fields, featurePositiveAppearance, featureAppearance)
	}

	dumpIndices(dig, group, cfg)
	calulateFeatureMutalInformation(cfg, featurePositiveAppearance, featureAppearance, dig)
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

func storageByColumes(fields []string, featurePositiveAppearance map[string]int, featureAppearance map[string]int) {
	/*
		Stoage sample labels and features by colume.

		input:
			positiveFeatureAppearance        a hash of feature and its appearance times of positive samples
			featureAppearance                number of feature's appearance time.

		outpu:
			fulfilled featureAppearance and featurePositiveAppearance
	*/

	label, err := strconv.ParseInt(fields[0], 10, 32)
	check(err)

	for _, s := range fields[1:] {
		if v, ok := featureAppearance[s]; ok {
			featureAppearance[s] = v + 1
		} else {
			featureAppearance[s] = 1
		}

		if label > 0 {
			if v, ok := featurePositiveAppearance[s]; ok {
				featurePositiveAppearance[s] = v + 1
			} else {
				featurePositiveAppearance[s] = 1
			}
		}
	}
}

type MutalInformation struct {
	featureName string
	mi          float64
}

type FeatureColume struct {
	featureName string
	appearance  int
}

func calulateFeatureMutalInformation(cfg Configuration, featurePositiveAppearance map[string]int, featureAppearance map[string]int, dig DiagnoseIndices) {
	/*
		Calulate the mutal infermation between given feature and sample labels

		input:
			cfg                              configurations
			featurePositiveAppearance        a hash of feature and its appearance times of positive samples
			featureAppearance                number of feature's appearance time.
			dig                              storage of indices

		output:
			a file named "feature.mutal.info" which contains feature-names and
			its mutal information writen in the same line.

	*/

	if cfg.enable_mi == false {
		return
	}

	mutal, err := os.Create("feature.mutal.info")
	check(err)

	var barrier sync.WaitGroup
	runtime.GOMAXPROCS(runtime.NumCPU())

	var miChannel = make(chan MutalInformation, runtime.NumCPU())
	var feChannel = make(chan FeatureColume, runtime.NumCPU())

	var total = float64(dig.rowCnt)

	/* start "cocurrency" number of calculation goroutine */
	cocurrency := runtime.NumCPU() + 2
	for i := 0; i < cocurrency; i++ {
		barrier.Add(1)
		go func() {
			for feature := range feChannel {
				var TT, TF, FT, FF int = 0, 0, 0, 0
				if value, ok := featurePositiveAppearance[feature.featureName]; ok {
					TT = value
				}

				var mi float64 = 0.0
				if TT > 0 {
					mi += float64(TT) * math.Log(total*float64(TT)/float64(feature.appearance*dig.positive))
				}

				if TF = feature.appearance - TT; TF > 0 {
					mi += float64(TF) * math.Log(total*float64(TF)/float64(feature.appearance*dig.negative))
				}

				var noShow int = dig.rowCnt - feature.appearance

				if FT = dig.positive - TT; FT > 0 {
					mi += float64(FT) * math.Log(total*float64(FT)/float64(noShow*dig.positive))
				}

				if FF = dig.negative - TF; FF > 0 {
					mi += float64(FF) * math.Log(total*float64(FF)/float64(noShow*dig.negative))
				}

				miChannel <- MutalInformation{feature.featureName, mi / float64(total) * float64(10000)}
			}
			barrier.Done()
		}()
	}

	/* Just 1 routine provides payloads for calculation */
	barrier.Add(1)
	go func() {
		for k, v := range featureAppearance {
			feChannel <- FeatureColume{k, v}
		}
		close(feChannel)
		barrier.Done()
	}()

	/* another routine spys on calculations routines; and Close miChannel after they all finished. */
	go func() {
		barrier.Wait()
		close(miChannel)
	}()

	/* [main routine] receives calculated mutal information and write into file system */
	for mutalInfo := range miChannel {
		mutal.WriteString(fmt.Sprintf("%s\t%.3f\n", mutalInfo.featureName, mutalInfo.mi))
	}

	mutal.Close()
}
