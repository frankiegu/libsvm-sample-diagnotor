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

func check(value interface{}) {
	/*
		This is an assert function;

		input:
			<value>        value can be bool, error or string

		output:
			check will panic when value is false(bool) or error(not nil) or string(not empty),
			however NOTHING will happen when value is nil

	*/

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
		panic(v)
	}
}

func main() {
	var config = Load("thresholds.json")
	basic_statistics(config)
}
