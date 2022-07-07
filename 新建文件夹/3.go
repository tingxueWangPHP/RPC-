package main

import (
	"fmt"
	"net"
)

/*ype data struct {
	Age  int
	Name string
}

type header struct {
	Method string
}*/

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

			buf := make([]byte, 5)
			for {
				if err := forRead(buf, conn); err != nil {
					//fmt.Println(err)
					break
				} else {
					fmt.Println(string(buf))
				}
			}
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

	str := []byte("12345")

	for i := 0; i < 10; i++ {
		if err := forWrite(str, conn); err != nil {
			fmt.Println(err)
			break
		}
	}

}

func forWrite(str []byte, conn net.Conn) error {
	if n, err := conn.Write(str); err != nil {
		return err
	} else if n < len(str) {
		return forWrite(str[n:], conn)
	} else {
		return nil
	}
}

func forRead(buf []byte, conn net.Conn) error {
	if n, err := conn.Read(buf); err != nil {
		return err
	} else if n < len(buf) {
		return forRead(buf[n:], conn)
	} else {
		return nil
	}
}
