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
			--mutal-max         specify mutal-infermation uppper bound for report
			--mutal-min         specify mutal-infermation lower bound for report
			--group-tag         specify tags of feature group. seperated by comma

		output:
			use any manual specified value instead of value readed from configuration file;

			Notice: All value specified must be in range [0,0, 1.0]
	*/

	var widthMax = FloatValue(*(flag.String("feature-max", "-1", "Threshold for width")))
	var widthMin = FloatValue(*(flag.String("feature-min", "-1", "Threshold for width")))
	var coverMax = FloatValue(*(flag.String("cover-max", "-1", "Threshold for cover")))
	var coverMin = FloatValue(*(flag.String("cover-min", "-1", "Threshold for cover")))
	var mutalMax = FloatValue(*(flag.String("mutal-max", "-1", "Threshold for mutal")))
	var mutalMin = FloatValue(*(flag.String("mutal-min", "-1", "Threshold for mutal")))

	var groupTag = flag.String("group-tag", "", "feature group tags, seperated by comma")

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

	if mutalMax > 0 {
		cfg.thresholds["mutal_max"] = mutalMax
	}

	if mutalMin > 0 {
		cfg.thresholds["mutal_min"] = mutalMin
	}

	if len(*groupTag) > 0 {
		var fields []string = strings.Split(*groupTag, ",")
		for _, tag := range fields {
			cfg.groupTags = append(cfg.groupTags, strings.Trim(tag, " \""))
		}
	}

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

	var cfg Configuration =  Configuration{}

	contents, err := ioutil.ReadFile(cfgName)
	check(err)
	err = json.Unmarshal(contents, &(cfg.thresholds))
	check(err)

	return Specified(cfg)
}
