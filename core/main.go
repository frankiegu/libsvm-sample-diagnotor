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

import "encoding/json"
import ioutil "io/ioutil"

func check(value interface{}) {
	if value == nil {
		return
	}

	switch v := value.(type) {
	case error:
		e, _ := value.(error)
		if e != nil {
			panic(e)
		}
	case bool:
		b, _ := value.(bool)
		if b == false {
			panic("Abort!")
		}
	case string:
		s, _ := value.(string)
		if len(s) > 0 {
			panic(s)
		}
	default:
		/* fmt.Println(value); */
		panic(v)
	}
}

/* load 配置文件, 返回包含所有值的map */
func load(cfg_name string) map[string]interface{} {
	cfg := make(map[string]interface{}, 1)
	contents, err := ioutil.ReadFile(cfg_name)
	check(err)
	err = json.Unmarshal(contents, &cfg)
	check(err)
	return cfg
}

func main() {
	diagnose(load("./thresholds.json"))
}
