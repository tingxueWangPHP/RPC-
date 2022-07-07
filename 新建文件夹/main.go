package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

type data struct {
	Age  int
	Name string
}

type header struct {
	Method string
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	defer listen.Close()

	go func() {
		client()
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error")
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			objHeader := new(header)
			obj := new(data)

			aa := ""

			decoder := gob.NewDecoder(conn)

			/*if err := decoder.Decode(objHeader); err != nil {
				panic(err)
			}
			if err := decoder.Decode(obj); err != nil {
				panic(err)
			}

			fmt.Println(objHeader.Method)
			fmt.Println(obj.Age)
			fmt.Println(obj.Name)*/

			//for i := 0; i < 10000; i++ {

			if err := decoder.Decode(&aa); err != nil {
				panic(err)
			}

			if err := decoder.Decode(objHeader); err != nil {
				panic(err)
			}
			if err := decoder.Decode(obj); err != nil {
				panic(err)
			}

			fmt.Println(objHeader.Method)
			fmt.Println(obj.Age)
			fmt.Println(obj.Name)

			fmt.Println(aa)
			//fmt.Println(i)
			//fmt.Println("---------")
			//}
		}(conn)
	}

}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}

	defer conn.Close() // 关闭连接

	d := data{Age: 30, Name: "zhangsan"}
	f := header{Method: "test.Say"}

	encoder := gob.NewEncoder(conn)

	//encoder.Encode("444455")

	/*if err := encoder.Encode(f); err != nil {
		panic(err)
	}

	if err := encoder.Encode(d); err != nil {
		panic(err)
	}*/

	//for i := 0; i < 10000; i++ {

	if err := encoder.Encode("99999"); err != nil {
		panic(err)
	}

	if err := encoder.Encode(f); err != nil {
		panic(err)
	}

	if err := encoder.Encode(d); err != nil {
		panic(err)
	}

	//}

	//encoder.Encode("444455")

}
