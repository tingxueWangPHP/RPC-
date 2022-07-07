package main

import (
	"fmt"
	"reflect"
)

type data struct {
	age  int
	name string
}

func main() {
	run(test, 1, "abc")
	run(test2, 8)
	run(test3, data{80, "liming"})
}

func test(a int, b string) error {
	fmt.Println(a)
	fmt.Println(b)

	return nil
}

func test2(a int) {
	fmt.Println(a)
}

func test3(d data) {
	fmt.Println(d.name)
}

func run(method interface{}, params ...interface{}) {

	v := reflect.ValueOf(method)

	/*t := reflect.TypeOf(method)
	fmt.Println(t.NumMethod())*/

	s := make([]reflect.Value, len(params))

	for k, item := range params {
		s[k] = reflect.ValueOf(item)
	}

	v.Call(s)
}
