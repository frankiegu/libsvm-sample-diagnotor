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

import "flag"
import "encoding/json"
import "io/ioutil"
import "strings"
import "strconv"

type Configuration struct {
	thresholds map[string]float32 /* thresholds from configration */
	groupTags  []string           /* options specified by command line */
	enable_mi  bool               /* whether evaluate mutual infermation between feature and label */
	inputFile  string             /* input file name of smaples */
}

func FloatValue(str string) float32 {
	/*
		function of atof, convert string into float

		input:
			str        a string value of number

		output:
			return a float32 converted from <str>

			Notice: Abort on any error.
	*/
	f64, err := strconv.ParseFloat(str, 32)
	check(err)
	return float32(f64)
}

func Specified(cfg Configuration) Configuration {
	/*
		load Config file as a map[key(string)]value(float32)

		input:
			--feature-max       specify feature number upper bound for report
			--feature-min       specify feature number lower bound for report
			--cover-max         specify feature coverage upperbound for report
			--cover-min         specify feature coverage lower bound for report
			--mutual-max        specify mutual-infermation uppper bound for report
			--mutual-min        specify mutual-infermation lower bound for report
			--group-tag         specify tags of feature group. seperated by comma
			--enable-mi         enable mutual-infermation calculation, default is false

		output:
			use any manual specified value instead of value readed from configuration file;

			Notice: All value specified must be in range [0,0, 1.0]
	*/

	var widthMax = FloatValue(*(flag.String("feature-max", "-1", "Threshold for width")))
	var widthMin = FloatValue(*(flag.String("feature-min", "-1", "Threshold for width")))
	var coverMax = FloatValue(*(flag.String("cover-max", "-1", "Threshold for cover")))
	var coverMin = FloatValue(*(flag.String("cover-min", "-1", "Threshold for cover")))
	var mutualMax = FloatValue(*(flag.String("mutual-max", "-1", "Threshold for mutual")))
	var mutualMin = FloatValue(*(flag.String("mutual-min", "-1", "Threshold for mutual")))

	var groupTag = flag.String("group-tag", "", "feature group tags, seperated by comma")

	var enable_mi = flag.Bool("enable-mi", false, "Whether enable mutual infermation evaluation, this is expensive, default is false")

	flag.Parse()

	if widthMax > 0 {
		cfg.thresholds["width_max"] = widthMax
	}

	if widthMin > 0 {
		cfg.thresholds["width_min"] = widthMin
	}

	if coverMax > 0 {
		cfg.thresholds["cover_max"] = coverMax
	}

	if coverMin > 0 {
		cfg.thresholds["cover_min"] = coverMin
	}

	if mutualMax > 0 {
		cfg.thresholds["mutual_max"] = mutualMax
	}

	if mutualMin > 0 {
		cfg.thresholds["mutual_min"] = mutualMin
	}

	if len(*groupTag) > 0 {
		var fields []string = strings.Split(*groupTag, ",")
		for _, tag := range fields {
			cfg.groupTags = append(cfg.groupTags, strings.Trim(tag, " \""))
		}
	}

	if flag.NArg() > 0 {
		cfg.inputFile = flag.Args()[0]
	}

	cfg.enable_mi = *enable_mi
	return cfg
}

func Load(cfgName string) Configuration {
	/*
		load Config file as a map[key(string)]value(float32)

		input:
			cfgName       file name of configuration file

		output:
			return a map[string]float32 which contains all
			key/values loaded from config file

			Notice: Abort on any error!
	*/

	var cfg Configuration = Configuration{}

	contents, err := ioutil.ReadFile(cfgName)
	check(err)
	err = json.Unmarshal(contents, &(cfg.thresholds))
	check(err)

	return Specified(cfg)
}
