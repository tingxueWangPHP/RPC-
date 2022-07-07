package main

import (
	"fmt"
	rpc "rpc/codec"
	"sync"
	"time"
)

type Person struct{}

func (p Person) Say1(a int) int {
	fmt.Println(a)
	return 3
}

func (p Person) Say2(name1, name2, name3 string) string {
	fmt.Println(name1, name2, name3)
	return "test"
}

func (p Person) Say3() Test {
	//fmt.Println(t.Name)
	/*return Test{Name: "zhangsan", Data: []int{1, 2, 3}, Data2: map[string][]int{
		"name": {4, 5, 6},
	}}*/

	return Test{Name: "zhangsan", Data: []int{1, 2, 3}}
}

func (p Person) Say4(m map[string]interface{}) []int {
	//fmt.Println(m)

	/*a := []int{5}

	return map[string][]int{"name": a}*/
	//a := []int{66}

	return []int{1, 2, 3}
}

func (p Person) Say5(m map[string]interface{}) map[string][]int {
	return map[string][]int{
		"zhangsan": {1, 2, 3},
		"lisi":     {9, 9, 9},
	}
}

func (p Person) Say6(m map[string]interface{}) map[string]string {
	return map[string]string{
		"zhangsan": "one",
		"lisi":     "two",
	}
}

func (p Person) Say7(m map[string]interface{}) map[string]map[string]string {
	return map[string]map[string]string{
		"zhangsan": {
			"a": "one",
		},
	}
}

func (p Person) Say8(m map[string]interface{}) []map[int]string {
	return []map[int]string{
		{1: "one"},
		{2: "two"},
	}
}

type Test struct {
	Name string
	Data []int
	//Data2 map[string][]int
	//Test2
}

type Test2 struct {
	age int
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		a := time.Now()
		client := rpc.DialServer("tcp", "127.0.0.1:20000")

		//var ret []int
		//var ret = make(map[string][]int, 0)

		//var ret map[string][]int

		//var ret map[string]string

		//var ret map[string]map[string]string

		//var ret []map[int]string

		//var ret string

		//fmt.Println(client.Call("Person.Say2", "11", "22", "33"))
		//client.Call("Person.Say7", &ret, map[string]interface{}{"name": nil})
		//client.Call("Person.Say2", &ret, "a", "b", "c")

		var ret = new(Test)

		fmt.Println(client.Call("Person.Say3", time.Second*10, ret))

		fmt.Println(ret)

		fmt.Println(time.Since(a))
		//fmt.Println(time.Now().Unix() - a)
	}()

	go func() {
		defer wg.Done()
		rpc.Serv.Register(Person{})

		rpc.ListenAndServe("tcp", "127.0.0.1:20000")
	}()

	wg.Wait()
	fmt.Println("end")
}
