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
	var (
		wg       = sync.WaitGroup{}
		wgServer = sync.WaitGroup{}

		serverList = []string{}

		lock sync.Mutex

		ch = make(chan struct{})
	)
	wg.Add(3)
	wgServer.Add(2)

	go func() {
		defer wg.Done()
		//a := time.Now()
		<-ch
		servers := rpc.NewMultiServersDiscovery()
		servers.Update(serverList)

		wg2 := sync.WaitGroup{}

		for i := 0; i < 10; i++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				addr, _ := servers.Get(rpc.RandomSelect)
				fmt.Println(addr)
				client, _ := rpc.DialServer("tcp", addr)
				defer func() {
					client.Close()
				}()
				var ret = new(Test)
				client.Call(client, "Person.Say3", time.Second*10, ret)
				//fmt.Println(ret)
			}()
		}

		wg2.Wait()

		/*addr, _ := servers.Get(rpc.RoundRobinSelect)

		client, _ := rpc.DialServer("tcp", addr)

		defer func() {
			client.Conn.Close()
		}()*/

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

		/*var ret = new(Test)

		fmt.Println(client.Call(client, "Person.Say3", time.Second*10, ret))

		fmt.Println(ret)

		fmt.Println(time.Since(a))*/
		//fmt.Println(time.Now().Unix() - a)
	}()

	go func() {
		defer func() {
			wg.Done()
		}()
		rpc.Serv.Register(Person{})

		var address string = "127.0.0.1:20000"
		lock.Lock()
		serverList = append(serverList, address)
		lock.Unlock()
		wgServer.Done()
		rpc.ListenAndServe("tcp", address)
	}()

	go func() {
		defer func() {
			wg.Done()
		}()
		rpc.Serv.Register(Person{})

		var address string = "127.0.0.1:20001"
		lock.Lock()
		serverList = append(serverList, address)
		lock.Unlock()
		wgServer.Done()
		rpc.ListenAndServe("tcp", address)
	}()

	wgServer.Wait()
	ch <- struct{}{}
	wg.Wait()
	fmt.Println("end")
}
