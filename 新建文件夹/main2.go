package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", "127.0.0.1:20000")
		if err != nil {
			fmt.Println("err :", err)
			return
		}

		defer conn.Close() // 关闭连接

		var args = []interface{}{}
		header := "Person.Say"
		body := 33
		args = append(args, body)

		encoder := gob.NewEncoder(conn)

		if err := encoder.Encode(header); err != nil {
			panic(err)
		}

		if err := encoder.Encode(args); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		listen, err := net.Listen("tcp", "127.0.0.1:20000")
		if err != nil {
			fmt.Println("listen failed, err:", err)
			return
		}

		defer listen.Close()

		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("accept error")
				continue
			}

			go func(conn net.Conn) {
				defer conn.Close()
				register(Person{})
				decoder := gob.NewDecoder(conn)

				var (
					header string
				)

				if err := decoder.Decode(&header); err != nil {
					panic(err)
				}

				name1, name2 := parseHeader(header)

				s := relation[name1]

				t := reflect.ValueOf(s)

				fn := t.MethodByName(name2)

				temp := []interface{}{}

				if err := decoder.Decode(&temp); err != nil {
					panic(err)
				}

				res := []reflect.Value{}

				for i := 0; i < len(temp); i++ {
					res = append(res, reflect.ValueOf(temp[i]))
				}

				fn.Call(res)

			}(conn)
		}
	}()

	wg.Wait()
}

type Person struct{}

func (p Person) Say(a int) {
	fmt.Println(a)
}

func (p Person) Say2() {
	fmt.Println("say2")
}

var relation = map[string]interface{}{}

func register(param interface{}) {
	relation[reflect.TypeOf(param).Name()] = param
}

func parseHeader(header string) (string, string) {
	s := strings.Split(header, ".")
	return s[0], s[1]
}
