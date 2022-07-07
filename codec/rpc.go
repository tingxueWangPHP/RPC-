package codec

import (
	"context"
	"errors"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

type DataFormat struct {
	Header string
	Body   []interface{}
}

type (
	Server struct{}

	Client struct {
		conn net.Conn
	}
)

var (
	relation = map[string]interface{}{}
	Serv     *Server

	lock sync.Mutex
)

func (s *Server) Register(structItem interface{}) {
	lock.Lock()
	defer lock.Unlock()
	relation[reflect.TypeOf(structItem).Name()] = structItem
}

func (s *Server) parseHeader(header string) (string, string) {
	ss := strings.Split(header, ".")
	return ss[0], ss[1]
}

func ListenAndServe(protocol, address string) error {
	listen, err := net.Listen(protocol, address)
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()

			data, _ := Decode(conn)
			structName, methodName := Serv.parseHeader(data.Header)
			fn := reflect.ValueOf(relation[structName]).MethodByName(methodName)

			//构建参数
			tempParam := []reflect.Value{}
			for i := 0; i < len(data.Body); i++ {
				tempParam = append(tempParam, reflect.ValueOf(data.Body[i]))
			}
			values := fn.Call(tempParam)

			d := DataFormat{
				Header: "reply",
			}
			for _, item := range values {
				d.Body = append(d.Body, item.Interface())
			}

			Encode(conn, d)
		}(conn)
	}

	return nil
}

func DialServer(protocol, address string) (*Client, error) {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil

}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Call(cli *Client, method string, timeout time.Duration, ret interface{}, args ...interface{}) error {
	if reflect.TypeOf(ret).Kind() != reflect.Ptr {
		return errors.New("type error")
	}

	//防止协程因为没有消费者而阻塞 产生泄露
	resCh := make(chan *DataFormat, 1)
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	go func() {
		d := DataFormat{
			Header: method,
			Body:   args,
		}
		Encode(cli.conn, d)
		res, _ := Decode(cli.conn)
		resCh <- res
	}()

	select {
	case res := <-resCh:
		for _, item := range res.Body {
			if err := recursionGet2(item, ret); err != nil {
				return err
			}
		}
		close(resCh)
	case <-ctx.Done():
		return errors.New("time out")
	}

	return nil
}

func recursionGet2(src, des interface{}) error {
	var (
		desElem  = reflect.ValueOf(des).Elem()
		srcValue = reflect.ValueOf(src)
		srcType  = reflect.TypeOf(src)

		result reflect.Value
	)

	//判断是否是指针
	if !desElem.CanSet() {
		return errors.New("cannot set")
	}

	switch srcValue.Kind() {
	case reflect.Slice:
		newSlice := reflect.MakeSlice(srcType, srcValue.Len(), srcValue.Cap())
		for i := 0; i < newSlice.Len(); i++ {
			p := reflect.New(srcValue.Index(i).Type())
			recursionGet2(srcValue.Index(i).Interface(), p.Interface())
			newSlice.Index(i).Set(p.Elem())
		}
		result = newSlice
	case reflect.Map:
		newMap := reflect.MakeMap(srcType)
		for _, v := range srcValue.MapKeys() {
			p := reflect.New(srcValue.MapIndex(v).Type())
			recursionGet2(srcValue.MapIndex(v).Interface(), p.Interface())
			newMap.SetMapIndex(v, p.Elem())
		}
		result = newMap
	default:
		result = srcValue
	}

	desElem.Set(result)

	return nil
}
